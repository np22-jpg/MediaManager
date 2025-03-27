import os
from typing import Literal

from pydantic import BaseModel


class TmdbConfig(BaseModel):
    api_key: str = os.getenv("TMDB_API_KEY") or None


class BasicConfig(BaseModel):
    storage_directory: str = os.getenv("STORAGE_FILE_PATH") or "."

class ProwlarrConfig(BaseModel):
    enabled: bool = bool(os.getenv("PROWLARR_ENABLED") or True)
    api_key: str = os.getenv("PROWLARR_API_KEY")
    url: str = os.getenv("PROWLARR_URL")


class AuthConfig(BaseModel):
    # to get a signing key run:
    # openssl rand -hex 32
    _jwt_signing_key: str = os.getenv("JWT_SIGNING_KEY")
    jwt_signing_algorithm: str = "HS256"
    jwt_access_token_lifetime: int = int(os.getenv("JWT_ACCESS_TOKEN_LIFETIME") or 60 * 24 * 30)

    @property
    def jwt_signing_key(self):
        return self._jwt_signing_key


class QbittorrentConfig(BaseModel):
    host: str = os.getenv("QBITTORRENT_HOST") or "localhost"
    port: int = os.getenv("QBITTORRENT_PORT") or 8080
    username: str = os.getenv("QBITTORRENT_USERNAME") or "admin"
    password: str = os.getenv("QBITTORRENT_PASSWORD") or "adminadmin"


class DownloadClientConfig(BaseModel):
    client: Literal['qbit'] = os.getenv("DOWNLOAD_CLIENT") or "qbit"


class MachineLearningConfig(BaseModel):
    model_name: str = os.getenv("OLLAMA_MODEL_NAME") or "qwen2.5:0.5b"


def get_db_config() -> DbConfig:
    return DbConfig()
