import os

import tvdb_v4_official
import logging
from fastapi import APIRouter
from .cache import cache_response

log = logging.getLogger(__name__)

router = APIRouter(prefix="/tvdb", tags=["TVDB"])

tvdb_client = tvdb_v4_official.TVDB(os.getenv("TVDB_API_KEY"))


@router.get("/tv/trending")
@cache_response("tvdb_tv_trending", ttl=7200)
async def get_tvdb_trending_tv():
    return tvdb_client.get_all_series()


@router.get("/tv/search")
@cache_response("tvdb_tv_search", ttl=14400)
async def search_tvdb_tv(query: str, page: int = 1):
    return tvdb_client.search(query)


@router.get("/tv/shows/{show_id}")
@cache_response("tvdb_tv_show", ttl=28800)
async def get_tvdb_show(show_id: int):
    return tvdb_client.get_series_extended(show_id)


@router.get("/tv/seasons/{season_id}")
@cache_response("tvdb_tv_season", ttl=28800)
async def get_tvdb_season(season_id: int):
    return tvdb_client.get_season_extended(season_id)


@router.get("/movies/trending")
@cache_response("tvdb_movies_trending", ttl=7200)
async def get_tvdb_trending_movies():
    return tvdb_client.get_all_movies()


@router.get("/movies/search")
@cache_response("tvdb_movies_search", ttl=14400)
async def search_tvdb_movies(query: str):
    return tvdb_client.search(query)


@router.get("/movies/{movie_id}")
@cache_response("tvdb_movie", ttl=28800)
async def get_tvdb_movie(movie_id: int):
    return tvdb_client.get_movie_extended(movie_id)