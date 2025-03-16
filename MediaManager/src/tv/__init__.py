import logging
import uuid
from uuid import UUID

import tmdbsimple as tmdb
from sqlalchemy import UniqueConstraint, ForeignKeyConstraint
from sqlmodel import Field, SQLModel

from config import TvConfig

class Show(SQLModel, table=True):
    __table_args__ = (UniqueConstraint("external_id", "metadata_provider"),)
    id: UUID = Field(primary_key=True, default_factory=uuid.uuid4)
    external_id: int
    metadata_provider: str
    name: str
    overview: str

class Season(SQLModel, table=True):
    show_id: UUID = Field(foreign_key="show.id", primary_key=True, default_factory=uuid.uuid4)
    number: int = Field(primary_key=True)
    requested: bool = Field(default=False)
    external_id: int
    name: str
    overview: str

class Episode(SQLModel, table=True):
    __table_args__ = (
        ForeignKeyConstraint(['show_id', 'season_number'], ['season.show_id', 'season.number']),
    )
    show_id: UUID = Field(primary_key=True)
    season_number: int = Field( primary_key=True)
    number: int = Field(primary_key=True)
    external_id: int
    title: str

config = TvConfig()
log = logging.getLogger(__name__)

tmdb.API_KEY = config.api_key
