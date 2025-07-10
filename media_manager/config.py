import logging
import os
from pathlib import Path
from typing import Type, Tuple

from pydantic import AnyHttpUrl
from pydantic_settings import (
    BaseSettings,
    SettingsConfigDict,
    PydanticBaseSettingsSource,
    TomlConfigSettingsSource,
)

from media_manager.auth.config import AuthConfig
from media_manager.database.config import DbConfig
from media_manager.indexer.config import IndexerConfig
from media_manager.metadataProvider.config import MetadataProviderConfig
from media_manager.notification.config import NotificationConfig
from media_manager.torrent.config import TorrentConfig

log = logging.getLogger(__name__)
config_path = os.getenv("CONFIG_FILE")

if config_path is None:
    log.info("No CONFIG_FILE environment variable set, using default config file path.")
    config_path = Path(__file__).parent.parent / "data" / "config.toml"
else:
    config_path = Path(config_path)
print("SERVAS CONFIG PATH: ", config_path)
log.info("Using config file path: %s", config_path)


class BasicConfig(BaseSettings):
    image_directory: Path = Path(__file__).parent.parent / "data" / "images"
    tv_directory: Path = Path(__file__).parent.parent / "data" / "tv"
    movie_directory: Path = Path(__file__).parent.parent / "data" / "movies"
    torrent_directory: Path = Path(__file__).parent.parent / "data" / "torrents"

    frontend_url: AnyHttpUrl = "http://localhost:3000/"
    cors_urls: list[str] = []
    development: bool = False
    api_base_path: str = "/api/v1"


class AllEncompassingConfig(BaseSettings):
    model_config = SettingsConfigDict(
        toml_file=config_path,
        case_sensitive=False,
    )
    """
    This class is used to load all configurations from the environment variables.
    It combines the BasicConfig with any additional configurations needed.
    """
    misc: BasicConfig
    torrents: TorrentConfig
    notifications: NotificationConfig
    metadata: MetadataProviderConfig
    indexers: IndexerConfig
    database: DbConfig
    auth: AuthConfig

    @classmethod
    def settings_customise_sources(
        cls,
        settings_cls: Type[BaseSettings],
        init_settings: PydanticBaseSettingsSource,
        env_settings: PydanticBaseSettingsSource,
        dotenv_settings: PydanticBaseSettingsSource,
        file_secret_settings: PydanticBaseSettingsSource,
    ) -> Tuple[PydanticBaseSettingsSource, ...]:
        return (TomlConfigSettingsSource(settings_cls),)
