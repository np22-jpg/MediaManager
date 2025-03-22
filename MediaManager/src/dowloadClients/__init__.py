from abc import ABC
from typing import Literal
from uuid import UUID

from sqlmodel import Field, SQLModel

from config import DownloadClientConfig
from dowloadClients.qbittorrent import QbittorrentClient


class Torrents(ABC):
    status: Literal["downloading", "finished", "error"] = Field(default="ownloading")
    url: str = Field(default=None)
    id: UUID


class SeasonTorrents(SQLModel, Torrents):
    id: UUID = Field(primary_key=True, foreign_key="season.show_id")
    season_number: int = Field(primary_key=True, foreign_key="season.number")


# TODO: implement MovieTorrents
# class MovieTorrents(SQLModel, Torrents):
#    id: UUID = Field(primary_key=True, foreign_key="movie.show_id")

config = DownloadClientConfig()

# TODO: add more elif when implementing more download clients
if config.client == "qbit":
    client = QbittorrentClient()
else:
    client = QbittorrentClient()
