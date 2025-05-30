from pathlib import Path

from pydantic import AnyHttpUrl
from pydantic_settings import BaseSettings


class BasicConfig(BaseSettings):
    image_directory: Path = "./data"
    tv_directory: Path = "./tv"
    movie_directory: Path = "./movie"
    torrent_directory: Path = "./torrent"
    FRONTEND_URL: AnyHttpUrl
    CORS_URLS: str = ""
    DEVELOPMENT: bool = False
