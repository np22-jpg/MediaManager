import logging
import pprint
from typing import Literal, List, Any
from uuid import UUID, uuid4

import requests
from pydantic import BaseModel

from config import TvConfig
from database import PgDatabase
import tmdbsimple as tmdb


# NOTE: use tmdbsimple for api calls

class Episode(BaseModel):
    number: int
    title: str


class Season(BaseModel):
    number: int
    episodes: List[Episode]

    def get_episode_count(self) -> int:
        return self.episodes.__len__()


class Show(BaseModel):
    id: UUID = uuid4()
    external_id: int
    metadata_provider: str
    name: str
    seasons: List[Season]  = []

    def get_season_count(self) -> int:
        return self.seasons.__len__()

    def get_episode_count(self) -> int:
        episode_count = 0
        for season in self.seasons:
            episode_count += season.get_episode_count()
        return episode_count

    def save_show(self) -> None:
        with PgDatabase() as db:
            db.connection.execute("""
            INSERT INTO tv_show (
                id,
                external_id,
                metadata_provider,
                name,
                episode_count,
                season_count
                )VALUES(%s,%s,%s,%s,%s,%s);
                """,
                                  (self.id,
                                   self.external_id,
                                   self.metadata_provider,
                                   self.name,
                                   self.get_episode_count(),
                                   self.get_season_count(),
                                   )
                                  )
        log.info("added show: " + self.__str__())

    def get_data_from_tmdb(self) -> None:
        data = tmdb.TV(self.external_id).info()
        log.debug("data from tmdb: " + pprint.pformat(data))
        self.name = data["original_name"]
        self.metadata_provider = "tmdb"

    def add_season(self, season_number: int) -> None:
        data = tmdb.TV_Seasons(self.external_id, season_number).info()
        log.debug("data from tmdb: " + pprint.pformat(data))

        episodes: List[Episode] = []
        for episode in data["episodes"]:
            episodes.append(Episode(title=episode["name"],number=episode["episode_number"]))

        season = Season(number=season_number, episodes=episodes)

        self.seasons.append(season)

    def add_seasons(self, season_numbers: List[int]) -> None:
        for season_number in season_numbers:
            self.add_season(season_number)

def get_all_shows() -> List[Show]:
    with PgDatabase() as db:
        result = db.connection.execute("""
        SELECT * FROM tv_show
        """).fetchall()
        return result

def get_show(id: UUID) -> Show:
    with PgDatabase() as db:
        result = db.connection.execute("""
        SELECT * FROM tv_show WHERE id = %s
        """, (id,)).fetchone()
    return  Show(**result)

config = TvConfig()
log = logging.getLogger(__name__)

tmdb.API_KEY = config.api_key
