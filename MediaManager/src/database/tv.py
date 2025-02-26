from typing import Literal, List
from uuid import UUID, uuid4

from pydantic import BaseModel

from database import PgDatabase, log


# NOTE: use tmdbsimple for api calls

class Episode(BaseModel):
    number: int
    title: str

class Season(BaseModel):
    number: int
    episodes: List[Episode]

    def get_episode_count(self)-> int:
        return self.episodes.__len__()

class Show(BaseModel):
    id: UUID = uuid4()
    external_id: int
    indexer: Literal["tmdb"]
    name: str
    seasons: List[Season]

    def get_season_count(self)-> int:
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
                indexer,
                name,
                episode_count,
                season_count
                )VALUES(%s,%s,%s,%s,%s,%s);
                """,
                                  (self.id,
                                   self.external_id,
                                   self.indexer,
                                   self.name,
                                    self.get_episode_count(),
                                   self.get_season_count(),
                                   )
                                  )
            log.info("added show: " + self.__str__())


# TODO: add NOT NULL and default values to DB

def init_table():
    with PgDatabase() as db:
        db.connection.execute("""
            CREATE TABLE IF NOT EXISTS tv_show (
                id UUID PRIMARY KEY,
                external_id TEXT,
                indexer TEXT,
                name TEXT,
                episode_count INTEGER,
                season_count INTEGER
            );""")
        log.info("tv_show Table initialized successfully")
        db.connection.execute("""
            CREATE TABLE IF NOT EXISTS tv_season (
                show_id UUID  REFERENCES tv_show(id),
                season_number INTEGER,
                episode_count INTEGER,
                CONSTRAINT PK_season PRIMARY KEY (show_id,season_number)

            );""")
        log.info("tv_seasonTable initialized successfully")
        db.connection.execute("""
            CREATE TABLE IF NOT EXISTS tv_episode (
                season  INTEGER,
                show_id uuid,
                episode_number INTEGER,
                title TEXT,
                CONSTRAINT PK_episode PRIMARY KEY (season,show_id,episode_number),
                FOREIGN KEY (season, show_id) REFERENCES tv_season(season_number,show_id)

            );""")
        log.info("tv_episode Table initialized successfully")
