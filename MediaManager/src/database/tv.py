from typing import Literal
from uuid import UUID, uuid4

from pydantic import BaseModel

from database import PgDatabase, log


# NOTE: use tmdbsimple for api calls

class Show(BaseModel):
    id: UUID = uuid4()
    external_id: int
    indexer: Literal["tmdb"]
    name: str
    number_of_episodes: int
    number_of_seasons: int
    origin_country: list[str]
    original_language: str
    status: str
    first_air_date: str

def save_show(show: Show) -> None:
    with PgDatabase() as db:
        db.connection.execute("""
            INSERT INTO tv_shows (
                id,
                external_id,
                indexer,
                name,
                number_of_episodes,
                number_of_seasons,
                origin_country,
                original_language,
                status,
                first_air_date
                )VALUES(%s,%s,%s,%s,%s,%s,%s,%s,%s,%s);
                """,
                              (show.id,
                               show.external_id,
                               show.indexer,
                               show.name,
                               show.number_of_episodes,
                               show.number_of_seasons,
                               show.origin_country,
                               show.original_language,
                               show.status,
                               show.first_air_date
                               )
                              )
        log.info("added show: "+show.__str__())


def init_table():
    with PgDatabase() as db:
        db.connection.execute("""
            CREATE TABLE IF NOT EXISTS tv_shows (
                id UUID PRIMARY KEY,
                external_id NUMERIC,
                indexer TEXT,
                name TEXT,
                number_of_episodes INTEGER,
                number_of_seasons INTEGER,
                origin_country TEXT[],
                original_language TEXT,
                status TEXT,
                first_air_date TEXT
            );""")
    log.info("tv_shows Table initialized successfully")
