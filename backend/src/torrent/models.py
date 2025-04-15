import re
from typing import Literal
from uuid import UUID

from pydantic import computed_field
from sqlalchemy.orm import Mapped, mapped_column

from database import Base
from torrent.schemas import Quality


class QualityMixin:
    title: str

    @computed_field
    @property
    def quality(self) -> Quality:
        high_quality_pattern = r'\b(4k|4K)\b'
        medium_quality_pattern = r'\b(1080p|1080P)\b'
        low_quality_pattern = r'\b(720p|720P)\b'
        very_low_quality_pattern = r'\b(480p|480P|360p|360P)\b'

        if re.search(high_quality_pattern, self.title):
            return Quality.high
        elif re.search(medium_quality_pattern, self.title):
            return Quality.medium
        elif re.search(low_quality_pattern, self.title):
            return Quality.low
        elif re.search(very_low_quality_pattern, self.title):
            return Quality.very_low
        else:
            return Quality.unknown


class Torrent(Base):
    __tablename__ = "torrent"

    id: Mapped[UUID] = mapped_column(primary_key=True)
    status: Mapped[Literal["downloading", "finished", "error"] | None]
    title: Mapped[str]
    quality: Mapped[Quality | None]
