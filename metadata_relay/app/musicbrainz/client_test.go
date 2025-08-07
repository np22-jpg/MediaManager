package musicbrainz

import (
	"net/url"
	"testing"
)

func TestMakeRequest(t *testing.T) {
	// Test basic MusicBrainz API request without caching
	params := url.Values{}
	params.Set("query", "The Beatles")
	params.Set("limit", "5")

	result, err := makeRequest("/artist", params)
	if err != nil {
		t.Fatalf("makeRequest failed: %v", err)
	}

	if result == nil {
		t.Fatal("makeRequest returned nil result")
	}

	t.Logf("makeRequest result type: %T", result)

	// Try to assert it's a map and has expected structure
	if resultMap, ok := result.(map[string]any); ok {
		if artists, hasArtists := resultMap["artists"]; hasArtists {
			t.Logf("Found artists in response: %T", artists)
		} else {
			t.Logf("Response keys: %v", getKeys(resultMap))
		}
	}
}

func TestMakeRequestReleaseGroup(t *testing.T) {
	// Test release group search
	params := url.Values{}
	params.Set("query", "Abbey Road")
	params.Set("limit", "5")

	result, err := makeRequest("/release-group", params)
	if err != nil {
		t.Fatalf("makeRequest for release-group failed: %v", err)
	}

	if result == nil {
		t.Fatal("makeRequest returned nil result")
	}

	t.Logf("makeRequest release-group result type: %T", result)
}

// Helper function to get map keys for debugging
func getKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
