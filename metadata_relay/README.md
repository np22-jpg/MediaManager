# MetadataRelay

This is a service that provides metadata for movies, TV shows, and music. It caches the metadata to not overload the TMDB, TVDB, and MusicBrainz APIs. This service is used by MediaManager to fetch metadata for movies, TV shows, and music. I (the developer) run a public instance of this service at https://metadata-relay.maxid.me, but you can also run your own instance.

## Supported APIs

- **TMDB** - Movies and TV Shows metadata
- **TVDB** - Movies and TV Shows metadata (alternative source)
- **MusicBrainz** - Music metadata (artists, albums, releases, recordings)

## Example Docker Compose Configuration

````yaml
services:
  valkey:
    image: valkey/valkey:latest
    restart: always
    container_name: valkey
    expose:
      - 6379
  metadata_relay:
    image: ghcr.io/maxdorninger/mediamanager/metadata_relay:latest
    restart: always
    environment:
      - TMDB_API_KEY=  # you need not provide a TMDB API key, if you only want to use TVDB metadata, or the other way around
      - TVDB_API_KEY=
      - VALKEY_HOST=valkey
    container_name: metadata_relay
    ports:
      - 8000:8000
````

## Environment Variables

| Name         | Default Value  | Description         |
| ------------ | -------------- | ------------------- | 
| VALKEY_HOST  | localhost      | Address/URL to DB   |
| VALKEY_PORT  | 6379           | Port to DB          |
| VALKEY_DB    | 0              | DB Name             |
| TMDB_API_KEY | *unset*        | API Key for TMDB    |
| TVDB_API_KEY | *unset*        | API Key for TVDB    |
| LOG_LEVEL    | info           | Log Verbosity       |
| PORT         | 8000           | Service port        |

## metadata_relay-specific Roadmap

- [x] port metadata relay to go
- [x] enable musicbrainz support in metadata relay
- [ ] enable AniDB support in metadata relay
  - [ ] enable SeaDex support in metadata relay
- [ ] add support for new metadata sources in backend
  - [ ] _maybe_ enable automated PGO/BOLT building
- [ ] expand E2E metadata relay testing
- [ ] expand Linting and formatting in metadata relay
- [ ] Update metadata relay docs
- [ ] Create a grafana dashboard for metadata relay

## API Endpoints

### MusicBrainz Endpoints

#### Artists
- `GET /musicbrainz/artists/search?query={query}&limit={limit}` - Search for artists
- `GET /musicbrainz/artists/search/advanced?artist={name}&area={area}&begin={date}&end={date}&limit={limit}` - Advanced artist search
- `GET /musicbrainz/artists/{mbid}` - Get artist by MBID
- `GET /musicbrainz/artists/{mbid}/release-groups?limit={limit}` - Browse artist's release groups

#### Release Groups (Albums)
- `GET /musicbrainz/release-groups/search?query={query}&limit={limit}` - Search for release groups
- `GET /musicbrainz/release-groups/{mbid}` - Get release group by MBID
- `GET /musicbrainz/release-groups/{mbid}/releases?limit={limit}` - Browse releases in a release group

#### Releases
- `GET /musicbrainz/releases/search?query={query}&limit={limit}` - Search for releases
- `GET /musicbrainz/releases/{mbid}` - Get release by MBID

#### Recordings (Tracks)
- `GET /musicbrainz/recordings/search?query={query}&limit={limit}` - Search for recordings
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

### Notes

- **TMDB**: Supports pagination with `page` parameter for search endpoints
- **TVDB**: No pagination parameters needed for search endpoints
- **MusicBrainz**: Supports `limit` parameter (1-100) for controlling result count
- All endpoints return JSON responses