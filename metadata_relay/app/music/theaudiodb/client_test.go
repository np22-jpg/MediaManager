package theaudiodb

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockHTTP implements HTTPDoer for tests
type mockHTTP struct{ handler http.Handler }

func (m mockHTTP) Do(req *http.Request) (*http.Response, error) {
	rr := httptest.NewRecorder()
	m.handler.ServeHTTP(rr, req)
	return rr.Result(), nil
}

func TestSearchArtist_SimplifiesResult(t *testing.T) {
	// Mock TheAudioDB response
	mux := http.NewServeMux()
	mux.HandleFunc("/2/search.php", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{
			"artists": []map[string]any{{
				"idArtist":       "123",
				"strArtist":      "Test Artist",
				"strBiographyEN": "Bio",
				"strWebsite":     "https://example.com",
				"strGenre":       "Rock",
				"intFormedYear":  "2001",
			}},
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	c := &Client{BaseURL: "http://example", APIKey: "2", HTTP: mockHTTP{handler: mux}}
	got, err := c.SearchArtist(context.Background(), "Test")
	if err != nil {
		t.Fatalf("SearchArtist error: %v", err)
	}
	if got["name"] != "Test Artist" {
		t.Fatalf("unexpected name: %v", got["name"])
	}
	if got["genre"] != "Rock" {
		t.Fatalf("unexpected genre: %v", got["genre"])
	}
}

func TestSearchArtistByMBID_SimplifiesResult(t *testing.T) {
	// Mock TheAudioDB MBID artist response
	mux := http.NewServeMux()
	mux.HandleFunc("/2/artist-mb.php", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{
			"artists": []map[string]any{{
				"idArtist":         "123",
				"strArtist":        "Test Artist",
				"strBiographyEN":   "Bio",
				"strWebsite":       "https://example.com",
				"strGenre":         "Rock",
				"intFormedYear":    "2001",
				"strMusicBrainzID": "5b11f4ce-a62d-471e-81fc-a69a8278c7da",
				"strArtistThumb":   "https://example.com/thumb.jpg",
				"strArtistFanart":  "https://example.com/fanart.jpg",
				"strCountry":       "United States",
			}},
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	c := &Client{BaseURL: "http://example", APIKey: "2", HTTP: mockHTTP{handler: mux}}
	got, err := c.SearchArtistByMBID(context.Background(), "5b11f4ce-a62d-471e-81fc-a69a8278c7da")
	if err != nil {
		t.Fatalf("SearchArtistByMBID error: %v", err)
	}
	if got["name"] != "Test Artist" {
		t.Fatalf("unexpected name: %v", got["name"])
	}
	if got["mbid"] != "5b11f4ce-a62d-471e-81fc-a69a8278c7da" {
		t.Fatalf("unexpected mbid: %v", got["mbid"])
	}
	if got["country"] != "United States" {
		t.Fatalf("unexpected country: %v", got["country"])
	}
}

func TestSearchAlbumByMBID_SimplifiesResult(t *testing.T) {
	// Mock TheAudioDB MBID album response
	mux := http.NewServeMux()
	mux.HandleFunc("/2/album-mb.php", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{
			"album": []map[string]any{{
				"idAlbum":          "456",
				"idMBAlbum":        "d7d2c46c-c9b6-4ff1-bd0a-b3e36f57e0d0",
				"strAlbum":         "Test Album",
				"strArtist":        "Test Artist",
				"intYearReleased":  "2005",
				"strGenre":         "Rock",
				"strAlbumThumb":    "https://example.com/album_thumb.jpg",
				"strDescriptionEN": "A great test album",
			}},
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	c := &Client{BaseURL: "http://example", APIKey: "2", HTTP: mockHTTP{handler: mux}}
	got, err := c.SearchAlbumByMBID(context.Background(), "d7d2c46c-c9b6-4ff1-bd0a-b3e36f57e0d0")
	if err != nil {
		t.Fatalf("SearchAlbumByMBID error: %v", err)
	}
	if got["name"] != "Test Album" {
		t.Fatalf("unexpected name: %v", got["name"])
	}
	if got["mbid"] != "d7d2c46c-c9b6-4ff1-bd0a-b3e36f57e0d0" {
		t.Fatalf("unexpected mbid: %v", got["mbid"])
	}
	if got["artist"] != "Test Artist" {
		t.Fatalf("unexpected artist: %v", got["artist"])
	}
	if got["year"] != "2005" {
		t.Fatalf("unexpected year: %v", got["year"])
	}
}

func TestSearchArtistByMBID_NotFound(t *testing.T) {
	// Mock empty response
	mux := http.NewServeMux()
	mux.HandleFunc("/2/artist-mb.php", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{"artists": nil}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	c := &Client{BaseURL: "http://example", APIKey: "2", HTTP: mockHTTP{handler: mux}}
	got, err := c.SearchArtistByMBID(context.Background(), "nonexistent-mbid")
	if err != nil {
		t.Fatalf("SearchArtistByMBID error: %v", err)
	}
	if got["error"] != "artist not found" {
		t.Fatalf("expected error message, got: %v", got)
	}
}

func TestSearchTrackByMBID_SimplifiesResult(t *testing.T) {
	// Mock TheAudioDB MBID track response
	mux := http.NewServeMux()
	mux.HandleFunc("/2/track-mb.php", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{
			"track": []map[string]any{{
				"idTrack":          "789",
				"idMBTrack":        "f1b10b1e-c2c6-4ff1-bd0a-b3e36f57e0d1",
				"strTrack":         "Test Track",
				"strArtist":        "Test Artist",
				"strAlbum":         "Test Album",
				"intDuration":      "210000",
				"strGenre":         "Rock",
				"strDescriptionEN": "A great test track",
				"strTrackThumb":    "https://example.com/track_thumb.jpg",
			}},
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	c := &Client{BaseURL: "http://example", APIKey: "2", HTTP: mockHTTP{handler: mux}}
	got, err := c.SearchTrackByMBID(context.Background(), "f1b10b1e-c2c6-4ff1-bd0a-b3e36f57e0d1")
	if err != nil {
		t.Fatalf("SearchTrackByMBID error: %v", err)
	}
	if got["name"] != "Test Track" {
		t.Fatalf("unexpected name: %v", got["name"])
	}
	if got["mbid"] != "f1b10b1e-c2c6-4ff1-bd0a-b3e36f57e0d1" {
		t.Fatalf("unexpected mbid: %v", got["mbid"])
	}
	if got["artist"] != "Test Artist" {
		t.Fatalf("unexpected artist: %v", got["artist"])
	}
	if got["duration"] != "210000" {
		t.Fatalf("unexpected duration: %v", got["duration"])
	}
}
