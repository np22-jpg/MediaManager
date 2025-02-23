import logging
import os

import psycopg
from psycopg.rows import dict_row

log = logging.getLogger(__name__)


class PgDatabase:
    """PostgreSQL Database context manager using psycopg"""

    def __init__(self) -> None:
        self.driver = psycopg

    def connect_to_database(self):
        return self.driver.connect(
            autocommit=True,
            host=os.getenv("DB_HOST"),
            port=os.getenv("DB_PORT"),
            user=os.getenv("DB_USERNAME"),
            password=os.getenv("DB_PASSWORD"),
            dbname=os.getenv("DB_NAME"),
            row_factory=dict_row
        )

    def __enter__(self):
        self.connection = self.connect_to_database()
        return self

    def __exit__(self, exception_type, exc_val, traceback):
        self.connection.close()


def init_db():
    log.info("Initializing database")

    from database import tv, users
    users.init_table()
    tv.init_table()

    log.info("Tables initialized successfully")
