package seadex

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSeaDexWithMockServer(t *testing.T) {
	// Mock response data
	searchResponse := map[string]any{
		"page":       1,
		"perPage":    30,
		"totalItems": 1,
		"totalPages": 1,
		"items": []map[string]any{
			{
				"id":               "test_entry_1",
				"anilist_id":       12345,
				"collection_id":    "test_collection",
				"collection_name":  "Test Anime Collection",
				"comparisons":      []string{},
				"created":          "2024-01-01T00:00:00Z",
				"is_incomplete":    false,
				"notes":            "Test anime entry",
				"size":             1073741824,
				"theoretical_best": nil,
				"torrents": []map[string]any{
					{
						"id":              "torrent_1",
						"collection_id":   "test_collection",
						"collection_name": "Test Anime Collection",
						"created":         "2024-01-01T00:00:00Z",
						"files": []map[string]any{
							{
								"name": "Test_Anime_E01.mkv",
								"size": 536870912,
							},
						},
						"grouped_url":   nil,
						"infohash":      "abc123def456",
						"is_best":       true,
						"is_dual_audio": false,
						"release_group": "TestGroup",
						"size":          536870912,
						"tags":          []string{"bluray", "1080p"},
						"tracker":       "nyaa",
						"updated":       "2024-01-01T00:00:00Z",
						"url":           "https://test.tracker/torrent/1",
					},
				},
				"updated": "2024-01-01T00:00:00Z",
				"url":     "https://releases.moe/entry/test_entry_1",
			},
		},
	}

	entryResponse := map[string]any{
		"id":              "test_entry_1",
		"anilist_id":      12345,
		"collection_id":   "test_collection",
		"collection_name": "Test Anime Collection",
		"created":         "2024-01-01T00:00:00Z",
		"is_incomplete":   false,
		"notes":           "Detailed test anime entry",
		"size":            1073741824,
		"torrents": []map[string]any{
			{
				"id":            "torrent_1",
				"release_group": "TestGroup",
				"is_best":       true,
				"tracker":       "nyaa",
			},
		},
	}

	// Mock server setup
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/api/collections/entries/records":
			w.WriteHeader(200)
			_ = json.NewEncoder(w).Encode(searchResponse)
		case "/api/collections/entries/records/test_entry_1":
			w.WriteHeader(200)
			_ = json.NewEncoder(w).Encode(entryResponse)
		default:
			w.WriteHeader(404)
			_, _ = w.Write([]byte("Not Found"))
		}
	}))
	defer server.Close()

	// Initialize SeaDex with mock server
	InitSeaDex(server.URL)

	ctx := context.Background()

	t.Run("SearchEntries", func(t *testing.T) {
		result, err := SearchEntries(ctx, "test anime", 1, 30)
		if err != nil {
			t.Errorf("SearchEntries() error = %v", err)
			return
		}
		if result == nil {
			t.Error("Expected result to not be nil")
			return
		}

		// Check result structure
		resultMap, ok := result.(map[string]any)
		if !ok {
			t.Fatalf("Result should be map[string]any, got %T", result)
		}

		if resultMap["page"] != float64(1) {
			t.Errorf("Expected page to be 1, got %v", resultMap["page"])
		}
		if resultMap["totalItems"] != float64(1) {
			t.Errorf("Expected totalItems to be 1, got %v", resultMap["totalItems"])
		}

		items, ok := resultMap["items"].([]any)
		if !ok {
			t.Fatalf("Items should be slice, got %T", resultMap["items"])
		}
		if len(items) != 1 {
			t.Errorf("Expected 1 item, got %d", len(items))
		}

		// Check first entry
		entry, ok := items[0].(map[string]any)
		if !ok {
			t.Fatalf("Entry should be map, got %T", items[0])
		}
		if entry["id"] != "test_entry_1" {
			t.Errorf("Expected id to be 'test_entry_1', got %v", entry["id"])
		}
		if entry["anilist_id"] != float64(12345) {
			t.Errorf("Expected anilist_id to be 12345, got %v", entry["anilist_id"])
		}
	})

	t.Run("GetEntryByID", func(t *testing.T) {
		result, err := GetEntryByID(ctx, "test_entry_1")
		if err != nil {
			t.Errorf("GetEntryByID() error = %v", err)
			return
		}
		if result == nil {
			t.Error("Expected result to not be nil")
			return
		}

		// Check result structure
		entryMap, ok := result.(map[string]any)
		if !ok {
			t.Fatalf("Result should be a map, got %T", result)
		}

		if entryMap["id"] != "test_entry_1" {
			t.Errorf("Expected id to be 'test_entry_1', got %v", entryMap["id"])
		}
		if entryMap["anilist_id"] != float64(12345) {
			t.Errorf("Expected anilist_id to be 12345, got %v", entryMap["anilist_id"])
		}
		if entryMap["collection_name"] != "Test Anime Collection" {
			t.Errorf("Expected collection_name to be 'Test Anime Collection', got %v", entryMap["collection_name"])
		}
	})

	t.Run("GetTrendingEntries", func(t *testing.T) {
		result, err := GetTrendingEntries(ctx, 50)
		if err != nil {
			t.Errorf("GetTrendingEntries() error = %v", err)
			return
		}
		if result == nil {
			t.Error("Expected result to not be nil")
		}
	})
}

func TestSeaDexErrorHandling(t *testing.T) {
	// Test with no base URL (should use default)
	InitSeaDex("")

	ctx := context.Background()

	t.Run("NoMockServerRunning", func(t *testing.T) {
		// This should fail because no server is running on the default URL
		result, err := SearchEntries(ctx, "test", 1, 30)
		if err == nil {
			t.Error("Expected error when no server running")
		}
		if result != nil {
			t.Error("Expected nil result when no server running")
		}
	})
}

func TestSeaDexServerErrors(t *testing.T) {
	// Mock server that returns server error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		_, _ = w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	// Initialize SeaDex with mock server
	InitSeaDex(server.URL)

	ctx := context.Background()

	t.Run("ServerError", func(t *testing.T) {
		result, err := SearchEntries(ctx, "test", 1, 30)
		if err == nil {
			t.Error("Expected error for server error")
		}
		if result != nil {
			t.Error("Expected nil result for server error")
		}
	})
}
