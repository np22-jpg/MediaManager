import logging
from typing import Annotated

from fastapi import APIRouter, Depends, status
from fastapi.responses import JSONResponse

import tv.repository
import tv.service
from auth.db import User
from auth.schemas import UserRead
from auth.users import current_active_user, current_superuser
from backend.src.database import DbSessionDependency
from indexer.schemas import PublicIndexerQueryResult, IndexerQueryResultId
from metadataProvider.schemas import MetaDataProviderShowSearchResult
from torrent.schemas import Torrent
from tv import log
from tv.exceptions import MediaAlreadyExists
from tv.schemas import Show, SeasonRequest, ShowId, RichShowTorrent, PublicShow, PublicSeasonFile, SeasonNumber, \
    CreateSeasonRequest, SeasonRequestId, UpdateSeasonRequest, RichSeasonRequest

router = APIRouter()


# --------------------------------
# CREATE AND DELETE SHOWS
# --------------------------------

@router.post("/shows", status_code=status.HTTP_201_CREATED, dependencies=[Depends(current_active_user)],
             responses={status.HTTP_201_CREATED: {"model": Show, "description": "Successfully created show"},
                        status.HTTP_409_CONFLICT: {"model": str, "description": "Show already exists"}, })
def add_a_show(db: DbSessionDependency, show_id: int, metadata_provider: str = "tmdb"):
    try:
        show = tv.service.add_show(db=db, external_id=show_id, metadata_provider=metadata_provider, )
    except MediaAlreadyExists as e:
        return JSONResponse(status_code=status.HTTP_409_CONFLICT, content={"message": str(e)})
    return show


@router.delete("/shows/{show_id}", status_code=status.HTTP_200_OK, dependencies=[Depends(current_active_user)])
def delete_a_show(db: DbSessionDependency, show_id: ShowId):
    db.delete(db.get(Show, show_id))
    db.commit()


# --------------------------------
# GET SHOW INFORMATION
# --------------------------------

@router.get("/shows", dependencies=[Depends(current_active_user)], response_model=list[Show])
def get_all_shows(db: DbSessionDependency, external_id: int = None, metadata_provider: str = "tmdb"):
    if external_id is not None:
        return tv.service.get_show_by_external_id(db=db, external_id=external_id, metadata_provider=metadata_provider)
    else:
        return tv.service.get_all_shows(db=db)


@router.get("/shows/torrents", dependencies=[Depends(current_active_user)], response_model=list[RichShowTorrent])
def get_shows_with_torrents(db: DbSessionDependency):
    """
    get all shows that are associated with torrents
    :return: A list of shows with all their torrents
    """
    result = tv.service.get_all_shows_with_torrents(db=db)
    return result


@router.get("/shows/{show_id}", dependencies=[Depends(current_active_user)], response_model=PublicShow)
def get_a_show(db: DbSessionDependency, show_id: ShowId):
    return tv.service.get_public_show_by_id(db=db, show_id=show_id)


@router.get("/shows/{show_id}/torrents", dependencies=[Depends(current_active_user)], response_model=RichShowTorrent)
def get_a_shows_torrents(db: DbSessionDependency, show_id: ShowId):
    return tv.service.get_torrents_for_show(db=db, show=tv.service.get_show_by_id(db=db, show_id=show_id))


# TODO: replace by route with season_id rather than show_id and season_number
@router.get("/shows/{show_id}/{season_number}/files", status_code=status.HTTP_200_OK,
            dependencies=[Depends(current_active_user)])
def get_season_files(db: DbSessionDependency, season_number: SeasonNumber, show_id: ShowId) -> list[PublicSeasonFile]:
    return tv.service.get_public_season_files_by_season_number(db=db, season_number=season_number, show_id=show_id)


# --------------------------------
# MANAGE REQUESTS
# --------------------------------

@router.post("/seasons/requests", status_code=status.HTTP_204_NO_CONTENT)
def request_a_season(db: DbSessionDependency, user: Annotated[User, Depends(current_active_user)],
                     season_request: CreateSeasonRequest):
    """
    adds request flag to a season
    """
    request: SeasonRequest = SeasonRequest.model_validate(season_request)
    request.requested_by = UserRead.model_validate(user)
    tv.service.add_season_request(db=db, season_request=request)
    return


@router.get("/seasons/requests", status_code=status.HTTP_200_OK, dependencies=[Depends(current_active_user)],
            response_model=list[RichSeasonRequest])
def get_season_requests(db: DbSessionDependency) -> list[RichSeasonRequest]:
    return tv.service.get_all_season_requests(db=db)


@router.delete("/seasons/requests/{request_id}", status_code=status.HTTP_204_NO_CONTENT,
               dependencies=[Depends(current_active_user)])
def delete_season_request(db: DbSessionDependency, request_id: SeasonRequestId):
    tv.service.delete_season_request(db=db, season_request_id=request_id)
    return



@router.patch("/seasons/requests/{season_request_id}", status_code=status.HTTP_204_NO_CONTENT)
def authorize_request(db: DbSessionDependency, user: Annotated[User, Depends(current_superuser)],
                      season_request_id: SeasonRequestId, authorized_status: bool = False):
    """
    updates the request flag to true
    """
    season_request: SeasonRequest = tv.repository.get_season_request(db=db, season_request_id=season_request_id)
    season_request.authorized_by = UserRead.model_validate(user)
    season_request.authorized = authorized_status
    tv.service.update_season_request(db=db, season_request=season_request)
    return


@router.put("/seasons/requests", status_code=status.HTTP_204_NO_CONTENT)
def update_request(db: DbSessionDependency, user: Annotated[User, Depends(current_superuser)],
                   season_request: UpdateSeasonRequest):
    season_request: SeasonRequest = SeasonRequest.model_validate(season_request)
    season_request.requested_by = UserRead.model_validate(user)
    tv.service.update_season_request(db=db, season_request=season_request)
    return

# --------------------------------
# MANAGE TORRENTS
# --------------------------------

# 1 is the default for season_number because it returns multi season torrents
@router.get("/torrents", status_code=status.HTTP_200_OK, dependencies=[Depends(current_superuser)],
            response_model=list[PublicIndexerQueryResult])
def get_torrents_for_a_season(db: DbSessionDependency, show_id: ShowId, season_number: int = 1,
                              search_query_override: str = None):
    return tv.service.get_all_available_torrents_for_a_season(db=db, season_number=season_number, show_id=show_id,
                                                              search_query_override=search_query_override)


# download a torrent
@router.post("/torrents", status_code=status.HTTP_200_OK, response_model=Torrent,
             dependencies=[Depends(current_superuser)])
def download_a_torrent(db: DbSessionDependency, public_indexer_result_id: IndexerQueryResultId, show_id: ShowId,
                       override_file_path_suffix: str = ""):
    return tv.service.download_torrent(db=db, public_indexer_result_id=public_indexer_result_id, show_id=show_id,
                                       override_show_file_path_suffix=override_file_path_suffix)


# --------------------------------
# SEARCH SHOWS ON METADATA PROVIDERS
# --------------------------------

@router.get("/search", dependencies=[Depends(current_active_user)],
            response_model=list[MetaDataProviderShowSearchResult])
def search_metadata_providers_for_a_show(db: DbSessionDependency, query: str, metadata_provider: str = "tmdb"):
    return tv.service.search_for_show(query=query, metadata_provider=metadata_provider, db=db)
