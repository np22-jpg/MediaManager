from fastapi import APIRouter
from fastapi import status
from fastapi.params import Depends

from media_manager.auth.users import current_active_user
from media_manager.torrent.dependencies import torrent_service_dep, torrent_dep
from media_manager.torrent.schemas import Torrent

router = APIRouter()


@router.get("/{torrent_id}", status_code=status.HTTP_200_OK, response_model=Torrent)
def get_torrent(service: torrent_service_dep, torrent: torrent_dep):
    return service.get_torrent_by_id(torrent_id=torrent.id)


@router.get(
    "",
    status_code=status.HTTP_200_OK,
    dependencies=[Depends(current_active_user)],
    response_model=list[Torrent],
)
def get_all_torrents(service: torrent_service_dep):
    return service.get_all_torrents()
