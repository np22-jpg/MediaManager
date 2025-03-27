import os
from typing import Literal

from pydantic import BaseModel


class BasicConfig(BaseModel):
    storage_directory: str = os.getenv("STORAGE_FILE_PATH") or "."

class ProwlarrConfig(BaseModel):
    enabled: bool = bool(os.getenv("PROWLARR_ENABLED") or True)
    api_key: str = os.getenv("PROWLARR_API_KEY")
    url: str = os.getenv("PROWLARR_URL")





class QbittorrentConfig(BaseModel):
    host: str = os.getenv("QBITTORRENT_HOST") or "localhost"
    port: int = os.getenv("QBITTORRENT_PORT") or 8080
    username: str = os.getenv("QBITTORRENT_USERNAME") or "admin"
    password: str = os.getenv("QBITTORRENT_PASSWORD") or "adminadmin"


class DownloadClientConfig(BaseModel):
    client: Literal['qbit'] = os.getenv("DOWNLOAD_CLIENT") or "qbit"


class MachineLearningConfig(BaseModel):
    model_name: str = os.getenv("OLLAMA_MODEL_NAME") or "qwen2.5:0.5b"