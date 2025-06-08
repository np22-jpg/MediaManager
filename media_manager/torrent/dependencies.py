from typing import Annotated

from fastapi import Depends

from media_manager.database import DbSessionDependency
from media_manager.torrent.service import TorrentService
from media_manager.torrent.repository import TorrentRepository


def get_torrent_repository(db: DbSessionDependency) -> TorrentRepository:
    return TorrentRepository(db=db)

TorrentRepositoryDependency = Annotated[TorrentRepository, Depends(get_torrent_repository)]

def get_torrent_service(torrent_repository: TorrentRepositoryDependency) -> TorrentService:
    return TorrentService(torrent_repository=torrent_repository)


TorrentServiceDependency = Annotated[TorrentService, Depends(get_torrent_service)]
