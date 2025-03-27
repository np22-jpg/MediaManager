import logging
from typing import Annotated, Any, Generator

from fastapi import Depends
from sqlmodel import SQLModel, Session, create_engine

import config
from config import DbConfig

log = logging.getLogger(__name__)
config: DbConfig = config.get_db_config()

db_url = "postgresql+psycopg" + "://" + config.user + ":" + config.password + "@" + config.host + ":" + str(
    config.port) + "/" + config.dbname

engine = create_engine(db_url, echo=False)


def init_db() -> None:
    SQLModel.metadata.create_all(engine)


def get_session() -> Generator[Session, Any, None]:
    with Session(engine) as session:
        yield session

SessionDependency = Annotated[Session, Depends(get_session)]

