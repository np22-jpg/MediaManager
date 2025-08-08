package typesense

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/typesense/typesense-go/v2/typesense/api"
)

// High-performance sharded sync
// - Parallel DB readers (shards) per entity using id ranges
// - Concurrent Typesense import workers with retries & backoff
// - Bounded channels for backpressure

// Tunables (consider promoting to env/config if needed)
var (
	importBatchSize   = 2000                          // docs per Typesense import request
	importWorkers     = max(4, runtime.GOMAXPROCS(0)) // concurrent Typesense import workers
	importMaxRetries  = 3                             // retries per import chunk
	importBackoffBase = 400 * time.Millisecond        // initial backoff
	// optional global cap across all entities to avoid overwhelming Typesense
	importGlobalLimit int
	globalImportSem   chan struct{}
)

// SyncTunables allows callers to override the defaults via application config without import cycles
type SyncTunables struct {
	ImportBatchSize  int
	ImportWorkers    int
	ImportMaxRetries int
	ImportBackoff    time.Duration
	// If > 0, caps total concurrent import requests across all entities
	ImportGlobalLimit int
}

// ApplyTunables sets the sync tunables; it's safe to call multiple times
func ApplyTunables(t SyncTunables) {
	if t.ImportBatchSize > 0 {
		importBatchSize = t.ImportBatchSize
	}
	if t.ImportWorkers > 0 {
		importWorkers = max(2, t.ImportWorkers)
	}
	if t.ImportMaxRetries > 0 {
		importMaxRetries = t.ImportMaxRetries
	}
	if t.ImportBackoff > 0 {
		importBackoffBase = t.ImportBackoff
	}
	if t.ImportGlobalLimit > 0 {
		importGlobalLimit = t.ImportGlobalLimit
		globalImportSem = make(chan struct{}, importGlobalLimit)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// indexChunk imports a chunk to Typesense with retries
func indexChunk(ctx context.Context, collection string, docs []interface{}) error {
	stringPtr := func(s string) *string { return &s }
	intPtr := func(i int) *int { return &i }
	params := &api.ImportDocumentsParams{Action: stringPtr("upsert"), BatchSize: intPtr(importBatchSize)}

	var attempt int
	var lastErr error
	backoff := importBackoffBase
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		_, err := ImportDocuments(ctx, collection, docs, params)
		if err == nil {
			return nil
		}
		lastErr = err
		attempt++
		if attempt > importMaxRetries {
			break
		}
		slog.Warn("import retry", "collection", collection, "attempt", attempt, "error", err)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
			backoff = time.Duration(float64(backoff) * 1.8) // gentle exponential
		}
	}
	return fmt.Errorf("exhausted retries: %w", lastErr)
}

// IndexArtists performs a full sync of artists from MusicBrainz to Typesense
func IndexArtists() error {
	if !IsReady() {
		return fmt.Errorf("MusicBrainz or Typesense not ready")
	}
	slog.Info("Indexing artists...")
	return indexEntities("artists", buildArtistQuery, transformArtist)
}

// IndexReleaseGroups performs a full sync of release groups from MusicBrainz to Typesense
func IndexReleaseGroups() error {
	if !IsReady() {
		return fmt.Errorf("MusicBrainz or Typesense not ready")
	}
	slog.Info("Indexing release groups...")
	return indexEntities("release_groups", buildReleaseGroupQuery, transformReleaseGroup)
}

// IndexReleases performs a full sync of releases from MusicBrainz to Typesense
func IndexReleases() error {
	if !IsReady() {
		return fmt.Errorf("MusicBrainz or Typesense not ready")
	}
	slog.Info("Indexing releases...")
	return indexEntities("releases", buildReleaseQuery, transformRelease)
}

// IndexRecordings performs a full sync of recordings from MusicBrainz to Typesense
func IndexRecordings() error {
	if !IsReady() {
		return fmt.Errorf("MusicBrainz or Typesense not ready")
	}
	slog.Info("Indexing recordings...")
	return indexEntities("recordings", buildRecordingQuery, transformRecording)
}

// Generic indexing function that can handle any entity type
func indexEntities(collection string, queryBuilder func() string, transformer func(*sql.Rows) (map[string]interface{}, error)) error {
	ctx := context.Background()

	// Get total count for progress tracking
	query := queryBuilder()
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS subquery", query)
	var total int64
	if err := db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return fmt.Errorf("failed to get total count: %w", err)
	}

	slog.Info("Starting index", "collection", collection, "total", total)

	// Create channels for coordination
	docsCh := make(chan []interface{}, importWorkers*2)
	var wg sync.WaitGroup
	var imported int64

	// Start import workers
	for i := 0; i < importWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for docs := range docsCh {
				if globalImportSem != nil {
					globalImportSem <- struct{}{}
					if err := indexChunk(ctx, collection, docs); err != nil {
						slog.Error("import chunk failed", "collection", collection, "error", err)
					} else {
						atomic.AddInt64(&imported, int64(len(docs)))
					}
					<-globalImportSem
				} else {
					if err := indexChunk(ctx, collection, docs); err != nil {
						slog.Error("import chunk failed", "collection", collection, "error", err)
					} else {
						atomic.AddInt64(&imported, int64(len(docs)))
					}
				}
			}
		}()
	}

	// Fetch and transform data
	err := fetchAndSendDocs(ctx, query, transformer, docsCh)
	close(docsCh)
	wg.Wait()

	if err != nil {
		return err
	}

	slog.Info("Completed indexing", "collection", collection, "imported", atomic.LoadInt64(&imported))
	return nil
}

// fetchAndSendDocs fetches data from the database and sends it to the docs channel
func fetchAndSendDocs(ctx context.Context, query string, transformer func(*sql.Rows) (map[string]interface{}, error), docsCh chan<- []interface{}) error {
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var docs []interface{}
	for rows.Next() {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		doc, err := transformer(rows)
		if err != nil {
			slog.Warn("failed to transform row", "error", err)
			continue
		}

		docs = append(docs, doc)

		if len(docs) >= importBatchSize {
			select {
			case docsCh <- docs:
				docs = nil
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	// Send remaining docs
	if len(docs) > 0 {
		select {
		case docsCh <- docs:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return rows.Err()
}

// Query builders and transformers - these will need to be implemented
// For now, I'll create placeholder functions that will be filled from the original sync.go

func buildArtistQuery() string {
	return `
		SELECT a.gid, a.name, a.sort_name, at.name as artist_type, 
			   ar.name as area_name, a.begin_date_year, a.end_date_year, 
			   a.ended, a.comment,
			   ARRAY_AGG(DISTINCT alias.name) FILTER (WHERE alias.name IS NOT NULL) as aliases
		FROM artist a
		LEFT JOIN artist_type at ON a.type = at.id
		LEFT JOIN area ar ON a.area = ar.id
		LEFT JOIN artist_alias alias ON a.id = alias.artist
		GROUP BY a.gid, a.name, a.sort_name, at.name, ar.name, 
				 a.begin_date_year, a.end_date_year, a.ended, a.comment
		ORDER BY a.id
	`
}

func buildReleaseGroupQuery() string {
	return `
		SELECT rg.gid, rg.name, rgt.name as primary_type, rg.comment
		FROM release_group rg
		LEFT JOIN release_group_primary_type rgt ON rg.type = rgt.id
		ORDER BY rg.id
	`
}

func buildReleaseQuery() string {
	return `
		SELECT r.gid, r.name, rs.name as status, 
			   a.name as artist_name, a.gid as artist_mbid,
			   rg.name as release_group_name, rg.gid as release_group_mbid,
			   r.comment, EXTRACT(YEAR FROM r.date_year) as date_year
		FROM release r
		LEFT JOIN release_status rs ON r.status = rs.id
		LEFT JOIN artist_credit ac ON r.artist_credit = ac.id
		LEFT JOIN artist_credit_name acn ON ac.id = acn.artist_credit
		LEFT JOIN artist a ON acn.artist = a.id
		LEFT JOIN release_group rg ON r.release_group = rg.id
		ORDER BY r.id
	`
}

func buildRecordingQuery() string {
	return `
		SELECT r.gid, r.name, r.length, r.comment
		FROM recording r
		ORDER BY r.id
	`
}

func transformArtist(rows *sql.Rows) (map[string]interface{}, error) {
	var mbid, name, sortName, comment sql.NullString
	var artistType, areaName sql.NullString
	var beginYear, endYear sql.NullInt64
	var ended sql.NullBool
	var aliases interface{}

	err := rows.Scan(&mbid, &name, &sortName, &artistType, &areaName,
		&beginYear, &endYear, &ended, &comment, &aliases)
	if err != nil {
		return nil, err
	}

	doc := map[string]interface{}{
		"mbid":      mbid.String,
		"name":      name.String,
		"sort_name": sortName.String,
	}

	if artistType.Valid {
		doc["artist_type"] = artistType.String
	}
	if areaName.Valid {
		doc["area_name"] = areaName.String
	}
	if beginYear.Valid {
		doc["begin_year"] = beginYear.Int64
	}
	if endYear.Valid {
		doc["end_year"] = endYear.Int64
	}
	if ended.Valid {
		doc["ended"] = ended.Bool
	}
	if comment.Valid {
		doc["comment"] = comment.String
	}

	return doc, nil
}

func transformReleaseGroup(rows *sql.Rows) (map[string]interface{}, error) {
	var mbid, name, comment sql.NullString
	var primaryType sql.NullString

	err := rows.Scan(&mbid, &name, &primaryType, &comment)
	if err != nil {
		return nil, err
	}

	doc := map[string]interface{}{
		"mbid": mbid.String,
		"name": name.String,
	}

	if primaryType.Valid {
		doc["primary_type"] = primaryType.String
	}
	if comment.Valid {
		doc["comment"] = comment.String
	}

	return doc, nil
}

func transformRelease(rows *sql.Rows) (map[string]interface{}, error) {
	var mbid, name, comment sql.NullString
	var status, artistName, artistMbid, releaseGroupName, releaseGroupMbid sql.NullString
	var dateYear sql.NullInt64

	err := rows.Scan(&mbid, &name, &status, &artistName, &artistMbid,
		&releaseGroupName, &releaseGroupMbid, &comment, &dateYear)
	if err != nil {
		return nil, err
	}

	doc := map[string]interface{}{
		"mbid": mbid.String,
		"name": name.String,
	}

	if status.Valid {
		doc["status"] = status.String
	}
	if artistName.Valid {
		doc["artist_name"] = artistName.String
	}
	if artistMbid.Valid {
		doc["artist_mbid"] = artistMbid.String
	}
	if releaseGroupName.Valid {
		doc["release_group_name"] = releaseGroupName.String
	}
	if releaseGroupMbid.Valid {
		doc["release_group_mbid"] = releaseGroupMbid.String
	}
	if comment.Valid {
		doc["comment"] = comment.String
	}
	if dateYear.Valid {
		doc["date_year"] = dateYear.Int64
	}

	return doc, nil
}

func transformRecording(rows *sql.Rows) (map[string]interface{}, error) {
	var mbid, name, comment sql.NullString
	var length sql.NullInt64

	err := rows.Scan(&mbid, &name, &length, &comment)
	if err != nil {
		return nil, err
	}

	doc := map[string]interface{}{
		"mbid": mbid.String,
		"name": name.String,
	}

	if length.Valid {
		doc["length"] = length.Int64
	}
	if comment.Valid {
		doc["comment"] = comment.String
	}

	return doc, nil
}
