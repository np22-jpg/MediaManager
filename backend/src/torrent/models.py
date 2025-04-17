from typing import Literal
from uuid import UUID

from sqlalchemy.orm import Mapped, mapped_column

from database import Base
from torrent.schemas import Quality


class Torrent(Base):
    __tablename__ = "torrent"

    id: Mapped[UUID] = mapped_column(primary_key=True)
    status: Mapped[Literal["downloading", "finished", "error"] | None]
    title: Mapped[str]
    quality: Mapped[Quality]
