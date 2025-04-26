from fastapi import APIRouter
from fastapi import status
from fastapi.params import Depends

from auth.users import current_active_user
from torrent.dependencies import TorrentServiceDependency
from torrent.schemas import TorrentId, Torrent

router = APIRouter()


@router.get("/", status_code=status.HTTP_200_OK, dependencies=[Depends(current_active_user)],
            response_model=list[Torrent])
def get_all_torrents(service: TorrentServiceDependency, ):
    return service.get_all_torrents()


@router.post("/{torrent_id}", status_code=status.HTTP_200_OK, dependencies=[Depends(current_active_user)],
             response_model=Torrent)
def import_torrent(service: TorrentServiceDependency, torrent_id: TorrentId):
    return service.import_torrent(service.get_torrent_by_id(id=torrent_id))


@router.post("/", status_code=status.HTTP_200_OK, dependencies=[Depends(current_active_user)],
             response_model=list[Torrent])
def import_all_torrents(service: TorrentServiceDependency):
    return service.import_all_torrents()


@router.get("/{torrent_id}", status_code=status.HTTP_200_OK, response_model=Torrent)
def get_torrent(service: TorrentServiceDependency, torrent_id: TorrentId):
    return service.get_torrent_by_id(id=torrent_id)


@router.delete("/torrents", status_code=status.HTTP_200_OK, dependencies=[Depends(current_active_user)])
def delete_torrent(service: TorrentServiceDependency, torrent_id: TorrentId):
    service.delete_torrent(torrent_id=torrent_id)
