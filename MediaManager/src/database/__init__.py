import logging
import os
from logging import getLogger

import psycopg
from psycopg.rows import dict_row

log = logging.getLogger(__name__)

log.debug("servas")

class PgDatabase():
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
    from database import user
    user.init_db()
    log.info("Tables initialized successfully")

init_db()
def drop_tables() -> None:
    with PgDatabase() as db:
        db.connection.execute("DROP TABLE IF EXISTS users CASCADE;")
        log.info("User Table dropped")

