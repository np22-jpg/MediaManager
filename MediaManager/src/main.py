from fastapi import FastAPI

import database
from routers import users
from auth import password

app = FastAPI()
app.include_router(users.router, tags=["users"])
app.include_router(password.app, tags=["authentication"])
database.__init__()