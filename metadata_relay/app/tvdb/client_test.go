package tvdb

import (
	"context"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTVDBWithMockServer(t *testing.T) {
	// Activate HTTP mock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Initialize TVDB with mock server
	mockBaseURL := "http://mock-tvdb-api.com/v4"
	InitTVDB("test_api_key", mockBaseURL)

	// Mock login endpoint
	httpmock.RegisterResponder("POST", mockBaseURL+"/login",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"status": "success",
			"data": map[string]any{
				"token": "mock_jwt_token_here",
			},
		}))

	// Mock series endpoint (GetTrendingTV)
	httpmock.RegisterResponder("GET", mockBaseURL+"/series",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"status": "success",
			"data": []map[string]any{
				{
					"id":            123,
					"name":          "Test TV Series",
					"overview":      "A test TV series",
					"firstAired":    "2024-01-01",
					"image":         "https://test.com/poster.jpg",
					"averageRating": 8.5,
				},
				{
					"id":            456,
					"name":          "Another Test Series",
					"overview":      "Another test series",
					"firstAired":    "2024-02-01",
					"image":         "https://test.com/poster2.jpg",
					"averageRating": 7.8,
				},
			},
		}))

	// Mock search TV endpoint
	httpmock.RegisterResponder("GET", mockBaseURL+"/search",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"status": "success",
			"data": []map[string]any{
				{
					"id":         789,
					"name":       "Searched Series",
					"overview":   "A searched series",
					"firstAired": "2024-03-01",
					"image":      "https://test.com/searched.jpg",
				},
			},
		}))

	// Mock series details endpoint (extended)
	httpmock.RegisterResponder("GET", mockBaseURL+"/series/123/extended",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"status": "success",
			"data": map[string]any{
				"id":            123,
				"name":          "Test TV Series",
				"overview":      "A detailed test TV series",
				"firstAired":    "2024-01-01",
				"averageRating": 8.5,
				"genres": []map[string]any{
					{"id": 1, "name": "Drama"},
					{"id": 2, "name": "Action"},
				},
			},
		}))

	// Mock movies endpoint (GetTrendingMovies)
	httpmock.RegisterResponder("GET", mockBaseURL+"/movies",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"status": "success",
			"data": []map[string]any{
				{
					"id":            101,
					"name":          "Test Movie",
					"overview":      "A test movie",
					"releaseDate":   "2024-01-01",
					"image":         "https://test.com/movie.jpg",
					"averageRating": 9.0,
				},
			},
		}))

	// Mock season details endpoint (extended)
	httpmock.RegisterResponder("GET", mockBaseURL+"/seasons/555/extended",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"status": "success",
			"data": map[string]any{
				"id":       555,
				"name":     "Season 1",
				"overview": "First season",
				"number":   1,
				"episodes": []map[string]any{
					{
						"id":   1001,
						"name": "Episode 1",
					},
					{
						"id":   1002,
						"name": "Episode 2",
					},
				},
			},
		}))

	ctx := context.Background()

	t.Run("GetTrendingTV", func(t *testing.T) {
		result, err := GetTrendingTV(ctx)
		require.NoError(t, err)
		assert.NotNil(t, result)

		// Check result structure - functions return map[string]any
		resultMap, ok := result.(map[string]any)
		require.True(t, ok, "Result should be map[string]any")

		assert.Equal(t, "success", resultMap["status"])

		data, ok := resultMap["data"].([]any)
		require.True(t, ok, "Data should be slice")
		assert.Len(t, data, 2)

		// Check first series
		series, ok := data[0].(map[string]any)
		require.True(t, ok, "Series should be map")
		assert.Equal(t, float64(123), series["id"])
		assert.Equal(t, "Test TV Series", series["name"])
	})

	t.Run("SearchTV", func(t *testing.T) {
		result, err := SearchTV(ctx, "test")
		require.NoError(t, err)
		assert.NotNil(t, result)

		// Check result structure
		resultMap, ok := result.(map[string]any)
		require.True(t, ok, "Result should be map[string]any")

		assert.Equal(t, "success", resultMap["status"])

		data, ok := resultMap["data"].([]any)
		require.True(t, ok, "Data should be slice")
		assert.Len(t, data, 1)

		// Check searched series
		series, ok := data[0].(map[string]any)
		require.True(t, ok, "Series should be map")
		assert.Equal(t, float64(789), series["id"])
		assert.Equal(t, "Searched Series", series["name"])
	})

	t.Run("GetTVShow", func(t *testing.T) {
		result, err := GetTVShow(ctx, 123)
		require.NoError(t, err)
		assert.NotNil(t, result)

		// Check result structure
		resultMap, ok := result.(map[string]any)
		require.True(t, ok, "Result should be map[string]any")

		assert.Equal(t, "success", resultMap["status"])

		data, ok := resultMap["data"].(map[string]any)
		require.True(t, ok, "Data should be map")
		assert.Equal(t, float64(123), data["id"])
		assert.Equal(t, "Test TV Series", data["name"])

		// Check genres
		genres, ok := data["genres"].([]any)
		require.True(t, ok, "Genres should be slice")
		assert.Len(t, genres, 2)

		genre, ok := genres[0].(map[string]any)
		require.True(t, ok, "Genre should be map")
		assert.Equal(t, "Drama", genre["name"])
	})

	t.Run("GetTrendingMovies", func(t *testing.T) {
		result, err := GetTrendingMovies(ctx)
		require.NoError(t, err)
		assert.NotNil(t, result)

		// Check result structure
		resultMap, ok := result.(map[string]any)
		require.True(t, ok, "Result should be map[string]any")

		assert.Equal(t, "success", resultMap["status"])

		data, ok := resultMap["data"].([]any)
		require.True(t, ok, "Data should be slice")
		assert.Len(t, data, 1)

		// Check movie
		movie, ok := data[0].(map[string]any)
		require.True(t, ok, "Movie should be map")
		assert.Equal(t, float64(101), movie["id"])
		assert.Equal(t, "Test Movie", movie["name"])
	})

	t.Run("GetTVSeason", func(t *testing.T) {
		result, err := GetTVSeason(ctx, 555)
		require.NoError(t, err)
		assert.NotNil(t, result)

		// Check result structure
		resultMap, ok := result.(map[string]any)
		require.True(t, ok, "Result should be map[string]any")

		assert.Equal(t, "success", resultMap["status"])

		data, ok := resultMap["data"].(map[string]any)
		require.True(t, ok, "Data should be map")
		assert.Equal(t, float64(555), data["id"])
		assert.Equal(t, "Season 1", data["name"])
		assert.Equal(t, float64(1), data["number"])

		// Check episodes
		episodes, ok := data["episodes"].([]any)
		require.True(t, ok, "Episodes should be slice")
		assert.Len(t, episodes, 2)
	})

	// Verify expected calls were made (login + 5 API calls)
	assert.Equal(t, 6, httpmock.GetTotalCallCount())
}

func TestTVDBErrorHandling(t *testing.T) {
	// Test with no API key - this should fail during authentication
	InitTVDB("", "http://localhost:99999") // Use a port that's definitely not listening

	ctx := context.Background()

	t.Run("NoAPIKey", func(t *testing.T) {
		result, err := GetTrendingTV(ctx)
		assert.Error(t, err)
		assert.Nil(t, result)
		// Since we have no API key, the error should be about authentication or connection
		assert.True(t,
			strings.Contains(err.Error(), "API key not configured") ||
				strings.Contains(err.Error(), "connection refused") ||
				strings.Contains(err.Error(), "invalid port") ||
				strings.Contains(err.Error(), "failed to authenticate") ||
				strings.Contains(err.Error(), "failed to make request"),
			"Error should be about API key, connection, or authentication: %v", err)
	})
}

func TestTVDBServerErrors(t *testing.T) {
	// Activate HTTP mock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Initialize TVDB with mock server
	mockBaseURL := "http://mock-tvdb-api.com/v4"
	InitTVDB("test_api_key", mockBaseURL)

	// Mock login failure
	httpmock.RegisterResponder("POST", mockBaseURL+"/login",
		httpmock.NewStringResponder(401, "Unauthorized"))

	ctx := context.Background()

	t.Run("AuthenticationFailure", func(t *testing.T) {
		result, err := GetTrendingTV(ctx)
		assert.Error(t, err)
		assert.Nil(t, result)
		// The error may be about authentication failure or the subsequent API call failure
		assert.True(t,
			strings.Contains(err.Error(), "authentication failed") ||
				strings.Contains(err.Error(), "no responder found") ||
				strings.Contains(err.Error(), "status 401"),
			"Error should be about authentication failure: %v", err)
	})
}
