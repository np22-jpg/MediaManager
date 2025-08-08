package anidb

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClient_GetAnimeInfo(t *testing.T) {
	// XML test data
	animeXML := `<?xml version="1.0" encoding="UTF-8"?>
<anime id="1">
	<type>TV Series</type>
	<episodecount>13</episodecount>
	<startdate>1999-01-03</startdate>
	<enddate>1999-03-28</enddate>
	<description>This is a test anime description.</description>
</anime>`

	hotAnimeXML := `<?xml version="1.0" encoding="UTF-8"?>
<hotanime>
	<anime id="2">
		<type>Movie</type>
		<episodecount>1</episodecount>
		<startdate>2000-01-01</startdate>
		<enddate>2000-01-01</enddate>
		<description>This is a hot anime description.</description>
	</anime>
</hotanime>`

	// Mock server setup
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryType := r.URL.Query().Get("request")

		switch queryType {
		case "anime":
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(200)
			_, _ = w.Write([]byte(animeXML))
		case "hotanime":
			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(200)
			_, _ = w.Write([]byte(hotAnimeXML))
		default:
			w.WriteHeader(400)
			_, _ = w.Write([]byte("Bad Request"))
		}
	}))
	defer server.Close()

	// Initialize AniDB with mock server
	InitAniDB(server.URL, "testclient", "1")

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

		anime, ok := result.(AnimeInfo)
		if !ok {
			t.Fatalf("Result should be AnimeInfo, got %T", result)
		}
		if anime.ID != 1 {
			t.Errorf("Expected ID 1, got %d", anime.ID)
		}
		if anime.Type != "TV Series" {
			t.Errorf("Expected Type 'TV Series', got %s", anime.Type)
		}
		if anime.EpisodeCount != 13 {
			t.Errorf("Expected EpisodeCount 13, got %d", anime.EpisodeCount)
		}
		if anime.StartDate != "1999-01-03" {
			t.Errorf("Expected StartDate '1999-01-03', got %s", anime.StartDate)
		}
		if anime.EndDate != "1999-03-28" {
			t.Errorf("Expected EndDate '1999-03-28', got %s", anime.EndDate)
		}
		if !strings.Contains(anime.Description, "test anime") {
			t.Errorf("Expected Description to contain 'test anime', got %s", anime.Description)
		}
	})

	t.Run("GetHotAnime", func(t *testing.T) {
		result, err := GetHotAnime(ctx)
		if err != nil {
			t.Errorf("GetHotAnime() error = %v", err)
			return
		}
		if result == nil {
			t.Error("Expected result to not be nil")
			return
		}

		animeList, ok := result.([]AnimeInfo)
		if !ok {
			t.Fatalf("Result should be []AnimeInfo, got %T", result)
		}
		if len(animeList) == 0 {
			t.Error("Expected at least one anime in the hot anime list")
			return
		}
		if animeList[0].ID != 2 {
			t.Errorf("Expected ID 2, got %d", animeList[0].ID)
		}
		if animeList[0].Type != "Movie" {
			t.Errorf("Expected Type 'Movie', got %s", animeList[0].Type)
		}
	})
}

func TestAniDBErrorHandling(t *testing.T) {
	// Test with no client configured
	InitAniDB("", "", "")

	ctx := context.Background()

	t.Run("NoClientConfigured", func(t *testing.T) {
		result, err := GetAnimeByID(ctx, 1)
		if err == nil {
			t.Error("Expected error when no client configured")
		}
		if result != nil {
			t.Error("Expected nil result when no client configured")
		}
		if !strings.Contains(err.Error(), "client") {
			t.Errorf("Expected error to contain 'client', got: %s", err.Error())
		}
	})
}

func TestAniDBServerErrors(t *testing.T) {
	// Mock server that returns server error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		_, _ = w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	// Initialize AniDB with mock server
	InitAniDB(server.URL, "testclient", "1")

	ctx := context.Background()

	t.Run("ServerError", func(t *testing.T) {
		result, err := GetAnimeByID(ctx, 1)
		if err == nil {
			t.Error("Expected error for server error")
		}
		if result != nil {
			t.Error("Expected nil result for server error")
		}
	})
}

func TestAniDBBannedError(t *testing.T) {
	// Mock banned error XML
	bannedXML := `<?xml version="1.0"?><error>Banned</error>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(bannedXML))
	}))
	defer server.Close()

	// Initialize AniDB with mock server
	InitAniDB(server.URL, "testclient", "1")

	ctx := context.Background()

	t.Run("BannedError", func(t *testing.T) {
		result, err := GetAnimeByID(ctx, 1)
		if err == nil {
			t.Error("Expected error for banned response")
		}
		if result != nil {
			t.Error("Expected nil result for banned response")
		}
		if !strings.Contains(err.Error(), "Banned") {
			t.Errorf("Expected error to contain 'Banned', got: %s", err.Error())
		}
	})
}
