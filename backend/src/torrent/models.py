from uuid import UUID

from sqlalchemy.orm import Mapped, mapped_column, relationship

from backend.src.database import Base
from torrent.schemas import Quality, TorrentStatus


class Torrent(Base):
    __tablename__ = "torrent"
    id: Mapped[UUID] = mapped_column(primary_key=True)
    status: Mapped[TorrentStatus]
    title: Mapped[str]
    quality: Mapped[Quality]
    imported: Mapped[bool]
    hash: Mapped[str]

    season_files = relationship("SeasonFile", back_populates="torrent")
