package seadex

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSeaDexWithMockServer(t *testing.T) {
	// Activate HTTP mock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Initialize SeaDex with mock server
	mockBaseURL := "http://mock-seadex-api.com"
	InitSeaDex(mockBaseURL)

	// Mock search entries endpoint
	httpmock.RegisterResponder("GET", mockBaseURL+"/api/collections/entries/records",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
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
		}))

	// Mock get entry by ID endpoint
	httpmock.RegisterResponder("GET", mockBaseURL+"/api/collections/entries/records/test_entry_1",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
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
		}))

	ctx := context.Background()

	t.Run("SearchEntries", func(t *testing.T) {
		result, err := SearchEntries(ctx, "test anime", 1, 30)
		require.NoError(t, err)
		assert.NotNil(t, result)

		// Check result structure
		resultMap, ok := result.(map[string]any)
		require.True(t, ok, "Result should be map[string]any")

		assert.Equal(t, float64(1), resultMap["page"])
		assert.Equal(t, float64(1), resultMap["totalItems"])

		items, ok := resultMap["items"].([]any)
		require.True(t, ok, "Items should be slice")
		assert.Len(t, items, 1)

		// Check first entry
		entry, ok := items[0].(map[string]any)
		require.True(t, ok, "Entry should be map")
		assert.Equal(t, "test_entry_1", entry["id"])
		assert.Equal(t, float64(12345), entry["anilist_id"])
	})

	t.Run("GetEntryByID", func(t *testing.T) {
		result, err := GetEntryByID(ctx, "test_entry_1")
		require.NoError(t, err)
		assert.NotNil(t, result)

		// Check result structure
		entryMap, ok := result.(map[string]any)
		require.True(t, ok, "Result should be a map")

		assert.Equal(t, "test_entry_1", entryMap["id"])
		assert.Equal(t, float64(12345), entryMap["anilist_id"])
		assert.Equal(t, "Test Anime Collection", entryMap["collection_name"])
	})

	t.Run("GetTrendingEntries", func(t *testing.T) {
		result, err := GetTrendingEntries(ctx, 50)
		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	// Verify expected calls were made
	assert.GreaterOrEqual(t, httpmock.GetTotalCallCount(), 3)
}

func TestSeaDexErrorHandling(t *testing.T) {
	// Test with no base URL (should use default)
	InitSeaDex("")

	ctx := context.Background()

	t.Run("NoMockServerRunning", func(t *testing.T) {
		// This should fail because no server is running on the default URL
		result, err := SearchEntries(ctx, "test", 1, 30)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestSeaDexServerErrors(t *testing.T) {
	// Activate HTTP mock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Initialize SeaDex with mock server
	mockBaseURL := "http://mock-seadx-api.com"
	InitSeaDex(mockBaseURL)

	// Mock server error
	httpmock.RegisterResponder("GET", mockBaseURL+"/api/collections/entries/records",
		httpmock.NewStringResponder(500, "Internal Server Error"))

	ctx := context.Background()

	t.Run("ServerError", func(t *testing.T) {
		result, err := SearchEntries(ctx, "test", 1, 30)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
