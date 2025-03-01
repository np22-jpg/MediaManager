from typing import Annotated
from uuid import UUID

import psycopg.errors
from fastapi import APIRouter, Depends

import auth
from database.users import User, UserInternal
from tv import Show, get_all_shows, tmdb, log, get_show

router = APIRouter(
    prefix="/tv",
)


@router.post("/show",  status_code=201,dependencies=[Depends(auth.get_current_user)])
def post_add_show_route(show_id: int, metadata_provider: str = "tmdb"):
    show: Show = Show(external_id=show_id, metadata_provider=metadata_provider, name="temp_name_set_in_post_show_route")
    show.get_data_from_tmdb()

    try:
        show.save_show()
    except psycopg.errors.UniqueViolation:
        log.info("Show already exists " + show.__str__())
        return show
    return show

@router.post("/{show_id}/{season}", status_code=201,dependencies=[Depends(auth.get_current_user)])
def post_add_season_route(show_id: UUID, season: int):
    show = get_show(show_id)
    show.add_season(season)
    show.save_show()
    return get_show(show_id)

@router.get("/show")
def get_shows_route():
    return get_all_shows()

@router.get("/search")
def search_show_route(query: str):
    search = tmdb.Search()
    return search.tv(query=query)

