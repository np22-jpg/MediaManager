import typing
import uuid
from abc import ABC
from enum import Enum

from pydantic import ConfigDict

TorrentId = typing.NewType("TorrentId", uuid.UUID)

class Quality(Enum):
    high = 1
    medium = 2
    low = 3
    very_low = 4
    unknown = 5


class TorrentStatus(Enum):
    finished = 1
    downloading = 2
    error = 3
    unknown = 4


class TorrentBase(ABC):
    model_config = ConfigDict(from_attributes=True)

    id: TorrentId
    status: TorrentStatus
    title: str
    quality: Quality
    imported: bool
