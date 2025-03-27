import logging
import sys
from logging.config import dictConfig

logging.basicConfig(level=logging.DEBUG,
                    format="%(asctime)s - %(levelname)s - %(name)s - %(funcName)s(): %(message)s",
                    stream=sys.stdout,
                    )
log = logging.getLogger(__name__)

import uvicorn
from fastapi import FastAPI

import database.users
import tv.router
from auth import password
from users import routers

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
app.include_router(routers.router, tags=["users"])
app.include_router(password.router, tags=["authentication"])
app.include_router(tv.router.router, tags=["tv"])

if __name__ == "__main__":
    uvicorn.run(app, host="127.0.0.1", port=5049, log_config=LOGGING_CONFIG)
