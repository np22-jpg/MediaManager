import re
from enum import Enum
from typing import Literal
from uuid import UUID, uuid4

from pydantic import computed_field
from sqlalchemy import Column, String
from sqlmodel import Field, SQLModel


class Quality(Enum):
    high = 1
    medium = 2
    low = 3
    very_low = 4
    unknown = 5


# TODO: make system to detect quality more sophisticated
class QualityMixin:
    title: str

    @property
    @computed_field
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


class TorrentMixin:
    torrent_id: UUID | None = Field(default=None, foreign_key="torrent.id")


class Torrent(SQLModel, QualityMixin, table=True):
    id: UUID = Field(default_factory=uuid4, primary_key=True)
    torrent_status: Literal["downloading", "finished", "error"] | None = Field(default=None,
                                                                               sa_column=Column(String))
    torrent_title: str = Field(default=None)

    @property
    @computed_field
    def torrent_filepath(self) -> str:
        return f"{self.id}.torrent"
