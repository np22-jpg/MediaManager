import uuid
from typing import Literal
from uuid import UUID

from sqlalchemy import Column, ForeignKeyConstraint, String, UniqueConstraint
from sqlmodel import Field, Relationship, SQLModel


class TorrentMixin:
    torrent_status: Literal["downloading", "finished", "error"] | None = Field(default=None,
                                                                               sa_column=Column(String))
    torrent_filepath: str | None = Field(default=None)
    requested: bool = Field(default=False)
    id: UUID


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


class Season(SQLModel, TorrentMixin, table=True):
    __table_args__ = (UniqueConstraint("show_id", "number"),)
    id: UUID = Field(primary_key=True, default_factory=uuid.uuid4)

    show_id: UUID = Field(foreign_key="show.id", ondelete="CASCADE")
    number: int = Field()

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
