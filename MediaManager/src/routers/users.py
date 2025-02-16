import logging

from fastapi import APIRouter
from fastapi import Depends
from pydantic import BaseModel
from starlette.responses import JSONResponse

import database
from auth.password import authenticate_user, get_password_hash
from database import User, UserInternal

router = APIRouter(
    prefix="/users",
)


class Message(BaseModel):
    message: str

class CreateUser(User):
    """"
    The Usermodel, but with an additional non-hashed password. attribute
    """
    password: str

log = logging.getLogger(__name__)
log.level = logging.DEBUG

@router.post("/",status_code=201, responses = {
    409: {"model":  Message, "description": "User with provided email already exists"},
    201:{"model": UserInternal, "description": "User  created successfully"}
})
async def create_user(
        user: CreateUser = Depends(CreateUser),
):
    internal_user = UserInternal(name=user.name, lastname=user.lastname, email=user.email,
                                 hashed_password=get_password_hash(user.password))
    if database.create_user(internal_user):
        log.info("Created new user",internal_user.model_dump())
        return internal_user
    else:
        log.warning("Failed to create new user, User with this email already exists,",internal_user.model_dump())
        return JSONResponse(status_code=409, content={"message": "User with this email already exists"})



@router.get("/me", response_model=User)
async def read_users_me(
        current_user: User = Depends(authenticate_user),
):
    return current_user


@router.get("/me/items")
async def read_own_items(
        current_user: User = Depends(authenticate_user),
):
    return [{"item_id": "Foo", "owner": current_user.username}]
