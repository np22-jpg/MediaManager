import logging
import sys
from logging.config import dictConfig

from pythonjsonlogger.json import JsonFormatter

import auth.users

LOGGING_CONFIG = {
    "version": 1,
    "disable_existing_loggers": False,
    "formatters": {
        "default": {
            "format": "%(asctime)s - %(levelname)s - %(name)s - %(funcName)s(): %(message)s"
        },
        "json": {
            "()": JsonFormatter,
        }

    },
    "handlers": {
        "console": {
            "class": "logging.StreamHandler",
            "formatter": "default",
            "stream": sys.stdout,
        },
        "file": {
            "class": "logging.handlers.RotatingFileHandler",
            "formatter": "json",
            "filename": "./log.txt",
            "maxBytes": 10485760,
            "backupCount": 5,
        }
    },
    "loggers": {
        "uvicorn": {"handlers": ["console", "file"], "level": "DEBUG"},
        "uvicorn.access": {"handlers": ["console", "file"], "level": "DEBUG"},
        "fastapi": {"handlers": ["console", "file"], "level": "DEBUG"},
    },
}
dictConfig(LOGGING_CONFIG)

logging.basicConfig(level=logging.DEBUG,
                    format="%(asctime)s - %(levelname)s - %(name)s - %(funcName)s(): %(message)s",
                    stream=sys.stdout,
                    )
log = logging.getLogger(__name__)

import database
from auth.schemas import UserCreate, UserRead, UserUpdate
from auth.users import bearer_auth_backend, fastapi_users, cookie_auth_backend
from config import BasicConfig
from auth.users import oauth_client
import uvicorn
from fastapi import FastAPI

import tv.router
import torrent.router

basic_config = BasicConfig()
if basic_config.DEVELOPMENT:
    basic_config.torrent_directory.mkdir(parents=True, exist_ok=True)
    basic_config.tv_directory.mkdir(parents=True, exist_ok=True)
    basic_config.movie_directory.mkdir(parents=True, exist_ok=True)
    basic_config.image_directory.mkdir(parents=True, exist_ok=True)
    log.warning("Development Mode activated!")
else:
    log.info("Development Mode not activated!")

database.init_db()

app = FastAPI(root_path="/api/v1")
# Standard Auth Routers
app.include_router(
    fastapi_users.get_auth_router(bearer_auth_backend),
    prefix="/auth/jwt",
    tags=["auth"]
)
app.include_router(
    fastapi_users.get_auth_router(cookie_auth_backend),
    prefix="/auth/cookie",
    tags=["auth"]
)
app.include_router(
    fastapi_users.get_register_router(UserRead, UserCreate),
    prefix="/auth",
    tags=["auth"],
)
app.include_router(
    fastapi_users.get_reset_password_router(),
    prefix="/auth",
    tags=["auth"],
)
app.include_router(
    fastapi_users.get_verify_router(UserRead),
    prefix="/auth",
    tags=["auth"],
)
# User Router
app.include_router(
    fastapi_users.get_users_router(UserRead, UserUpdate),
    prefix="/users",
    tags=["users"],
)
# OAuth2 Routers
if oauth_client is not None:
    app.include_router(
        fastapi_users.get_oauth_router(oauth_client,
                                       bearer_auth_backend,
                                       auth.users.SECRET,
                                       associate_by_email=True,
                                       is_verified_by_default=True
                                       ),
        prefix=f"/auth/jwt/{oauth_client.name}",
        tags=["oauth"],
    )
    app.include_router(
        fastapi_users.get_oauth_router(oauth_client,
                                       cookie_auth_backend,
                                       auth.users.SECRET,
                                       associate_by_email=True,
                                       is_verified_by_default=True
                                       ),
        prefix=f"/auth/cookie/{oauth_client.name}",
        tags=["oauth"],

    )

app.include_router(
    tv.router.router,
    prefix="/tv",
    tags=["tv"]
)
app.include_router(
    torrent.router.router,
    prefix="/torrent",
    tags=["torrent"]
)
if __name__ == "__main__":
    uvicorn.run(app, host="127.0.0.1", port=5049, log_config=LOGGING_CONFIG)
