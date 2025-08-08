package typesense

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"relay/app/cache"

	"github.com/typesense/typesense-go/v2/typesense"
	"github.com/typesense/typesense-go/v2/typesense/api"
)

var typesenseClient *typesense.Client
var db *sql.DB

// SetDB sets the database connection for the typesense package
func SetDB(database *sql.DB) {
	db = database
}

// ImportDocuments imports documents to a Typesense collection
func ImportDocuments(ctx context.Context, collection string, docs []interface{}, params *api.ImportDocumentsParams) (interface{}, error) {
	if typesenseClient == nil {
		return nil, fmt.Errorf("typesense client not initialized")
	}
	return typesenseClient.Collection(collection).Documents().Import(ctx, docs, params)
}

// IsReady checks if both MusicBrainz and Typesense are configured and ready
func IsReady() bool {
	return db != nil && typesenseClient != nil
}

// InitTypesense initializes the Typesense client and creates collections
func InitTypesense(host, port, apiKey string, timeout time.Duration) error {
	if timeout <= 0 {
		timeout = 60 * time.Second
	}
	typesenseClient = typesense.NewClient(
		typesense.WithServer(fmt.Sprintf("http://%s:%s", host, port)),
		typesense.WithAPIKey(apiKey),
		typesense.WithConnectionTimeout(timeout),
		typesense.WithCircuitBreakerMaxRequests(50),
		typesense.WithCircuitBreakerInterval(2*time.Minute),
		typesense.WithCircuitBreakerTimeout(1*time.Minute),
	)

	// Test connection
	ctx := context.Background()
	_, err := typesenseClient.Collections().Retrieve(ctx)
	if err != nil {
		typesenseClient = nil // Reset client on connection failure
		return fmt.Errorf("failed to connect to Typesense: %w", err)
	}

	slog.Info("Connected to Typesense successfully")

	// Create collections if they don't exist
	if err := createCollections(); err != nil {
		return fmt.Errorf("failed to create collections: %w", err)
	}

	return nil
}

// createCollections creates the necessary Typesense collections for MusicBrainz data
func createCollections() error {
	ctx := context.Background()

	// Helper function to create bool pointers
	boolPtr := func(b bool) *bool { return &b }
	stringPtr := func(s string) *string { return &s }

	// Artist collection schema
	artistSchema := &api.CollectionSchema{
		Name: "artists",
		Fields: []api.Field{
			{Name: "mbid", Type: "string", Facet: boolPtr(false), Index: boolPtr(true)},
			{Name: "name", Type: "string", Facet: boolPtr(false), Index: boolPtr(true), Sort: boolPtr(true)},
			{Name: "sort_name", Type: "string", Facet: boolPtr(false), Index: boolPtr(true), Sort: boolPtr(true)},
			{Name: "artist_type", Type: "string", Facet: boolPtr(true), Index: boolPtr(true), Optional: boolPtr(true)},
			{Name: "area_name", Type: "string", Facet: boolPtr(true), Index: boolPtr(true), Optional: boolPtr(true)},
			{Name: "begin_year", Type: "int32", Facet: boolPtr(true), Index: boolPtr(true), Optional: boolPtr(true)},
			{Name: "end_year", Type: "int32", Facet: boolPtr(true), Index: boolPtr(true), Optional: boolPtr(true)},
			{Name: "ended", Type: "bool", Facet: boolPtr(true), Index: boolPtr(true), Optional: boolPtr(true)},
			{Name: "comment", Type: "string", Facet: boolPtr(false), Index: boolPtr(true), Optional: boolPtr(true)},
			{Name: "aliases", Type: "string[]", Facet: boolPtr(false), Index: boolPtr(true), Optional: boolPtr(true)},
		},
		DefaultSortingField: stringPtr("name"),
	}

	// Release Group collection schema
	releaseGroupSchema := &api.CollectionSchema{
		Name: "release_groups",
		Fields: []api.Field{
			{Name: "mbid", Type: "string", Facet: boolPtr(false), Index: boolPtr(true)},
			{Name: "name", Type: "string", Facet: boolPtr(false), Index: boolPtr(true), Sort: boolPtr(true)},
			{Name: "primary_type", Type: "string", Facet: boolPtr(true), Index: boolPtr(true), Optional: boolPtr(true)},
			{Name: "comment", Type: "string", Facet: boolPtr(false), Index: boolPtr(true), Optional: boolPtr(true)},
		},
		DefaultSortingField: stringPtr("name"),
	}

	// Release collection schema
	releaseSchema := &api.CollectionSchema{
		Name: "releases",
		Fields: []api.Field{
			{Name: "mbid", Type: "string", Facet: boolPtr(false), Index: boolPtr(true)},
			{Name: "name", Type: "string", Facet: boolPtr(false), Index: boolPtr(true), Sort: boolPtr(true)},
			{Name: "status", Type: "string", Facet: boolPtr(true), Index: boolPtr(true), Optional: boolPtr(true)},
			{Name: "artist_name", Type: "string", Facet: boolPtr(false), Index: boolPtr(true)},
			{Name: "artist_mbid", Type: "string", Facet: boolPtr(false), Index: boolPtr(true)},
			{Name: "release_group_name", Type: "string", Facet: boolPtr(false), Index: boolPtr(true)},
			{Name: "release_group_mbid", Type: "string", Facet: boolPtr(false), Index: boolPtr(true)},
			{Name: "comment", Type: "string", Facet: boolPtr(false), Index: boolPtr(true), Optional: boolPtr(true)},
			{Name: "date_year", Type: "int32", Facet: boolPtr(true), Index: boolPtr(true), Optional: boolPtr(true)},
		},
		DefaultSortingField: stringPtr("name"),
	}

	// Recording collection schema
	recordingSchema := &api.CollectionSchema{
		Name: "recordings",
		Fields: []api.Field{
			{Name: "mbid", Type: "string", Facet: boolPtr(false), Index: boolPtr(true)},
			{Name: "name", Type: "string", Facet: boolPtr(false), Index: boolPtr(true), Sort: boolPtr(true)},
			{Name: "length", Type: "int64", Facet: boolPtr(true), Index: boolPtr(true), Optional: boolPtr(true)},
			{Name: "comment", Type: "string", Facet: boolPtr(false), Index: boolPtr(true), Optional: boolPtr(true)},
		},
		DefaultSortingField: stringPtr("name"),
	}

	schemas := []*api.CollectionSchema{artistSchema, releaseGroupSchema, releaseSchema, recordingSchema}

	for _, schema := range schemas {
		// Check if collection exists
		_, err := typesenseClient.Collection(schema.Name).Retrieve(ctx)
		if err != nil {
			// Collection doesn't exist, create it
			_, err := typesenseClient.Collections().Create(ctx, schema)
			if err != nil {
				return fmt.Errorf("failed to create collection %s: %w", schema.Name, err)
			}
			slog.Info("Created Typesense collection", "collection", schema.Name)
		} else {
			// Collection exists, check if we need to recreate it due to schema changes
			// For now, just log that it exists
			slog.Info("Typesense collection already exists", "collection", schema.Name)
		}
	}

	return nil
}

// SearchArtistsTypesense searches artists using Typesense with prioritized fields
func SearchArtistsTypesense(ctx context.Context, query string, limit int) (any, error) {
	return cache.NewCache("typesense_artist_search").TTL(24*time.Hour).Wrap(func() (any, error) {
		if typesenseClient == nil {
			return nil, fmt.Errorf("typesense client not initialized")
		}

		intPtr := func(i int) *int { return &i }
		stringPtr := func(s string) *string { return &s }

		searchParams := &api.SearchCollectionParams{
			Q:              stringPtr(query),
			QueryBy:        stringPtr("name,sort_name,aliases"),
			QueryByWeights: stringPtr("3,2,1"), // Prioritize: name > sort_name > aliases
			SortBy:         stringPtr("_text_match:desc,name:asc"),
			PerPage:        intPtr(limit),
			Page:           intPtr(1),
		}

		result, err := typesenseClient.Collection("artists").Documents().Search(ctx, searchParams)
		if err != nil {
			return nil, fmt.Errorf("failed to search artists in Typesense: %w", err)
		}

		var artists []map[string]any
		if result.Hits != nil {
			for _, hit := range *result.Hits {
				doc := hit.Document
				artist := map[string]any{
					"id":        (*doc)["mbid"],
					"name":      (*doc)["name"],
					"sort-name": (*doc)["sort_name"],
					"score":     int(*hit.TextMatch * 100),
				}

				if artistType, ok := (*doc)["artist_type"]; ok {
					artist["type"] = artistType
				}
				if areaName, ok := (*doc)["area_name"]; ok {
					artist["area"] = areaName
					artist["begin-area"] = areaName
				}
				if ended, ok := (*doc)["ended"]; ok {
					artist["ended"] = ended
				}
				if comment, ok := (*doc)["comment"]; ok {
					artist["disambiguation"] = comment
				}

				// Handle life-span
				if beginYear, ok := (*doc)["begin_year"]; ok {
					lifeSpan := map[string]any{"begin": fmt.Sprintf("%.0f", beginYear.(float64))}
					if endYear, ok := (*doc)["end_year"]; ok {
						lifeSpan["end"] = fmt.Sprintf("%.0f", endYear.(float64))
					}
					artist["life-span"] = lifeSpan
				}

				artists = append(artists, artist)
			}
		}

		return map[string]any{
			"artists": artists,
			"count":   len(artists),
		}, nil
	})(ctx, query, limit)
}

// SearchReleaseGroupsTypesense searches release groups using Typesense with prioritized fields
func SearchReleaseGroupsTypesense(ctx context.Context, query string, limit int) (any, error) {
	return cache.NewCache("typesense_release_group_search").TTL(24*time.Hour).Wrap(func() (any, error) {
		if typesenseClient == nil {
			return nil, fmt.Errorf("typesense client not initialized")
		}

		intPtr := func(i int) *int { return &i }
		stringPtr := func(s string) *string { return &s }

		searchParams := &api.SearchCollectionParams{
			Q:              stringPtr(query),
			QueryBy:        stringPtr("name,comment"),
			QueryByWeights: stringPtr("3,1"), // Prioritize: name > description
			SortBy:         stringPtr("_text_match:desc,name:asc"),
			PerPage:        intPtr(limit),
			Page:           intPtr(1),
		}

		result, err := typesenseClient.Collection("release_groups").Documents().Search(ctx, searchParams)
		if err != nil {
			return nil, fmt.Errorf("failed to search release groups in Typesense: %w", err)
		}

		var releaseGroups []map[string]any
		if result.Hits != nil {
			for _, hit := range *result.Hits {
				doc := hit.Document
				releaseGroup := map[string]any{
					"id":    (*doc)["mbid"],
					"title": (*doc)["name"],
					"score": int(*hit.TextMatch * 100),
				}

				if primaryType, ok := (*doc)["primary_type"]; ok {
					releaseGroup["primary-type"] = primaryType
				}
				if comment, ok := (*doc)["comment"]; ok {
					releaseGroup["disambiguation"] = comment
				}

				releaseGroups = append(releaseGroups, releaseGroup)
			}
		}

		return map[string]any{
			"release-groups": releaseGroups,
			"count":          len(releaseGroups),
		}, nil
	})(ctx, query, limit)
}

// SearchReleasesTypesense searches releases using Typesense with prioritized fields
func SearchReleasesTypesense(ctx context.Context, query string, limit int) (any, error) {
	return cache.NewCache("typesense_release_search").TTL(24*time.Hour).Wrap(func() (any, error) {
		if typesenseClient == nil {
			return nil, fmt.Errorf("typesense client not initialized")
		}

		intPtr := func(i int) *int { return &i }
		stringPtr := func(s string) *string { return &s }

		// Use fields that exist in the releases schema
		searchParams := &api.SearchCollectionParams{
			Q:              stringPtr(query),
			QueryBy:        stringPtr("name,artist_name,release_group_name,comment"),
			QueryByWeights: stringPtr("4,3,2,1"), // title > artist > release group > comment
			SortBy:         stringPtr("_text_match:desc,name:asc"),
			PerPage:        intPtr(limit),
			Page:           intPtr(1),
		}

		result, err := typesenseClient.Collection("releases").Documents().Search(ctx, searchParams)
		if err != nil {
			return nil, fmt.Errorf("failed to search releases in Typesense: %w", err)
		}

		var releases []map[string]any
		if result.Hits != nil {
			for _, hit := range *result.Hits {
				doc := hit.Document
				release := map[string]any{
					"id":    (*doc)["mbid"],
					"title": (*doc)["name"],
					"score": int(*hit.TextMatch * 100),
				}

				if status, ok := (*doc)["status"]; ok {
					release["status"] = status
				}
				// artist-credit block
				if an, ok := (*doc)["artist_name"]; ok {
					release["artist-credit"] = []map[string]any{
						{
							"name": an,
						},
					}
					if am, ok := (*doc)["artist_mbid"]; ok {
						release["artist-credit"].([]map[string]any)[0]["artist"] = map[string]any{"id": am, "name": an}
					}
				}
				// release-group block
				if rgn, ok := (*doc)["release_group_name"]; ok {
					rg := map[string]any{"title": rgn}
					if rgm, ok := (*doc)["release_group_mbid"]; ok {
						rg["id"] = rgm
					}
					release["release-group"] = rg
				}
				if comment, ok := (*doc)["comment"]; ok {
					release["disambiguation"] = comment
				}
				if yr, ok := (*doc)["date_year"]; ok {
					release["date-year"] = yr
				}

				releases = append(releases, release)
			}
		}

		return map[string]any{
			"releases": releases,
			"count":    len(releases),
		}, nil
	})(ctx, query, limit)
}

// SearchRecordingsTypesense searches recordings using Typesense with prioritized fields
func SearchRecordingsTypesense(ctx context.Context, query string, limit int) (any, error) {
	return cache.NewCache("typesense_recording_search").TTL(24*time.Hour).Wrap(func() (any, error) {
		if typesenseClient == nil {
			return nil, fmt.Errorf("typesense client not initialized")
		}

		intPtr := func(i int) *int { return &i }
		stringPtr := func(s string) *string { return &s }

		searchParams := &api.SearchCollectionParams{
			Q:              stringPtr(query),
			QueryBy:        stringPtr("name,comment"),
			QueryByWeights: stringPtr("3,1"), // Prioritize: song name > description
			SortBy:         stringPtr("_text_match:desc,name:asc"),
			PerPage:        intPtr(limit),
			Page:           intPtr(1),
		}

		result, err := typesenseClient.Collection("recordings").Documents().Search(ctx, searchParams)
		if err != nil {
			return nil, fmt.Errorf("failed to search recordings in Typesense: %w", err)
		}

		var recordings []map[string]any
		if result.Hits != nil {
			for _, hit := range *result.Hits {
				doc := hit.Document
				recording := map[string]any{
					"id":    (*doc)["mbid"],
					"title": (*doc)["name"],
					"score": int(*hit.TextMatch * 100),
				}

				if length, ok := (*doc)["length"]; ok {
					recording["length"] = int64(length.(float64))
				}
				if comment, ok := (*doc)["comment"]; ok {
					recording["disambiguation"] = comment
				}

				recordings = append(recordings, recording)
			}
		}

		return map[string]any{
			"recordings": recordings,
			"count":      len(recordings),
		}, nil
	})(ctx, query, limit)
}

// GetArtistTypesense gets a specific artist by MBID from Typesense
func GetArtistTypesense(ctx context.Context, mbid string) (any, error) {
	return cache.NewCache("typesense_artist_mbid").TTL(7*24*time.Hour).Wrap(func() (any, error) {
		if typesenseClient == nil {
			return nil, fmt.Errorf("typesense client not initialized")
		}

		intPtr := func(i int) *int { return &i }
		stringPtr := func(s string) *string { return &s }

		searchParams := &api.SearchCollectionParams{
			Q:       stringPtr(mbid),
			QueryBy: stringPtr("mbid"),
			PerPage: intPtr(1),
			Page:    intPtr(1),
		}

		result, err := typesenseClient.Collection("artists").Documents().Search(ctx, searchParams)
		if err != nil {
			return nil, fmt.Errorf("failed to search artist by MBID in Typesense: %w", err)
		}

		if result.Hits == nil || len(*result.Hits) == 0 {
			return nil, fmt.Errorf("artist not found")
		}

		hit := (*result.Hits)[0]
		doc := hit.Document

		artist := map[string]any{
			"id":        (*doc)["mbid"],
			"name":      (*doc)["name"],
			"sort-name": (*doc)["sort_name"],
		}

		if artistType, ok := (*doc)["artist_type"]; ok {
			artist["type"] = artistType
		}
		if areaName, ok := (*doc)["area_name"]; ok {
			artist["area"] = areaName
		}
		if ended, ok := (*doc)["ended"]; ok {
			artist["ended"] = ended
		}
		if comment, ok := (*doc)["comment"]; ok {
			artist["disambiguation"] = comment
		}

		if beginYear, ok := (*doc)["begin_year"]; ok {
			lifeSpan := map[string]any{"begin": fmt.Sprintf("%.0f", beginYear.(float64))}
			if endYear, ok := (*doc)["end_year"]; ok {
				lifeSpan["end"] = fmt.Sprintf("%.0f", endYear.(float64))
			}
			artist["life-span"] = lifeSpan
		}

		return artist, nil
	})(ctx, mbid)
}
