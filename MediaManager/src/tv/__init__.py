import logging
import uuid
from uuid import UUID

import tmdbsimple as tmdb
from sqlalchemy import UniqueConstraint, ForeignKeyConstraint
from sqlmodel import Field, SQLModel, Relationship

from config import TvConfig

class Show(SQLModel, table=True):
    __table_args__ = (UniqueConstraint("external_id", "metadata_provider"),)
    id: UUID = Field(primary_key=True, default_factory=uuid.uuid4)
    external_id: int
    metadata_provider: str
    name: str
    overview: str

    seasons: list["Season"] = Relationship(back_populates="show", cascade_delete=True)

class Season(SQLModel, table=True):
    show_id: UUID = Field(foreign_key="show.id", primary_key=True, default_factory=uuid.uuid4, ondelete="CASCADE")
    number: int = Field(primary_key=True)
    requested: bool = Field(default=False)
    external_id: int
    name: str
    overview: str

    show: Show = Relationship(back_populates="seasons")
    episodes: list["Episode"] = Relationship(back_populates="season", cascade_delete=True)

class Episode(SQLModel, table=True):
    __table_args__ = (
        ForeignKeyConstraint(['show_id', 'season_number'], ['season.show_id', 'season.number'], ondelete="CASCADE"),
    )
    show_id: UUID = Field(primary_key=True)
    season_number: int = Field( primary_key=True)
    number: int = Field(primary_key=True)
    external_id: int
    title: str

    season: Season = Relationship(back_populates="episodes")

config = TvConfig()
log = logging.getLogger(__name__)

tmdb.API_KEY = config.api_key
