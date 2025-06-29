# Backend

These variables configure the core backend application, database connections, authentication, and integrations.

<note>
    <include from="notes.topic" element-id="list-format"/>
</note>

## General Settings

### `API_BASE_PATH`

The url base of the backend. Default is `/api/v1`.

### `CORS_URLS`

Enter a list of origins you are going to access the api from. Example: `["https://mm.example"]`.

## Database Settings

### `DB_HOST`

Hostname or IP of the PostgreSQL server. Default is `localhost`.

### `DB_PORT`

Port number of the PostgreSQL server. Default is `5432`.

### `DB_USER`

Username for PostgreSQL connection. Default is `MediaManager`.

### `DB_PASSWORD`

Password for the PostgreSQL user. Default is `MediaManager`.

### `DB_DBNAME`

Name of the PostgreSQL database. Default is `MediaManager`.

## Download Client Settings

Currently, only qBittorrent is supported as a download client. But support for other clients isn't unlikely in the
future.

### `QBITTORRENT_HOST`

Host of the QBittorrent API. Default is `localhost`. Example: `qbit.example.com`.

### `QBITTORRENT_PORT`

Port of the QBittorrent API. Default is `8080`. Example: `443`.

### `QBITTORRENT_USER`

Username for QBittorrent. Default is `admin`.

### `QBITTORRENT_PASSWORD`

Password for QBittorrent. Default is `admin`.

## Metadata Provider Settings


<note>
   Note the lack of a trailing slash in some env vars like <code>TMDB_RELAY_URL</code>. This is important.
</note>



These settings configure the integrations with external metadata providers like The Movie Database (TMDB) and The TVDB.

### TMDB (The Movie Database)

TMDB is the primary metadata provider for MediaManager. It provides detailed information about movies and TV shows.


<tip>
    Other software like Jellyfin use TMDB as well, so there won't be any metadata discrepancies.
</tip>

#### `TMDB_RELAY_URL`

If you want use your own TMDB relay service, set this to the URL of your own MetadataRelay. Otherwise, don't set it to
use the default relay.

Default: `https://metadata-relay.maxid.me/tmdb`.

### TVDB (The TVDB)

<warning>
    The TVDB might provide false metadata, also it doesn't support some features of MediaManager like to show overviews, therfore TMDB is the preferred metadata provider.
</warning>

#### `TVDB_RELAY_URL`

If you want use your own TVDB relay service, set this to the URL of your own MetadataRelay. Otherwise, don't set it to
use the default relay.

Default: `https://metadata-relay.maxid.me/tvdb`.

### MetadataRelay

<note>
  To use MediaManager <strong>you don't need to set up your own MetadataRelay</strong>, as the default relay which is hosted by me, the dev of MediaManager, should be sufficient for most purposes.
</note>

The MetadataRelay is a service that provides metadata for MediaManager. It acts as a proxy for TMDB and TVDB, allowing
you to use your own API keys, but not strictly needing your own because only me, the developer, needs to create accounts
for API keys.
You might want to use it if you want to avoid rate limits, to protect your privacy, or other reasons.
If you know Sonarr's Skyhook, this is similar to that.

#### Where to get API keys

Get an API key from [The Movie Database](https://www.themoviedb.org/settings/api). You can create
an account and generate a free API key in your account settings.

Get an API key from [The TVDB](https://thetvdb.com/auth/register). You can create an account and
generate a free API key in your account settings.

<tip>
    If you want to use your own MetadataRelay, you can set the  <code>TMDB_RELAY_URL</code> and/or  <code>TVDB_RELAY_URL</code> to your own relay service.
</tip>

## Directory Settings

<note>
    Normally you don't need to change these, as the default mountpoints are usually sufficient. In your <code>docker-compose.yaml</code>, you can just mount <code>/any/directory</code> to <code>/data/torrents</code>.
</note>

### `IMAGE_DIRECTORY`

Media images (posters, backdrops) will be stored here. Default is `/data/images`.

### `TV_DIRECTORY`

Location of TV show files. Default is `/data/tv`.

### `MOVIE_DIRECTORY`

Location of movie files. Default is `/data/movies`.

### `TORRENT_DIRECTORY`

Location of torrent files and downloads. Default is `/data/torrents`.

## Build Arguments (Dockerfile)

### `VERSION`

Labels the Docker image with a version. Passed during build (e.g., by GitHub Actions). Frontend uses this as
`PUBLIC_VERSION`. Example (in build command): `docker build --build-arg VERSION=1.2.3 .`
