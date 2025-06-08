from fastapi import APIRouter
from fastapi import status
from fastapi.params import Depends

from media_manager.auth.users import current_active_user, current_superuser
from media_manager.torrent.dependencies import tv_service_dep
from media_manager.torrent.schemas import TorrentId, Torrent

router = APIRouter()


@router.get("/{torrent_id}", status_code=status.HTTP_200_OK, response_model=Torrent)
def get_torrent(service: tv_service_dep, torrent_id: TorrentId):
    return service.get_torrent_by_id(id=torrent_id)


@router.get(
    "/",
    status_code=status.HTTP_200_OK,
    dependencies=[Depends(current_active_user)],
    response_model=list[Torrent],
)
def get_all_torrents(service: tv_service_dep):
    return service.get_all_torrents()


@router.post(
    "/{torrent_id}",
    status_code=status.HTTP_200_OK,
    dependencies=[Depends(current_active_user)],
    response_model=Torrent,
)
def import_torrent(service: tv_service_dep, torrent_id: TorrentId):
    return service.import_torrent(service.get_torrent_by_id(id=torrent_id))


@router.post(
    "/",
    status_code=status.HTTP_200_OK,
    dependencies=[Depends(current_active_user)],
    response_model=list[Torrent],
)
def import_all_torrents(service: tv_service_dep):
    return service.import_all_torrents()


@router.delete(
    "/{torrent_id}",
    status_code=status.HTTP_200_OK,
    dependencies=[Depends(current_superuser)],
)
def delete_torrent(service: tv_service_dep, torrent_id: TorrentId):
    service.delete_torrent(torrent_id=torrent_id)
