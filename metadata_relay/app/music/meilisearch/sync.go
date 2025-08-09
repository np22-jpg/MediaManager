package meilisearch

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// Configuration variables for sync performance tuning
var (
	dbPageSize       = 8000 // Rows fetched per DB page per shard
	shardParallelism = 0    // Parallel DB reader shards per entity
	importBatchSize  = 2000 // Documents per Meilisearch import call
	importWorkers    = 0    // Concurrent Meilisearch import workers per entity
	importMaxRetries = 3    // Retries per failed import chunk
	importBackoff    = 400 * time.Millisecond
	globalSemaphore  chan struct{} // Global rate limiter across all entities
)

// TunableConfig represents adjustable sync parameters
type TunableConfig struct {
	ImportBatchSize   int
	ImportWorkers     int
	ImportMaxRetries  int
	ImportBackoff     time.Duration
	ImportGlobalLimit int
}

// ApplyTunables sets global sync parameters for all sync operations
func ApplyTunables(config TunableConfig) {
	if config.ImportBatchSize > 0 {
		importBatchSize = config.ImportBatchSize
	}
	if config.ImportWorkers > 0 {
		importWorkers = config.ImportWorkers
	}
	if config.ImportMaxRetries > 0 {
		importMaxRetries = config.ImportMaxRetries
	}
	if config.ImportBackoff > 0 {
		importBackoff = config.ImportBackoff
	}
	if config.ImportGlobalLimit > 0 {
		globalSemaphore = make(chan struct{}, config.ImportGlobalLimit)
	}
	slog.Info("Applied sync tunables", "batchSize", importBatchSize, "workers", importWorkers,
		"maxRetries", importMaxRetries, "backoff", importBackoff,
		"globalLimit", config.ImportGlobalLimit)
}

// indexChunk uploads a batch of documents to Meilisearch with retry logic
func indexChunk(ctx context.Context, indexName string, docs []interface{}) error {
	if globalSemaphore != nil {
		globalSemaphore <- struct{}{}
		defer func() { <-globalSemaphore }()
	}

	var lastErr error
	for attempt := 0; attempt <= importMaxRetries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(attempt) * importBackoff
			slog.Debug("retrying import chunk", "attempt", attempt, "backoff", backoff)
			time.Sleep(backoff)
		}

		if err := ImportDocuments(ctx, indexName, docs); err != nil {
			lastErr = err
			slog.Warn("import chunk failed", "attempt", attempt, "error", err)
			continue
		}
		return nil
	}
	return fmt.Errorf("failed to import chunk after %d attempts: %w", importMaxRetries, lastErr)
}

// IndexArtists indexes all artists from the MusicBrainz database to Meilisearch
func IndexArtists() error {
	return indexEntities("artists", func() string {
		return `
			SELECT a.gid, a.name, a.sort_name, 
				   at.name as artist_type, area.name as area_name,
				   a.begin_date_year, a.end_date_year, a.ended, a.comment
			FROM artist a
			LEFT JOIN artist_type at ON a.type = at.id
			LEFT JOIN area ON a.area = area.id`
	}, func(rows *sql.Rows) (map[string]interface{}, error) {
		var artist map[string]interface{}
		var gid, name, sortName sql.NullString
		var artistType, areaName, comment sql.NullString
		var beginYear, endYear sql.NullInt32
		var ended sql.NullBool

		err := rows.Scan(&gid, &name, &sortName, &artistType, &areaName,
			&beginYear, &endYear, &ended, &comment)
		if err != nil {
			return nil, err
		}

		artist = map[string]interface{}{
			"mbid":      gid.String,
			"name":      name.String,
			"sort_name": sortName.String,
		}

		if artistType.Valid {
			artist["artist_type"] = artistType.String
		}
		if areaName.Valid {
			artist["area_name"] = areaName.String
		}
		if beginYear.Valid {
			artist["begin_year"] = beginYear.Int32
		}
		if endYear.Valid {
			artist["end_year"] = endYear.Int32
		}
		if ended.Valid {
			artist["ended"] = ended.Bool
		}
		if comment.Valid {
			artist["comment"] = comment.String
		}

		return artist, nil
	})
}

// IndexReleaseGroups indexes all release groups from the MusicBrainz database to Meilisearch
func IndexReleaseGroups() error {
	return indexEntities("release_groups", func() string {
		return `
			SELECT rg.gid, rg.name, rgt.name as primary_type, rg.comment
			FROM release_group rg
			LEFT JOIN release_group_primary_type rgt ON rg.type = rgt.id`
	}, func(rows *sql.Rows) (map[string]interface{}, error) {
		var releaseGroup map[string]interface{}
		var gid, name sql.NullString
		var primaryType, comment sql.NullString

		err := rows.Scan(&gid, &name, &primaryType, &comment)
		if err != nil {
			return nil, err
		}

		releaseGroup = map[string]interface{}{
			"mbid": gid.String,
			"name": name.String,
		}

		if primaryType.Valid {
			releaseGroup["primary_type"] = primaryType.String
		}
		if comment.Valid {
			releaseGroup["comment"] = comment.String
		}

		return releaseGroup, nil
	})
}

// IndexReleases indexes all releases from the MusicBrainz database to Meilisearch
func IndexReleases() error {
	return indexEntities("releases", func() string {
		return `
			SELECT r.gid, r.name, rs.name as status, r.comment,
				   ac.name as artist_name, a.gid as artist_mbid,
				   rg.name as release_group_name, rg.gid as release_group_mbid
			FROM release r
			LEFT JOIN release_status rs ON r.status = rs.id
			LEFT JOIN artist_credit ac ON r.artist_credit = ac.id
			LEFT JOIN artist_credit_name acn ON ac.id = acn.artist_credit AND acn.position = 0
			LEFT JOIN artist a ON acn.artist = a.id
			LEFT JOIN release_group rg ON r.release_group = rg.id`
	}, func(rows *sql.Rows) (map[string]interface{}, error) {
		var release map[string]interface{}
		var gid, name, status, comment sql.NullString
		var artistName, artistMbid sql.NullString
		var releaseGroupName, releaseGroupMbid sql.NullString

		err := rows.Scan(&gid, &name, &status, &comment,
			&artistName, &artistMbid, &releaseGroupName, &releaseGroupMbid)
		if err != nil {
			return nil, err
		}

		release = map[string]interface{}{
			"mbid": gid.String,
			"name": name.String,
		}

		if status.Valid {
			release["status"] = status.String
		}
		if comment.Valid {
			release["comment"] = comment.String
		}
		if artistName.Valid {
			release["artist_name"] = artistName.String
		}
		if artistMbid.Valid {
			release["artist_mbid"] = artistMbid.String
		}
		if releaseGroupName.Valid {
			release["release_group_name"] = releaseGroupName.String
		}
		if releaseGroupMbid.Valid {
			release["release_group_mbid"] = releaseGroupMbid.String
		}

		return release, nil
	})
}

// IndexRecordings indexes all recordings from the MusicBrainz database to Meilisearch
func IndexRecordings() error {
	return indexEntities("recordings", func() string {
		return `
			SELECT r.gid, r.name, r.length, r.comment
			FROM recording r`
	}, func(rows *sql.Rows) (map[string]interface{}, error) {
		var recording map[string]interface{}
		var gid, name sql.NullString
		var length sql.NullInt32
		var comment sql.NullString

		err := rows.Scan(&gid, &name, &length, &comment)
		if err != nil {
			return nil, err
		}

		recording = map[string]interface{}{
			"mbid": gid.String,
			"name": name.String,
		}

		if length.Valid {
			recording["length"] = length.Int32
		}
		if comment.Valid {
			recording["comment"] = comment.String
		}

		return recording, nil
	})
}

// indexEntities is a generic function to index any entity type to Meilisearch
func indexEntities(indexName string, queryBuilder func() string, transformer func(*sql.Rows) (map[string]interface{}, error)) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}
	if meilisearchClient == nil {
		return fmt.Errorf("meilisearch client not initialized")
	}

	slog.Info("Starting indexing", "entity", indexName)

	// Determine optimal parallelism
	workers := importWorkers
	if workers <= 0 {
		workers = runtime.NumCPU()
		if workers < 2 {
			workers = 2
		}
	}

	shards := shardParallelism
	if shards <= 0 {
		shards = runtime.NumCPU()
		if shards < 2 {
			shards = 2
		}
	}

	// Channel for work distribution
	workCh := make(chan []interface{}, workers*2)
	var wg sync.WaitGroup
	var successful, failed int64

	// Start import workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for batch := range workCh {
				if err := indexChunk(context.Background(), indexName, batch); err != nil {
					slog.Error("worker failed to import batch", "worker", workerID, "error", err)
					atomic.AddInt64(&failed, 1)
				} else {
					atomic.AddInt64(&successful, 1)
				}
			}
		}(i)
	}

	// Producer: fetch data and send to workers
	go func() {
		defer close(workCh)

		query := queryBuilder() + " ORDER BY 1 LIMIT $1 OFFSET $2"
		offset := 0

		for {
			rows, err := db.Query(query, dbPageSize, offset)
			if err != nil {
				slog.Error("failed to query database", "error", err)
				return
			}

			var batch []interface{}
			for rows.Next() {
				entity, err := transformer(rows)
				if err != nil {
					slog.Error("failed to transform row", "error", err)
					continue
				}
				batch = append(batch, entity)

				if len(batch) >= importBatchSize {
					workCh <- batch
					batch = nil
				}
			}
			rows.Close()

			// Send remaining items
			if len(batch) > 0 {
				workCh <- batch
			}

			if len(batch) < dbPageSize {
				break // No more data
			}

			offset += dbPageSize
			slog.Info("Processed page", "entity", indexName, "offset", offset)
		}
	}()

	// Wait for all workers to complete
	wg.Wait()

	slog.Info("Indexing completed", "entity", indexName, "successful", successful, "failed", failed)
	return nil
}
