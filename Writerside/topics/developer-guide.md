# Developer Guide

This section is for those who want to contribute to Media Manager or understand its internals.

### Source Code

- `media_manager/`: Backend FastAPI application
- `web/`: Frontend SvelteKit application
- `Writerside/`: Documentation
- `metadata_relay/`: Metadata relay service
-

### Backend Development

- Uses `uv` for dependency management
- Follows standard FastAPI project structure
- Database migrations are handled by Alembic

### Frontend Development

- Uses `npm` for package management
- SvelteKit with TypeScript

### Contributing

- Consider opening an issue to discuss significant changes before starting work

## Sequence Diagrams

```mermaid
sequenceDiagram
    title Step-by-step: going from adding a show to importing a torrent of one of its seasons
    
    User->>TV Router: Add a show (POST /tv/shows)
    TV Router->>TV Service: Receive Show Request
    TV Service->>MetadataProviderService: Get Metadata for Show
    MetadataProviderService->>File System: Save Poster Image
    TV Service->>Database: Store show information

    User->>TV Router: Get Available Torrents for a Season (GET /tv/torrents)
    TV Router->>TV Service: Receive Request
    TV Service->>Indexer Service: Search for torrents
    TV Service->>User: Returns Public Indexer Results

    User->>TV Router: Download Torrent (POST /tv/torrents)
    TV Router->>TV Service: Receive Request
    Note over Database: This associates a season with a torrent id and the file_path_suffix
    TV Service->>Database: Saves a SeasonFile object
    TV Service->>Torrent Service: Download Torrent
    Torrent Service->>File System: Save Torrentfile
    Torrent Service->>QBittorrent: Download Torrent

    Note over Scheduler: Hourly scheduler trigger
    Scheduler->>TV Service: auto_import_all_show_torrents()
    TV Service->>Database: Get all Shows and seasons which are associated with a torrent
    TV Service->>Torrent Service: Update Torrent download statuses
    Note over TV Service: if a torrent is finished downloading it will be imported
    TV Service->>Torrent Service: get all files in the torrents directory
    Note over Torrent Service: Extracts archives, guesses mimetype (Video/Subtitle/Other)
    Note over TV Service: filters files based on some regex and renames them
    TV Service->>File System: Move/Hardlink video and subtitle files

    Note over User: User can now access the show in e.g. Jellyfin


```

## Tech Stack

### Backend

- Python with FastAPI
- SQLAlchemy
- Pydantic and Pydantic-Settings

### Frontend

- TypeScript with SvelteKit
- Tailwind CSS
- shadcn-svelte

### CI/CD

- GitHub Actions
