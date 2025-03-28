import uuid
from uuid import UUID

from sqlalchemy import ForeignKey, ForeignKeyConstraint, Integer, String, UniqueConstraint
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column, relationship


class Base(DeclarativeBase):
    pass


class Show(Base):
    __tablename__ = "show"
    __table_args__ = (UniqueConstraint("external_id", "metadata_provider", "version"),)

    id: Mapped[UUID] = mapped_column(primary_key=True, default=uuid.uuid4)
    external_id: Mapped[int] = mapped_column(Integer, nullable=False)
    metadata_provider: Mapped[str] = mapped_column(String, nullable=False)
    name: Mapped[str] = mapped_column(String, nullable=False)
    overview: Mapped[str] = mapped_column(String, nullable=False)
    year: Mapped[int | None] = mapped_column(Integer, nullable=True)
    version: Mapped[str] = mapped_column(String, default="")

    seasons: Mapped[list["Season"]] = relationship(back_populates="show", cascade="all, delete")


class Season(Base):
    __tablename__ = "season"
    __table_args__ = (UniqueConstraint("show_id", "number"),)

    id: Mapped[UUID] = mapped_column(primary_key=True, default=uuid.uuid4)
    show_id: Mapped[UUID] = mapped_column(ForeignKey(column="show.id", ondelete="CASCADE"), nullable=False)
    number: Mapped[int] = mapped_column(Integer, nullable=False)
    external_id: Mapped[int] = mapped_column(Integer, nullable=False)
    name: Mapped[str] = mapped_column(String, nullable=False)
    overview: Mapped[str] = mapped_column(String, nullable=False)
    torrent_id: Mapped[UUID] = mapped_column(ForeignKey(column="torrent.id"), nullable=False)

    show: Mapped[Show] = relationship(back_populates="seasons")
    episodes: Mapped[list["Episode"]] = relationship(back_populates="season", cascade="all, delete")


class Episode(Base):
    __tablename__ = "episode"
    __table_args__ = (
        ForeignKeyConstraint(columns=["show_id", "season_number"],
                             refcolumns=["season.show_id", "season.number"],
                             ondelete="CASCADE"),
    )

    show_id: Mapped[UUID] = mapped_column(ForeignKey("show.id"), primary_key=True)
    season_number: Mapped[int] = mapped_column(Integer, primary_key=True)
    number: Mapped[int] = mapped_column(Integer, primary_key=True)
    external_id: Mapped[int] = mapped_column(Integer, nullable=False)
    title: Mapped[str] = mapped_column(String, nullable=False)

    season: Mapped[Season] = relationship(back_populates="episodes")
