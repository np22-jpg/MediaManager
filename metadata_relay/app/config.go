package app

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
)

// Config holds all application configuration settings loaded from environment variables.
type Config struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
	Port     string `env:"PORT" envDefault:"8000"`

	CacheHost string `env:"VALKEY_HOST" envDefault:"localhost"`
	CachePort int    `env:"VALKEY_PORT" envDefault:"6379"`
	CacheDB   int    `env:"VALKEY_DB" envDefault:"0"`

	TMDBAPIKey  string `env:"TMDB_API_KEY"`
	TMDBBaseURL string `env:"TMDB_BASE_URL" envDefault:"https://api.themoviedb.org/3"`
	TVDBAPIKey  string `env:"TVDB_API_KEY"`
	TVDBBaseURL string `env:"TVDB_BASE_URL" envDefault:"https://api4.thetvdb.com/v4"`

	// TheAudioDB Configuration
	TheAudioDBAPIKey  string `env:"THEAUDIODB_API_KEY"`
	TheAudioDBBaseURL string `env:"THEAUDIODB_BASE_URL" envDefault:"https://www.theaudiodb.com/api/v1/json"`

	// Spotify (for images)
	SpotifyClientID     string `env:"SPOTIFY_CLIENT_ID"`
	SpotifyClientSecret string `env:"SPOTIFY_CLIENT_SECRET"`

	// LRCLib (lyrics)
	LRCLibBaseURL string `env:"LRCLIB_BASE_URL" envDefault:"https://lrclib.net/api"`

	// SeaDx (anime torrents)
	SeaDxBaseURL string `env:"SEADX_BASE_URL" envDefault:"https://releases.moe/api"`

	// AniDB (anime database)
	AniDBBaseURL   string `env:"ANIDB_BASE_URL" envDefault:"http://api.anidb.info:9001/httpapi"`
	AniDBClient    string `env:"ANIDB_CLIENT"`
	AniDBClientVer string `env:"ANIDB_CLIENT_VER" envDefault:"1"`

	// Media storage directory (on disk, do not cache images in Redis)
	MediaDir string `env:"MEDIA_DIR" envDefault:"./media"`

	// MusicBrainz PostgreSQL Configuration
	MusicBrainzDBHost     string `env:"MUSICBRAINZ_DB_HOST"`
	MusicBrainzDBPort     string `env:"MUSICBRAINZ_DB_PORT" envDefault:"5432"`
	MusicBrainzDBUser     string `env:"MUSICBRAINZ_DB_USER" envDefault:"musicbrainz"`
	MusicBrainzDBPassword string `env:"MUSICBRAINZ_DB_PASSWORD" envDefault:"musicbrainz"`
	MusicBrainzDBName     string `env:"MUSICBRAINZ_DB_NAME"`

	// Typesense Configuration
	TypesenseHost   string `env:"TYPESENSE_HOST" envDefault:"localhost"`
	TypesensePort   string `env:"TYPESENSE_PORT" envDefault:"8108"`
	TypesenseAPIKey string `env:"TYPESENSE_API_KEY"`
	// HTTP timeout for Typesense client operations (e.g., bulk imports)
	TypesenseTimeout string `env:"TYPESENSE_TIMEOUT" envDefault:"60s"`

	// Sync Configuration
	SyncInterval string `env:"SYNC_INTERVAL" envDefault:"24h"` // How often to sync data to Typesense
	SyncEnabled  bool   `env:"SYNC_ENABLED" envDefault:"true"` // Toggle background sync scheduler
	// Comma-separated list of entities to sync when target is "all" or scheduler runs
	// Allowed values: artists, release-groups, releases, recordings
	SyncEntities string `env:"SYNC_ENTITIES" envDefault:"artists,release-groups,releases,recordings"`
	// If true, attempts to skip unchanged documents during sync using a content fingerprint cache
	SyncSkipUnchanged bool `env:"SYNC_SKIP_UNCHANGED" envDefault:"true"`

	// Sync Performance Tunables
	SyncDBPageSize       int    `env:"SYNC_DB_PAGE_SIZE" envDefault:"8000"`
	SyncShardParallelism int    `env:"SYNC_SHARD_PARALLELISM"` // default derived from CPU if 0
	SyncImportBatchSize  int    `env:"SYNC_IMPORT_BATCH_SIZE" envDefault:"2000"`
	SyncImportWorkers    int    `env:"SYNC_IMPORT_WORKERS"` // default derived from CPU if 0
	SyncImportMaxRetries int    `env:"SYNC_IMPORT_MAX_RETRIES" envDefault:"3"`
	SyncImportBackoff    string `env:"SYNC_IMPORT_BACKOFF" envDefault:"400ms"`
	// Global cap across all entities for concurrent Typesense import requests
	SyncImportGlobalLimit int `env:"SYNC_IMPORT_GLOBAL_LIMIT"`
}

var AppConfig Config

// LoadConfig parses environment variables and loads the application configuration.
func LoadConfig() error {
	if err := env.Parse(&AppConfig); err != nil {
		return err
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

// IsAniDBConfigured checks if AniDB is configured (requires client name).
func (c *Config) IsAniDBConfigured() bool {
	return c.AniDBClient != ""
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
