
import uvicorn
from fastapi import FastAPI

from routers import users
from auth import password

app = FastAPI()

app.include_router(users.router, tags=["users"])
app.include_router(password.app, tags=["authentication"])



if __name__ == "__main__":
    uvicorn.run(app, host="127.0.0.1", port=5049)