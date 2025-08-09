package musicbrainz

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"relay/app/cache"
	"relay/app/music"
	"relay/app/music/meilisearch"

	_ "github.com/lib/pq"
)

var db *sql.DB
var spotifyClient *music.SpotifyClient
var lyricsClient *music.LRCLibClient
var mediaDir string

// Meilisearch-related variables and functions that reference the meilisearch package
var meilisearchClient interface{} // This will be set by InitMeilisearch

// IsReady checks if both MusicBrainz and Meilisearch are configured and ready
func IsReady() bool {
	return meilisearch.IsReady()
}

// InitMeilisearch initializes the Meilisearch client and creates indexes
func InitMeilisearch(host, port, apiKey string, timeout time.Duration) error {
	err := meilisearch.InitMeilisearch(host, port, apiKey, timeout)
	if err == nil {
		meilisearchClient = true // Just a marker that it's initialized
	}
	return err
}

// ApplyTunables sets the sync tunables for Meilisearch operations
func ApplyTunables(t meilisearch.SyncTunables) {
	meilisearch.ApplyTunables(t)
}

// Index functions that delegate to the meilisearch package
func IndexArtists() error {
	return meilisearch.IndexArtists()
}

func IndexReleaseGroups() error {
	return meilisearch.IndexReleaseGroups()
}

func IndexReleases() error {
	return meilisearch.IndexReleases()
}

func IndexRecordings() error {
	return meilisearch.IndexRecordings()
}

// Search functions that delegate to the meilisearch package
func SearchArtistsMeilisearch(ctx context.Context, query string, limit int) (any, error) {
	return meilisearch.SearchArtistsMeilisearch(ctx, query, limit)
}

func SearchReleaseGroupsMeilisearch(ctx context.Context, query string, limit int) (any, error) {
	return meilisearch.SearchReleaseGroupsMeilisearch(ctx, query, limit)
}

func SearchReleasesMeilisearch(ctx context.Context, query string, limit int) (any, error) {
	return meilisearch.SearchReleasesMeilisearch(ctx, query, limit)
}

func SearchRecordingsMeilisearch(ctx context.Context, query string, limit int) (any, error) {
	return meilisearch.SearchRecordingsMeilisearch(ctx, query, limit)
}

func GetArtistMeilisearch(ctx context.Context, mbid string) (any, error) {
	return meilisearch.GetArtistMeilisearch(ctx, mbid)
}

// SetSpotifyClient injects a Spotify client used to fetch images stored on disk.
func SetSpotifyClient(c *music.SpotifyClient) { spotifyClient = c }

// SetLyricsClient injects an LRCLib client used to fetch lyrics stored on disk.
func SetLyricsClient(c *music.LRCLibClient) { lyricsClient = c }

// SetMediaDir configures on-disk media directory used for images and lyrics.
func SetMediaDir(dir string) { mediaDir = dir }

var httpClient = &http.Client{Timeout: 6 * time.Second}

// Database structures for URL relationships
type ArtistURLRelation struct {
	URL      string `db:"url"`
	LinkType string `db:"link_type"`
	TypeName string `db:"type_name"`
	Ended    bool   `db:"ended"`
}

// WikidataEntity represents a Wikidata entity response
type WikidataEntity struct {
	Entities map[string]struct {
		SiteLinks map[string]struct {
			Site  string `json:"site"`
			Title string `json:"title"`
			URL   string `json:"url"`
		} `json:"sitelinks"`
	} `json:"entities"`
}

var wikidataIDRegex = regexp.MustCompile(`/wiki/(Q\d+)$`)

// fetchArtistURLsFromDB retrieves URL relationships for an artist from PostgreSQL database
func fetchArtistURLsFromDB(ctx context.Context, mbid string) ([]ArtistURLRelation, error) {
	if mbid == "" {
		return nil, fmt.Errorf("empty MBID")
	}

	query := `
		SELECT 
			u.url,
			lt.gid as link_type,
			lt.name as type_name,
			lar.ended
		FROM artist a
		JOIN l_artist_url lar ON a.id = lar.entity0
		JOIN url u ON lar.entity1 = u.id
		JOIN link l ON lar.link = l.id
		JOIN link_type lt ON l.link_type = lt.id
		WHERE a.gid = $1
		ORDER BY lt.name
	`

	rows, err := db.QueryContext(ctx, query, mbid)
	if err != nil {
		return nil, fmt.Errorf("failed to query URL relationships: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var relations []ArtistURLRelation
	for rows.Next() {
		var rel ArtistURLRelation
		err := rows.Scan(&rel.URL, &rel.LinkType, &rel.TypeName, &rel.Ended)
		if err != nil {
			return nil, fmt.Errorf("failed to scan URL relationship: %w", err)
		}
		relations = append(relations, rel)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating URL relationships: %w", err)
	}

	return relations, nil
}

// extractWikidataID extracts Wikidata ID from URL relationships
func extractWikidataID(relations []ArtistURLRelation) string {
	for _, relation := range relations {
		if relation.TypeName == "wikidata" && !relation.Ended {
			matches := wikidataIDRegex.FindStringSubmatch(relation.URL)
			if len(matches) > 1 {
				return matches[1]
			}
		}
	}
	return ""
}

// fetchWikipediaPageFromWikidata retrieves Wikipedia page title from Wikidata ID
func fetchWikipediaPageFromWikidata(ctx context.Context, wikidataID string) (string, error) {
	if wikidataID == "" {
		return "", fmt.Errorf("empty Wikidata ID")
	}

	endpoint := fmt.Sprintf("https://www.wikidata.org/wiki/Special:EntityData/%s.json", wikidataID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "MetadataRelay/1.0 (https://github.com/your-org/metadata-relay)")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("wikidata API returned status %d", resp.StatusCode)
	}

	var entity WikidataEntity
	if err := json.NewDecoder(resp.Body).Decode(&entity); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Look for English Wikipedia sitelink
	if entityData, ok := entity.Entities[wikidataID]; ok {
		if enwiki, ok := entityData.SiteLinks["enwiki"]; ok {
			return enwiki.Title, nil
		}
	}

	return "", fmt.Errorf("no English Wikipedia page found")
}

// fetchWikipediaFromMusicBrainz fetches Wikipedia info via PostgreSQL URL relationships
func fetchWikipediaFromMusicBrainz(ctx context.Context, mbid string) (summary string, pageURL string, err error) {
	result, err := cache.NewCache("musicbrainz_wikipedia").TTL(7*24*time.Hour).Wrap(func() (any, error) {
		// Get URL relationships from database
		relations, err := fetchArtistURLsFromDB(ctx, mbid)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch URL relationships: %w", err)
		}

		// Extract Wikidata ID
		wikidataID := extractWikidataID(relations)
		if wikidataID == "" {
			return nil, fmt.Errorf("no Wikidata link found")
		}

		// Get Wikipedia page title from Wikidata
		pageTitle, err := fetchWikipediaPageFromWikidata(ctx, wikidataID)
		if err != nil {
			return nil, fmt.Errorf("failed to get Wikipedia page from Wikidata: %w", err)
		}

		// Fetch Wikipedia summary
		sum, link, err := fetchWikipediaSummary(ctx, pageTitle)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch Wikipedia summary: %w", err)
		}

		return map[string]string{
			"summary": sum,
			"url":     link,
		}, nil
	})(ctx, mbid)

	if err != nil {
		return "", "", err
	}

	data := result.(map[string]string)
	return data["summary"], data["url"], nil
}

// fetchWikipediaSummary best-effort retrieves a summary and canonical URL for a title from Wikipedia.
func fetchWikipediaSummary(ctx context.Context, title string) (summary string, pageURL string, err error) {
	if title == "" {
		return "", "", fmt.Errorf("empty title")
	}
	t := strings.ReplaceAll(title, " ", "_")
	endpoint := "https://en.wikipedia.org/api/rest_v1/page/summary/" + url.PathEscape(t)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	req.Header.Set("Accept", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("wikipedia http %d", resp.StatusCode)
	}
	var out struct {
		Extract     string `json:"extract"`
		ContentURLs struct {
			Desktop struct {
				Page string `json:"page"`
			} `json:"desktop"`
		} `json:"content_urls"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", "", err
	}
	return out.Extract, out.ContentURLs.Desktop.Page, nil
}

// InitMusicBrainz initializes the MusicBrainz PostgreSQL connection with proper
// connection pooling and timeout settings.
func InitMusicBrainz(connStr string) {
	// Check if required connection parameters are missing
	if connStr == "host= port=5432 user=musicbrainz password=musicbrainz dbname= sslmode=disable" {
		fmt.Printf("WARNING: MUSICBRAINZ_DB_HOST and MUSICBRAINZ_DB_NAME environment variables are not set.\n")
		fmt.Printf("Music endpoints will not be available.\n")
		return
	}

	var err error

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		slog.Error("failed to connect to MusicBrainz database", "error", err)
		return
	}

	// Set the database connection for the meilisearch package
	meilisearch.SetDB(db)

	// Test the connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		slog.Error("failed to ping MusicBrainz database", "error", err)
		return
	}

	// Set connection pool settings for optimal performance
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	slog.Info("Connected to MusicBrainz PostgreSQL database")
}

// GetArtist gets a specific artist by MBID
func GetArtist(ctx context.Context, mbid string) (any, error) {
	return cache.NewCache("musicbrainz_artist").TTL(7*24*time.Hour).Wrap(func() (any, error) {
		query := `
			SELECT 
				a.gid,
				a.name,
				a.sort_name,
				a.type as artist_type,
				ac.name as area_name,
				a.begin_date_year,
				a.end_date_year,
				a.ended,
				a.comment
			FROM artist a
			LEFT JOIN area ac ON a.area = ac.id
			WHERE a.gid = $1
		`

		row := db.QueryRowContext(ctx, query, mbid)

		var gid, name, sortName, artistType, areaName, comment sql.NullString
		var beginYear, endYear sql.NullInt32
		var ended sql.NullBool

		err := row.Scan(&gid, &name, &sortName, &artistType, &areaName, &beginYear, &endYear, &ended, &comment)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("artist not found")
			}
			return nil, fmt.Errorf("failed to get artist: %w", err)
		}

		artist := map[string]any{
			"id":             gid.String,
			"name":           name.String,
			"sort-name":      sortName.String,
			"type":           artistType.String,
			"area":           areaName.String,
			"ended":          ended.Bool,
			"disambiguation": comment.String,
		}

		if beginYear.Valid {
			artist["life-span"] = map[string]any{
				"begin": fmt.Sprintf("%d", beginYear.Int32),
			}
			if endYear.Valid {
				artist["life-span"].(map[string]any)["end"] = fmt.Sprintf("%d", endYear.Int32)
			}
		}

		// Wikipedia summary via MusicBrainz URL relationships (cached)
		if gid.String != "" {
			if sum, link, err := fetchWikipediaFromMusicBrainz(ctx, gid.String); err == nil {
				if sum != "" {
					artist["wikipedia-summary"] = sum
				}
				if link != "" {
					artist["wikipedia-url"] = link
				}
			} else {
				// Fallback to name-based Wikipedia search
				if name.String != "" {
					if sum, link, err := fetchWikipediaSummary(ctx, name.String); err == nil {
						if sum != "" {
							artist["wikipedia-summary"] = sum
						}
						if link != "" {
							artist["wikipedia-url"] = link
						}
					}
				}
			}
		}

		// Spotify image (stored on disk under /media)
		if spotifyClient != nil && mediaDir != "" && name.String != "" {
			if path, err := spotifyClient.DownloadArtistImage(ctx, name.String, mediaDir); err == nil && path != "" {
				artist["image-url"] = "/media/spotify/artists/" + url.PathEscape(name.String) + ".jpg"
			}
		}

		return artist, nil
	})(ctx, mbid)
}

// GetReleaseGroup gets a specific release group by MBID
func GetReleaseGroup(ctx context.Context, mbid string) (any, error) {
	return cache.NewCache("musicbrainz_release_group").TTL(7*24*time.Hour).Wrap(func() (any, error) {
		query := `
			SELECT 
				rg.gid,
				rg.name,
				rg.type as release_group_type,
				a.name as artist_name,
				a.gid as artist_id,
				rg.comment
			FROM release_group rg
			JOIN artist_credit_name acn ON rg.artist_credit = acn.artist_credit
			JOIN artist a ON acn.artist = a.id
			WHERE rg.gid = $1
		`

		row := db.QueryRowContext(ctx, query, mbid)

		var rgGid, rgName, rgType, artistName, artistGid, comment sql.NullString

		err := row.Scan(&rgGid, &rgName, &rgType, &artistName, &artistGid, &comment)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("release group not found")
			}
			return nil, fmt.Errorf("failed to get release group: %w", err)
		}

		releaseGroup := map[string]any{
			"id":           rgGid.String,
			"title":        rgName.String,
			"primary-type": rgType.String,
			"artist-credit": []map[string]any{
				{
					"name": artistName.String,
					"artist": map[string]any{
						"id":   artistGid.String,
						"name": artistName.String,
					},
				},
			},
			"disambiguation": comment.String,
		}

		return releaseGroup, nil
	})(ctx, mbid)
}

// GetRelease gets a specific release by MBID
func GetRelease(ctx context.Context, mbid string) (any, error) {
	return cache.NewCache("musicbrainz_release").TTL(7*24*time.Hour).Wrap(func() (any, error) {
		query := `
			SELECT 
				r.gid,
				r.name,
				r.status,
				a.name as artist_name,
				a.gid as artist_id,
				rg.name as release_group_name,
				rg.gid as release_group_id,
				r.comment
			FROM release r
			JOIN artist_credit_name acn ON r.artist_credit = acn.artist_credit
			JOIN artist a ON acn.artist = a.id
			LEFT JOIN release_group rg ON r.release_group = rg.id
			WHERE r.gid = $1
		`

		row := db.QueryRowContext(ctx, query, mbid)

		var releaseGid, releaseName, releaseStatus, artistName, artistGid, rgName, rgGid, comment sql.NullString

		err := row.Scan(&releaseGid, &releaseName, &releaseStatus, &artistName, &artistGid, &rgName, &rgGid, &comment)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("release not found")
			}
			return nil, fmt.Errorf("failed to get release: %w", err)
		}

		release := map[string]any{
			"id":     releaseGid.String,
			"title":  releaseName.String,
			"status": releaseStatus.String,
			"artist-credit": []map[string]any{
				{
					"name": artistName.String,
					"artist": map[string]any{
						"id":   artistGid.String,
						"name": artistName.String,
					},
				},
			},
			"release-group": map[string]any{
				"id":    rgGid.String,
				"title": rgName.String,
			},
			"disambiguation": comment.String,
		}

		return release, nil
	})(ctx, mbid)
}

// GetRecording gets a specific recording by MBID
func GetRecording(ctx context.Context, mbid string) (any, error) {
	return cache.NewCache("musicbrainz_recording").TTL(7*24*time.Hour).Wrap(func() (any, error) {
		query := `
			SELECT 
				rec.gid,
				rec.name,
				rec.length,
				a.name as artist_name,
				a.gid as artist_id,
				rec.comment
			FROM recording rec
			JOIN artist_credit_name acn ON rec.artist_credit = acn.artist_credit
			JOIN artist a ON acn.artist = a.id
			WHERE rec.gid = $1
		`

		row := db.QueryRowContext(ctx, query, mbid)

		var recGid, recName, artistName, artistGid, comment sql.NullString
		var length sql.NullInt64

		err := row.Scan(&recGid, &recName, &length, &artistName, &artistGid, &comment)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("recording not found")
			}
			return nil, fmt.Errorf("failed to get recording: %w", err)
		}

		recording := map[string]any{
			"id":    recGid.String,
			"title": recName.String,
			"artist-credit": []map[string]any{
				{
					"name": artistName.String,
					"artist": map[string]any{
						"id":   artistGid.String,
						"name": artistName.String,
					},
				},
			},
			"disambiguation": comment.String,
		}

		if length.Valid {
			recording["length"] = length.Int64
		}

		// Lyrics via LRCLib (stored on disk)
		if lyricsClient != nil && mediaDir != "" && artistName.String != "" && recName.String != "" {
			if path, err := lyricsClient.FetchLyrics(ctx, artistName.String, recName.String, mediaDir); err == nil && path != "" {
				// public URL pointing to static media
				recording["lyrics-url"] = "/media/lyrics/" + url.PathEscape(artistName.String+" - "+recName.String) + ".lrc"
			}
		}

		return recording, nil
	})(ctx, mbid)
}

// BrowseArtistReleaseGroups browses release groups for a specific artist
func BrowseArtistReleaseGroups(ctx context.Context, artistMbid string, limit int) (any, error) {
	return cache.NewCache("musicbrainz_artist_release_groups").TTL(24*time.Hour).Wrap(func() (any, error) {
		query := `
			SELECT DISTINCT
				rg.gid,
				rg.name,
				rg.type as release_group_type,
				a.name as artist_name,
				a.gid as artist_id,
				rg.comment
			FROM release_group rg
			JOIN artist_credit_name acn ON rg.artist_credit = acn.artist_credit
			JOIN artist a ON acn.artist = a.id
			WHERE a.gid = $1
			ORDER BY rg.name
			LIMIT $2
		`

		rows, err := db.QueryContext(ctx, query, artistMbid, limit)
		if err != nil {
			return nil, fmt.Errorf("failed to browse artist release groups: %w", err)
		}
		defer func() {
			if closeErr := rows.Close(); closeErr != nil {
				slog.Warn("failed to close rows", "error", closeErr)
			}
		}()

		var releaseGroups []map[string]any
		for rows.Next() {
			var rgGid, rgName, rgType, artistName, artistGid, comment sql.NullString

			err := rows.Scan(&rgGid, &rgName, &rgType, &artistName, &artistGid, &comment)
			if err != nil {
				continue
			}

			releaseGroup := map[string]any{
				"id":           rgGid.String,
				"title":        rgName.String,
				"primary-type": rgType.String,
				"artist-credit": []map[string]any{
					{
						"name": artistName.String,
						"artist": map[string]any{
							"id":   artistGid.String,
							"name": artistName.String,
						},
					},
				},
				"disambiguation": comment.String,
			}

			releaseGroups = append(releaseGroups, releaseGroup)
		}

		return map[string]any{
			"release-groups": releaseGroups,
			"count":          len(releaseGroups),
		}, nil
	})(ctx, artistMbid, limit)
}

// BrowseReleaseGroupReleases browses releases for a specific release group
func BrowseReleaseGroupReleases(ctx context.Context, releaseGroupMbid string, limit int) (any, error) {
	return cache.NewCache("musicbrainz_release_group_releases").TTL(24*time.Hour).Wrap(func() (any, error) {
		query := `
			SELECT DISTINCT
				r.gid,
				r.name,
				r.status,
				a.name as artist_name,
				a.gid as artist_id,
				rg.name as release_group_name,
				rg.gid as release_group_id,
				r.comment
			FROM release r
			JOIN artist_credit_name acn ON r.artist_credit = acn.artist_credit
			JOIN artist a ON acn.artist = a.id
			JOIN release_group rg ON r.release_group = rg.id
			WHERE rg.gid = $1
			ORDER BY r.name
			LIMIT $2
		`

		rows, err := db.QueryContext(ctx, query, releaseGroupMbid, limit)
		if err != nil {
			return nil, fmt.Errorf("failed to browse release group releases: %w", err)
		}
		defer func() {
			if closeErr := rows.Close(); closeErr != nil {
				slog.Warn("failed to close rows", "error", closeErr)
			}
		}()

		var releases []map[string]any
		for rows.Next() {
			var releaseGid, releaseName, releaseStatus, artistName, artistGid, rgName, rgGid, comment sql.NullString

			err := rows.Scan(&releaseGid, &releaseName, &releaseStatus, &artistName, &artistGid, &rgName, &rgGid, &comment)
			if err != nil {
				continue
			}

			release := map[string]any{
				"id":     releaseGid.String,
				"title":  releaseName.String,
				"status": releaseStatus.String,
				"artist-credit": []map[string]any{
					{
						"name": artistName.String,
						"artist": map[string]any{
							"id":   artistGid.String,
							"name": artistName.String,
						},
					},
				},
				"release-group": map[string]any{
					"id":    rgGid.String,
					"title": rgName.String,
				},
				"disambiguation": comment.String,
			}

			releases = append(releases, release)
		}

		return map[string]any{
			"releases": releases,
			"count":    len(releases),
		}, nil
	})(ctx, releaseGroupMbid, limit)
}

// AdvancedSearchArtists performs an advanced artist search with field-specific queries
func AdvancedSearchArtists(ctx context.Context, artistName, area, beginDate, endDate string, limit int) (any, error) {
	cacheKey := fmt.Sprintf("musicbrainz_artist_advanced_search_%s_%s_%s_%s_%d",
		artistName, area, beginDate, endDate, limit)

	return cache.NewCache(cacheKey).TTL(24*time.Hour).Wrap(func() (any, error) {
		var whereClauses []string
		var args []any
		argIndex := 1

		if artistName != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("(to_tsvector('english', COALESCE(a.name, '')) @@ plainto_tsquery('english', $%d) OR to_tsvector('english', COALESCE(a.sort_name, '')) @@ plainto_tsquery('english', $%d))", argIndex, argIndex))
			args = append(args, artistName)
			argIndex++
		}

		if area != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("ac.name ILIKE $%d", argIndex))
			args = append(args, "%"+area+"%")
			argIndex++
		}

		if beginDate != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("a.begin_date_year = $%d", argIndex))
			args = append(args, beginDate)
			argIndex++
		}

		if endDate != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("a.end_date_year = $%d", argIndex))
			args = append(args, endDate)
			argIndex++
		}

		if len(whereClauses) == 0 {
			return map[string]any{"artists": []map[string]any{}, "count": 0}, nil
		}

		query := fmt.Sprintf(`
			SELECT DISTINCT
				a.gid,
				a.name,
				a.sort_name,
				a.type as artist_type,
				ac.name as area_name,
				a.begin_date_year,
				a.end_date_year,
				a.ended,
				a.comment
			FROM artist a
			LEFT JOIN area ac ON a.area = ac.id
			WHERE %s
			ORDER BY a.name
			LIMIT $%d
		`, strings.Join(whereClauses, " AND "), argIndex)

		args = append(args, limit)

		rows, err := db.QueryContext(ctx, query, args...)
		if err != nil {
			return nil, fmt.Errorf("failed to perform advanced artist search: %w", err)
		}
		defer func() {
			if closeErr := rows.Close(); closeErr != nil {
				slog.Warn("failed to close rows", "error", closeErr)
			}
		}()

		var artists []map[string]any
		for rows.Next() {
			var gid, name, sortName, artistType, areaName, comment sql.NullString
			var beginYear, endYear sql.NullInt32
			var ended sql.NullBool

			err := rows.Scan(&gid, &name, &sortName, &artistType, &areaName, &beginYear, &endYear, &ended, &comment)
			if err != nil {
				continue
			}

			artist := map[string]any{
				"id":             gid.String,
				"name":           name.String,
				"sort-name":      sortName.String,
				"type":           artistType.String,
				"area":           areaName.String,
				"ended":          ended.Bool,
				"disambiguation": comment.String,
			}

			if beginYear.Valid {
				artist["life-span"] = map[string]any{
					"begin": fmt.Sprintf("%d", beginYear.Int32),
				}
				if endYear.Valid {
					artist["life-span"].(map[string]any)["end"] = fmt.Sprintf("%d", endYear.Int32)
				}
			}

			artists = append(artists, artist)
		}

		return map[string]any{
			"artists": artists,
			"count":   len(artists),
		}, nil
	})(ctx, artistName, area, beginDate, endDate, limit)
}
