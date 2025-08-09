package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"relay/app"
	"relay/app/cache"
	"relay/app/music"
	"relay/app/music/musicbrainz"
	"relay/app/music/theaudiodb"
	"relay/app/music/typesense"
	"relay/app/seadex"
	sched "relay/app/sync"
	"relay/app/tmdb"
	"relay/app/tvdb"

	"github.com/VictoriaMetrics/metrics"
)

// runSyncCommand handles the sync command for indexing MusicBrainz data to Typesense.
// This is a standalone operation that can be run independently of the web server.
// Supports targeted sync operations with optional entity type arguments.
func runSyncCommand() {
	// Load configuration
	if err := app.LoadConfig(); err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Initialize cache (required for some database functions)
	cache.InitCache(app.AppConfig.CacheHost, app.AppConfig.CachePort, app.AppConfig.CacheDB)

	// Check if MusicBrainz is configured
	if !app.AppConfig.IsMusicBrainzConfigured() {
		slog.Error("MusicBrainz is not configured - cannot run sync")
		os.Exit(1)
	}

	// Initialize MusicBrainz database connection
	musicbrainz.InitMusicBrainz(app.AppConfig.GetMusicBrainzConnStr())

	// Check if Typesense is configured
	if !app.AppConfig.IsTypesenseConfigured() {
		slog.Error("Typesense is not configured - cannot run sync")
		os.Exit(1)
	}

	// Initialize Typesense client
	err := musicbrainz.InitTypesense(app.AppConfig.TypesenseHost, app.AppConfig.TypesensePort, app.AppConfig.TypesenseAPIKey, app.AppConfig.GetTypesenseTimeout())
	if err != nil {
		slog.Error("failed to initialize Typesense", "error", err)
		os.Exit(1)
	}

	// Apply sync tunables from config
	musicbrainz.ApplyTunables(typesense.SyncTunables{
		ImportBatchSize:   app.AppConfig.GetSyncImportBatchSize(),
		ImportWorkers:     app.AppConfig.GetSyncImportWorkers(),
		ImportMaxRetries:  app.AppConfig.GetSyncImportMaxRetries(),
		ImportBackoff:     app.AppConfig.GetSyncImportBackoff(),
		ImportGlobalLimit: app.AppConfig.GetSyncImportGlobalLimit(),
	})

	// Determine what to sync based on command arguments
	var syncTarget string
	if len(os.Args) > 2 {
		syncTarget = os.Args[2]
	} else {
		syncTarget = "all"
	}

	slog.Info("Starting data sync to Typesense", "target", syncTarget)

	switch syncTarget {
	case "artists":
		slog.Info("Indexing artists...")
		if err := musicbrainz.IndexArtists(); err != nil {
			slog.Error("failed to index artists", "error", err)
			os.Exit(1)
		}
		slog.Info("✓ Artists indexed")

	case "release-groups":
		slog.Info("Indexing release groups...")
		if err := musicbrainz.IndexReleaseGroups(); err != nil {
			slog.Error("failed to index release groups", "error", err)
			os.Exit(1)
		}
		slog.Info("✓ Release groups indexed")

	case "releases":
		slog.Info("Indexing releases...")
		if err := musicbrainz.IndexReleases(); err != nil {
			slog.Error("failed to index releases", "error", err)
			os.Exit(1)
		}
		slog.Info("✓ Releases indexed")

	case "recordings":
		slog.Info("Indexing recordings...")
		if err := musicbrainz.IndexRecordings(); err != nil {
			slog.Error("failed to index recordings", "error", err)
			os.Exit(1)
		}
		slog.Info("✓ Recordings indexed")

	case "all":
		// Build the list of entities from config (SYNC_ENTITIES)
		entities := app.AppConfig.GetSyncEntities()
		slog.Info("Indexing entities in parallel...", "entities", entities)
		errCh := make(chan error, len(entities))
		var wg sync.WaitGroup
		wg.Add(len(entities))

		run := func(name string, fn func() error) {
			defer wg.Done()
			if err := fn(); err != nil {
				errCh <- fmt.Errorf("%s: %w", name, err)
			} else {
				slog.Info("✓ indexed", "entity", name)
			}
		}

		for _, e := range entities {
			switch e {
			case "artists":
				go run("artists", musicbrainz.IndexArtists)
			case "release-groups":
				go run("release-groups", musicbrainz.IndexReleaseGroups)
			case "releases":
				go run("releases", musicbrainz.IndexReleases)
			case "recordings":
				go run("recordings", musicbrainz.IndexRecordings)
			default:
				wg.Done()
				slog.Warn("Unknown entity in SYNC_ENTITIES - skipping", "entity", e)
			}
		}
		wg.Wait()
		close(errCh)
		if err, ok := <-errCh; ok {
			slog.Error("failed to index one or more entities", "error", err)
			os.Exit(1)
		}
		slog.Info("✓ All requested entities indexed")

	default:
		slog.Error("Invalid sync target", "target", syncTarget, "valid_targets", []string{"artists", "release-groups", "releases", "recordings", "all"})
		os.Exit(1)
	}

	slog.Info("🎉 Sync completed successfully!", "target", syncTarget)
}

func main() {
	// Check if this is a sync command
	if len(os.Args) > 1 && os.Args[1] == "sync" {
		runSyncCommand()
		return
	}

	// Load configuration
	if err := app.LoadConfig(); err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Initialize cache
	cache.InitCache(app.AppConfig.CacheHost, app.AppConfig.CachePort, app.AppConfig.CacheDB)

	// Initialize API clients
	tmdb.InitTMDB(app.AppConfig.TMDBAPIKey, app.AppConfig.TMDBBaseURL)
	tvdb.InitTVDB(app.AppConfig.TVDBAPIKey, app.AppConfig.TVDBBaseURL)

	// Initialize TheAudioDB (optional)
	if app.AppConfig.TheAudioDBAPIKey != "" {
		base := app.AppConfig.TheAudioDBBaseURL
		if base == "" {
			base = "https://www.theaudiodb.com/api/v1/json"
		}
		audClient := theaudiodb.New(base, app.AppConfig.TheAudioDBAPIKey)
		theaudiodb.SetClient(audClient)
	}

	// Media directory
	musicbrainz.SetMediaDir(app.AppConfig.MediaDir)

	// Initialize Spotify (optional)
	if app.AppConfig.SpotifyClientID != "" && app.AppConfig.SpotifyClientSecret != "" {
		sp := music.NewSpotify(app.AppConfig.SpotifyClientID, app.AppConfig.SpotifyClientSecret)
		musicbrainz.SetSpotifyClient(sp)
	}

	// Initialize LRCLib
	lyrics := music.NewLRCLib(app.AppConfig.LRCLibBaseURL)
	musicbrainz.SetLyricsClient(lyrics)

	// Initialize SeaDx conditionally
	var seadexEnabled bool
	if app.AppConfig.IsSeaDxConfigured() {
		seadex.InitSeaDex(app.AppConfig.SeaDxBaseURL)
		seadexEnabled = true
		slog.Info("SeaDx initialized successfully")
	} else {
		seadexEnabled = false
		slog.Info("SeaDx not configured - skipping")
	}

	// Initialize MusicBrainz conditionally
	var musicBrainzEnabled bool
	// Always try to initialize MusicBrainz (similar to TMDB/TVDB)
	musicbrainz.InitMusicBrainz(app.AppConfig.GetMusicBrainzConnStr())

	if !app.AppConfig.IsMusicBrainzConfigured() {
		musicBrainzEnabled = false
	} else {
		musicBrainzEnabled = true
		slog.Info("MusicBrainz initialized successfully")
	}

	// Initialize Typesense conditionally (only if MusicBrainz is also enabled)
	if musicBrainzEnabled {
		if !app.AppConfig.IsTypesenseConfigured() {
			slog.Warn("Typesense is not configured - search will not be available")
		} else {
			err := musicbrainz.InitTypesense(app.AppConfig.TypesenseHost, app.AppConfig.TypesensePort, app.AppConfig.TypesenseAPIKey, app.AppConfig.GetTypesenseTimeout())
			if err != nil {
				slog.Error("failed to initialize Typesense (configured but connection failed)", "error", err)
				os.Exit(1)
			}
			slog.Info("Typesense initialized successfully")

			// Start sync scheduler if both MusicBrainz and Typesense are available and sync is enabled
			if musicbrainz.IsReady() && app.AppConfig.IsSyncEnabled() {
				// Ensure tunables applied for scheduler-run syncs
				musicbrainz.ApplyTunables(typesense.SyncTunables{
					ImportBatchSize:   app.AppConfig.GetSyncImportBatchSize(),
					ImportWorkers:     app.AppConfig.GetSyncImportWorkers(),
					ImportMaxRetries:  app.AppConfig.GetSyncImportMaxRetries(),
					ImportBackoff:     app.AppConfig.GetSyncImportBackoff(),
					ImportGlobalLimit: app.AppConfig.GetSyncImportGlobalLimit(),
				})
				syncScheduler := sched.NewScheduler(app.AppConfig.GetSyncInterval())
				syncScheduler.Start()

				// Register cleanup for scheduler on shutdown
				defer syncScheduler.Stop()
			} else if !app.AppConfig.IsSyncEnabled() {
				slog.Info("Background sync scheduler is disabled via SYNC_ENABLED=false")
			}
		}
	}

	// Set up structured logger with configurable level
	logLevel := app.AppConfig.GetLogLevel()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	slog.Info("starting server")

	// Create HTTP router
	router := app.NewRouter()

	// Add middleware
	router.Use(app.LoggingMiddleware())
	router.Use(app.RecoveryMiddleware())
	router.Use(app.MetricsMiddleware())

	// AniList is always enabled as it requires no configuration
	anilistEnabled := true
	slog.Info("AniList enabled (no configuration required)")

	// Mount app routes
	app.RegisterRoutes(router, musicBrainzEnabled, seadexEnabled, anilistEnabled)

	// Serve media directory statically (images, lyrics files)
	// Ensure media dir exists
	if app.AppConfig.MediaDir != "" {
		if err := os.MkdirAll(app.AppConfig.MediaDir, 0o755); err != nil {
			slog.Warn("failed to create media dir", "dir", app.AppConfig.MediaDir, "error", err)
		}
		router.Handle("/media/", http.StripPrefix("/media/", http.FileServer(http.Dir(app.AppConfig.MediaDir))))
	}

	// Create metrics server
	metricsRouter := http.NewServeMux()
	metricsRouter.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		metrics.WritePrometheus(w, true)
	})

	// Graceful shutdown setup
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Main server
	srv := &http.Server{
		Addr:    app.AppConfig.GetServerAddr(),
		Handler: router,
	}

	// Metrics server
	metricsSrv := &http.Server{
		Addr:    app.AppConfig.GetMetricsAddr(),
		Handler: metricsRouter,
	}

	// Start servers in goroutines
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start main server", "error", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := metricsSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start metrics server", "error", err)
		}
	}()

	slog.Info("Servers started", "main_address", app.AppConfig.GetServerAddr(), "metrics_address", app.AppConfig.GetMetricsAddr())

	// Wait for interrupt signal to gracefully shutdown the servers
	<-ctx.Done()

	// The context is canceled, now attempt graceful shutdown
	slog.Info("Servers are shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown both servers
	go func() {
		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.Error("Main server forced to shutdown", "error", err)
		}
	}()

	go func() {
		if err := metricsSrv.Shutdown(shutdownCtx); err != nil {
			slog.Error("Metrics server forced to shutdown", "error", err)
		}
	}()

	// Wait for both servers to finish
	wg.Wait()

	slog.Info("Servers exited")
}
