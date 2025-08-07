package musicbrainz

import (
	"context"
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
		}
	}()

	cache.InitCache(cacheHost, 6379, 0)
}

func TestSearchArtists(t *testing.T) {
	// Initialize the database connection if not already done
	InitMusicBrainz()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test basic artist search
	result, err := SearchArtists(ctx, "Beatles", 5)
	if err != nil {
		t.Fatalf("SearchArtists failed: %v", err)
	}

	if result == nil {
		t.Fatal("SearchArtists returned nil result")
	}

	// Check if result is properly structured
	if resultMap, ok := result.(map[string]any); ok {
		if artists, hasArtists := resultMap["artists"]; hasArtists {
			t.Logf("Found artists in response: %T", artists)
			if artistList, ok := artists.([]map[string]any); ok && len(artistList) > 0 {
				t.Logf("First artist: %+v", artistList[0])
			}
		} else {
			t.Logf("Response keys: %v", getKeys(resultMap))
		}
		if count, hasCount := resultMap["count"]; hasCount {
			t.Logf("Artist count: %v", count)
		}
	} else {
		t.Fatalf("SearchArtists returned unexpected type: %T", result)
	}
}

func TestGetArtist(t *testing.T) {
	// Initialize the database connection if not already done
	InitMusicBrainz()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

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
	// Initialize the database connection if not already done
	InitMusicBrainz()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Test release group search
	result, err := SearchReleaseGroups(ctx, "Abbey Road", 5)
	if err != nil {
		t.Fatalf("SearchReleaseGroups failed: %v", err)
	}

	if result == nil {
		t.Fatal("SearchReleaseGroups returned nil result")
	}

	// Check if result is properly structured
	if resultMap, ok := result.(map[string]any); ok {
		if releaseGroups, hasReleaseGroups := resultMap["release-groups"]; hasReleaseGroups {
			t.Logf("Found release groups in response: %T", releaseGroups)
			if rgList, ok := releaseGroups.([]map[string]any); ok && len(rgList) > 0 {
				t.Logf("First release group: %+v", rgList[0])
			}
		} else {
			t.Logf("Response keys: %v", getKeys(resultMap))
		}
		if count, hasCount := resultMap["count"]; hasCount {
			t.Logf("Release group count: %v", count)
		}
	} else {
		t.Fatalf("SearchReleaseGroups returned unexpected type: %T", result)
	}
}

func TestSearchReleases(t *testing.T) {
	// Initialize the database connection if not already done
	InitMusicBrainz()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test release search
	result, err := SearchReleases(ctx, "Abbey Road", 5)
	if err != nil {
		t.Fatalf("SearchReleases failed: %v", err)
	}

	if result == nil {
		t.Fatal("SearchReleases returned nil result")
	}

	// Check if result is properly structured
	if resultMap, ok := result.(map[string]any); ok {
		if releases, hasReleases := resultMap["releases"]; hasReleases {
			t.Logf("Found releases in response: %T", releases)
			if releaseList, ok := releases.([]map[string]any); ok && len(releaseList) > 0 {
				t.Logf("First release: %+v", releaseList[0])
			}
		} else {
			t.Logf("Response keys: %v", getKeys(resultMap))
		}
		if count, hasCount := resultMap["count"]; hasCount {
			t.Logf("Release count: %v", count)
		}
	} else {
		t.Fatalf("SearchReleases returned unexpected type: %T", result)
	}
}

func TestSearchRecordings(t *testing.T) {
	// Initialize the database connection if not already done
	InitMusicBrainz()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Test recording search
	result, err := SearchRecordings(ctx, "Come Together", 5)
	if err != nil {
		t.Fatalf("SearchRecordings failed: %v", err)
	}

	if result == nil {
		t.Fatal("SearchRecordings returned nil result")
	}

	// Check if result is properly structured
	if resultMap, ok := result.(map[string]any); ok {
		if recordings, hasRecordings := resultMap["recordings"]; hasRecordings {
			t.Logf("Found recordings in response: %T", recordings)
			if recordingList, ok := recordings.([]map[string]any); ok && len(recordingList) > 0 {
				t.Logf("First recording: %+v", recordingList[0])
			}
		} else {
			t.Logf("Response keys: %v", getKeys(resultMap))
		}
		if count, hasCount := resultMap["count"]; hasCount {
			t.Logf("Recording count: %v", count)
		}
	} else {
		t.Fatalf("SearchRecordings returned unexpected type: %T", result)
	}
}

func TestBrowseArtistReleaseGroups(t *testing.T) {
	// Initialize the database connection if not already done
	InitMusicBrainz()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

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
	// Initialize the database connection if not already done
	InitMusicBrainz()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

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

// Helper function to get map keys for debugging
func getKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
