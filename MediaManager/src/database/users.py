from uuid import uuid4

import psycopg
from pydantic import BaseModel

from database import PgDatabase, log


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
    log.debug(f"Inserted User: " + user.model_dump().__str__())
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
        user = UserInternal(id=result["id"].__str__(), name=result["name"], lastname=result["lastname"],
                            email=result["email"],
                            hashed_password=result["hashed_password"])
        log.debug(f"Retrieved User successfully:  {user.model_dump()} ")
        return user


def init_table():
    with PgDatabase() as db:
        db.connection.execute("""
            CREATE TABLE IF NOT EXISTS users (
                id UUID NOT NULL PRIMARY KEY,
                lastname TEXT,
                name TEXT NOT NULL,
                email TEXT NOT NULL UNIQUE,
                hashed_password TEXT NOT NULL
            );
        """)
    log.info("users Table initialized successfully")
