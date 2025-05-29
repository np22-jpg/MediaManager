from typing import Annotated

from fastapi import Depends

from media_manager.database import DbSessionDependency
from media_manager.torrent.service import TorrentService


def get_torrent_service(db: DbSessionDependency) -> TorrentService:
    return TorrentService(db=db)


TorrentServiceDependency = Annotated[TorrentService, Depends(get_torrent_service)]
