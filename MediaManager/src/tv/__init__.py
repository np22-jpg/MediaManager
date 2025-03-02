import logging
import pprint
from typing import List
from uuid import UUID, uuid4

import tmdbsimple as tmdb
from pydantic import BaseModel

from config import TvConfig
from database import PgDatabase



from sqlmodel import Field, Session, SQLModel, create_engine, select


# NOTE: use tmdbsimple for api calls

class Episode(SQLModel):
    show_id: int = Field(foreign_key="show.id")
    season_number: int = Field(foreign_key="season.number")
    number: int
    title: str


class Season(SQLModel, table=True):
    show_id: UUID = Field(foreign_key="show.id")
    number: int


class Show(SQLModel, table=True):
    id: UUID = Field(primary_key=True)
    external_id: int
    metadata_provider: str
    name: str

#   def get_season_count(self) -> int:
#       return self.seasons.__len__()

#   def get_episode_count(self) -> int:
#       episode_count = 0
#       for season in self.seasons:
#           episode_count += season.get_episode_count()
#       return episode_count

#   def save_show(self) -> None:
#       with PgDatabase() as db:
#           db.connection.execute("""
#           INSERT INTO tv_show (
#               id,
#               external_id,
#               metadata_provider,
#               name,
#               episode_count,
#               season_count
#               )VALUES(%s,%s,%s,%s,%s,%s);
#               """,
#                                 (self.id,
#                                  self.external_id,
#                                  self.metadata_provider,
#                                  self.name,
#                                  self.get_episode_count(),
#                                  self.get_season_count(),
#                                  )
#                                 )
#       log.info("added show: " + self.__str__())

#   def get_data_from_tmdb(self) -> None:
#       data = tmdb.TV(self.external_id).info()
#       log.debug("data from tmdb: " + pprint.pformat(data))
#       self.name = data["original_name"]
#       self.metadata_provider = "tmdb"

#   def add_season(self, season_number: int) -> None:
#       data = tmdb.TV_Seasons(self.external_id, season_number).info()
#       log.debug("data from tmdb: " + pprint.pformat(data))

#       episodes: List[Episode] = []
#       for episode in data["episodes"]:
#           episodes.append(Episode(title=episode["name"], number=episode["episode_number"]))

#       season = Season(number=season_number, episodes=episodes)

#       self.seasons.append(season)

#   def add_seasons(self, season_numbers: List[int]) -> None:
#       for season_number in season_numbers:
#           self.add_season(season_number)


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
    return Show(**result)


config = TvConfig()
log = logging.getLogger(__name__)

tmdb.API_KEY = config.api_key
