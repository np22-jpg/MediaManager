import uuid
from abc import ABC
from typing import Literal
from uuid import UUID

from sqlalchemy import ForeignKeyConstraint, UniqueConstraint
from sqlmodel import Field, Relationship, SQLModel


class Torrents(ABC):
    torrent_status: Literal["downloading", "finished", "error"] | None = Field(default="downloading")
    torrent_url: str | None = Field(default=None)


class Show(SQLModel, table=True):
    __table_args__ = (UniqueConstraint("external_id", "metadata_provider"),)
    id: UUID = Field(primary_key=True, default_factory=uuid.uuid4)
    external_id: int
    metadata_provider: str
    name: str
    overview: str
    # For some shows the first_air_date isn't known, therefore it needs to be nullable
    year: int | None

    seasons: list["Season"] = Relationship(back_populates="show", cascade_delete=True)


class Season(SQLModel, Torrents, table=True):
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
    season_number: int = Field(primary_key=True)
    number: int = Field(primary_key=True)
    external_id: int
    title: str

    season: Season = Relationship(back_populates="episodes")
