package app

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all application configuration settings loaded from environment variables.
type Config struct {
	LogLevel string
	Port     string

	// Metrics server configuration
	MetricsPort string

	CacheHost string
	CachePort int
	CacheDB   int

	TMDBAPIKey  string
	TMDBBaseURL string
	TVDBAPIKey  string
	TVDBBaseURL string

	// TheAudioDB Configuration
	TheAudioDBAPIKey  string
	TheAudioDBBaseURL string

	// Spotify (for images)
	SpotifyClientID     string
	SpotifyClientSecret string

	// LRCLib (lyrics)
	LRCLibBaseURL string

	// SeaDx (anime)
	SeaDxBaseURL string

	// AniList Configuration
	AniListGraphQLURL string
	AniListUserAgent  string

	// Media storage directory (on disk, do not cache images in Redis)
	MediaDir string

	// MusicBrainz PostgreSQL Configuration
	MusicBrainzDBHost     string
	MusicBrainzDBPort     string
	MusicBrainzDBUser     string
	MusicBrainzDBPassword string
	MusicBrainzDBName     string

	// Typesense Configuration
	TypesenseHost   string
	TypesensePort   string
	TypesenseAPIKey string
	// HTTP timeout for Typesense client operations (e.g., bulk imports)
	TypesenseTimeout string

	// Sync Configuration
	SyncInterval string // How often to sync data to Typesense
	SyncEnabled  bool   // Toggle background sync scheduler
	// Comma-separated list of entities to sync when target is "all" or scheduler runs
	// Allowed values: artists, release-groups, releases, recordings
	SyncEntities string
	// If true, attempts to skip unchanged documents during sync using a content fingerprint cache
	SyncSkipUnchanged bool

	// Sync Performance Tunables
	SyncDBPageSize       int
	SyncShardParallelism int
	SyncImportBatchSize  int
	SyncImportWorkers    int
	SyncImportMaxRetries int
	SyncImportBackoff    string
	// Global cap across all entities for concurrent Typesense import requests
	SyncImportGlobalLimit int
}

var AppConfig Config

// getEnv returns environment variable value or default if not set
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt returns environment variable as int or default if not set or invalid
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

// getEnvBool returns environment variable as bool or default if not set or invalid
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

// LoadConfig parses environment variables and loads the application configuration.
func LoadConfig() error {
	AppConfig = Config{
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		Port:        getEnv("PORT", "8000"),
		MetricsPort: getEnv("METRICS_PORT", "9090"),

		CacheHost: getEnv("VALKEY_HOST", "localhost"),
		CachePort: getEnvInt("VALKEY_PORT", 6379),
		CacheDB:   getEnvInt("VALKEY_DB", 0),

		TMDBAPIKey:  getEnv("TMDB_API_KEY", ""),
		TMDBBaseURL: getEnv("TMDB_BASE_URL", "https://api.themoviedb.org/3"),
		TVDBAPIKey:  getEnv("TVDB_API_KEY", ""),
		TVDBBaseURL: getEnv("TVDB_BASE_URL", "https://api4.thetvdb.com/v4"),

		TheAudioDBAPIKey:  getEnv("THEAUDIODB_API_KEY", ""),
		TheAudioDBBaseURL: getEnv("THEAUDIODB_BASE_URL", "https://www.theaudiodb.com/api/v1/json"),

		SpotifyClientID:     getEnv("SPOTIFY_CLIENT_ID", ""),
		SpotifyClientSecret: getEnv("SPOTIFY_CLIENT_SECRET", ""),

		LRCLibBaseURL: getEnv("LRCLIB_BASE_URL", "https://lrclib.net/api"),

		SeaDxBaseURL: getEnv("SEADX_BASE_URL", "https://releases.moe/api"),

		AniListGraphQLURL: getEnv("ANILIST_GRAPHQL_URL", "https://graphql.anilist.co"),
		AniListUserAgent:  getEnv("ANILIST_USER_AGENT", "MediaManager-Relay/1.0"),

		MediaDir: getEnv("MEDIA_DIR", "./media"),

		MusicBrainzDBHost:     getEnv("MUSICBRAINZ_DB_HOST", ""),
		MusicBrainzDBPort:     getEnv("MUSICBRAINZ_DB_PORT", "5432"),
		MusicBrainzDBUser:     getEnv("MUSICBRAINZ_DB_USER", "musicbrainz"),
		MusicBrainzDBPassword: getEnv("MUSICBRAINZ_DB_PASSWORD", "musicbrainz"),
		MusicBrainzDBName:     getEnv("MUSICBRAINZ_DB_NAME", ""),

		TypesenseHost:    getEnv("TYPESENSE_HOST", "localhost"),
		TypesensePort:    getEnv("TYPESENSE_PORT", "8108"),
		TypesenseAPIKey:  getEnv("TYPESENSE_API_KEY", ""),
		TypesenseTimeout: getEnv("TYPESENSE_TIMEOUT", "60s"),

		SyncInterval:          getEnv("SYNC_INTERVAL", "24h"),
		SyncEnabled:           getEnvBool("SYNC_ENABLED", true),
		SyncEntities:          getEnv("SYNC_ENTITIES", "artists,release-groups,releases,recordings"),
		SyncSkipUnchanged:     getEnvBool("SYNC_SKIP_UNCHANGED", true),
		SyncDBPageSize:        getEnvInt("SYNC_DB_PAGE_SIZE", 8000),
		SyncShardParallelism:  getEnvInt("SYNC_SHARD_PARALLELISM", 0),
		SyncImportBatchSize:   getEnvInt("SYNC_IMPORT_BATCH_SIZE", 2000),
		SyncImportWorkers:     getEnvInt("SYNC_IMPORT_WORKERS", 0),
		SyncImportMaxRetries:  getEnvInt("SYNC_IMPORT_MAX_RETRIES", 3),
		SyncImportBackoff:     getEnv("SYNC_IMPORT_BACKOFF", "400ms"),
		SyncImportGlobalLimit: getEnvInt("SYNC_IMPORT_GLOBAL_LIMIT", 0),
	}
	return nil
}

// GetLogLevel converts string log level to slog.Level with case-insensitive matching.
func (c *Config) GetLogLevel() slog.Level {
	verbosity := strings.ToLower(c.LogLevel)
	switch verbosity {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// GetServerAddr returns the server address in host:port format.
func (c *Config) GetServerAddr() string {
	return ":" + c.Port
}

// GetMetricsAddr returns the metrics server address in host:port format.
func (c *Config) GetMetricsAddr() string {
	return ":" + c.MetricsPort
}

// GetCacheAddr returns the cache address in host:port format.
func (c *Config) GetCacheAddr() string {
	return fmt.Sprintf("%s:%d", c.CacheHost, c.CachePort)
}

// GetMusicBrainzConnStr returns the MusicBrainz database connection string.
func (c *Config) GetMusicBrainzConnStr() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.MusicBrainzDBHost, c.MusicBrainzDBPort, c.MusicBrainzDBUser, c.MusicBrainzDBPassword, c.MusicBrainzDBName)
}

// GetTypesenseURL returns the Typesense server URL.
func (c *Config) GetTypesenseURL() string {
	return fmt.Sprintf("http://%s:%s", c.TypesenseHost, c.TypesensePort)
}

// GetTypesenseTimeout parses the Typesense HTTP client timeout.
func (c *Config) GetTypesenseTimeout() time.Duration {
	d, err := time.ParseDuration(c.TypesenseTimeout)
	if err != nil || d <= 0 {
		return 60 * time.Second
	}
	return d
}

// IsMusicBrainzConfigured checks if MusicBrainz is configured (requires both host and database name).
func (c *Config) IsMusicBrainzConfigured() bool {
	return c.MusicBrainzDBHost != "" && c.MusicBrainzDBName != ""
}

// IsTypesenseConfigured checks if Typesense is configured (requires API key).
func (c *Config) IsTypesenseConfigured() bool {
	return c.TypesenseAPIKey != ""
}

// IsSeaDxConfigured checks if SeaDx is configured (has base URL).
func (c *Config) IsSeaDxConfigured() bool {
	return c.SeaDxBaseURL != ""
}

// GetSyncInterval returns the sync interval duration, defaulting to 24 hours if parsing fails.
func (c *Config) GetSyncInterval() time.Duration {
	duration, err := time.ParseDuration(c.SyncInterval)
	if err != nil {
		return 24 * time.Hour // Default to 24 hours
	}
	return duration
}

// IsSyncEnabled returns whether the background sync scheduler should run.
func (c *Config) IsSyncEnabled() bool {
	return c.SyncEnabled
}

// ShouldSkipUnchanged controls whether the sync tries to avoid sending unchanged docs
func (c *Config) ShouldSkipUnchanged() bool { return c.SyncSkipUnchanged }

// GetSyncEntities parses SYNC_ENTITIES into a normalized, de-duplicated slice.
// Supported values: artists, release-groups, releases, recordings
func (c *Config) GetSyncEntities() []string {
	raw := strings.TrimSpace(c.SyncEntities)
	if raw == "" {
		return []string{"artists", "release-groups", "releases", "recordings"}
	}
	parts := strings.Split(raw, ",")
	seen := map[string]bool{}
	var out []string
	norm := func(s string) string {
		s = strings.TrimSpace(strings.ToLower(s))
		switch s {
		case "release_groups", "releasegroup", "releasegroups":
			return "release-groups"
		case "artist", "artists":
			return "artists"
		case "release", "releases":
			return "releases"
		case "recording", "recordings":
			return "recordings"
		case "release-groups":
			return "release-groups"
		default:
			return "" // ignore unknowns
		}
	}
	for _, p := range parts {
		v := norm(p)
		if v == "" || seen[v] {
			continue
		}
		seen[v] = true
		out = append(out, v)
	}
	if len(out) == 0 {
		return []string{"artists", "release-groups", "releases", "recordings"}
	}
	return out
}

// Sync performance getters (with sane fallbacks)
func (c *Config) GetSyncDBPageSize() int {
	if c.SyncDBPageSize <= 0 {
		return 8000
	}
	return c.SyncDBPageSize
}

func (c *Config) GetSyncShardParallelism() int {
	if c.SyncShardParallelism <= 0 {
		// default to number of CPUs, at least 2
		n := 0
		// avoid importing runtime here; choose a conservative default
		// caller may override from their environment
		n = 4
		if n < 2 {
			n = 2
		}
		return n
	}
	return c.SyncShardParallelism
}

func (c *Config) GetSyncImportBatchSize() int {
	if c.SyncImportBatchSize <= 0 {
		return 2000
	}
	return c.SyncImportBatchSize
}

func (c *Config) GetSyncImportWorkers() int {
	if c.SyncImportWorkers <= 0 {
		// default to a conservative 4
		return 4
	}
	return c.SyncImportWorkers
}

func (c *Config) GetSyncImportMaxRetries() int {
	if c.SyncImportMaxRetries <= 0 {
		return 3
	}
	return c.SyncImportMaxRetries
}

func (c *Config) GetSyncImportBackoff() time.Duration {
	d, err := time.ParseDuration(c.SyncImportBackoff)
	if err != nil || d <= 0 {
		return 400 * time.Millisecond
	}
	return d
}

func (c *Config) GetSyncImportGlobalLimit() int {
	if c.SyncImportGlobalLimit <= 0 {
		return 0 // disabled by default
	}
	return c.SyncImportGlobalLimit
}
