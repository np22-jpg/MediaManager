from pydantic_settings import BaseSettings, SettingsConfigDict


class QbittorrentConfig(BaseSettings):
    model_config = SettingsConfigDict(env_prefix="QBITTORRENT_")
    host: str = "localhost"
    port: int = 8080
    username: str = "admin"
    password: str = "admin"
    enabled: bool = False


class SabnzbdConfig(BaseSettings):
    model_config = SettingsConfigDict(env_prefix="SABNZBD_")
    host: str = "localhost"
    port: int = 8080
    api_key: str = ""
    enabled: bool = False


class TorrentConfig(BaseSettings):
    qbittorrent: QbittorrentConfig
    sabnzbd: SabnzbdConfig
