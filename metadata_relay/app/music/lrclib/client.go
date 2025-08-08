package lrclib

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	BaseURL string
	http    HTTPDoer
}

func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = "https://lrclib.net/api"
	}
	return &Client{
		BaseURL: baseURL,
		http:    &http.Client{Timeout: 8 * time.Second},
	}
}

// FetchLyrics retrieves lyrics for an artist/title pair from LRCLib API.
// If found, saves the lyrics to disk in LRC format and returns the file path.
func (l *Client) FetchLyrics(ctx context.Context, artist, title, mediaDir string) (string, error) {
	u := fmt.Sprintf("%s/get?artist_name=%s&track_name=%s",
		l.BaseURL,
		url.QueryEscape(artist),
		url.QueryEscape(title))

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	resp, err := l.http.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("lrclib http %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Join(mediaDir, "lyrics"), 0o755); err != nil {
		return "", err
	}

	// Save lyrics to disk
	filename := url.PathEscape(artist+" - "+title) + ".lrc"
	file := filepath.Join(mediaDir, "lyrics", filename)
	if err := os.WriteFile(file, b, 0o644); err != nil {
		return "", err
	}

	return file, nil
}
