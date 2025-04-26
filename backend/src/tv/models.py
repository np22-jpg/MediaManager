from uuid import UUID

from sqlalchemy import ForeignKey, PrimaryKeyConstraint, UniqueConstraint
from sqlalchemy.orm import Mapped, mapped_column, relationship

from database import Base
from torrent.models import Quality


class Show(Base):
    __tablename__ = "show"
    __table_args__ = (UniqueConstraint("external_id", "metadata_provider"),)

    id: Mapped[UUID] = mapped_column(primary_key=True)
    external_id: Mapped[int]
    metadata_provider: Mapped[str]
    name: Mapped[str]
    overview: Mapped[str]
    year: Mapped[int | None]

    seasons: Mapped[list["Season"]] = relationship(back_populates="show", cascade="all, delete")


class Season(Base):
    __tablename__ = "season"
    __table_args__ = (UniqueConstraint("show_id", "number"),)

    id: Mapped[UUID] = mapped_column(primary_key=True)
    show_id: Mapped[UUID] = mapped_column(ForeignKey(column="show.id", ondelete="CASCADE"), )
    number: Mapped[int]
    external_id: Mapped[int]
    name: Mapped[str]
    overview: Mapped[str]

    show: Mapped["Show"] = relationship(back_populates="seasons")
    episodes: Mapped[list["Episode"]] = relationship(back_populates="season", cascade="all, delete")

    season_files = relationship("SeasonFile", back_populates="season", cascade="all, delete")


class Episode(Base):
    __tablename__ = "episode"
    __table_args__ = (
        UniqueConstraint("season_id", "number"),
    )
    id: Mapped[UUID] = mapped_column(primary_key=True)
    season_id: Mapped[UUID] = mapped_column(ForeignKey("season.id", ondelete="CASCADE"), )
    number: Mapped[int]
    external_id: Mapped[int]
    title: Mapped[str]

    season: Mapped["Season"] = relationship(back_populates="episodes")


class SeasonFile(Base):
    __tablename__ = "season_file"
    __table_args__ = (
        PrimaryKeyConstraint("season_id", "file_path_suffix"),
    )
    season_id: Mapped[UUID] = mapped_column(ForeignKey(column="season.id", ondelete="CASCADE"), )
    torrent_id: Mapped[UUID | None] = mapped_column(ForeignKey(column="torrent.id", ondelete="SET NULL"), )
    file_path_suffix: Mapped[str]
    quality: Mapped[Quality]

    torrent = relationship("Torrent", back_populates="season_files", uselist=False)
    season = relationship("Season", back_populates="season_files", uselist=False)


class SeasonRequest(Base):
    __tablename__ = "season_request"
    __table_args__ = (
        PrimaryKeyConstraint("season_id", "wanted_quality"),
    )
    season_id: Mapped[UUID] = mapped_column(ForeignKey(column="season.id", ondelete="CASCADE"), )
    wanted_quality: Mapped[Quality]
    min_quality: Mapped[Quality]
