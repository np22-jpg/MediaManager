from fastapi import APIRouter, Depends
from fastapi import status
from sqlalchemy import select

from auth.config import OAuth2Config
from auth.db import User
from auth.schemas import UserRead
from auth.users import current_superuser
from database import DbSessionDependency
from auth.users import oauth_client

users_router = APIRouter()
auth_metadata_router = APIRouter()
oauth_enabled = oauth_client is not None
if oauth_enabled:
    oauth_config = OAuth2Config()


@users_router.get("/users/all", status_code=status.HTTP_200_OK, dependencies=[Depends(current_superuser)])
def get_all_users(db: DbSessionDependency) -> list[UserRead]:
    stmt = select(User)
    result = db.execute(stmt).scalars().unique()
    return [UserRead.model_validate(user) for user in result]


@auth_metadata_router.get("/auth/metadata", status_code=status.HTTP_200_OK)
def get_auth_metadata() -> dict:
    if oauth_enabled:
        return {
            "oauth_name": oauth_config.name,
        }
    else:
        return {"oauth_name": None}
