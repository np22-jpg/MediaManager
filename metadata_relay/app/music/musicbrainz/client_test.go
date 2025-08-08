package musicbrainz

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"relay/app/cache"
)

func init() {
	// Initialize cache for tests - use environment variables or defaults
	cacheHost := os.Getenv("CACHE_HOST")
	if cacheHost == "" {
		cacheHost = "localhost"
	}

	// Try to initialize cache, but don't fail if Redis is not available in test environment
	defer func() {
		if r := recover(); r != nil {
			// Cache initialization failed, that's ok for tests
			// The cache will be disabled and functions will execute directly
			_ = r // Mark as used to satisfy linter
		}
	}()

	cache.InitCache(cacheHost, 6379, 0)
}

// getTestDBConnStr returns the connection string for testing
func getTestDBConnStr() string {
	host := os.Getenv("MUSICBRAINZ_DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("MUSICBRAINZ_DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("MUSICBRAINZ_DB_USER")
	if user == "" {
		user = "musicbrainz"
	}
	password := os.Getenv("MUSICBRAINZ_DB_PASSWORD")
	if password == "" {
		password = "musicbrainz"
	}
	dbname := os.Getenv("MUSICBRAINZ_DB_NAME")
	if dbname == "" {
		dbname = "musicbrainz_db"
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

func TestSearchArtists(t *testing.T) {
	t.Skip("skipping: SearchArtists deprecated test; Typesense-only search")
}

func TestGetArtist(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTS") != "true" {
		t.Skip("skipping integration test: set INTEGRATION_TESTS=true to enable")
	}

	// Initialize the database connection if not already done
	InitMusicBrainz(getTestDBConnStr())

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if db == nil {
		t.Skip("skipping: MusicBrainz DB not configured")
	}
	if err := db.PingContext(ctx); err != nil {
		t.Skipf("skipping: MusicBrainz DB not reachable: %v", err)
	}

	// Test getting a specific artist by MBID - using a well-known Beatles MBID
	// This is The Beatles' MBID in MusicBrainz
	beatlesMbid := "b10bbbfc-cf9e-42e0-be17-e2c3e1d2600d"

	result, err := GetArtist(ctx, beatlesMbid)
	if err != nil {
		t.Logf("GetArtist failed (expected if Beatles not in local DB): %v", err)
		// This might fail if the local database doesn't have The Beatles
		// Let's skip this test if the artist is not found
		return
	}

	if result == nil {
		t.Fatal("GetArtist returned nil result")
	}

	if artistMap, ok := result.(map[string]any); ok {
		t.Logf("Artist details: %+v", artistMap)
		if name, hasName := artistMap["name"]; hasName {
			t.Logf("Artist name: %v", name)
		}
	} else {
		t.Fatalf("GetArtist returned unexpected type: %T", result)
	}
}

func TestSearchReleaseGroups(t *testing.T) {
	t.Skip("skipping: deprecated search tests; Typesense-only search")
}

func TestSearchReleases(t *testing.T) {
	t.Skip("skipping: deprecated search tests; Typesense-only search")
}

func TestSearchRecordings(t *testing.T) {
	t.Skip("skipping: deprecated search tests; Typesense-only search")
}

func TestBrowseArtistReleaseGroups(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTS") != "true" {
		t.Skip("skipping integration test: set INTEGRATION_TESTS=true to enable")
	}

	// Initialize the database connection if not already done
	InitMusicBrainz(getTestDBConnStr())

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if db == nil {
		t.Skip("skipping: MusicBrainz DB not configured")
	}
	if err := db.PingContext(ctx); err != nil {
		t.Skipf("skipping: MusicBrainz DB not reachable: %v", err)
	}

	// Test browsing release groups for an artist
	// Using The Beatles' MBID
	beatlesMbid := "b10bbbfc-cf9e-42e0-be17-e2c3e1d2600d"

	result, err := BrowseArtistReleaseGroups(ctx, beatlesMbid, 10)
	if err != nil {
		t.Logf("BrowseArtistReleaseGroups failed (expected if Beatles not in local DB): %v", err)
		// This might fail if the local database doesn't have The Beatles
		return
	}

	if result == nil {
		t.Fatal("BrowseArtistReleaseGroups returned nil result")
	}

	// Check if result is properly structured
	if resultMap, ok := result.(map[string]any); ok {
		if releaseGroups, hasReleaseGroups := resultMap["release-groups"]; hasReleaseGroups {
			t.Logf("Found release groups for artist: %T", releaseGroups)
			if rgList, ok := releaseGroups.([]map[string]any); ok && len(rgList) > 0 {
				t.Logf("First release group: %+v", rgList[0])
			}
		}
		if count, hasCount := resultMap["count"]; hasCount {
			t.Logf("Release group count for artist: %v", count)
		}
	} else {
		t.Fatalf("BrowseArtistReleaseGroups returned unexpected type: %T", result)
	}
}

func TestAdvancedSearchArtists(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTS") != "true" {
		t.Skip("skipping integration test: set INTEGRATION_TESTS=true to enable")
	}

	// Initialize the database connection if not already done
	InitMusicBrainz(getTestDBConnStr())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if db == nil {
		t.Skip("skipping: MusicBrainz DB not configured")
	}
	if err := db.PingContext(ctx); err != nil {
		t.Skipf("skipping: MusicBrainz DB not reachable: %v", err)
	}

	// Test advanced artist search
	result, err := AdvancedSearchArtists(ctx, "Beatles", "", "", "", 5)
	if err != nil {
		t.Fatalf("AdvancedSearchArtists failed: %v", err)
	}

	if result == nil {
		t.Fatal("AdvancedSearchArtists returned nil result")
	}

	// Check if result is properly structured
	if resultMap, ok := result.(map[string]any); ok {
		if artists, hasArtists := resultMap["artists"]; hasArtists {
			t.Logf("Found artists in advanced search: %T", artists)
			if artistList, ok := artists.([]map[string]any); ok && len(artistList) > 0 {
				t.Logf("First artist from advanced search: %+v", artistList[0])
			}
		}
		if count, hasCount := resultMap["count"]; hasCount {
			t.Logf("Advanced search artist count: %v", count)
		}
	} else {
		t.Fatalf("AdvancedSearchArtists returned unexpected type: %T", result)
	}
}
