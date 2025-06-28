import logging
import os

import tmdbsimple
from tmdbsimple import TV, TV_Seasons, Movies, Trending, Search
from fastapi import APIRouter
from .cache import cache_response

log = logging.getLogger(__name__)

tmdb_api_key = os.getenv("TMDB_API_KEY")
router = APIRouter(prefix="/tmdb", tags=["TMDB"])

if tmdb_api_key is None:
    log.warning("TMDB_API_KEY environment variable is not set.")
else:
    tmdbsimple.API_KEY = tmdb_api_key

    @router.get("/tv/trending")
    @cache_response("tmdb_tv_trending", ttl=7200)
    async def get_tmdb_trending_tv():
        return Trending(media_type="tv").info()

    @router.get("/tv/search")
    @cache_response("tmdb_tv_search", ttl=14400)
    async def search_tmdb_tv(query: str, page: int = 1):
        return Search().tv(page=page, query=query, include_adult=True)

    @router.get("/tv/shows/{show_id}")
    @cache_response("tmdb_tv_show", ttl=28800)
    async def get_tmdb_show(show_id: int):
        return TV(show_id).info()

    @router.get("/tv/shows/{show_id}/{season_number}")
    @cache_response("tmdb_tv_season", ttl=28800)
    async def get_tmdb_season(season_number: int, show_id: int):
        return TV_Seasons(season_number=season_number, tv_id=show_id).info()

    @router.get("/movies/trending")
    @cache_response("tmdb_movies_trending", ttl=7200)
    async def get_tmdb_trending_movies():
        return Trending(media_type="movie").info()

    @router.get("/movies/search")
    @cache_response("tmdb_movies_search", ttl=14400)
    async def search_tmdb_movies(query: str, page: int = 1):
        return Search().movie(page=page, query=query, include_adult=True)

    @router.get("/movies/{movie_id}")
    @cache_response("tmdb_movie", ttl=28800)
    async def get_tmdb_movie(movie_id: int):
        return Movies(movie_id).info()



