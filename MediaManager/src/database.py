import logging
import os
from abc import ABC, abstractmethod
from logging import getLogger
from uuid import uuid4

import psycopg
from annotated_types.test_cases import cases
from psycopg.rows import dict_row
from pydantic import BaseModel

log = getLogger(__name__)


class User(BaseModel):
    """
    User model
    """
    name: str
    lastname: str
    email: str


class UserInternal(User):
    """"
    Internal user model, assumes the password is already hashed, when a new instance is created
    """
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
            autocommit=True,
            host=os.getenv("DB_HOST"),
            port=os.getenv("DB_PORT"),
            user=os.getenv("DB_USERNAME"),
            password=os.getenv("DB_PASSWORD"),
            dbname=os.getenv("DB_NAME"),
            row_factory=dict_row
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
    """

    :param user: user to create, password must already be hashed
    :return:  True if user was created, False otherwise
    """
    with PgDatabase() as db:
        try:
            db.connection.execute(
                """
                INSERT INTO users (id, name, lastname, email, hashed_password)
                VALUES (%s, %s, %s, %s, %s)
                """,
                (user.id, user.name, user.lastname, user.email, user.hashed_password)
            )
        except psycopg.errors.UniqueViolation as e:
            log.error(e)
            return False

    log.info("User inserted successfully")
    log.debug(f"Inserted following User: "+ user.model_dump())
    return True


def get_user(email: str = None, uid: str = None) -> UserInternal | None:
    """
    either specify the email or uid to search for the user, if both parameters are provided the uid will be used


    :param email:  the users email address
    :param uid:  the users id
    :return:  if user was found it a UserInternal object is returned, otherwise None
    """
    with PgDatabase() as db:
        if email is not None and uid is None:
            result = db.connection.execute(
                    "SELECT id, name, lastname, email, hashed_password FROM users WHERE email=%s",
                    (email,)
                ).fetchone()
        if uid is not None:
            result = db.connection.execute(
                "SELECT id, name, lastname, email, hashed_password FROM users WHERE id=%s",
                (uid,)
            ).fetchone()

        if result is None:
            return None
        user = UserInternal(id = result["id"], name = result["name"], lastname = result["lastname"], email = result["email"], hashed_password = result["hashed_password"])
        log.debug(f"Retrieved User succesfully:  {user.model_dump()} ")
        return user
