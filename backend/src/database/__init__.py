import logging
from contextvars import ContextVar
from typing import Annotated, Any, Generator

from fastapi import Depends
from sqlalchemy import create_engine
from sqlalchemy.orm import Session, declarative_base, sessionmaker

from .config import DbConfig

log = logging.getLogger(__name__)
config = DbConfig()

db_url = "postgresql+psycopg" + "://" + config.USER + ":" + config.PASSWORD + "@" + config.HOST + ":" + str(
    config.PORT) + "/" + config.DBNAME

engine = create_engine(db_url, echo=False)
Base = declarative_base()
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)



def init_db() -> None:
    Base.metadata.create_all(engine)


def get_session() -> Generator[Session, Any, None]:
    db = SessionLocal()
    try:
        yield db
        db.commit()
    except Exception as e:
        db.rollback()
        log.critical(f"error occurred: {e}")
    finally:
        db.close()


db_session: ContextVar[Session] = ContextVar('db_session')


DbSessionDependency = Annotated[Session, Depends(get_session)]
