package sync

import (
	"context"
	"log/slog"
	"time"

	"relay/app"
	"relay/app/music/musicbrainz"
)

// Scheduler manages periodic sync operations to keep Typesense index updated.
type Scheduler struct {
	interval time.Duration
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewScheduler creates a new sync scheduler with the specified interval.
func NewScheduler(interval time.Duration) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		interval: interval,
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Start begins the periodic sync operations in a background goroutine.
// Includes an initial delay before first sync and then runs at regular intervals.
func (s *Scheduler) Start() {
	slog.Info("Starting sync scheduler", "interval", s.interval)

	// Run initial sync after a short delay to allow system startup
	go func() {
		select {
		case <-time.After(5 * time.Minute): // Wait 5 minutes after startup
			s.runSync()
		case <-s.ctx.Done():
			return
		}

		// Then run periodically at the configured interval
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.runSync()
			case <-s.ctx.Done():
				slog.Info("Sync scheduler stopped")
				return
			}
		}
	}()
}

// Stop stops the sync scheduler and cancels any running operations.
func (s *Scheduler) Stop() {
	slog.Info("Stopping sync scheduler")
	s.cancel()
}

// runSync executes a full sync operation to update the Typesense search index.
func (s *Scheduler) runSync() {
	slog.Info("Starting scheduled sync operation")
	start := time.Now()

	// Only run if both MusicBrainz and Typesense are configured and ready
	if !musicbrainz.IsReady() {
		slog.Warn("Skipping sync - MusicBrainz or Typesense not ready")
		return
	}

	entities := app.AppConfig.GetSyncEntities()
	slog.Info("Scheduled sync entities", "entities", entities)

	for _, e := range entities {
		switch e {
		case "artists":
			if err := musicbrainz.IndexArtists(); err != nil {
				slog.Error("Failed to sync artists", "error", err)
			} else {
				slog.Info("Successfully synced artists")
			}
		case "release-groups":
			if err := musicbrainz.IndexReleaseGroups(); err != nil {
				slog.Error("Failed to sync release groups", "error", err)
			} else {
				slog.Info("Successfully synced release groups")
			}
		case "releases":
			if err := musicbrainz.IndexReleases(); err != nil {
				slog.Error("Failed to sync releases", "error", err)
			} else {
				slog.Info("Successfully synced releases")
			}
		case "recordings":
			if err := musicbrainz.IndexRecordings(); err != nil {
				slog.Error("Failed to sync recordings", "error", err)
			} else {
				slog.Info("Successfully synced recordings")
			}
		default:
			slog.Warn("Unknown entity in SYNC_ENTITIES - skipping", "entity", e)
		}
	}

	duration := time.Since(start)
	slog.Info("Completed scheduled sync operation", "duration", duration)
}
