package musicbrainz

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"relay/app/cache"

	_ "github.com/lib/pq"
)

var db *sql.DB

// InitMusicBrainz initializes the MusicBrainz PostgreSQL connection
func InitMusicBrainz() {
	var err error

	// Get configuration from environment variables with defaults
	host := getEnvOrDefault("MUSICBRAINZ_DB_HOST", "192.168.10.202")
	port := getEnvOrDefault("MUSICBRAINZ_DB_PORT", "5432")
	user := getEnvOrDefault("MUSICBRAINZ_DB_USER", "musicbrainz")
	password := getEnvOrDefault("MUSICBRAINZ_DB_PASSWORD", "musicbrainz")
	dbname := getEnvOrDefault("MUSICBRAINZ_DB_NAME", "musicbrainz_db")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		slog.Error("failed to connect to MusicBrainz database", "error", err)
		return
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		slog.Error("failed to ping MusicBrainz database", "error", err)
		return
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	fmt.Printf("Connected to MusicBrainz PostgreSQL database at %s:%s\n", host, port)
}

// getEnvOrDefault returns the environment variable value or a default value if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} // SearchArtists searches for artists using PostgreSQL full-text search
func SearchArtists(ctx context.Context, query string, limit int) (any, error) {
	return cache.NewCache("musicbrainz_artist_search").TTL(24*time.Hour).Wrap(func() (any, error) {
		// Simplified search query with better performance
		searchQuery := `
			SELECT DISTINCT
				a.gid,
				a.name,
				a.sort_name,
				a.type as artist_type,
				ac.name as area_name,
				a.begin_date_year,
				a.end_date_year,
				a.ended,
				a.comment,
				CASE 
					WHEN a.name ILIKE $1 || '%' THEN 1
					WHEN a.name ILIKE '%' || $1 || '%' THEN 2
					WHEN a.sort_name ILIKE $1 || '%' THEN 3
					WHEN aa.name ILIKE $1 || '%' THEN 4
					ELSE 5
				END as sort_order
			FROM artist a
			LEFT JOIN area ac ON a.area = ac.id
			LEFT JOIN artist_alias aa ON aa.artist = a.id
			WHERE (
				a.name ILIKE '%' || $1 || '%' OR
				a.sort_name ILIKE '%' || $1 || '%' OR
				aa.name ILIKE '%' || $1 || '%'
			)
			ORDER BY sort_order, a.name
			LIMIT $2
		`

		rows, err := db.QueryContext(ctx, searchQuery, query, limit)
		if err != nil {
			return nil, fmt.Errorf("failed to search artists: %w", err)
		}
		defer rows.Close()

		var artists []map[string]any
		for rows.Next() {
			var gid, name, sortName, artistType, areaName, comment sql.NullString
			var beginYear, endYear sql.NullInt32
			var ended sql.NullBool
			var sortOrder int

			err := rows.Scan(&gid, &name, &sortName, &artistType, &areaName, &beginYear, &endYear, &ended, &comment, &sortOrder)
			if err != nil {
				continue
			}

			artist := map[string]any{
				"id":             gid.String,
				"name":           name.String,
				"sort-name":      sortName.String,
				"type":           artistType.String,
				"area":           areaName.String,
				"begin-area":     areaName.String,
				"ended":          ended.Bool,
				"disambiguation": comment.String,
				"score":          100, // Static score for simplified query
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
	})(ctx, query, limit)
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

		return artist, nil
	})(ctx, mbid)
}

// SearchReleaseGroups searches for release groups using PostgreSQL full-text search
func SearchReleaseGroups(ctx context.Context, query string, limit int) (any, error) {
	return cache.NewCache("musicbrainz_release_group_search").TTL(24*time.Hour).Wrap(func() (any, error) {
		// Simplified query with better performance
		searchQuery := `
			SELECT DISTINCT
				rg.gid,
				rg.name,
				rg.type as release_group_type,
				a.name as artist_name,
				a.gid as artist_id,
				rg.comment,
				CASE 
					WHEN rg.name ILIKE $1 || '%' THEN 1
					WHEN rg.name ILIKE '%' || $1 || '%' THEN 2
					WHEN a.name ILIKE $1 || '%' THEN 3
					ELSE 4
				END as sort_order
			FROM release_group rg
			JOIN artist_credit_name acn ON rg.artist_credit = acn.artist_credit
			JOIN artist a ON acn.artist = a.id
			WHERE (
				rg.name ILIKE '%' || $1 || '%' OR
				a.name ILIKE '%' || $1 || '%'
			)
			ORDER BY sort_order, rg.name
			LIMIT $2
		`

		rows, err := db.QueryContext(ctx, searchQuery, query, limit)
		if err != nil {
			return nil, fmt.Errorf("failed to search release groups: %w", err)
		}
		defer rows.Close()

		var releaseGroups []map[string]any
		for rows.Next() {
			var rgGid, rgName, rgType, artistName, artistGid, comment sql.NullString
			var sortOrder int

			err := rows.Scan(&rgGid, &rgName, &rgType, &artistName, &artistGid, &comment, &sortOrder)
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
				"score":          100, // Static score for simplified query
			}

			releaseGroups = append(releaseGroups, releaseGroup)
		}

		return map[string]any{
			"release-groups": releaseGroups,
			"count":          len(releaseGroups),
		}, nil
	})(ctx, query, limit)
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

// SearchReleases searches for releases using PostgreSQL full-text search
func SearchReleases(ctx context.Context, query string, limit int) (any, error) {
	return cache.NewCache("musicbrainz_release_search").TTL(24*time.Hour).Wrap(func() (any, error) {
		// Priority: 1. Release name, 2. Artist name, 3. Release group name
		searchQuery := `
			SELECT DISTINCT
				r.gid,
				r.name,
				r.status,
				a.name as artist_name,
				a.gid as artist_id,
				rg.name as release_group_name,
				rg.gid as release_group_id,
				r.comment,
				ts_rank_cd(
					setweight(to_tsvector('english', COALESCE(r.name, '')), 'A') ||
					setweight(to_tsvector('english', COALESCE(a.name, '')), 'B') ||
					setweight(to_tsvector('english', COALESCE(rg.name, '')), 'C'),
					plainto_tsquery('english', $1)
				) as rank
			FROM release r
			JOIN artist_credit_name acn ON r.artist_credit = acn.artist_credit
			JOIN artist a ON acn.artist = a.id
			LEFT JOIN release_group rg ON r.release_group = rg.id
			WHERE (
				to_tsvector('english', COALESCE(r.name, '')) @@ plainto_tsquery('english', $1) OR
				to_tsvector('english', COALESCE(a.name, '')) @@ plainto_tsquery('english', $1) OR
				to_tsvector('english', COALESCE(rg.name, '')) @@ plainto_tsquery('english', $1)
			)
			GROUP BY r.id, r.gid, r.name, r.status, a.name, a.gid, rg.name, rg.gid, r.comment
			ORDER BY rank DESC, r.name
			LIMIT $2
		`

		rows, err := db.QueryContext(ctx, searchQuery, query, limit)
		if err != nil {
			return nil, fmt.Errorf("failed to search releases: %w", err)
		}
		defer rows.Close()

		var releases []map[string]any
		for rows.Next() {
			var releaseGid, releaseName, releaseStatus, artistName, artistGid, rgName, rgGid, comment sql.NullString
			var rank float64

			err := rows.Scan(&releaseGid, &releaseName, &releaseStatus, &artistName, &artistGid, &rgName, &rgGid, &comment, &rank)
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
				"score":          int(rank * 100),
			}

			releases = append(releases, release)
		}

		return map[string]any{
			"releases": releases,
			"count":    len(releases),
		}, nil
	})(ctx, query, limit)
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

// SearchRecordings searches for recordings using PostgreSQL full-text search
func SearchRecordings(ctx context.Context, query string, limit int) (any, error) {
	return cache.NewCache("musicbrainz_recording_search").TTL(24*time.Hour).Wrap(func() (any, error) {
		// Simplified query with better performance
		searchQuery := `
			SELECT DISTINCT
				rec.gid,
				rec.name,
				rec.length,
				a.name as artist_name,
				a.gid as artist_id,
				rec.comment,
				CASE 
					WHEN rec.name ILIKE $1 || '%' THEN 1
					WHEN rec.name ILIKE '%' || $1 || '%' THEN 2
					WHEN a.name ILIKE $1 || '%' THEN 3
					ELSE 4
				END as sort_order
			FROM recording rec
			JOIN artist_credit_name acn ON rec.artist_credit = acn.artist_credit
			JOIN artist a ON acn.artist = a.id
			WHERE (
				rec.name ILIKE '%' || $1 || '%' OR
				a.name ILIKE '%' || $1 || '%'
			)
			ORDER BY sort_order, rec.name
			LIMIT $2
		`

		rows, err := db.QueryContext(ctx, searchQuery, query, limit)
		if err != nil {
			return nil, fmt.Errorf("failed to search recordings: %w", err)
		}
		defer rows.Close()

		var recordings []map[string]any
		for rows.Next() {
			var recGid, recName, artistName, artistGid, comment sql.NullString
			var length sql.NullInt64
			var sortOrder int

			err := rows.Scan(&recGid, &recName, &length, &artistName, &artistGid, &comment, &sortOrder)
			if err != nil {
				continue
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
				"score":          100, // Static score for simplified query
			}

			if length.Valid {
				recording["length"] = length.Int64
			}

			recordings = append(recordings, recording)
		}

		return map[string]any{
			"recordings": recordings,
			"count":      len(recordings),
		}, nil
	})(ctx, query, limit)
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
		defer rows.Close()

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
		defer rows.Close()

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
		defer rows.Close()

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
