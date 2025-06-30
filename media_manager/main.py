import logging
import os
import sys
from logging.config import dictConfig
from pathlib import Path

from apscheduler.jobstores.sqlalchemy import SQLAlchemyJobStore
from pythonjsonlogger.json import JsonFormatter

import media_manager.database

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

from media_manager.database import init_db  # noqa: E402
from media_manager.config import BasicConfig  # noqa: E402
import media_manager.torrent.router as torrent_router  # noqa: E402
import media_manager.movies.router as movies_router  # noqa: E402
import media_manager.tv.router as tv_router  # noqa: E402
from media_manager.tv.service import (  # noqa: E402
    auto_download_all_approved_season_requests,
    import_all_show_torrents,
    update_all_non_ended_shows_metadata,
)
from media_manager.movies.service import (  # noqa: E402
    import_all_movie_torrents,
    update_all_movies_metadata,
)
import uvicorn  # noqa: E402
from fastapi.staticfiles import StaticFiles  # noqa: E402
from media_manager.auth.users import openid_client  # noqa: E402
from media_manager.auth.users import SECRET as AUTH_USERS_SECRET  # noqa: E402
from media_manager.auth.router import users_router as custom_users_router  # noqa: E402
from media_manager.auth.router import auth_metadata_router  # noqa: E402
from media_manager.auth.schemas import UserCreate, UserRead, UserUpdate  # noqa: E402
from media_manager.auth.oauth import get_oauth_router  # noqa: E402

from media_manager.auth.users import (  # noqa: E402
    bearer_auth_backend,
    fastapi_users,
    cookie_auth_backend,
    openid_cookie_auth_backend,
)
from media_manager.exceptions import (  # noqa: E402
    NotFoundError,
    not_found_error_exception_handler,
    MediaAlreadyExists,
    media_already_exists_exception_handler,
    InvalidConfigError,
    invalid_config_error_exception_handler,
)


import shutil  # noqa: E402
from fastapi import FastAPI  # noqa: E402
from fastapi.middleware.cors import CORSMiddleware  # noqa: E402
from datetime import datetime  # noqa: E402
from contextlib import asynccontextmanager  # noqa: E402
from apscheduler.schedulers.background import BackgroundScheduler  # noqa: E402
from apscheduler.triggers.cron import CronTrigger  # noqa: E402

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
    log.info(f"Hourly tasks are running at {datetime.now()}")
    auto_download_all_approved_season_requests()
    import_all_show_torrents()
    import_all_movie_torrents()


def weekly_tasks():
    log.info(f"Weekly tasks are running at {datetime.now()}")
    update_all_non_ended_shows_metadata()
    update_all_movies_metadata()


jobstores = {"default": SQLAlchemyJobStore(engine=media_manager.database.engine)}

scheduler = BackgroundScheduler(jobstores=jobstores)
trigger = CronTrigger(minute=0, hour="*")
weekly_trigger = CronTrigger(
    day_of_week="mon", hour=0, minute=0, jitter=60 * 60 * 24 * 2
)
scheduler.add_job(hourly_tasks, trigger)
scheduler.add_job(weekly_tasks, weekly_trigger)
scheduler.start()


@asynccontextmanager
async def lifespan(app: FastAPI):
    yield
    scheduler.shutdown()


base_path = os.getenv("API_BASE_PATH") or "/api/v1"
log.info("Base Path for API: %s", base_path)
app = FastAPI(root_path=base_path, lifespan=lifespan)

origins = basic_config.CORS_URLS
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


# ----------------------------
# Standard Auth Routers
# ----------------------------

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

# ----------------------------
# User Management Routers
# ----------------------------

app.include_router(custom_users_router, tags=["users"])
app.include_router(
    fastapi_users.get_users_router(UserRead, UserUpdate),
    prefix="/users",
    tags=["users"],
)

# ----------------------------
# OpenID Connect Routers
# ----------------------------

app.include_router(auth_metadata_router, tags=["openid"])

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
app.include_router(movies_router.router, prefix="/movies", tags=["movie"])
app.mount(
    "/static/image",
    StaticFiles(directory=basic_config.image_directory),
    name="static-images",
)

# ----------------------------
# Custom Exception Handlers
# ----------------------------

app.add_exception_handler(NotFoundError, not_found_error_exception_handler)
app.add_exception_handler(MediaAlreadyExists, media_already_exists_exception_handler)
app.add_exception_handler(InvalidConfigError, invalid_config_error_exception_handler)
log.info("Hello World!")

# ----------------------------
# Startup filesystem checks
# ----------------------------
try:
    test_dir = basic_config.tv_directory / Path(".media_manager_test_dir")
    test_dir.mkdir(parents=True, exist_ok=True)
    test_dir.rmdir()
    log.info(f"Successfully created test dir in TV directory at: {test_dir}")

    test_dir = basic_config.movie_directory / Path(".media_manager_test_dir")
    test_dir.mkdir(parents=True, exist_ok=True)
    test_dir.rmdir()
    log.info(f"Successfully created test dir in Movie directory at: {test_dir}")

    test_dir = basic_config.image_directory / Path(".media_manager_test_dir")
    test_dir.touch()
    test_dir.unlink()
    log.info(f"Successfully created test file in Image directory at: {test_dir}")

    # check if hardlink creation works
    test_dir = basic_config.tv_directory / Path(".media_manager_test_dir")
    test_dir.mkdir(parents=True, exist_ok=True)

    torrent_dir = basic_config.torrent_directory / Path(".media_manager_test_dir")
    torrent_dir.mkdir(parents=True, exist_ok=True)

    test_torrent_file = torrent_dir / Path(".media_manager.test.torrent")
    test_torrent_file.touch()

    test_hardlink = test_dir / Path(".media_manager.test.hardlink")
    try:
        test_hardlink.hardlink_to(test_torrent_file)
        if not test_hardlink.samefile(test_torrent_file):
            log.critical("Hardlink creation failed!")
        log.info("Successfully created test hardlink in TV directory")
    except OSError as e:
        log.error(
            f"Hardlink creation failed, falling back to copying files. Error: {e}"
        )
        shutil.copy(src=test_torrent_file, dst=test_hardlink)
    finally:
        test_hardlink.unlink()
        test_torrent_file.unlink()
        torrent_dir.rmdir()
        test_dir.rmdir()

except Exception as e:
    log.error(f"Error creating test directory: {e}")
    raise


@app.get("/", response_model={"message": str, "version": str})
async def hello_world():
    """
    A simple endpoint to check if the API is running.
    """
    return {"message": "Hello World!", "version": os.getenv("PUBLIC_VERSION")}


if __name__ == "__main__":
    uvicorn.run(app, host="127.0.0.1", port=5049, log_config=LOGGING_CONFIG)
