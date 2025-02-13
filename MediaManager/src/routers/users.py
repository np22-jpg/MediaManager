import logging

from fastapi import APIRouter
from fastapi import Depends

import database
from auth.password import authenticate_user, get_password_hash
from database import User, UserInternal

router = APIRouter(
    prefix="/users",
)


class CreateUser(User):
    password: str


log = logging.getLogger(__file__)


@router.post("/", response_model=User)
async def create_user(user: CreateUser):
    internal_user = UserInternal(name=user.name, lastname=user.lastname, email=user.email,
                                 hashed_password=get_password_hash(user.password))
    database.create_user(internal_user)
    return user


@router.get("/me/", response_model=User)
async def read_users_me(
        current_user: User = Depends(authenticate_user),
):
    return current_user


@router.get("/me/items/")
async def read_own_items(
        current_user: User = Depends(authenticate_user),
):
    return [{"item_id": "Foo", "owner": current_user.username}]
