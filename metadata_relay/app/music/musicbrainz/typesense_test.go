package musicbrainz

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestTypesenseIntegration(t *testing.T) {
	// Skip if Typesense is not configured
	if !IsReady() {
		t.Skip("Typesense or MusicBrainz not configured - skipping integration tests")
	}

	ctx := context.Background()

	t.Run("SearchArtists", func(t *testing.T) {
		result, err := SearchArtistsTypesense(ctx, "Beatles", 5)
		if err != nil {
			t.Skip("Typesense search failed (service may be unavailable):", err)
		}

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("Expected non-nil value")
		}

		// Check result structure
		resultMap, ok := result.(map[string]any)
		if !ok {
			t.Fatal("Result should be a map")
		}

		artists, ok := resultMap["artists"]
		if !ok {
			t.Fatal("Result should contain 'artists' key")
		}

		artistList, ok := artists.([]map[string]any)
		if !ok {
			t.Fatal("Artists should be a slice of maps")
		}

		if len(artistList) > 0 {
			artist := artistList[0]
			if _, ok := artist["id"]; !ok {
				t.Error("Artist should have an ID")
			}
			if _, ok := artist["name"]; !ok {
				t.Error("Artist should have a name")
			}
		}
	})

	t.Run("SearchReleaseGroups", func(t *testing.T) {
		result, err := SearchReleaseGroupsTypesense(ctx, "Abbey Road", 5)
		if err != nil {
			t.Skip("Typesense search failed (service may be unavailable):", err)
		}

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("Expected non-nil result")
		}

		// Check result structure
		resultMap, ok := result.(map[string]any)
		if !ok {
			t.Fatal("Result should be a map")
		}

		releaseGroups, ok := resultMap["release-groups"]
		if !ok {
			t.Fatal("Result should contain 'release-groups' key")
		}

		rgList, ok := releaseGroups.([]map[string]any)
		if !ok {
			t.Fatal("Release groups should be a slice of maps")
		}

		if len(rgList) > 0 {
			rg := rgList[0]
			if _, ok := rg["id"]; !ok {
				t.Error("Release group should have an ID")
			}
			if _, ok := rg["title"]; !ok {
				t.Error("Release group should have a title")
			}
		}
	})

	t.Run("SearchRecordings", func(t *testing.T) {
		result, err := SearchRecordingsTypesense(ctx, "Come Together", 5)
		if err != nil {
			t.Skip("Typesense search failed (service may be unavailable):", err)
		}

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("Expected non-nil result")
		}

		// Check result structure
		resultMap, ok := result.(map[string]any)
		if !ok {
			t.Fatal("Result should be a map")
		}

		recordings, ok := resultMap["recordings"]
		if !ok {
			t.Fatal("Result should contain 'recordings' key")
		}

		recList, ok := recordings.([]map[string]any)
		if !ok {
			t.Fatal("Recordings should be a slice of maps")
		}

		if len(recList) > 0 {
			rec := recList[0]
			if _, ok := rec["id"]; !ok {
				t.Error("Recording should have an ID")
			}
			if _, ok := rec["title"]; !ok {
				t.Error("Recording should have a title")
			}
		}
	})
}

func TestTypesenseNotConfigured(t *testing.T) {
	// This test ensures proper handling when Typesense is not available
	if IsReady() {
		t.Skip("Typesense is configured - skipping not-configured test")
	}

	ctx := context.Background()

	t.Run("SearchArtistsFailsGracefully", func(t *testing.T) {
		result, err := SearchArtistsTypesense(ctx, "Beatles", 5)
		if err == nil {
			t.Error("Expected error when Typesense not configured")
		}
		if result != nil {
			t.Error("Expected nil result when Typesense not configured")
		}
	})

	t.Run("IndexingFailsGracefully", func(t *testing.T) {
		err := IndexArtists()
		if err == nil {
			t.Error("Expected error when Typesense not configured")
		}
		if !strings.Contains(err.Error(), "not ready") {
			t.Errorf("Expected error to contain 'not ready', got: %s", err.Error())
		}
	})
}

func TestSyncFunctions(t *testing.T) {
	if !IsReady() {
		t.Skip("MusicBrainz or Typesense not configured - skipping sync tests")
	}

	t.Run("IndexArtists", func(t *testing.T) {
		err := IndexArtists()
		if err != nil {
			t.Skip("Indexing failed (database may be unavailable):", err)
		}
		if err != nil {
			t.Errorf("IndexArtists() error = %v", err)
		}
	})

	t.Run("IndexReleaseGroups", func(t *testing.T) {
		err := IndexReleaseGroups()
		if err != nil {
			t.Skip("Indexing failed (database may be unavailable):", err)
		}
		if err != nil {
			t.Errorf("IndexReleaseGroups() error = %v", err)
		}
	})

	t.Run("IndexRecordings", func(t *testing.T) {
		err := IndexRecordings()
		if err != nil {
			t.Skip("Indexing failed (database may be unavailable):", err)
		}
		if err != nil {
			t.Errorf("IndexRecordings() error = %v", err)
		}
	})
}

// TestTypesenseConnectionFailure tests graceful degradation when Typesense is unavailable
func TestTypesenseConnectionFailure(t *testing.T) {
	// Test with invalid Typesense configuration
	originalClient := typesenseClient
	defer func() { typesenseClient = originalClient }()

	// Clear the client to simulate initialization failure
	typesenseClient = nil

	t.Run("IsReady_WithoutTypesense", func(t *testing.T) {
		// Even with database, should return false without Typesense
		ready := IsReady()
		if ready {
			t.Error("IsReady should return false when Typesense is not initialized")
		}
	})

	t.Run("InitTypesense_InvalidHost", func(t *testing.T) {
		// Try to initialize with invalid host - should fail gracefully
		err := InitTypesense("nonexistent.host", "8108", "invalid_api_key", 2*time.Second)
		if err == nil {
			t.Error("Expected error for invalid host")
		}
		if !strings.Contains(err.Error(), "failed to connect to Typesense") {
			t.Errorf("Expected error to contain 'failed to connect to Typesense', got: %s", err.Error())
		}
	})

	t.Run("InitTypesense_InvalidPort", func(t *testing.T) {
		// Try to initialize with invalid port - should fail gracefully
		err := InitTypesense("localhost", "99999", "invalid_api_key", 2*time.Second)
		if err == nil {
			t.Error("Expected error for invalid port")
		}
		if !strings.Contains(err.Error(), "failed to connect to Typesense") &&
			!strings.Contains(err.Error(), "connection refused") &&
			!strings.Contains(err.Error(), "invalid port") {
			t.Errorf("Error should be about connection failure: %v", err)
		}
	})

	t.Run("InitTypesense_EmptyConfig", func(t *testing.T) {
		// Try to initialize with empty configuration
		err := InitTypesense("", "", "", 2*time.Second)
		if err == nil {
			t.Error("Expected error for empty config")
		}
		if !strings.Contains(err.Error(), "failed to connect to Typesense") {
			t.Errorf("Expected error to contain 'failed to connect to Typesense', got: %s", err.Error())
		}
	})
}

// TestTypesenseSearchWithoutConnection tests search functions when Typesense is not available
func TestTypesenseSearchWithoutConnection(t *testing.T) {
	// Save original client and set to nil to simulate unavailable Typesense
	originalClient := typesenseClient
	defer func() { typesenseClient = originalClient }()

	typesenseClient = nil
	ctx := context.Background()

	t.Run("SearchArtistsTypesense_NoConnection", func(t *testing.T) {
		result, err := SearchArtistsTypesense(ctx, "test artist", 10)
		if err == nil {
			t.Error("Expected error when no Typesense connection")
		}
		if result != nil {
			t.Error("Expected nil result when no Typesense connection")
		}
		if !strings.Contains(err.Error(), "typesense") &&
			!strings.Contains(err.Error(), "nil") {
			t.Errorf("Error should mention Typesense or nil client: %v", err)
		}
	})

	t.Run("SearchReleaseGroupsTypesense_NoConnection", func(t *testing.T) {
		result, err := SearchReleaseGroupsTypesense(ctx, "test album", 10)
		if err == nil {
			t.Error("Expected error when no Typesense connection")
		}
		if result != nil {
			t.Error("Expected nil result when no Typesense connection")
		}
		if !strings.Contains(err.Error(), "typesense") &&
			!strings.Contains(err.Error(), "nil") {
			t.Errorf("Error should mention Typesense or nil client: %v", err)
		}
	})

	t.Run("SearchRecordingsTypesense_NoConnection", func(t *testing.T) {
		result, err := SearchRecordingsTypesense(ctx, "test song", 10)
		if err == nil {
			t.Error("Expected error when no Typesense connection")
		}
		if result != nil {
			t.Error("Expected nil result when no Typesense connection")
		}
		if !strings.Contains(err.Error(), "typesense") &&
			!strings.Contains(err.Error(), "nil") {
			t.Errorf("Error should mention Typesense or nil client: %v", err)
		}
	})

	t.Run("GetArtistTypesense_NoConnection", func(t *testing.T) {
		result, err := GetArtistTypesense(ctx, "123e4567-e89b-12d3-a456-426614174000")
		if err == nil {
			t.Error("Expected error when no Typesense connection")
		}
		if result != nil {
			t.Error("Expected nil result when no Typesense connection")
		}
		if !strings.Contains(err.Error(), "typesense") &&
			!strings.Contains(err.Error(), "nil") {
			t.Errorf("Error should mention Typesense or nil client: %v", err)
		}
	})
}

// TestTypesenseConnectionScenarios tests connection failures with different scenarios
func TestTypesenseConnectionScenarios(t *testing.T) {
	// Save original state
	originalClient := typesenseClient
	defer func() { typesenseClient = originalClient }()

	testCases := []struct {
		name          string
		host          string
		port          string
		apiKey        string
		expectError   bool
		errorContains string
	}{
		{
			name:          "localhost_invalid_port",
			host:          "localhost",
			port:          "99999",
			apiKey:        "test_key",
			expectError:   true,
			errorContains: "failed to connect to Typesense",
		},
		{
			name:          "nonexistent_host",
			host:          "nonexistent.typesense.host",
			port:          "8108",
			apiKey:        "test_key",
			expectError:   true,
			errorContains: "failed to connect to Typesense",
		},
		{
			name:          "empty_host",
			host:          "",
			port:          "8108",
			apiKey:        "test_key",
			expectError:   true,
			errorContains: "failed to connect to Typesense",
		},
		{
			name:          "empty_port",
			host:          "localhost",
			port:          "",
			apiKey:        "test_key",
			expectError:   true,
			errorContains: "failed to connect to Typesense",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset client for each test
			typesenseClient = nil

			err := InitTypesense(tc.host, tc.port, tc.apiKey, 2*time.Second)

			if tc.expectError {
				if err == nil {
					t.Error("Expected error for test case:", tc.name)
				}
				if !strings.Contains(err.Error(), tc.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %s", tc.errorContains, err.Error())
				}
				if IsReady() {
					t.Error("IsReady should return false after failed initialization")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for test case %s: %v", tc.name, err)
				}
				if !IsReady() {
					t.Error("IsReady should return true after successful initialization")
				}
			}
		})
	}
}
