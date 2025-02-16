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
    """"
    The Usermodel, but with an additional non-hashed password. attribute
    """
    password: str

log = logging.getLogger(__name__)
log.level = logging.DEBUG

@router.post("/")
async def create_user(
        user: CreateUser = Depends(CreateUser),
):
    internal_user = UserInternal(name=user.name, lastname=user.lastname, email=user.email,
                                 hashed_password=get_password_hash(user.password))
    if database.create_user(internal_user):
        log.info("Created new user",internal_user.model_dump())
        return user
    else:
        log.warning("Failed to create new user", internal_user.model_dump())
        return {"error": "Failed to create new user"}



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
