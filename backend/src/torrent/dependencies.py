from typing import Annotated

from fastapi import Depends

from backend.src.database import DbSessionDependency
from torrent.service import TorrentService


def get_torrent_service(db: DbSessionDependency) -> TorrentService:
    return TorrentService(db=db)


TorrentServiceDependency = Annotated[TorrentService, Depends(get_torrent_service)]
