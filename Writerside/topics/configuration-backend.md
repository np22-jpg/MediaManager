# Backend

These variables configure the core backend application, database connections, authentication, and integrations.

<note>
    <include from="notes.topic" element-id="list-format"/>
</note>

## General Settings

| Variable        | Description                                                                            | Default   |
|-----------------|----------------------------------------------------------------------------------------|-----------|
| `API_BASE_PATH` | The url base of the backend                                                            | `/api/v1` |
| `CORS_URLS`     | Enter a list of origins you are going to access the api from (e.g. https://mm.example) | -         |

## Database Settings

| Variable      | Description                              | Default        | Example      |
|---------------|------------------------------------------|----------------|--------------|
| `DB_HOST`     | Hostname or IP of the PostgreSQL server. | `localhost`    | `postgres`   |
| `DB_PORT`     | Port number of the PostgreSQL server.    | `5432`         | `5432`       |
| `DB_USER`     | Username for PostgreSQL connection.      | `MediaManager` | `myuser`     |
| `DB_PASSWORD` | Password for the PostgreSQL user.        | `MediaManager` | `mypassword` |
| `DB_DBNAME`   | Name of the PostgreSQL database.         | `MediaManager` | `mydatabase` |

## Download Client Settings

Currently, only qBittorrent is supported as a download client. But support for other clients isn't unlikely in the
future.

| Variable           | Description                 | Default     | Example            |
|--------------------|-----------------------------|-------------|--------------------|
| `QBITTORRENT_HOST` | Host of the QBittorrent API | `localhost` | `qbit.example.com` |
| `QBITTORRENT_PORT` | Port of the QBittorrent API | `8080`      | `443`              |
| `QBITTORRENT_USER` | Username for QBittorrent    | `admin`     | -                  |
| `QBITTORRENT_PASS` | Password for QBittorrent    | `admin`     | -                  |

## Metadata Provider Settings

These settings configure the integrations with external metadata providers like The Movie Database (TMDB) and The TVDB.

### TMDB (The Movie Database)

TMDB is the primary metadata provider for MediaManager. It provides detailed information about movies and TV shows.
Get an API key from [The Movie Database](https://www.themoviedb.org/settings/api) to use this provider. You can create
an account and generate a free API key in your account settings.

<tip>
    Other software like Jellyfin use TMDB as well, so there won't be any metadata discrepancies.
</tip>

| Variable       | Default | Example                               |
|----------------|---------|---------------------------------------|
| `TMDB_API_KEY` | None    | `TMDB_API_KEY=your_tmdb_api_key_here` |

### TVDB (The TVDB)

<warning>
    The TVDB might provide false metadata, also it doesn't support some features of MediaManager like to show overviews, therfore TMDB is the preferred metadata provider. 
</warning>

Get an API key from [The TVDB](https://thetvdb.com/auth/register) to use this provider. You can create an account and
generate a free API key in your account settings.

| Variable       | Default | Example                               |
|----------------|---------|---------------------------------------|
| `TVDB_API_KEY` | None    | `TVDB_API_KEY=your_tvdb_api_key_here` |

## Directory Settings

<note>
    Normally you don't need to change these, as the default mountpoints are usually sufficient. In your <code>docker-compose.yaml</code>, you can just mount <code>/any/directory</code> to <code>/data/torrents</code>.
</note>

| Variable            | Description                                       | Default          |
|---------------------|---------------------------------------------------|------------------|
| `IMAGE_DIRECTORY`   | media images (posters, backdrops) will be stored. | `/data/images`   |
| `TV_DIRECTORY`      | location of TV show files                         | `/data/tv`       |
| `MOVIE_DIRECTORY`   | location of movie files                           | `/data/movies`   |
| `TORRENT_DIRECTORY` | location of torrent files and downloads           | `/data/torrents` |

## Build Arguments (Dockerfile)

| Argument  | Description                                                                                                                    | Example (in build command)                 |
|-----------|--------------------------------------------------------------------------------------------------------------------------------|--------------------------------------------|
| `VERSION` | Labels the Docker image with a version. Passed during build (e.g., by GitHub Actions). Frontend uses this as `PUBLIC_VERSION`. | `docker build --build-arg VERSION=1.2.3 .` |

