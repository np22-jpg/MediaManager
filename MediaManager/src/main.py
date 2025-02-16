import logging

import uvicorn
from fastapi import FastAPI, Depends
from pydantic import BaseModel

import database
from fastapi.testclient import TestClient
from routers import users
from auth import password
from routers.users import CreateUser

app = FastAPI()

logging.info("OIDA")
app.include_router(users.router, tags=["users"])
app.include_router(password.app, tags=["authentication"])



if __name__ == "__main__":
    uvicorn.run(app, host="127.0.0.1", port=5049)