package music

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestLRCLib_FetchLyrics_WritesFile(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("[00:00.00] Hello"))
	}))
	defer srv.Close()

	c := NewLRCLib(srv.URL)
	dir := t.TempDir()
	path, err := c.FetchLyrics(context.Background(), "Artist", "Title", dir)
	if err != nil {
		t.Fatalf("FetchLyrics error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("missing file: %v", err)
	}
	if filepath.Dir(path) != filepath.Join(dir, "lyrics") {
		t.Fatalf("unexpected dir: %s", filepath.Dir(path))
	}
}
