# Backend

These variables configure the core backend application, database connections, authentication, and integrations.

<note>
    <include from="notes.topic" element-id="list-format"/>
</note>

## General Settings

### `API_BASE_PATH`

The url base of the backend. Default is `/api/v1`.

### `CORS_URLS`

Enter a list of origins you are going to access the api from. Example: `https://mm.example`.

## Database Settings

### `DB_HOST`

Hostname or IP of the PostgreSQL server. Default is `localhost`. Example: `postgres`.

### `DB_PORT`

Port number of the PostgreSQL server. Default is `5432`. Example: `5432`.

### `DB_USER`

Username for PostgreSQL connection. Default is `MediaManager`. Example: `myuser`.

### `DB_PASSWORD`

Password for the PostgreSQL user. Default is `MediaManager`. Example: `mypassword`.

### `DB_DBNAME`

Name of the PostgreSQL database. Default is `MediaManager`. Example: `mydatabase`.

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

These settings configure the integrations with external metadata providers like The Movie Database (TMDB) and The TVDB.

### TMDB (The Movie Database)

TMDB is the primary metadata provider for MediaManager. It provides detailed information about movies and TV shows.
Get an API key from [The Movie Database](https://www.themoviedb.org/settings/api) to use this provider. You can create
an account and generate a free API key in your account settings.

<tip>
    Other software like Jellyfin use TMDB as well, so there won't be any metadata discrepancies.
</tip>

#### `TMDB_API_KEY`

Your TMDB API key. Example: `your_tmdb_api_key_here`.

### TVDB (The TVDB)

<warning>
    The TVDB might provide false metadata, also it doesn't support some features of MediaManager like to show overviews, therfore TMDB is the preferred metadata provider.
</warning>

Get an API key from [The TVDB](https://thetvdb.com/auth/register) to use this provider. You can create an account and
generate a free API key in your account settings.

#### `TVDB_API_KEY`

Your TVDB API key. Example: `your_tvdb_api_key_here`.

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
