package jikan

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestJikan_GetAnimeByID(t *testing.T) {
	// Mock Jikan API response
	mockResponse := `{
		"data": {
			"mal_id": 1,
			"title": "Cowboy Bebop",
			"title_japanese": "カウボーイビバップ",
			"type": "TV",
			"episodes": 26,
			"status": "Finished Airing",
			"score": 8.78,
			"synopsis": "In the year 2071, humanity has colonized several of the planets and moons of the solar system...",
			"images": {
				"jpg": {
					"image_url": "https://cdn.myanimelist.net/images/anime/4/19644.jpg"
				}
			},
			"genres": [
				{"mal_id": 1, "name": "Action"},
				{"mal_id": 46, "name": "Award Winning"}
			]
		}
	}`

	// Mock server setup
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Initialize Jikan with mock server
	InitJikan(server.URL)

	ctx := context.Background()

	t.Run("GetAnimeByID", func(t *testing.T) {
		result, err := GetAnimeByID(ctx, 1)
		if err != nil {
			t.Errorf("GetAnimeByID() error = %v", err)
			return
		}
		if result == nil {
			t.Error("Expected result to not be nil")
			return
		}

		// The result should be the entire API response
		response, ok := result.(map[string]any)
		if !ok {
			t.Fatalf("Result should be map[string]any, got %T", result)
		}

		data, ok := response["data"].(map[string]any)
		if !ok {
			t.Fatalf("Expected data field, got %T", response["data"])
		}

		if data["mal_id"].(float64) != 1 {
			t.Errorf("Expected mal_id 1, got %v", data["mal_id"])
		}
		if data["title"].(string) != "Cowboy Bebop" {
			t.Errorf("Expected title 'Cowboy Bebop', got %s", data["title"])
		}
	})
}

func TestJikan_RateLimiting(t *testing.T) {
	// This test ensures we're respecting rate limits
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"data": {"mal_id": 1, "title": "Test"}}`))
	}))
	defer server.Close()

	InitJikan(server.URL)

	// Make multiple requests and ensure they're properly spaced
	start := time.Now()
	for i := 0; i < 3; i++ {
		_, err := makeRequest("/anime/1")
		if err != nil {
			t.Errorf("Request %d failed: %v", i, err)
		}
	}
	duration := time.Since(start)

	// Should take at least 700ms for 3 requests (350ms spacing * 2 gaps)
	if duration < 700*time.Millisecond {
		t.Errorf("Requests completed too quickly, may not be respecting rate limits: %v", duration)
	}
}
