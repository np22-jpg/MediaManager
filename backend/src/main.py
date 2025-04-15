import logging
import sys
from logging.config import dictConfig

import database
from auth.schemas import UserCreate, UserRead, UserUpdate
from auth.users import bearer_auth_backend, fastapi_users

logging.basicConfig(level=logging.DEBUG,
                    format="%(asctime)s - %(levelname)s - %(name)s - %(funcName)s(): %(message)s",
                    stream=sys.stdout,
                    )
log = logging.getLogger(__name__)

import uvicorn
from fastapi import FastAPI

import tv.router

LOGGING_CONFIG = {
    "version": 1,
    "disable_existing_loggers": True,
    "formatters": {
        "default": {
            "format": "%(asctime)s - %(levelname)s - %(name)s - %(funcName)s(): %(message)s"
        }
    },
    "handlers": {
        "console": {
            "class": "logging.StreamHandler",
            "formatter": "default",
            "stream": sys.stdout,
        },
    },
    "loggers": {
        "uvicorn": {"handlers": ["console"], "level": "DEBUG"},
        "uvicorn.access": {"handlers": ["console"], "level": "DEBUG"},
        "fastapi": {"handlers": ["console"], "level": "DEBUG"},
        "__main__": {"handlers": ["console"], "level": "DEBUG"},
    },
}

# Apply logging config
dictConfig(LOGGING_CONFIG)

database.init_db()

app = FastAPI(root_path="/api/v1")
app.include_router(
    fastapi_users.get_auth_router(bearer_auth_backend),
    prefix="/auth/jwt",
    tags=["auth"]
)
app.include_router(
    fastapi_users.get_auth_router(bearer_auth_backend),
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
app.include_router(
    fastapi_users.get_users_router(UserRead, UserUpdate),
    prefix="/users",
    tags=["users"],
)

app.include_router(
    tv.router.router,
    prefix="/tv",
    tags=["tv"]
)

if __name__ == "__main__":
    uvicorn.run(app, host="127.0.0.1", port=5049, log_config=LOGGING_CONFIG)
