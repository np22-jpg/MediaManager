import logging
import sys
from logging.config import dictConfig

from pythonjsonlogger.json import JsonFormatter

LOGGING_CONFIG = {
    "version": 1,
    "disable_existing_loggers": False,
    "formatters": {
        "default": {
            "format": "%(asctime)s - %(levelname)s - %(name)s - %(funcName)s(): %(message)s"
        },
        "json": {
            "()": JsonFormatter,
        },
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
        },
    },
    "loggers": {
        "uvicorn": {"handlers": ["console", "file"], "level": "DEBUG"},
        "uvicorn.access": {"handlers": ["console", "file"], "level": "DEBUG"},
        "fastapi": {"handlers": ["console", "file"], "level": "DEBUG"},
    },
}
dictConfig(LOGGING_CONFIG)

logging.basicConfig(
    level=logging.DEBUG,
    format="%(asctime)s - %(levelname)s - %(name)s - %(funcName)s(): %(message)s",
    stream=sys.stdout,
)
log = logging.getLogger(__name__)

from media_manager.database import init_db
import media_manager.tv.router as tv_router
import media_manager.torrent.router as torrent_router

init_db()
log.info("Database initialized")

from media_manager.config import BasicConfig
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

basic_config = BasicConfig()
if basic_config.DEVELOPMENT:
    basic_config.torrent_directory.mkdir(parents=True, exist_ok=True)
    basic_config.tv_directory.mkdir(parents=True, exist_ok=True)
    basic_config.movie_directory.mkdir(parents=True, exist_ok=True)
    basic_config.image_directory.mkdir(parents=True, exist_ok=True)
    log.warning("Development Mode activated!")
else:
    log.info("Development Mode not activated!")

app = FastAPI(root_path="/api/v1")

if basic_config.DEVELOPMENT:
    origins = [
        "*",
    ]
else:
    origins = basic_config.CORS_URLS.split(",")
    log.info("CORS URLs activated for following origins:")
    for origin in origins:
        log.info(f" - {origin}")

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

import uvicorn
from fastapi.staticfiles import StaticFiles
from media_manager.auth.users import openid_client
from media_manager.auth.users import SECRET as AUTH_USERS_SECRET
from media_manager.auth.router import users_router as custom_users_router
from media_manager.auth.router import auth_metadata_router
from media_manager.auth.schemas import UserCreate, UserRead, UserUpdate
from media_manager.auth.users import (
    bearer_auth_backend,
    fastapi_users,
    cookie_auth_backend,
    openid_cookie_auth_backend,
)


# Standard Auth Routers
app.include_router(
    fastapi_users.get_auth_router(bearer_auth_backend),
    prefix="/auth/jwt",
    tags=["auth"],
)
app.include_router(
    fastapi_users.get_auth_router(cookie_auth_backend),
    prefix="/auth/cookie",
    tags=["auth"],
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
# All users route router
app.include_router(custom_users_router, tags=["users"])
# OAuth Metadata Router
app.include_router(auth_metadata_router, tags=["openid"])
# User Router
app.include_router(
    fastapi_users.get_users_router(UserRead, UserUpdate),
    prefix="/users",
    tags=["users"],
)
# OAuth2 Routers
if openid_client is not None:
    app.include_router(
        fastapi_users.get_oauth_router(
            openid_client,
            openid_cookie_auth_backend,
            AUTH_USERS_SECRET,
            associate_by_email=True,
            is_verified_by_default=True,
        ),
        prefix=f"/auth/cookie/{openid_client.name}",
        tags=["openid"],
    )

app.include_router(tv_router.router, prefix="/tv", tags=["tv"])
app.include_router(torrent_router.router, prefix="/torrent", tags=["torrent"])

# static file routers
app.mount(
    "/static/image",
    StaticFiles(directory=basic_config.image_directory),
    name="static-images",
)

if __name__ == "__main__":
    uvicorn.run(app, host="127.0.0.1", port=5049, log_config=LOGGING_CONFIG)
