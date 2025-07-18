from pydantic_settings import BaseSettings


class TmdbConfig(BaseSettings):
    tmdb_relay_url: str = "https://metadata-relay.maxid.me/tmdb"


class TvdbConfig(BaseSettings):
    tvdb_relay_url: str = "https://metadata-relay.maxid.me/tvdb"


class MetadataProviderConfig(BaseSettings):
    tvdb: TvdbConfig = TvdbConfig()
    tmdb: TmdbConfig = TmdbConfig()
