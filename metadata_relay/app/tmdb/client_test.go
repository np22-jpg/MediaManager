package tmdb

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTMDBWithMockServer(t *testing.T) {
	// Activate HTTP mock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Initialize TMDB with mock server
	mockBaseURL := "http://mock-tmdb-api.com/3"
	InitTMDB("test_api_key", mockBaseURL)

	// Mock trending movies endpoint
	httpmock.RegisterResponder("GET", mockBaseURL+"/trending/movie/week",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"page":          1,
			"total_pages":   1,
			"total_results": 2,
			"results": []map[string]any{
				{
					"id":            123,
					"title":         "Test Movie 1",
					"overview":      "A test movie",
					"release_date":  "2024-01-01",
					"poster_path":   "/test1.jpg",
					"backdrop_path": "/test1_backdrop.jpg",
					"vote_average":  8.5,
				},
				{
					"id":            456,
					"title":         "Test Movie 2",
					"overview":      "Another test movie",
					"release_date":  "2024-02-01",
					"poster_path":   "/test2.jpg",
					"backdrop_path": "/test2_backdrop.jpg",
					"vote_average":  7.8,
				},
			},
		}))

	// Mock search movies endpoint
	httpmock.RegisterResponder("GET", mockBaseURL+"/search/movie",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"page":          1,
			"total_pages":   1,
			"total_results": 1,
			"results": []map[string]any{
				{
					"id":           789,
					"title":        "Searched Movie",
					"overview":     "A searched movie",
					"release_date": "2024-03-01",
					"poster_path":  "/searched.jpg",
					"vote_average": 9.0,
				},
			},
		}))

	// Mock movie details endpoint
	httpmock.RegisterResponder("GET", mockBaseURL+"/movie/123",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"id":           123,
			"title":        "Test Movie 1",
			"overview":     "A detailed test movie",
			"release_date": "2024-01-01",
			"runtime":      120,
			"vote_average": 8.5,
			"genres": []map[string]any{
				{"id": 28, "name": "Action"},
				{"id": 12, "name": "Adventure"},
			},
		}))

	// Mock trending TV shows endpoint
	httpmock.RegisterResponder("GET", mockBaseURL+"/trending/tv/week",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"page":          1,
			"total_pages":   1,
			"total_results": 1,
			"results": []map[string]any{
				{
					"id":             101,
					"name":           "Test TV Show",
					"overview":       "A test TV show",
					"first_air_date": "2024-01-01",
					"poster_path":    "/test_tv.jpg",
					"vote_average":   8.0,
				},
			},
		}))

	ctx := context.Background()

	t.Run("GetTrendingMovies", func(t *testing.T) {
		result, err := GetTrendingMovies(ctx)
		require.NoError(t, err)
		assert.NotNil(t, result)

		// Check result structure - functions return map[string]any, not typed structs
		resultMap, ok := result.(map[string]any)
		require.True(t, ok, "Result should be map[string]any")

		assert.Equal(t, float64(1), resultMap["page"])
		assert.Equal(t, float64(2), resultMap["total_results"])

		results, ok := resultMap["results"].([]any)
		require.True(t, ok, "Results should be slice")
		assert.Len(t, results, 2)

		// Check first movie
		movie, ok := results[0].(map[string]any)
		require.True(t, ok, "Movie should be map")
		assert.Equal(t, float64(123), movie["id"])
		assert.Equal(t, "Test Movie 1", movie["title"])
	})

	t.Run("SearchMovies", func(t *testing.T) {
		result, err := SearchMovies(ctx, "test", 1)
		require.NoError(t, err)
		assert.NotNil(t, result)

		// Check result structure - functions return map[string]any, not typed structs
		resultMap, ok := result.(map[string]any)
		require.True(t, ok, "Result should be map[string]any")

		assert.Equal(t, float64(1), resultMap["page"])
		assert.Equal(t, float64(1), resultMap["total_results"])

		results, ok := resultMap["results"].([]any)
		require.True(t, ok, "Results should be slice")
		assert.Len(t, results, 1)

		// Check searched movie
		movie, ok := results[0].(map[string]any)
		require.True(t, ok, "Movie should be map")
		assert.Equal(t, float64(789), movie["id"])
		assert.Equal(t, "Searched Movie", movie["title"])
	})

	t.Run("GetMovie", func(t *testing.T) {
		result, err := GetMovie(ctx, 123)
		require.NoError(t, err)
		assert.NotNil(t, result)

		// Check result structure
		movieMap, ok := result.(map[string]any)
		require.True(t, ok, "Result should be a map")

		assert.Equal(t, float64(123), movieMap["id"])
		assert.Equal(t, "Test Movie 1", movieMap["title"])
		assert.Equal(t, float64(120), movieMap["runtime"])

		// Check genres
		genres, ok := movieMap["genres"].([]any)
		require.True(t, ok, "Genres should be a slice of any")
		assert.Len(t, genres, 2)

		// Check first genre
		genre, ok := genres[0].(map[string]any)
		require.True(t, ok, "Genre should be map")
		assert.Equal(t, "Action", genre["name"])
	})

	t.Run("GetTrendingTV", func(t *testing.T) {
		result, err := GetTrendingTV(ctx)
		require.NoError(t, err)
		assert.NotNil(t, result)

		// Check result structure - functions return map[string]any, not typed structs
		resultMap, ok := result.(map[string]any)
		require.True(t, ok, "Result should be map[string]any")

		assert.Equal(t, float64(1), resultMap["page"])
		assert.Equal(t, float64(1), resultMap["total_results"])

		results, ok := resultMap["results"].([]any)
		require.True(t, ok, "Results should be slice")
		assert.Len(t, results, 1)

		// Check TV show
		show, ok := results[0].(map[string]any)
		require.True(t, ok, "Show should be map")
		assert.Equal(t, float64(101), show["id"])
		assert.Equal(t, "Test TV Show", show["name"])
	})

	// Verify all expected calls were made
	assert.Equal(t, 4, httpmock.GetTotalCallCount())
}

func TestTMDBErrorHandling(t *testing.T) {
	// Test with no API key
	InitTMDB("", "http://test.com")

	ctx := context.Background()

	t.Run("NoAPIKey", func(t *testing.T) {
		result, err := GetTrendingMovies(ctx)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "API key not configured")
	})
}

func TestTMDBServerErrors(t *testing.T) {
	// Activate HTTP mock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Initialize TMDB with mock server
	mockBaseURL := "http://mock-tmdb-api.com/3"
	InitTMDB("test_api_key", mockBaseURL)

	// Mock server error
	httpmock.RegisterResponder("GET", mockBaseURL+"/trending/movie/week",
		httpmock.NewStringResponder(500, "Internal Server Error"))

	ctx := context.Background()

	t.Run("ServerError", func(t *testing.T) {
		result, err := GetTrendingMovies(ctx)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
