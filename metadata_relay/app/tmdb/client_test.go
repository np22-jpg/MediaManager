package tmdb

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTMDBWithMockServer(t *testing.T) {
	// Create a test server that will handle all API requests
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case strings.Contains(r.URL.Path, "/trending/movie/week"):
			response := map[string]any{
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
			}
			_ = json.NewEncoder(w).Encode(response)

		case strings.Contains(r.URL.Path, "/search/movie"):
			response := map[string]any{
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
			}
			_ = json.NewEncoder(w).Encode(response)

		case strings.Contains(r.URL.Path, "/movie/123"):
			response := map[string]any{
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
			}
			_ = json.NewEncoder(w).Encode(response)

		case strings.Contains(r.URL.Path, "/trending/tv/week"):
			response := map[string]any{
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
			}
			_ = json.NewEncoder(w).Encode(response)

		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	// Initialize TMDB with test server URL
	InitTMDB("test_api_key", server.URL+"/3")

	ctx := context.Background()

	t.Run("GetTrendingMovies", func(t *testing.T) {
		result, err := GetTrendingMovies(ctx)
		if err != nil {
			t.Fatalf("GetTrendingMovies failed: %v", err)
		}
		if result == nil {
			t.Fatal("Result should not be nil")
		}

		// Check result structure - functions return map[string]any, not typed structs
		resultMap, ok := result.(map[string]any)
		if !ok {
			t.Fatal("Result should be map[string]any")
		}

		if resultMap["page"] != float64(1) {
			t.Errorf("Expected page 1, got %v", resultMap["page"])
		}
		if resultMap["total_results"] != float64(2) {
			t.Errorf("Expected total_results 2, got %v", resultMap["total_results"])
		}

		results, ok := resultMap["results"].([]any)
		if !ok {
			t.Fatal("Results should be slice")
		}
		if len(results) != 2 {
			t.Errorf("Expected 2 results, got %d", len(results))
		}

		// Check first movie
		movie, ok := results[0].(map[string]any)
		if !ok {
			t.Fatal("Movie should be map")
		}
		if movie["id"] != float64(123) {
			t.Errorf("Expected movie id 123, got %v", movie["id"])
		}
		if movie["title"] != "Test Movie 1" {
			t.Errorf("Expected movie title 'Test Movie 1', got %v", movie["title"])
		}
	})

	t.Run("SearchMovies", func(t *testing.T) {
		result, err := SearchMovies(ctx, "test", 1)
		if err != nil {
			t.Fatalf("SearchMovies failed: %v", err)
		}
		if result == nil {
			t.Fatal("Result should not be nil")
		}

		// Check result structure - functions return map[string]any, not typed structs
		resultMap, ok := result.(map[string]any)
		if !ok {
			t.Fatal("Result should be map[string]any")
		}

		if resultMap["page"] != float64(1) {
			t.Errorf("Expected page 1, got %v", resultMap["page"])
		}
		if resultMap["total_results"] != float64(1) {
			t.Errorf("Expected total_results 1, got %v", resultMap["total_results"])
		}

		results, ok := resultMap["results"].([]any)
		if !ok {
			t.Fatal("Results should be slice")
		}
		if len(results) != 1 {
			t.Errorf("Expected 1 result, got %d", len(results))
		}

		// Check searched movie
		movie, ok := results[0].(map[string]any)
		if !ok {
			t.Fatal("Movie should be map")
		}
		if movie["id"] != float64(789) {
			t.Errorf("Expected movie id 789, got %v", movie["id"])
		}
		if movie["title"] != "Searched Movie" {
			t.Errorf("Expected movie title 'Searched Movie', got %v", movie["title"])
		}
	})

	t.Run("GetMovie", func(t *testing.T) {
		result, err := GetMovie(ctx, 123)
		if err != nil {
			t.Fatalf("GetMovie failed: %v", err)
		}
		if result == nil {
			t.Fatal("Result should not be nil")
		}

		// Check result structure
		movieMap, ok := result.(map[string]any)
		if !ok {
			t.Fatal("Result should be a map")
		}

		if movieMap["id"] != float64(123) {
			t.Errorf("Expected movie id 123, got %v", movieMap["id"])
		}
		if movieMap["title"] != "Test Movie 1" {
			t.Errorf("Expected movie title 'Test Movie 1', got %v", movieMap["title"])
		}
		if movieMap["runtime"] != float64(120) {
			t.Errorf("Expected runtime 120, got %v", movieMap["runtime"])
		}

		// Check genres
		genres, ok := movieMap["genres"].([]any)
		if !ok {
			t.Fatal("Genres should be a slice of any")
		}
		if len(genres) != 2 {
			t.Errorf("Expected 2 genres, got %d", len(genres))
		}

		// Check first genre
		genre, ok := genres[0].(map[string]any)
		if !ok {
			t.Fatal("Genre should be map")
		}
		if genre["name"] != "Action" {
			t.Errorf("Expected genre name 'Action', got %v", genre["name"])
		}
	})

	t.Run("GetTrendingTV", func(t *testing.T) {
		result, err := GetTrendingTV(ctx)
		if err != nil {
			t.Fatalf("GetTrendingTV failed: %v", err)
		}
		if result == nil {
			t.Fatal("Result should not be nil")
		}

		// Check result structure - functions return map[string]any, not typed structs
		resultMap, ok := result.(map[string]any)
		if !ok {
			t.Fatal("Result should be map[string]any")
		}

		if resultMap["page"] != float64(1) {
			t.Errorf("Expected page 1, got %v", resultMap["page"])
		}
		if resultMap["total_results"] != float64(1) {
			t.Errorf("Expected total_results 1, got %v", resultMap["total_results"])
		}

		results, ok := resultMap["results"].([]any)
		if !ok {
			t.Fatal("Results should be slice")
		}
		if len(results) != 1 {
			t.Errorf("Expected 1 result, got %d", len(results))
		}

		// Check TV show
		show, ok := results[0].(map[string]any)
		if !ok {
			t.Fatal("Show should be map")
		}
		if show["id"] != float64(101) {
			t.Errorf("Expected show id 101, got %v", show["id"])
		}
		if show["name"] != "Test TV Show" {
			t.Errorf("Expected show name 'Test TV Show', got %v", show["name"])
		}
	})
}

func TestTMDBErrorHandling(t *testing.T) {
	// Test with no API key
	InitTMDB("", "http://test.com")

	ctx := context.Background()

	t.Run("NoAPIKey", func(t *testing.T) {
		result, err := GetTrendingMovies(ctx)
		if err == nil {
			t.Error("Expected error but got none")
		}
		if result != nil {
			t.Error("Expected nil result but got non-nil")
		}
		if !strings.Contains(err.Error(), "API key not configured") {
			t.Errorf("Expected error to contain 'API key not configured', got: %v", err)
		}
	})
}

func TestTMDBServerErrors(t *testing.T) {
	// Create a test server that returns 500 error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/trending/movie/week") {
			w.WriteHeader(500)
			_, _ = w.Write([]byte("Internal Server Error"))
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	// Initialize TMDB with test server
	InitTMDB("test_api_key", server.URL+"/3")

	ctx := context.Background()

	t.Run("ServerError", func(t *testing.T) {
		result, err := GetTrendingMovies(ctx)
		if err == nil {
			t.Error("Expected error but got none")
		}
		if result != nil {
			t.Error("Expected nil result but got non-nil")
		}
	})
}
