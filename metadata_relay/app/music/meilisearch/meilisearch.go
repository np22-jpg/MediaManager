package meilisearch

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"relay/app/cache"

	"github.com/meilisearch/meilisearch-go"
)

var meilisearchClient meilisearch.ServiceManager
var db *sql.DB

// SetDB sets the database connection for the meilisearch package
func SetDB(database *sql.DB) {
	db = database
}

// IsReady checks if both MusicBrainz and Meilisearch are configured and ready
func IsReady() bool {
	return db != nil && meilisearchClient != nil
}

// InitMeilisearch initializes the Meilisearch client and creates indexes
func InitMeilisearch(host, port, apiKey string, timeout time.Duration) error {
	url := fmt.Sprintf("http://%s:%s", host, port)
	meilisearchClient = meilisearch.New(url, meilisearch.WithAPIKey(apiKey))

	// Test connection by getting version
	_, err := meilisearchClient.Version()
	if err != nil {
		meilisearchClient = nil // Reset client on connection failure
		return fmt.Errorf("failed to connect to Meilisearch: %w", err)
	}

	slog.Info("Connected to Meilisearch successfully")

	// Create indexes if they don't exist
	if err := createIndexes(); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}

// createIndexes creates the necessary Meilisearch indexes for MusicBrainz data
func createIndexes() error {
	indexNames := []string{"artists", "release_groups", "releases", "recordings"}

	for _, indexName := range indexNames {
		// Check if index exists
		_, err := meilisearchClient.GetIndex(indexName)
		if err != nil {
			// Index doesn't exist, create it
			taskInfo, err := meilisearchClient.CreateIndex(&meilisearch.IndexConfig{
				Uid:        indexName,
				PrimaryKey: "mbid",
			})
			if err != nil {
				return fmt.Errorf("failed to create index %s: %w", indexName, err)
			}

			// Wait for index creation to complete
			_, err = meilisearchClient.WaitForTask(taskInfo.TaskUID, 10*time.Second)
			if err != nil {
				return fmt.Errorf("failed waiting for index creation %s: %w", indexName, err)
			}

			slog.Info("Created Meilisearch index", "index", indexName)
		} else {
			slog.Info("Meilisearch index already exists", "index", indexName)
		}

		// Configure index settings (simplified for now)
		index := meilisearchClient.Index(indexName)
		
		// Set searchable attributes based on index type
		var searchableAttrs []string
		switch indexName {
		case "artists":
			searchableAttrs = []string{"name", "sort_name", "aliases"}
		case "release_groups":
			searchableAttrs = []string{"name", "comment"}
		case "releases":
			searchableAttrs = []string{"name", "artist_name", "release_group_name", "comment"}
		case "recordings":
			searchableAttrs = []string{"name", "comment"}
		}
		
		if len(searchableAttrs) > 0 {
			task, err := index.UpdateSearchableAttributes(&searchableAttrs)
			if err != nil {
				slog.Warn("failed to set searchable attributes", "index", indexName, "error", err)
			} else {
				_, _ = meilisearchClient.WaitForTask(task.TaskUID, 5*time.Second)
			}
		}
	}

	return nil
}

// getValue safely extracts a value from JSON RawMessage
func getValue(raw json.RawMessage) interface{} {
	var value interface{}
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil
	}
	return value
}

// SearchArtistsMeilisearch searches artists using Meilisearch
func SearchArtistsMeilisearch(ctx context.Context, query string, limit int) (any, error) {
	return cache.NewCache("meilisearch_artist_search").TTL(24*time.Hour).Wrap(func() (any, error) {
		if meilisearchClient == nil {
			return nil, fmt.Errorf("meilisearch client not initialized")
		}

		index := meilisearchClient.Index("artists")
		searchRes, err := index.Search(query, &meilisearch.SearchRequest{
			Limit: int64(limit),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to search artists in Meilisearch: %w", err)
		}

		var artists []map[string]any
		for _, hit := range searchRes.Hits {
			artist := map[string]any{
				"id":        getValue(hit["mbid"]),
				"name":      getValue(hit["name"]),
				"sort-name": getValue(hit["sort_name"]),
				"score":     100, // Meilisearch doesn't provide detailed scores like Typesense
			}

			if artistType := getValue(hit["artist_type"]); artistType != nil {
				artist["type"] = artistType
			}
			if areaName := getValue(hit["area_name"]); areaName != nil {
				artist["area"] = areaName
				artist["begin-area"] = areaName
			}
			if ended := getValue(hit["ended"]); ended != nil {
				artist["ended"] = ended
			}
			if comment := getValue(hit["comment"]); comment != nil {
				artist["disambiguation"] = comment
			}

			// Handle life-span
			if beginYear := getValue(hit["begin_year"]); beginYear != nil {
				lifeSpan := map[string]any{"begin": fmt.Sprintf("%.0f", beginYear.(float64))}
				if endYear := getValue(hit["end_year"]); endYear != nil {
					lifeSpan["end"] = fmt.Sprintf("%.0f", endYear.(float64))
				}
				artist["life-span"] = lifeSpan
			}

			artists = append(artists, artist)
		}

		return map[string]any{
			"artists": artists,
			"count":   len(artists),
		}, nil
	})(ctx, query, limit)
}

// SearchReleaseGroupsMeilisearch searches release groups using Meilisearch
func SearchReleaseGroupsMeilisearch(ctx context.Context, query string, limit int) (any, error) {
	return cache.NewCache("meilisearch_release_group_search").TTL(24*time.Hour).Wrap(func() (any, error) {
		if meilisearchClient == nil {
			return nil, fmt.Errorf("meilisearch client not initialized")
		}

		index := meilisearchClient.Index("release_groups")
		searchRes, err := index.Search(query, &meilisearch.SearchRequest{
			Limit: int64(limit),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to search release groups in Meilisearch: %w", err)
		}

		var releaseGroups []map[string]any
		for _, hit := range searchRes.Hits {
			releaseGroup := map[string]any{
				"id":    getValue(hit["mbid"]),
				"title": getValue(hit["name"]),
				"score": 100,
			}

			if primaryType := getValue(hit["primary_type"]); primaryType != nil {
				releaseGroup["primary-type"] = primaryType
			}
			if comment := getValue(hit["comment"]); comment != nil {
				releaseGroup["disambiguation"] = comment
			}

			releaseGroups = append(releaseGroups, releaseGroup)
		}

		return map[string]any{
			"release-groups": releaseGroups,
			"count":          len(releaseGroups),
		}, nil
	})(ctx, query, limit)
}

// SearchReleasesMeilisearch searches releases using Meilisearch
func SearchReleasesMeilisearch(ctx context.Context, query string, limit int) (any, error) {
	return cache.NewCache("meilisearch_release_search").TTL(24*time.Hour).Wrap(func() (any, error) {
		if meilisearchClient == nil {
			return nil, fmt.Errorf("meilisearch client not initialized")
		}

		index := meilisearchClient.Index("releases")
		searchRes, err := index.Search(query, &meilisearch.SearchRequest{
			Limit: int64(limit),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to search releases in Meilisearch: %w", err)
		}

		var releases []map[string]any
		for _, hit := range searchRes.Hits {
			release := map[string]any{
				"id":    getValue(hit["mbid"]),
				"title": getValue(hit["name"]),
				"score": 100,
			}

			if status := getValue(hit["status"]); status != nil {
				release["status"] = status
			}
			// artist-credit block
			if artistName := getValue(hit["artist_name"]); artistName != nil {
				release["artist-credit"] = []map[string]any{
					{
						"name": artistName,
					},
				}
				if artistMbid := getValue(hit["artist_mbid"]); artistMbid != nil {
					release["artist-credit"].([]map[string]any)[0]["artist"] = map[string]any{
						"id":   artistMbid,
						"name": artistName,
					}
				}
			}
			// release-group block
			if releaseGroupName := getValue(hit["release_group_name"]); releaseGroupName != nil {
				rg := map[string]any{"title": releaseGroupName}
				if releaseGroupMbid := getValue(hit["release_group_mbid"]); releaseGroupMbid != nil {
					rg["id"] = releaseGroupMbid
				}
				release["release-group"] = rg
			}
			if comment := getValue(hit["comment"]); comment != nil {
				release["disambiguation"] = comment
			}

			releases = append(releases, release)
		}

		return map[string]any{
			"releases": releases,
			"count":    len(releases),
		}, nil
	})(ctx, query, limit)
}

// SearchRecordingsMeilisearch searches recordings using Meilisearch
func SearchRecordingsMeilisearch(ctx context.Context, query string, limit int) (any, error) {
	return cache.NewCache("meilisearch_recording_search").TTL(24*time.Hour).Wrap(func() (any, error) {
		if meilisearchClient == nil {
			return nil, fmt.Errorf("meilisearch client not initialized")
		}

		index := meilisearchClient.Index("recordings")
		searchRes, err := index.Search(query, &meilisearch.SearchRequest{
			Limit: int64(limit),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to search recordings in Meilisearch: %w", err)
		}

		var recordings []map[string]any
		for _, hit := range searchRes.Hits {
			recording := map[string]any{
				"id":    getValue(hit["mbid"]),
				"title": getValue(hit["name"]),
				"score": 100,
			}

			if length := getValue(hit["length"]); length != nil {
				recording["length"] = int64(length.(float64))
			}
			if comment := getValue(hit["comment"]); comment != nil {
				recording["disambiguation"] = comment
			}

			recordings = append(recordings, recording)
		}

		return map[string]any{
			"recordings": recordings,
			"count":      len(recordings),
		}, nil
	})(ctx, query, limit)
}

// GetArtistMeilisearch gets a specific artist by MBID from Meilisearch
func GetArtistMeilisearch(ctx context.Context, mbid string) (any, error) {
	return cache.NewCache("meilisearch_artist_mbid").TTL(7*24*time.Hour).Wrap(func() (any, error) {
		if meilisearchClient == nil {
			return nil, fmt.Errorf("meilisearch client not initialized")
		}

		index := meilisearchClient.Index("artists")
		searchRes, err := index.Search(mbid, &meilisearch.SearchRequest{
			Filter: fmt.Sprintf("mbid = %q", mbid),
			Limit:  1,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to search artist by MBID in Meilisearch: %w", err)
		}

		if len(searchRes.Hits) == 0 {
			return nil, fmt.Errorf("artist not found")
		}

		hit := searchRes.Hits[0]
		artist := map[string]any{
			"id":        getValue(hit["mbid"]),
			"name":      getValue(hit["name"]),
			"sort-name": getValue(hit["sort_name"]),
		}

		if artistType := getValue(hit["artist_type"]); artistType != nil {
			artist["type"] = artistType
		}
		if areaName := getValue(hit["area_name"]); areaName != nil {
			artist["area"] = areaName
		}
		if ended := getValue(hit["ended"]); ended != nil {
			artist["ended"] = ended
		}
		if comment := getValue(hit["comment"]); comment != nil {
			artist["disambiguation"] = comment
		}

		if beginYear := getValue(hit["begin_year"]); beginYear != nil {
			lifeSpan := map[string]any{"begin": fmt.Sprintf("%.0f", beginYear.(float64))}
			if endYear := getValue(hit["end_year"]); endYear != nil {
				lifeSpan["end"] = fmt.Sprintf("%.0f", endYear.(float64))
			}
			artist["life-span"] = lifeSpan
		}

		return artist, nil
	})(ctx, mbid)
}

// ImportDocuments imports documents to a Meilisearch index
func ImportDocuments(ctx context.Context, indexName string, docs []interface{}) error {
	if meilisearchClient == nil {
		return fmt.Errorf("meilisearch client not initialized")
	}
	
	index := meilisearchClient.Index(indexName)
	task, err := index.AddDocuments(docs, nil)
	if err != nil {
		return fmt.Errorf("failed to add documents: %w", err)
	}
	
	// Wait for the task to complete
	_, err = meilisearchClient.WaitForTask(task.TaskUID, 30*time.Second)
	return err
}