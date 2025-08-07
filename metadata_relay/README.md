# MetadataRelay

This is a service that provides metadata for movies and TV shows. It caches the metadata to not overload the TMDB and
TVDB APIs. This service is used by MediaManager to fetch metadata for movies and TV shows. I (the developer) run a
public instance of this service at https://metadata-relay.maxid.me, but you can also run your
own instance.

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