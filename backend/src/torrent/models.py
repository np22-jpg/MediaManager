from uuid import UUID

from sqlalchemy.orm import Mapped, mapped_column

from database import Base
from torrent.schemas import Quality, TorrentStatus


class TorrentBase(Base):
    __abstract__ = True

    id: Mapped[UUID] = mapped_column(primary_key=True)
    status: Mapped[TorrentStatus | None]
    title: Mapped[str]
    quality: Mapped[Quality]
    imported: Mapped[bool]
