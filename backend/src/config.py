from pathlib import Path

from pydantic_settings import BaseSettings


class BasicConfig(BaseSettings):
    storage_directory: Path = "./data"
    tv_directory: Path = "./tv"
    movie_directory: Path = "./movie"
    torrent_directory: Path = "./torrent"
    DEVELOPMENT: bool = False
