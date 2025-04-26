from fastapi import APIRouter, Depends, status
from fastapi.responses import JSONResponse

import metadataProvider
import tv.service
from auth.users import current_active_user
from database import DbSessionDependency
from indexer.schemas import PublicIndexerQueryResult, IndexerQueryResultId
from tv.exceptions import MediaAlreadyExists
from tv.schemas import Show, SeasonRequest, ShowId

router = APIRouter()


# --------------------------------
# CREATE AND DELETE SHOWS
# --------------------------------

@router.post("/shows", status_code=status.HTTP_201_CREATED, dependencies=[Depends(current_active_user)],
             responses={
                 status.HTTP_201_CREATED: {"model": Show, "description": "Successfully created show"},
                 status.HTTP_409_CONFLICT: {"model": str, "description": "Show already exists"},
             })
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
def get_all_shows(db: DbSessionDependency):
    """"""
    return tv.service.get_all_shows(db=db)


@router.get("/shows/{show_id}", dependencies=[Depends(current_active_user)], response_model=Show)
def get_a_show(db: DbSessionDependency, show_id: ShowId):
    """

    :param show_id:
    :type show_id:
    :return:
    :rtype:
    """

    return tv.service.get_show_by_id(db=db, show_id=show_id)


# --------------------------------
# MANAGE REQUESTS
# --------------------------------

@router.post("/season/request", status_code=status.HTTP_200_OK, dependencies=[Depends(current_active_user)])
def request_a_season(db: DbSessionDependency, season_request: SeasonRequest):
    """
    adds request flag to a season
    """
    tv.service.request_season(db=db, season_request=season_request)


@router.get("/season/request", status_code=status.HTTP_200_OK, dependencies=[Depends(current_active_user)])
def get_requested_seasons(db: DbSessionDependency) -> list[SeasonRequest]:
    return tv.service.get_all_requested_seasons(db=db)


@router.delete("/season/request", status_code=status.HTTP_200_OK, dependencies=[Depends(current_active_user)])
def unrequest_season(db: DbSessionDependency, request: SeasonRequest):
    tv.service.unrequest_season(db=db, season_request=request)


# --------------------------------
# MANAGE TORRENTS
# --------------------------------

# 1 is the default for season_number because it returns multi season torrents
@router.get("/torrents", status_code=status.HTTP_200_OK, dependencies=[Depends(current_active_user)],
            response_model=list[PublicIndexerQueryResult])
def get_torrents_for_a_season(db: DbSessionDependency, show_id: ShowId, season_number: int = 1):
    return tv.service.get_all_available_torrents_for_a_season(db=db, season_number=season_number, show_id=show_id)


# download a torrent
@router.post("/torrents", status_code=status.HTTP_200_OK, dependencies=[Depends(current_active_user)])
def download_a_torrent(db: DbSessionDependency, public_indexer_result_id: IndexerQueryResultId, show_id: ShowId,
                       override_file_path_suffix: str = ""):
    return tv.service.download_torrent(db=db, public_indexer_result_id=public_indexer_result_id, show_id=show_id,
                                       override_show_file_path_suffix=override_file_path_suffix)

# --------------------------------
# SEARCH SHOWS ON METADATA PROVIDERS
# --------------------------------

@router.get("/search", dependencies=[Depends(current_active_user)])
def search_metadata_providers_for_a_show(query: str, metadata_provider: str = "tmdb"):
    return metadataProvider.search_show(query, metadata_provider)
