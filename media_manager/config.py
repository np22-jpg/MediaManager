from pathlib import Path

from pydantic import AnyHttpUrl
from pydantic_settings import BaseSettings


class BasicConfig(BaseSettings):
    image_directory: Path = "/data/images"
    tv_directory: Path = "/data/tv"
    movie_directory: Path = "/data/movies"
    torrent_directory: Path = "/data/torrents"
    FRONTEND_URL: AnyHttpUrl = "http://localhost:3000"
    CORS_URLS: list[str] = []
    DEVELOPMENT: bool = False
