from typing import Annotated

from fastapi import APIRouter, Depends, status, HTTPException
from fastapi.responses import JSONResponse

from media_manager.auth.db import User
from media_manager.auth.schemas import UserRead
from media_manager.auth.users import current_active_user, current_superuser
from media_manager.indexer.schemas import PublicIndexerQueryResult, IndexerQueryResultId
from media_manager.metadataProvider.schemas import MetaDataProviderShowSearchResult
from media_manager.torrent.schemas import Torrent
from media_manager.tv import log
from media_manager.tv.exceptions import MediaAlreadyExists
from media_manager.tv.schemas import (
    Show,
    SeasonRequest,
    ShowId,
    RichShowTorrent,
    PublicShow,
    PublicSeasonFile,
    CreateSeasonRequest,
    SeasonRequestId,
    UpdateSeasonRequest,
    RichSeasonRequest,
)
from media_manager.tv.dependencies import (
    torrent_service_dep,
    season_dep,
    show_dep,
    tv_repository_dep,
)

router = APIRouter()


# --------------------------------
# CREATE AND DELETE SHOWS
# --------------------------------


@router.post(
    "/shows",
    status_code=status.HTTP_201_CREATED,
    dependencies=[Depends(current_active_user)],
    responses={
        status.HTTP_201_CREATED: {
            "model": Show,
            "description": "Successfully created show",
        },
        status.HTTP_409_CONFLICT: {"model": str, "description": "Show already exists"},
    },
)
def add_a_show(
    tv_service: torrent_service_dep, show_id: int, metadata_provider: str = "tmdb"
):
    try:
        show = tv_service.add_show(
            external_id=show_id,
            metadata_provider=metadata_provider,
        )
    except MediaAlreadyExists as e:
        return JSONResponse(
            status_code=status.HTTP_409_CONFLICT, content={"message": str(e)}
        )
    return show


@router.delete(
    "/shows/{show_id}",
    status_code=status.HTTP_204_NO_CONTENT,
    dependencies=[Depends(current_active_user)],
)
def delete_a_show(tv_repository: tv_repository_dep, show_id: ShowId):
    tv_repository.delete_show(show_id=show_id)


# --------------------------------
# GET SHOW INFORMATION
# --------------------------------


@router.get(
    "/shows", dependencies=[Depends(current_active_user)], response_model=list[Show]
)
def get_all_shows(
    tv_service: torrent_service_dep, external_id: int = None, metadata_provider: str = "tmdb"
):
    if external_id is not None:
        return tv_service.get_show_by_external_id(
            external_id=external_id, metadata_provider=metadata_provider
        )
    else:
        return tv_service.get_all_shows()


@router.get(
    "/shows/torrents",
    dependencies=[Depends(current_active_user)],
    response_model=list[RichShowTorrent],
)
def get_shows_with_torrents(tv_service: torrent_service_dep):
    """
    get all shows that are associated with torrents
    :return: A list of shows with all their torrents
    """
    result = tv_service.get_all_shows_with_torrents()
    return result


@router.get(
    "/shows/{show_id}",
    dependencies=[Depends(current_active_user)],
    response_model=PublicShow,
)
def get_a_show(show: show_dep, tv_service: torrent_service_dep) -> PublicShow:
    return tv_service.get_public_show_by_id(show_id=show.id)


@router.get(
    "/shows/{show_id}/torrents",
    dependencies=[Depends(current_active_user)],
    response_model=RichShowTorrent,
)
def get_a_shows_torrents(show: show_dep, tv_service: torrent_service_dep):
    return tv_service.get_torrents_for_show(show=show)


# --------------------------------
# MANAGE REQUESTS
# --------------------------------


@router.post("/seasons/requests", status_code=status.HTTP_204_NO_CONTENT)
def request_a_season(
    user: Annotated[User, Depends(current_active_user)],
    season_request: CreateSeasonRequest,
    tv_service: torrent_service_dep,
):
    """
    adds request flag to a season
    """
    log.info(f"Got season request: {season_request.model_dump()}")
    request: SeasonRequest = SeasonRequest.model_validate(season_request)
    request.requested_by = UserRead.model_validate(user)
    if user.is_superuser:
        request.authorized = True
        request.authorized_by = UserRead.model_validate(user)
    log.info(f"Adding season request: {request.model_dump()}")
    tv_service.add_season_request(request)
    return


@router.get(
    "/seasons/requests",
    status_code=status.HTTP_200_OK,
    dependencies=[Depends(current_active_user)],
    response_model=list[RichSeasonRequest],
)
def get_season_requests(tv_service: torrent_service_dep) -> list[RichSeasonRequest]:
    return tv_service.get_all_season_requests()


@router.delete(
    "/seasons/requests/{request_id}",
    status_code=status.HTTP_204_NO_CONTENT,
)
def delete_season_request(
    tv_service: torrent_service_dep,
    user: Annotated[User, Depends(current_active_user)],
    request_id: SeasonRequestId,
):
    request = tv_service.get_season_request_by_id(season_request_id=request_id)
    if user.is_superuser or request.requested_by.id == user.id:
        tv_service.delete_season_request(season_request_id=request_id)
        log.info(f"User {user.id} deleted season request {request_id}.")
        return None
    else:
        log.warning(
            f"User {user.id} tried to delete season request {request_id} but is not authorized."
        )
        return HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Not authorized to delete this request",
        )


@router.patch(
    "/seasons/requests/{season_request_id}", status_code=status.HTTP_204_NO_CONTENT
)
def authorize_request(
    tv_service: torrent_service_dep,
    user: Annotated[User, Depends(current_superuser)],
    season_request_id: SeasonRequestId,
    authorized_status: bool = False,
):
    """
    updates the request flag to true
    """
    season_request = tv_service.get_season_request_by_id(
        season_request_id=season_request_id
    )
    season_request.authorized_by = UserRead.model_validate(user)
    season_request.authorized = authorized_status
    if not authorized_status:
        season_request.authorized_by = None
    tv_service.update_season_request(season_request=season_request)
    return


@router.put("/seasons/requests", status_code=status.HTTP_204_NO_CONTENT)
def update_request(
    tv_service: torrent_service_dep,
    user: Annotated[User, Depends(current_active_user)],
    season_request: UpdateSeasonRequest,
):
    updated_season_request: SeasonRequest = SeasonRequest.model_validate(season_request)
    request = tv_service.get_season_request_by_id(
        season_request_id=updated_season_request.id
    )
    if request.requested_by.id == user.id or user.is_superuser:
        updated_season_request.requested_by = UserRead.model_validate(user)
        tv_service.update_season_request(season_request=updated_season_request)
    return


@router.get(
    "/seasons/{season_id}/files",
    dependencies=[Depends(current_active_user)],
    response_model=list[PublicSeasonFile],
)
def get_season_files(
    season: season_dep, tv_service: torrent_service_dep
) -> list[PublicSeasonFile]:
    return tv_service.get_public_season_files_by_season_id(season_id=season.id)


# --------------------------------
# MANAGE TORRENTS
# --------------------------------


# 1 is the default for season_number because it returns multi season torrents
@router.get(
    "/torrents",
    status_code=status.HTTP_200_OK,
    dependencies=[Depends(current_superuser)],
    response_model=list[PublicIndexerQueryResult],
)
def get_torrents_for_a_season(
    tv_service: torrent_service_dep,
    show_id: ShowId,
    season_number: int = 1,
    search_query_override: str = None,
):
    return tv_service.get_all_available_torrents_for_a_season(
        season_number=season_number,
        show_id=show_id,
        search_query_override=search_query_override,
    )


# download a torrent
@router.post(
    "/torrents",
    status_code=status.HTTP_200_OK,
    response_model=Torrent,
    dependencies=[Depends(current_superuser)],
)
def download_a_torrent(
    tv_service: torrent_service_dep,
    public_indexer_result_id: IndexerQueryResultId,
    show_id: ShowId,
    override_file_path_suffix: str = "",
):
    return tv_service.download_torrent(
        public_indexer_result_id=public_indexer_result_id,
        show_id=show_id,
        override_show_file_path_suffix=override_file_path_suffix,
    )


# --------------------------------
# SEARCH SHOWS ON METADATA PROVIDERS
# --------------------------------


@router.get(
    "/search",
    dependencies=[Depends(current_active_user)],
    response_model=list[MetaDataProviderShowSearchResult],
)
def search_metadata_providers_for_a_show(
    tv_service: torrent_service_dep, query: str, metadata_provider: str = "tmdb"
):
    return tv_service.search_for_show(query=query, metadata_provider=metadata_provider)


@router.get(
    "/recommended",
    dependencies=[Depends(current_active_user)],
    response_model=list[MetaDataProviderShowSearchResult],
)
def get_recommended_shows(tv_service: torrent_service_dep, metadata_provider: str = "tmdb"):
    return tv_service.get_popular_shows(metadata_provider=metadata_provider)
