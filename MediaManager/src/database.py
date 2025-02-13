import logging
import os
import sys
from abc import ABC, abstractmethod
from logging import getLogger
from uuid import uuid4

import psycopg
from pydantic import BaseModel

log = getLogger(__name__)
log.addHandler(logging.StreamHandler(sys.stdout))
log.level = logging.DEBUG


class User(BaseModel):
    name: str
    lastname: str
    email: str


class UserInternal(User):
    id: str = str(uuid4())
    hashed_password: str


class Database(ABC):
    """
    Database context manager
    """

    def __init__(self, driver) -> None:
        self.driver = driver

    @abstractmethod
    def connect_to_database(self):
        raise NotImplementedError()

    def __enter__(self):
        self.connection = self.connect_to_database()
        return self

    def __exit__(self, exception_type, exc_val, traceback):
        self.connection.close()


class PgDatabase(Database):
    """PostgreSQL Database context manager using psycopg"""

    def __init__(self) -> None:
        self.driver = psycopg
        super().__init__(self.driver)

    def connect_to_database(self):
        return self.driver.connect(
            host=os.getenv("DB_HOST"),
            port=os.getenv("DB_PORT"),
            user=os.getenv("DB_USERNAME"),
            password=os.getenv("DB_PASSWORD"),
            dbname=os.getenv("DB_NAME")
        )


def init_db():
    with PgDatabase() as db:
        db.connection.execute("""
            CREATE TABLE IF NOT EXISTS users (
                id TEXT NOT NULL PRIMARY KEY,
                lastname TEXT,
                name TEXT NOT NULL,
                email TEXT NOT NULL UNIQUE,
                hashed_password TEXT NOT NULL
            );
        """)
        log.info("User Table initialized successfully")


def drop_tables() -> None:
    with PgDatabase() as db:
        db.connection.execute("DROP TABLE IF EXISTS users CASCADE;")
        log.info("User Table dropped")


def create_user(user: UserInternal) -> bool:
    with PgDatabase() as db:
        try:
            db.connection.execute(
                """
                INSERT INTO users (id, name, lastname, email, hashed_password)
                VALUES (%s, %s, %s, %s, %s)
                """,
                (user.id, user.name, user.lastname, user.email, user.hashed_password)
            )
        except psycopg.errors.UniqueViolation:
            return False

    log.info("User inserted successfully")
    log.debug(f"User {user.model_dump()} created successfully")
    return True


def get_user(email: str) -> UserInternal:
    with PgDatabase() as db:
        result = db.connection.execute(
            "SELECT id, name, lastname, email, hashed_password FROM users WHERE email=%s",
            (email,)
        ).fetchone()

        if result is None:
            return None

        user = UserInternal.model_construct(**dict(zip(["id", "name", "lastname", "email", "hashed_password"], result)))
        log.debug(f"User {user.model_dump()} retrieved successfully")
        return user
