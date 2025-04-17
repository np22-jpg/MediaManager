import typing
import uuid
from uuid import UUID

from pydantic import BaseModel, Field, ConfigDict

from torrent.models import Quality

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
    year: int

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
    torrent_id: UUID
    file_path: str
