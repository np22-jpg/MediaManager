from fastapi import APIRouter
from fastapi import Depends
from pydantic import BaseModel
from starlette.responses import JSONResponse

import database
from auth import get_current_user
from auth.password import get_password_hash
from database.users import UserInternal, User
from routers import log

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


@router.post("/", status_code=201, responses={
    409: {"model": Message, "description": "User with provided email already exists"},
    201: {"model": UserInternal, "description": "User  created successfully"}
})
async def create_user(
        user: CreateUser = Depends(CreateUser),
):
    internal_user = UserInternal(name=user.name, lastname=user.lastname, email=user.email,
                                 hashed_password=get_password_hash(user.password))
    if database.users.create_user(internal_user):
        log.info("Created new user", internal_user.model_dump())
        return internal_user
    else:
        log.warning("Failed to create new user, User with this email already exists,", internal_user.model_dump())
        return JSONResponse(status_code=409, content={"message": "User with this email already exists"})


@router.get("/me")
async def read_users_me(
        current_user: UserInternal = Depends(get_current_user),
):
    return JSONResponse(status_code=200, content=current_user.model_dump())


@router.get("/me/items")
async def read_own_items(
        current_user: User = Depends(get_current_user),
):
    return [{"item_id": "Foo", "owner": current_user.name}]
