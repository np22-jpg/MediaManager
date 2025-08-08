package tvdb

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTVDBWithMockServer(t *testing.T) {
	// Create a test server that will handle all API requests
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == "POST" && strings.Contains(r.URL.Path, "/login"):
			response := map[string]any{
				"status": "success",
				"data": map[string]any{
					"token": "mock_jwt_token_here",
				},
			}
			_ = json.NewEncoder(w).Encode(response)

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/series") && !strings.Contains(r.URL.Path, "/extended"):
			response := map[string]any{
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
			}
			_ = json.NewEncoder(w).Encode(response)

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/search"):
			response := map[string]any{
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
			}
			_ = json.NewEncoder(w).Encode(response)

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/series/123/extended"):
			response := map[string]any{
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
			}
			_ = json.NewEncoder(w).Encode(response)

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/movies"):
			response := map[string]any{
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
			}
			_ = json.NewEncoder(w).Encode(response)

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/seasons/555/extended"):
			response := map[string]any{
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
			}
			_ = json.NewEncoder(w).Encode(response)

		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	// Initialize TVDB with test server URL
	InitTVDB("test_api_key", server.URL+"/v4")

	ctx := context.Background()

	t.Run("GetTrendingTV", func(t *testing.T) {
		result, err := GetTrendingTV(ctx)
		if err != nil {
			t.Fatalf("GetTrendingTV failed: %v", err)
		}
		if result == nil {
			t.Fatal("Result should not be nil")
		}

		// Check result structure - functions return map[string]any
		resultMap, ok := result.(map[string]any)
		if !ok {
			t.Fatal("Result should be map[string]any")
		}

		if resultMap["status"] != "success" {
			t.Errorf("Expected status 'success', got %v", resultMap["status"])
		}

		data, ok := resultMap["data"].([]any)
		if !ok {
			t.Fatal("Data should be slice")
		}
		if len(data) != 2 {
			t.Errorf("Expected 2 series, got %d", len(data))
		}

		// Check first series
		series, ok := data[0].(map[string]any)
		if !ok {
			t.Fatal("Series should be map")
		}
		if series["id"] != float64(123) {
			t.Errorf("Expected series id 123, got %v", series["id"])
		}
		if series["name"] != "Test TV Series" {
			t.Errorf("Expected series name 'Test TV Series', got %v", series["name"])
		}
	})

	t.Run("SearchTV", func(t *testing.T) {
		result, err := SearchTV(ctx, "test")
		if err != nil {
			t.Fatalf("SearchTV failed: %v", err)
		}
		if result == nil {
			t.Fatal("Result should not be nil")
		}

		// Check result structure
		resultMap, ok := result.(map[string]any)
		if !ok {
			t.Fatal("Result should be map[string]any")
		}

		if resultMap["status"] != "success" {
			t.Errorf("Expected status 'success', got %v", resultMap["status"])
		}

		data, ok := resultMap["data"].([]any)
		if !ok {
			t.Fatal("Data should be slice")
		}
		if len(data) != 1 {
			t.Errorf("Expected 1 series, got %d", len(data))
		}

		// Check searched series
		series, ok := data[0].(map[string]any)
		if !ok {
			t.Fatal("Series should be map")
		}
		if series["id"] != float64(789) {
			t.Errorf("Expected series id 789, got %v", series["id"])
		}
		if series["name"] != "Searched Series" {
			t.Errorf("Expected series name 'Searched Series', got %v", series["name"])
		}
	})

	t.Run("GetTVShow", func(t *testing.T) {
		result, err := GetTVShow(ctx, 123)
		if err != nil {
			t.Fatalf("GetTVShow failed: %v", err)
		}
		if result == nil {
			t.Fatal("Result should not be nil")
		}

		// Check result structure
		resultMap, ok := result.(map[string]any)
		if !ok {
			t.Fatal("Result should be map[string]any")
		}

		if resultMap["status"] != "success" {
			t.Errorf("Expected status 'success', got %v", resultMap["status"])
		}

		data, ok := resultMap["data"].(map[string]any)
		if !ok {
			t.Fatal("Data should be map")
		}
		if data["id"] != float64(123) {
			t.Errorf("Expected series id 123, got %v", data["id"])
		}
		if data["name"] != "Test TV Series" {
			t.Errorf("Expected series name 'Test TV Series', got %v", data["name"])
		}

		// Check genres
		genres, ok := data["genres"].([]any)
		if !ok {
			t.Fatal("Genres should be slice")
		}
		if len(genres) != 2 {
			t.Errorf("Expected 2 genres, got %d", len(genres))
		}

		genre, ok := genres[0].(map[string]any)
		if !ok {
			t.Fatal("Genre should be map")
		}
		if genre["name"] != "Drama" {
			t.Errorf("Expected genre name 'Drama', got %v", genre["name"])
		}
	})

	t.Run("GetTrendingMovies", func(t *testing.T) {
		result, err := GetTrendingMovies(ctx)
		if err != nil {
			t.Fatalf("GetTrendingMovies failed: %v", err)
		}
		if result == nil {
			t.Fatal("Result should not be nil")
		}

		// Check result structure
		resultMap, ok := result.(map[string]any)
		if !ok {
			t.Fatal("Result should be map[string]any")
		}

		if resultMap["status"] != "success" {
			t.Errorf("Expected status 'success', got %v", resultMap["status"])
		}

		data, ok := resultMap["data"].([]any)
		if !ok {
			t.Fatal("Data should be slice")
		}
		if len(data) != 1 {
			t.Errorf("Expected 1 movie, got %d", len(data))
		}

		// Check movie
		movie, ok := data[0].(map[string]any)
		if !ok {
			t.Fatal("Movie should be map")
		}
		if movie["id"] != float64(101) {
			t.Errorf("Expected movie id 101, got %v", movie["id"])
		}
		if movie["name"] != "Test Movie" {
			t.Errorf("Expected movie name 'Test Movie', got %v", movie["name"])
		}
	})

	t.Run("GetTVSeason", func(t *testing.T) {
		result, err := GetTVSeason(ctx, 555)
		if err != nil {
			t.Fatalf("GetTVSeason failed: %v", err)
		}
		if result == nil {
			t.Fatal("Result should not be nil")
		}

		// Check result structure
		resultMap, ok := result.(map[string]any)
		if !ok {
			t.Fatal("Result should be map[string]any")
		}

		if resultMap["status"] != "success" {
			t.Errorf("Expected status 'success', got %v", resultMap["status"])
		}

		data, ok := resultMap["data"].(map[string]any)
		if !ok {
			t.Fatal("Data should be map")
		}
		if data["id"] != float64(555) {
			t.Errorf("Expected season id 555, got %v", data["id"])
		}
		if data["name"] != "Season 1" {
			t.Errorf("Expected season name 'Season 1', got %v", data["name"])
		}
		if data["number"] != float64(1) {
			t.Errorf("Expected season number 1, got %v", data["number"])
		}

		// Check episodes
		episodes, ok := data["episodes"].([]any)
		if !ok {
			t.Fatal("Episodes should be slice")
		}
		if len(episodes) != 2 {
			t.Errorf("Expected 2 episodes, got %d", len(episodes))
		}
	})
}

func TestTVDBErrorHandling(t *testing.T) {
	// Test with no API key - this should fail during authentication
	InitTVDB("", "http://localhost:99999") // Use a port that's definitely not listening

	ctx := context.Background()

	t.Run("NoAPIKey", func(t *testing.T) {
		result, err := GetTrendingTV(ctx)
		if err == nil {
			t.Error("Expected error but got none")
		}
		if result != nil {
			t.Error("Expected nil result but got non-nil")
		}
		// Since we have no API key, the error should be about authentication or connection
		errorMessage := err.Error()
		if !strings.Contains(errorMessage, "API key not configured") &&
			!strings.Contains(errorMessage, "connection refused") &&
			!strings.Contains(errorMessage, "invalid port") &&
			!strings.Contains(errorMessage, "failed to authenticate") &&
			!strings.Contains(errorMessage, "failed to make request") {
			t.Errorf("Error should be about API key, connection, or authentication, got: %v", err)
		}
	})
}

func TestTVDBServerErrors(t *testing.T) {
	// Create a test server that returns 401 for login
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && strings.Contains(r.URL.Path, "/login") {
			w.WriteHeader(401)
			_, _ = w.Write([]byte("Unauthorized"))
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	// Initialize TVDB with test server that will fail authentication
	InitTVDB("test_api_key", server.URL+"/v4")

	ctx := context.Background()

	t.Run("AuthenticationFailure", func(t *testing.T) {
		result, err := GetTrendingTV(ctx)
		if err == nil {
			t.Error("Expected error but got none")
		}
		if result != nil {
			t.Error("Expected nil result but got non-nil")
		}
		// The error may be about authentication failure or the subsequent API call failure
		errorMessage := err.Error()
		if !strings.Contains(errorMessage, "authentication failed") &&
			!strings.Contains(errorMessage, "no responder found") &&
			!strings.Contains(errorMessage, "status 401") &&
			!strings.Contains(errorMessage, "status 404") {
			t.Errorf("Error should be about authentication failure, got: %v", err)
		}
	})
}
