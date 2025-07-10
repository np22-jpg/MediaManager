import os
from pathlib import Path

from pydantic import AnyHttpUrl
from pydantic_settings import (
    BaseSettings,
    SettingsConfigDict,
)

from media_manager.auth.config import AuthConfig
from media_manager.database.config import DbConfig
from media_manager.indexer.config import IndexerConfig
from media_manager.metadataProvider.config import MetadataProviderConfig
from media_manager.notification.config import NotificationConfig
from media_manager.torrent.config import TorrentConfig


class BasicConfig(BaseSettings):
    image_directory: Path = Path(__file__).parent.parent / "data" / "images"
    tv_directory: Path = Path(__file__).parent.parent / "data" / "tv"
    movie_directory: Path = Path(__file__).parent.parent / "data" / "movies"
    torrent_directory: Path = Path(__file__).parent.parent / "data" / "torrents"

    FRONTEND_URL: AnyHttpUrl = "http://localhost:3000/"
    CORS_URLS: list[str] = []
    DEVELOPMENT: bool = False
    api_base_path: str = "/api/v1"


class AllEncompassingConfig(BaseSettings):
    model_config = SettingsConfigDict(
        toml_file=os.getenv("CONFIG_FILE", "./config.toml")
    )
    """
    This class is used to load all configurations from the environment variables.
    It combines the BasicConfig with any additional configurations needed.
    """
    misc: BasicConfig = BasicConfig()
    torrents: TorrentConfig = TorrentConfig()
    notifications: NotificationConfig = NotificationConfig()
    metadata: MetadataProviderConfig = MetadataProviderConfig()
    indexers: IndexerConfig = IndexerConfig()
    database: DbConfig = DbConfig()
    auth: AuthConfig = AuthConfig()
