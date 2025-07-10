from pydantic_settings import BaseSettings

from media_manager.metadataProvider.tmdb import TmdbConfig
from media_manager.metadataProvider.tvdb import TvdbConfig


class MetadataProviderConfig(BaseSettings):
    tvdb: TvdbConfig = TvdbConfig()
    tmdb: TmdbConfig = TmdbConfig()