import typing
import uuid
from uuid import UUID

from pydantic import BaseModel, Field, ConfigDict

from torrent.models import Quality
from torrent.schemas import TorrentId, TorrentStatus

ShowId = typing.NewType("ShowId", UUID)
SeasonId = typing.NewType("SeasonId", UUID)
EpisodeId = typing.NewType("EpisodeId", UUID)

SeasonNumber = typing.NewType("SeasonNumber", int)
EpisodeNumber = typing.NewType("EpisodeNumber", int)


class Episode(BaseModel):
    model_config = ConfigDict(from_attributes=True)

    id: EpisodeId = Field(default_factory=uuid.uuid4)
    number: EpisodeNumber
    external_id: int
    title: str


class Season(BaseModel):
    model_config = ConfigDict(from_attributes=True)

    id: SeasonId = Field(default_factory=uuid.uuid4)
    number: SeasonNumber

    name: str
    overview: str

    external_id: int

    episodes: list[Episode]


class Show(BaseModel):
    model_config = ConfigDict(from_attributes=True)

    id: ShowId = Field(default_factory=uuid.uuid4)

    name: str
    overview: str
    year: int | None

    external_id: int
    metadata_provider: str

    seasons: list[Season]


class SeasonRequest(BaseModel):
    model_config = ConfigDict(from_attributes=True)

    season_id: SeasonId
    min_quality: Quality
    wanted_quality: Quality


class SeasonFile(BaseModel):
    model_config = ConfigDict(from_attributes=True)

    season_id: SeasonId
    quality: Quality
    torrent_id: TorrentId | None
    file_path_suffix: str

class RichSeasonTorrent(BaseModel):
    model_config = ConfigDict(from_attributes=True)

    torrent_id: TorrentId
    torrent_title: str
    status: TorrentStatus
    quality: Quality
    imported: bool

    file_path_suffix: str
    seasons: list[SeasonNumber]

class RichShowTorrent(BaseModel):
    model_config = ConfigDict(from_attributes=True)

    show_id: ShowId
    name: str
    year: int | None
    metadata_provider: str
    torrents: list[RichSeasonTorrent]


class PublicSeason(BaseModel):
    model_config = ConfigDict(from_attributes=True)

    id: SeasonId
    number: SeasonNumber

    downloaded: bool = False
    name: str
    overview: str

    external_id: int

    episodes: list[Episode]


class PublicShow(BaseModel):
    model_config = ConfigDict(from_attributes=True)

    id: ShowId

    name: str
    overview: str
    year: int | None

    external_id: int
    metadata_provider: str

    seasons: list[PublicSeason]
