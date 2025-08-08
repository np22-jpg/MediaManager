package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	ClientID     string
	ClientSecret string
	http         HTTPDoer
	token        string
	tokenExpiry  time.Time
}

func NewClient(clientID, clientSecret string) *Client {
	return &Client{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		http:         &http.Client{Timeout: 8 * time.Second},
	}
}

func (s *Client) ensureToken(ctx context.Context) error {
	if s == nil || s.ClientID == "" || s.ClientSecret == "" {
		return fmt.Errorf("spotify not configured")
	}
	if time.Now().Before(s.tokenExpiry.Add(-1*time.Minute)) && s.token != "" {
		return nil
	}

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "https://accounts.spotify.com/api/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(s.ClientID, s.ClientSecret)

	resp, err := s.http.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("spotify token http %d", resp.StatusCode)
	}

	var out struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return err
	}

	s.token = out.AccessToken
	s.tokenExpiry = time.Now().Add(time.Duration(out.ExpiresIn) * time.Second)
	return nil
}

// DownloadArtistImage searches for an artist on Spotify and downloads their image to disk.
// Returns the local file path if successful.
func (s *Client) DownloadArtistImage(ctx context.Context, artistName, mediaDir string) (string, error) {
	if err := s.ensureToken(ctx); err != nil {
		return "", err
	}

	q := url.QueryEscape(artistName)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.spotify.com/v1/search?type=artist&limit=1&q="+q, nil)
	req.Header.Set("Authorization", "Bearer "+s.token)

	resp, err := s.http.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("spotify search http %d", resp.StatusCode)
	}

	var out struct {
		Artists struct {
			Items []struct {
				Images []struct {
					URL string `json:"url"`
				} `json:"images"`
			} `json:"items"`
		} `json:"artists"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}

	if len(out.Artists.Items) == 0 || len(out.Artists.Items[0].Images) == 0 {
		return "", fmt.Errorf("no image found for artist: %s", artistName)
	}

	imgURL := out.Artists.Items[0].Images[0].URL

	// Download image
	ireq, _ := http.NewRequestWithContext(ctx, http.MethodGet, imgURL, nil)
	iresp, err := s.http.Do(ireq)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = iresp.Body.Close()
	}()

	if iresp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("image download http %d", iresp.StatusCode)
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Join(mediaDir, "spotify", "artists"), 0o755); err != nil {
		return "", err
	}

	// Save image to disk
	file := filepath.Join(mediaDir, "spotify", "artists", url.PathEscape(artistName)+".jpg")
	f, err := os.Create(file)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = f.Close()
	}()

	if _, err := io.Copy(f, iresp.Body); err != nil {
		return "", err
	}

	return file, nil
}
