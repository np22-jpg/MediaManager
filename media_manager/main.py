import logging
import os
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
from media_manager.tv.service import auto_download_all_approved_season_requests
import media_manager.torrent.router as torrent_router
from media_manager.config import BasicConfig
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from datetime import datetime
from contextlib import asynccontextmanager
from apscheduler.schedulers.background import BackgroundScheduler
from apscheduler.triggers.cron import CronTrigger
import media_manager.torrent.service
from media_manager.database import SessionLocal

init_db()
log.info("Database initialized")

basic_config = BasicConfig()
if basic_config.DEVELOPMENT:
    basic_config.torrent_directory.mkdir(parents=True, exist_ok=True)
    basic_config.tv_directory.mkdir(parents=True, exist_ok=True)
    basic_config.movie_directory.mkdir(parents=True, exist_ok=True)
    basic_config.image_directory.mkdir(parents=True, exist_ok=True)
    log.warning("Development Mode activated!")
else:
    log.info("Development Mode not activated!")


def hourly_tasks():
    log.info(f"Tasks are running at {datetime.now()}")
    auto_download_all_approved_season_requests()
    # media_manager.torrent.service.TorrentService(
    #    db=SessionLocal()
    #).import_all_torrents()


scheduler = BackgroundScheduler()
trigger = CronTrigger(second=0, hour="*")
scheduler.add_job(hourly_tasks, trigger)
scheduler.start()


@asynccontextmanager
async def lifespan(app: FastAPI):
    yield
    scheduler.shutdown()


base_path = os.getenv("API_BASE_PATH") or "/api/v1"
log.info("Base Path for API: %s", base_path)
app = FastAPI(root_path=base_path, lifespan=lifespan)

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
from media_manager.auth.oauth import get_oauth_router

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
        get_oauth_router(
            oauth_client=openid_client,
            backend=openid_cookie_auth_backend,
            get_user_manager=fastapi_users.get_user_manager,
            state_secret=AUTH_USERS_SECRET,
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

log.info("Hello World!")

if __name__ == "__main__":
    uvicorn.run(app, host="127.0.0.1", port=5049, log_config=LOGGING_CONFIG)
