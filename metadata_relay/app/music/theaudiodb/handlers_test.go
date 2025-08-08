package theaudiodb

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSearchArtistHandler_NoClient(t *testing.T) {
	// reset client
	client = nil
	r := gin.New()
	r.GET("/theaudiodb/artist", SearchArtistHandler)
	req := httptest.NewRequest(http.MethodGet, "/theaudiodb/artist?name=abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", w.Code)
	}
}

func TestSearchArtistHandler_BadRequest(t *testing.T) {
	// Set dummy client so we can test 400 path
	client = &Client{BaseURL: "x", APIKey: "y"}
	r := gin.New()
	r.GET("/theaudiodb/artist", SearchArtistHandler)
	req := httptest.NewRequest(http.MethodGet, "/theaudiodb/artist", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestSearchArtistHandler_OK(t *testing.T) {
	// mock client with custom HTTP
	mux := http.NewServeMux()
	mux.HandleFunc("/2/search.php", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"artists": []map[string]any{{"strArtist": "A"}}})
	})
	c := &Client{BaseURL: "http://example", APIKey: "2", HTTP: mockHTTP{handler: mux}}
	SetClient(c)

	r := gin.New()
	r.GET("/theaudiodb/artist", SearchArtistHandler)
	req := httptest.NewRequest(http.MethodGet, "/theaudiodb/artist?name=A", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetArtistByMBIDHandler_NoClient(t *testing.T) {
	// reset client
	client = nil
	r := gin.New()
	r.GET("/theaudiodb/artist/:mbid", GetArtistByMBIDHandler)
	req := httptest.NewRequest(http.MethodGet, "/theaudiodb/artist/12345", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", w.Code)
	}
}

func TestGetArtistByMBIDHandler_OK(t *testing.T) {
	// mock client with custom HTTP
	mux := http.NewServeMux()
	mux.HandleFunc("/2/artist-mb.php", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"artists": []map[string]any{{
				"strArtist":        "Test Artist",
				"strMusicBrainzID": "12345",
				"strBiographyEN":   "Test biography",
				"strGenre":         "Rock",
			}},
		})
	})
	c := &Client{BaseURL: "http://example", APIKey: "2", HTTP: mockHTTP{handler: mux}}
	SetClient(c)

	r := gin.New()
	r.GET("/theaudiodb/artist/:mbid", GetArtistByMBIDHandler)
	req := httptest.NewRequest(http.MethodGet, "/theaudiodb/artist/12345", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetAlbumByMBIDHandler_NoClient(t *testing.T) {
	// reset client
	client = nil
	r := gin.New()
	r.GET("/theaudiodb/album/:mbid", GetAlbumByMBIDHandler)
	req := httptest.NewRequest(http.MethodGet, "/theaudiodb/album/67890", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", w.Code)
	}
}

func TestGetAlbumByMBIDHandler_OK(t *testing.T) {
	// mock client with custom HTTP
	mux := http.NewServeMux()
	mux.HandleFunc("/2/album-mb.php", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"album": []map[string]any{{
				"strAlbum":         "Test Album",
				"strArtist":        "Test Artist",
				"idMBAlbum":        "67890",
				"strDescriptionEN": "Test description",
				"strGenre":         "Rock",
			}},
		})
	})
	c := &Client{BaseURL: "http://example", APIKey: "2", HTTP: mockHTTP{handler: mux}}
	SetClient(c)

	r := gin.New()
	r.GET("/theaudiodb/album/:mbid", GetAlbumByMBIDHandler)
	req := httptest.NewRequest(http.MethodGet, "/theaudiodb/album/67890", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetTrackByMBIDHandler_NoClient(t *testing.T) {
	// reset client
	client = nil
	r := gin.New()
	r.GET("/theaudiodb/track/:mbid", GetTrackByMBIDHandler)
	req := httptest.NewRequest(http.MethodGet, "/theaudiodb/track/99999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", w.Code)
	}
}

func TestGetTrackByMBIDHandler_OK(t *testing.T) {
	// mock client with custom HTTP
	mux := http.NewServeMux()
	mux.HandleFunc("/2/track-mb.php", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{
			"track": []map[string]any{{
				"strTrack":         "Test Track",
				"strArtist":        "Test Artist",
				"strAlbum":         "Test Album",
				"idMBTrack":        "99999",
				"strDescriptionEN": "Test track description",
				"intDuration":      "210000",
			}},
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	c := &Client{BaseURL: "http://example", APIKey: "2", HTTP: mockHTTP{handler: mux}}
	SetClient(c)

	r := gin.New()
	r.GET("/theaudiodb/track/:mbid", GetTrackByMBIDHandler)
	req := httptest.NewRequest(http.MethodGet, "/theaudiodb/track/99999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
