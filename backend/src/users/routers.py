from fastapi import APIRouter, Depends
from pydantic import BaseModel
from sqlalchemy.exc import IntegrityError
from starlette.responses import JSONResponse

from auth import get_current_user
from auth.password import get_password_hash
from database import SessionDependency
from database.users import User, UserCreate, UserPublic
from users import log

router = APIRouter(
    prefix="/users",
)

# TODO: remove from users.py
class Message(BaseModel):
    message: str



@router.post("/", status_code=201, responses={
    409: {"model": Message, "description": "User with provided email already exists"},
    201: {"model": UserPublic, "description": "User  created successfully"}
})
async def create_user(
        db: SessionDependency,
        user: UserCreate = Depends(UserCreate),
):
    internal_user = User(name=user.name, lastname=user.lastname, email=user.email,
                         hashed_password=get_password_hash(user.password))
    db.add(internal_user)
    try:
        db.commit()
    except IntegrityError as e:
        log.debug(e)
        log.warning("Failed to create new USER, User with this email already exists " + internal_user.model_dump().__str__())
        return JSONResponse(status_code=409, content={"message": "User with this email already exists"})
    log.info("Created new USER " + internal_user.email)
    return UserPublic(**internal_user.model_dump())


@router.get("/me")
async def read_users_me(
        current_user: User = Depends(get_current_user),
):
    return JSONResponse(status_code=200, content=current_user.model_dump())


@router.get("/me/items")
async def read_own_items(
        current_user: User = Depends(get_current_user),
):
    return [{"item_id": "Foo", "owner": current_user.name}]
