import logging
import sys

import uvicorn
from fastapi import FastAPI

import config
import database.users
import tv.router
from auth import password
from routers import users

logging.basicConfig(level=logging.DEBUG,
                    format="%(asctime)s - %(levelname)s - %(name)s - %(funcName)s(): %(message)s",
                    stream=sys.stdout)

config.load_config()
database.init_db()
app = FastAPI(root_path="/api/v1")
app.include_router(users.router, tags=["users"])
app.include_router(password.router, tags=["authentication"])
app.include_router(tv.router.router, tags=["tv"])




if __name__ == "__main__":
    uvicorn.run(app, host="127.0.0.1", port=5049)
