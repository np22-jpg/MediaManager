# IMPLEMENTATION NOTES

## Web Server
This application now uses the standard `net/http` library with a custom router implementation that supports path parameters and middleware chaining. Previously, Gin was used for better performance, but during migration testing, the performance difference was minimal in real-world scenarios. The custom router provides the needed functionality without external dependencies. If throughput becomes a problem in the future, moving to [fiber](https://gofiber.io/) might be considered.

## Metrics
For metrics, it was decided to go with `VictoriaMetrics/metrics` for its lightweight footprint and excellent performance. This replaced the previous `prometheus/client_golang` which was quite large. VictoriaMetrics provides a smaller binary size, lower memory usage, and faster metric operations while maintaining full Prometheus compatibility for scraping.

## Docker 
The Docker image consists of a static version of this program, which is then loaded into busybox. This provides a nice list of CLI utilities alongside the main program. It is distroless in the sense that there is no package manager, and it runs rootless. Further hardening should be considered, though.

## Metabrainz
This one is a doozy.

### Song search

Originally, it ran by directly querying the public metadata API; however, it has severe rate limits. The only way around this is to simply host a full DB mirror, in the same manner other programs do. 

PostgreSQL full-text search was initially tried but was complex and slow (over a minute per query). The current implementation uses Typesense for a faster search

Syncing is scheduled and done automatically without the use of a third-party program with `sync/scheduler.go` and `musicbrainz/sync.go`.

#### Sync engine (current)

- Sharded parallel DB scanning per entity (artists, release_groups, recordings):
	- Compute min/max id and split into contiguous id ranges (shards).
	- Each shard reads with cursor-based pagination: `WHERE id > $1 AND id <= $2 ORDER BY id LIMIT $3`.
	- Avoids OFFSET and expensive JOINs; queries only essential columns.
- Concurrent Typesense importers:
	- Buffered channel feeds import workers.
	- Upsert via `ImportDocuments` with chunking and retry/backoff.
	- Bounded queues preserve backpressure.
- Parallel entity sync in "sync all" mode: artists, release groups, and recordings run simultaneously.

If maintenance burden ends up being too much, it might be advisable to use a third-party program instead and do periodic batch syncs. 

The periodic schema changes might make it difficult to keep going. If it becomes a serious problem, it'd be a good idea to swap back to sir/solr and/or use all the upstream Metabrainz Docker software. In this scenario, we'd use their software to host the endpoints, then proxy those through our Redis caching server. That is:
PostgreSQL -> [sir -> solr] -> Their web API -> This caching service

That being said, their web API, even when self-hosted, is apparently quite slow, and even a distributed cache might not be able to keep up.


## Tests
Right now, the tests are set up to automatically pass, even when a metabrainz/typesense server is not available. To run full integration tests, set the `INTEGRATION_TESTS` env var to "true".

## Future
- hardcover support
- PGO
- get the debug packages out of the build