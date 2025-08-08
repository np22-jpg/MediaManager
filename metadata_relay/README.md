# MediaManager Metadata Relay

This is a service that provides metadata for movies, TV shows, and music. It caches the metadata to not overload the TMDB and TVDB APIs and uses a local MusicBrainz PostgreSQL mirror for music metadata. This service is used by MediaManager to fetch metadata for movies, TV shows, and music. I (the developer) run a public instance of this service at https://metadata-relay.maxid.me, but you can also run your own instance.

## Docker Compose Example

```yaml
version: '3.8'
services:
  valkey:
    image: valkey/valkey:latest
    restart: always
    container_name: valkey
    expose:
      - 6379
  metadata-relay:
    image: ghcr.io/maxdorninger/mediamanager/metadata_relay:latest
    restart: always
    environment:
      - CACHE_HOST=redis
      - TMDB_API_KEY=${TMDB_API_KEY} # you need not provide a TMDB API key, if you only want to use TVDB metadata, or the other way around
      - TVDB_API_KEY=${TVDB_API_KEY}
    ports:
      - "8000:8000"   # Main API server
      - "9090:9090"   # Metrics server
    depends_on:
      - valkey
```


## Environment Variables

| Name                      | Default                                    | Required | Description |
|---------------------------|--------------------------------------------|----------|-------------|
| LOG_LEVEL                 | info                                       | No       | Log verbosity (debug, info, warn, error) |
| PORT                      | 8000                                       | No       | Service port |
| METRICS_PORT              | 9090                                       | No       | Metrics server port |
| VALKEY_HOST               | localhost                                  | No       | Cache DB host |
| VALKEY_PORT               | 6379                                       | No       | Cache DB port |
| VALKEY_DB                 | 0                                          | No       | Cache DB number |
| TMDB_API_KEY              | unset                                      | No       | TMDB API key |
| TMDB_BASE_URL             | https://api.themoviedb.org/3               | No       | TMDB base URL |
| TVDB_API_KEY              | unset                                      | No       | TVDB API key |
| TVDB_BASE_URL             | https://api4.thetvdb.com/v4                | No       | TVDB base URL |
| SEADX_BASE_URL            | https://releases.moe/api                   | No       | SeaDx anime service index base URL |
| JIKAN_BASE_URL            | https://api.jikan.moe/v4                   | No       | Jikan base URL |
| THEAUDIODB_API_KEY        | unset                                      | No       | TheAudioDB API key (use "2" for testing) |
| THEAUDIODB_BASE_URL       | https://www.theaudiodb.com/api/v1/json     | No       | TheAudioDB base URL |
| MEDIA_DIR                 | ./media                                    | No       | On-disk directory to store images and lyrics; served at /media |
| SPOTIFY_CLIENT_ID         | unset                                      | No       | Spotify Client ID (for fetching artist images) |
| SPOTIFY_CLIENT_SECRET     | unset                                      | No       | Spotify Client Secret (for fetching artist images) |
| LRCLIB_BASE_URL           | https://lrclib.net/api                     | No       | LRCLib base URL (for fetching synced lyrics) |
| MUSICBRAINZ_DB_HOST       | unset                                      | No¹      | MusicBrainz PostgreSQL host |
| MUSICBRAINZ_DB_PORT       | 5432                                       | No       | MusicBrainz PostgreSQL port |
| MUSICBRAINZ_DB_USER       | musicbrainz                                | No       | MusicBrainz database username |
| MUSICBRAINZ_DB_PASSWORD   | musicbrainz                                | No       | MusicBrainz database password |
| MUSICBRAINZ_DB_NAME       | unset                                      | No¹      | MusicBrainz database name |
| TYPESENSE_HOST            | localhost                                  | No       | Typesense server host |
| TYPESENSE_PORT            | 8108                                       | No       | Typesense server port |
| TYPESENSE_API_KEY         | unset                                      | No²      | Typesense API key |
| TYPESENSE_TIMEOUT         | 60s                                        | No       | Typesense HTTP client timeout |
| SYNC_ENABLED              | true                                       | No       | Enable background sync scheduler |
| SYNC_INTERVAL             | 24h                                        | No       | How often the scheduler runs |
| SYNC_ENTITIES             | artists,release-groups,releases,recordings | No       | Entities to sync for `sync all` and scheduler (artists, release-groups, releases, recordings) |
| SYNC_SKIP_UNCHANGED       | true                                       | No       | Skip sending unchanged docs to Typesense (uses cache fingerprints) |
| SYNC_DB_PAGE_SIZE         | 8000                                       | No       | Rows fetched per DB page per shard |
| SYNC_SHARD_PARALLELISM    | (CPU)                                      | No       | Parallel DB reader shards per entity |
| SYNC_IMPORT_BATCH_SIZE    | 2000                                       | No       | Documents per Typesense import call |
| SYNC_IMPORT_WORKERS       | (CPU)                                      | No       | Concurrent Typesense import workers per entity |
| SYNC_IMPORT_MAX_RETRIES   | 3                                          | No       | Retries per failed import chunk |
| SYNC_IMPORT_BACKOFF       | 400ms                                      | No       | Initial backoff for retries |
| SYNC_IMPORT_GLOBAL_LIMIT  | unset                                      | No       | Global cap on concurrent import requests across entities |

**Notes:**
1. Both `MUSICBRAINZ_DB_HOST` and `MUSICBRAINZ_DB_NAME` are required to enable MusicBrainz endpoints. If either is missing, MusicBrainz endpoints will not be available (404).
2. `TYPESENSE_API_KEY` is required only if you want to enable search functionality. Without it, MusicBrainz endpoints return 503 for search operations.
3. `THEAUDIODB_API_KEY` can be set to `2` for public testing (rate-limited by TheAudioDB). Leaving it unset disables TheAudioDB endpoints. Note that image fetching via MBID *requires* a premium API key.
4. Jikan is enabled by default and requires no API key. 

## API Endpoints

### MusicBrainz Endpoints

#### Artists
- `GET /musicbrainz/artists/search?query={query}&limit={limit}` - Search for artists using typesense
- `GET /musicbrainz/artists/search/advanced?artist={name}&area={area}&begin={date}&end={date}&limit={limit}` - Advanced artist search with field-specific queries
- `GET /musicbrainz/artists/{mbid}` - Get artist by MBID
- `GET /musicbrainz/artists/{mbid}/release-groups?limit={limit}` - Browse artist's release groups

#### Release Groups (Albums)
- `GET /musicbrainz/release-groups/search?query={query}&limit={limit}` - Search for release groups (albums) using full-text search
- `GET /musicbrainz/release-groups/{mbid}` - Get release group by MBID
- `GET /musicbrainz/release-groups/{mbid}/releases?limit={limit}` - Browse releases in a release group

#### Releases
- `GET /musicbrainz/releases/search?query={query}&limit={limit}` - Search for releases using full-text search
- `GET /musicbrainz/releases/{mbid}` - Get release by MBID

#### Recordings (Tracks)
- `GET /musicbrainz/recordings/search?query={query}&limit={limit}` - Search for recordings (tracks) using full-text search
- `GET /musicbrainz/recordings/{mbid}` - Get recording by MBID

### TMDB Endpoints

#### TV Shows
- `GET /tmdb/tv/trending` - Get trending TV shows
- `GET /tmdb/tv/search?query={query}&page={page}` - Search for TV shows
- `GET /tmdb/tv/shows/{showId}` - Get TV show details by ID
- `GET /tmdb/tv/shows/{showId}/{seasonNumber}` - Get specific season of a TV show

#### Movies
- `GET /tmdb/movies/trending` - Get trending movies
- `GET /tmdb/movies/search?query={query}&page={page}` - Search for movies
- `GET /tmdb/movies/{movieId}` - Get movie details by ID

### TVDB Endpoints

#### TV Shows
- `GET /tvdb/tv/trending` - Get trending TV shows (all series)
- `GET /tvdb/tv/search?query={query}` - Search for TV shows
- `GET /tvdb/tv/shows/{showId}` - Get TV show details by ID
- `GET /tvdb/tv/seasons/{seasonId}` - Get season details by ID

#### Movies
- `GET /tvdb/movies/trending` - Get trending movies (all movies)
- `GET /tvdb/movies/search?query={query}` - Search for movies
- `GET /tvdb/movies/{movieId}` - Get movie details by ID

### SeaDx Endpoints

#### Anime
- `GET /seadx/search?query={query}&page={page}&perPage={perPage}` - Search anime entries
- `GET /seadx/entries/{id}` - Get anime entry details by ID
- `GET /seadx/anilist/{anilistId}` - Get anime entry by AniList ID
- `GET /seadx/trending?limit={limit}` - Get trending anime entries
- `GET /seadx/release-groups?group={group}` - Filter by release group/fansub
- `GET /seadx/trackers?tracker={tracker}` - Filter by tracker

### Jikan Endpoints

#### Anime Database
- `GET /jikan/anime/{id}` - Get anime details by MAL ID
- `GET /jikan/top` - Get top-rated anime
- `GET /jikan/seasonal` - Get current season anime
- `GET /jikan/search?q={query}&page={page}` - Search for anime by title
- `GET /jikan/anime/{id}/recommendations` - Get anime recommendations for specific anime
- `GET /jikan/random` - Get random anime recommendation

### TheAudioDB Endpoints

- `GET /theaudiodb/artist?name={name}` - Lookup artist info by name (uses TheAudioDB; requires THEAUDIODB_API_KEY; public test key is "2")
- `GET /theaudiodb/artist/{mbid}` - Lookup artist info by MusicBrainz ID (prioritized caching: 7-day TTL)
- `GET /theaudiodb/album/{mbid}` - Lookup album info by MusicBrainz release group ID (prioritized caching: 7-day TTL)
- `GET /theaudiodb/track/{mbid}` - Lookup track info by MusicBrainz recording ID (prioritized caching: 7-day TTL)

### System Endpoints
- `GET /` - Service status and information
- `GET /metrics` - Prometheus metrics (served on separate metrics port for security)
- `GET /media/*` - Static file serving for images and lyrics

**Note**: The `/metrics` endpoint is served on the dedicated metrics port (default 9090) to isolate monitoring traffic from the main API.