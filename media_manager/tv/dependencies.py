from typing import Annotated

from fastapi import Depends, Path

from media_manager.database import DbSessionDependency
from media_manager.tv.repository import TvRepository
from media_manager.tv.schemas import Show, ShowId, SeasonId, Season
from media_manager.tv.service import TvService


def get_tv_repository(db_session: DbSessionDependency) -> TvRepository:
    return TvRepository(db_session)


tv_repository_dep = Annotated[TvRepository, Depends(get_tv_repository)]


def get_tv_service(
    tv_repository: tv_repository_dep,
) -> TvService:
    return TvService(tv_repository)


tv_service_dep = Annotated[TvService, Depends(get_tv_service)]


def get_show_by_id(
    tv_service: tv_service_dep,
    show_id: ShowId = Path(..., description="The ID of the show"),
) -> Show:
    show = tv_service.get_show_by_id(show_id)
    return show


show_dep = Annotated[Show, Depends(get_show_by_id)]


def get_season_by_id(
    tv_service: tv_service_dep,
    season_id: SeasonId = Path(..., description="The ID of the season"),
) -> Season:
    return tv_service.get_season(season_id=season_id)


season_dep = Annotated[Season, Depends(get_season_by_id)]
