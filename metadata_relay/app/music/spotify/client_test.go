package spotify

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// mockHTTPClient implements HTTPDoer for testing
type mockHTTPClient struct {
	responses map[string]*http.Response
	requests  []*http.Request
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	m.requests = append(m.requests, req)

	// Create a key based on method and path
	var key string
	if strings.Contains(req.URL.String(), "accounts.spotify.com") {
		key = "POST token"
	} else if strings.Contains(req.URL.String(), "api.spotify.com/v1/search") {
		key = "GET search"
	} else if strings.Contains(req.URL.String(), "example.com") {
		key = "GET image"
	} else {
		key = req.Method + " " + req.URL.String()
	}

	if resp, exists := m.responses[key]; exists {
		return resp, nil
	}

	// Default 404 response
	return &http.Response{
		StatusCode: 404,
		Body:       io.NopCloser(strings.NewReader("Not Found")),
	}, nil
}

func (m *mockHTTPClient) addResponse(key string, statusCode int, body string) {
	if m.responses == nil {
		m.responses = make(map[string]*http.Response)
	}

	m.responses[key] = &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

func TestSpotifyClient(t *testing.T) {
	t.Run("NewClient", func(t *testing.T) {
		client := NewClient("test_id", "test_secret")
		if client.ClientID != "test_id" {
			t.Errorf("Expected ClientID 'test_id', got %s", client.ClientID)
		}
		if client.ClientSecret != "test_secret" {
			t.Errorf("Expected ClientSecret 'test_secret', got %s", client.ClientSecret)
		}
		if client.http == nil {
			t.Error("HTTP client should not be nil")
		}
	})

	t.Run("EnsureToken_Success", func(t *testing.T) {
		mockClient := &mockHTTPClient{}

		// Mock token response
		tokenResponse := map[string]any{
			"access_token": "test_token_123",
			"expires_in":   3600,
		}
		tokenBody, _ := json.Marshal(tokenResponse)
		mockClient.addResponse("POST token", 200, string(tokenBody))

		client := NewClient("test_id", "test_secret")
		client.http = mockClient

		ctx := context.Background()
		err := client.ensureToken(ctx)
		if err != nil {
			t.Fatalf("ensureToken failed: %v", err)
		}

		if client.token != "test_token_123" {
			t.Errorf("Expected token 'test_token_123', got %s", client.token)
		}

		// Token should be valid for approximately 1 hour
		expectedExpiry := time.Now().Add(3600 * time.Second)
		if client.tokenExpiry.Before(expectedExpiry.Add(-time.Minute)) ||
			client.tokenExpiry.After(expectedExpiry.Add(time.Minute)) {
			t.Errorf("Token expiry time is not as expected: %v", client.tokenExpiry)
		}

		// Verify request was made with correct auth
		if len(mockClient.requests) != 1 {
			t.Fatalf("Expected 1 request, got %d", len(mockClient.requests))
		}

		req := mockClient.requests[0]
		username, password, ok := req.BasicAuth()
		if !ok {
			t.Error("Expected basic auth to be set")
		}
		if username != "test_id" || password != "test_secret" {
			t.Errorf("Expected basic auth test_id:test_secret, got %s:%s", username, password)
		}
	})

	t.Run("EnsureToken_NoCredentials", func(t *testing.T) {
		client := NewClient("", "")
		ctx := context.Background()

		err := client.ensureToken(ctx)
		if err == nil {
			t.Error("Expected error for missing credentials")
		}
		if !strings.Contains(err.Error(), "spotify not configured") {
			t.Errorf("Expected error about configuration, got: %v", err)
		}
	})

	t.Run("EnsureToken_HTTPError", func(t *testing.T) {
		mockClient := &mockHTTPClient{}
		mockClient.addResponse("POST token", 401, "Unauthorized")

		client := NewClient("bad_id", "bad_secret")
		client.http = mockClient

		ctx := context.Background()
		err := client.ensureToken(ctx)
		if err == nil {
			t.Error("Expected error for HTTP 401")
		}
		if !strings.Contains(err.Error(), "spotify token http 401") {
			t.Errorf("Expected HTTP error, got: %v", err)
		}
	})

	t.Run("EnsureToken_ReuseValidToken", func(t *testing.T) {
		mockClient := &mockHTTPClient{}

		client := NewClient("test_id", "test_secret")
		client.http = mockClient
		client.token = "existing_token"
		client.tokenExpiry = time.Now().Add(10 * time.Minute) // Valid for 10 more minutes

		ctx := context.Background()
		err := client.ensureToken(ctx)
		if err != nil {
			t.Fatalf("ensureToken failed: %v", err)
		}

		// Should not make any HTTP requests since token is still valid
		if len(mockClient.requests) != 0 {
			t.Errorf("Expected 0 requests, got %d", len(mockClient.requests))
		}

		if client.token != "existing_token" {
			t.Errorf("Expected token 'existing_token', got %s", client.token)
		}
	})
}

func TestSpotifyDownloadArtistImage(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "spotify_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	t.Run("DownloadArtistImage_Success", func(t *testing.T) {
		mockClient := &mockHTTPClient{}

		// Mock token response
		tokenResponse := map[string]any{
			"access_token": "test_token_123",
			"expires_in":   3600,
		}
		tokenBody, _ := json.Marshal(tokenResponse)
		mockClient.addResponse("POST token", 200, string(tokenBody))

		// Mock search response
		searchResponse := map[string]any{
			"artists": map[string]any{
				"items": []map[string]any{
					{
						"images": []map[string]any{
							{
								"url": "https://example.com/artist_image.jpg",
							},
						},
					},
				},
			},
		}
		searchBody, _ := json.Marshal(searchResponse)
		mockClient.addResponse("GET search", 200, string(searchBody))

		// Mock image download response
		imageData := "fake_jpeg_data"
		mockClient.addResponse("GET image", 200, imageData)

		client := NewClient("test_id", "test_secret")
		client.http = mockClient

		ctx := context.Background()
		filePath, err := client.DownloadArtistImage(ctx, "Test Artist", tempDir)
		if err != nil {
			t.Fatalf("DownloadArtistImage failed: %v", err)
		}

		expectedPath := filepath.Join(tempDir, "spotify", "artists", "Test%20Artist.jpg")
		if filePath != expectedPath {
			t.Errorf("Expected file path %s, got %s", expectedPath, filePath)
		}

		// Check that file was created
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Expected file to be created at %s", filePath)
		}

		// Check file contents
		data, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read downloaded file: %v", err)
		}
		if string(data) != imageData {
			t.Errorf("Expected file content '%s', got '%s'", imageData, string(data))
		}

		// Verify requests
		if len(mockClient.requests) != 3 {
			t.Fatalf("Expected 3 requests (token, search, image), got %d", len(mockClient.requests))
		}

		// Check search request has authorization header
		searchReq := mockClient.requests[1]
		authHeader := searchReq.Header.Get("Authorization")
		if authHeader != "Bearer test_token_123" {
			t.Errorf("Expected Authorization header 'Bearer test_token_123', got %s", authHeader)
		}
	})

	t.Run("DownloadArtistImage_NoResults", func(t *testing.T) {
		mockClient := &mockHTTPClient{}

		// Mock token response
		tokenResponse := map[string]any{
			"access_token": "test_token_123",
			"expires_in":   3600,
		}
		tokenBody, _ := json.Marshal(tokenResponse)
		mockClient.addResponse("POST token", 200, string(tokenBody))

		// Mock empty search response
		searchResponse := map[string]any{
			"artists": map[string]any{
				"items": []map[string]any{},
			},
		}
		searchBody, _ := json.Marshal(searchResponse)
		mockClient.addResponse("GET search", 200, string(searchBody))

		client := NewClient("test_id", "test_secret")
		client.http = mockClient

		ctx := context.Background()
		filePath, err := client.DownloadArtistImage(ctx, "Unknown Artist", tempDir)
		if err == nil {
			t.Error("Expected error for no search results")
		}
		if filePath != "" {
			t.Errorf("Expected empty file path, got %s", filePath)
		}
		if !strings.Contains(err.Error(), "no image found for artist") {
			t.Errorf("Expected error about no image found, got: %v", err)
		}
	})

	t.Run("DownloadArtistImage_NoImages", func(t *testing.T) {
		mockClient := &mockHTTPClient{}

		// Mock token response
		tokenResponse := map[string]any{
			"access_token": "test_token_123",
			"expires_in":   3600,
		}
		tokenBody, _ := json.Marshal(tokenResponse)
		mockClient.addResponse("POST token", 200, string(tokenBody))

		// Mock search response with artist but no images
		searchResponse := map[string]any{
			"artists": map[string]any{
				"items": []map[string]any{
					{
						"images": []map[string]any{},
					},
				},
			},
		}
		searchBody, _ := json.Marshal(searchResponse)
		mockClient.addResponse("GET search", 200, string(searchBody))

		client := NewClient("test_id", "test_secret")
		client.http = mockClient

		ctx := context.Background()
		filePath, err := client.DownloadArtistImage(ctx, "No Image Artist", tempDir)
		if err == nil {
			t.Error("Expected error for no images")
		}
		if filePath != "" {
			t.Errorf("Expected empty file path, got %s", filePath)
		}
		if !strings.Contains(err.Error(), "no image found for artist") {
			t.Errorf("Expected error about no image found, got: %v", err)
		}
	})

	t.Run("DownloadArtistImage_SearchError", func(t *testing.T) {
		mockClient := &mockHTTPClient{}

		// Mock token response
		tokenResponse := map[string]any{
			"access_token": "test_token_123",
			"expires_in":   3600,
		}
		tokenBody, _ := json.Marshal(tokenResponse)
		mockClient.addResponse("POST token", 200, string(tokenBody))

		// Mock search error response
		mockClient.addResponse("GET search", 500, "Internal Server Error")

		client := NewClient("test_id", "test_secret")
		client.http = mockClient

		ctx := context.Background()
		filePath, err := client.DownloadArtistImage(ctx, "Error Artist", tempDir)
		if err == nil {
			t.Error("Expected error for search HTTP error")
		}
		if filePath != "" {
			t.Errorf("Expected empty file path, got %s", filePath)
		}
		if !strings.Contains(err.Error(), "spotify search http 500") {
			t.Errorf("Expected error about search HTTP error, got: %v", err)
		}
	})

	t.Run("DownloadArtistImage_ImageDownloadError", func(t *testing.T) {
		mockClient := &mockHTTPClient{}

		// Mock token response
		tokenResponse := map[string]any{
			"access_token": "test_token_123",
			"expires_in":   3600,
		}
		tokenBody, _ := json.Marshal(tokenResponse)
		mockClient.addResponse("POST token", 200, string(tokenBody))

		// Mock search response
		searchResponse := map[string]any{
			"artists": map[string]any{
				"items": []map[string]any{
					{
						"images": []map[string]any{
							{
								"url": "https://example.com/bad_image.jpg",
							},
						},
					},
				},
			},
		}
		searchBody, _ := json.Marshal(searchResponse)
		mockClient.addResponse("GET search", 200, string(searchBody))

		// Mock image download error
		mockClient.addResponse("GET image", 404, "Not Found")

		client := NewClient("test_id", "test_secret")
		client.http = mockClient

		ctx := context.Background()
		filePath, err := client.DownloadArtistImage(ctx, "Bad Image Artist", tempDir)
		if err == nil {
			t.Error("Expected error for image download error")
		}
		if filePath != "" {
			t.Errorf("Expected empty file path, got %s", filePath)
		}
		if !strings.Contains(err.Error(), "image download http 404") {
			t.Errorf("Expected error about image download error, got: %v", err)
		}
	})

	t.Run("DownloadArtistImage_TokenError", func(t *testing.T) {
		client := NewClient("", "") // No credentials

		ctx := context.Background()
		filePath, err := client.DownloadArtistImage(ctx, "Test Artist", tempDir)
		if err == nil {
			t.Error("Expected error for token failure")
		}
		if filePath != "" {
			t.Errorf("Expected empty file path, got %s", filePath)
		}
		if !strings.Contains(err.Error(), "spotify not configured") {
			t.Errorf("Expected error about configuration, got: %v", err)
		}
	})
}

func TestSpotifyWithRealHTTPServer(t *testing.T) {
	// Test with actual HTTP server to verify request formatting
	var serverURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/api/token":
			// Verify content type and body for token request
			if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
				t.Errorf("Expected Content-Type application/x-www-form-urlencoded, got %s", r.Header.Get("Content-Type"))
			}

			// Verify basic auth
			username, password, ok := r.BasicAuth()
			if !ok {
				t.Error("Expected basic auth to be set")
			}
			if username != "test_id" || password != "test_secret" {
				t.Errorf("Expected basic auth test_id:test_secret, got %s:%s", username, password)
			}

			response := map[string]any{
				"access_token": "server_token",
				"expires_in":   3600,
			}
			_ = json.NewEncoder(w).Encode(response)

		case "/v1/search":
			// Verify authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader != "Bearer server_token" {
				t.Errorf("Expected Authorization Bearer server_token, got %s", authHeader)
			}

			query := r.URL.Query().Get("q")
			if query != "Test Artist" {
				t.Errorf("Expected query 'Test Artist', got %s", query)
			}

			response := map[string]any{
				"artists": map[string]any{
					"items": []map[string]any{
						{
							"images": []map[string]any{
								{
									"url": serverURL + "/image.jpg",
								},
							},
						},
					},
				},
			}
			_ = json.NewEncoder(w).Encode(response)

		case "/image.jpg":
			w.Header().Set("Content-Type", "image/jpeg")
			_, _ = w.Write([]byte("fake_image_data"))

		default:
			w.WriteHeader(404)
		}
	}))
	defer server.Close()
	serverURL = server.URL

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "spotify_server_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	// Create client with custom HTTP client that points to test server
	client := NewClient("test_id", "test_secret")

	// Create custom transport to redirect requests to test server
	transport := &testTransport{serverURL: server.URL}
	client.http = &http.Client{
		Timeout:   8 * time.Second,
		Transport: transport,
	}

	ctx := context.Background()
	filePath, err := client.DownloadArtistImage(ctx, "Test Artist", tempDir)
	if err != nil {
		t.Fatalf("DownloadArtistImage failed: %v", err)
	}

	expectedPath := filepath.Join(tempDir, "spotify", "artists", "Test%20Artist.jpg")
	if filePath != expectedPath {
		t.Errorf("Expected file path %s, got %s", expectedPath, filePath)
	}

	// Verify file was created with correct content
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %v", err)
	}
	if string(data) != "fake_image_data" {
		t.Errorf("Expected file content 'fake_image_data', got '%s'", string(data))
	}
}

// Custom transport for rewriting URLs in tests
type testTransport struct {
	serverURL string
}

func (t *testTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Rewrite URLs to point to test server
	if strings.Contains(req.URL.String(), "accounts.spotify.com") {
		newURL, _ := url.Parse(t.serverURL + "/api/token")
		req.URL = newURL
	} else if strings.Contains(req.URL.String(), "api.spotify.com") {
		req.URL.Host = strings.TrimPrefix(t.serverURL, "http://")
		req.URL.Scheme = "http"
	}

	return http.DefaultTransport.RoundTrip(req)
}
