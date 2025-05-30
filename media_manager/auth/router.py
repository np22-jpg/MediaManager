from fastapi import APIRouter, Depends
from fastapi import status
from sqlalchemy import select

from media_manager.auth.config import OpenIdConfig
from media_manager.auth.db import User
from media_manager.auth.schemas import UserRead
from media_manager.auth.users import current_superuser
from media_manager.database import DbSessionDependency
from media_manager.auth.users import openid_client

users_router = APIRouter()
auth_metadata_router = APIRouter()
oauth_enabled = openid_client is not None
if oauth_enabled:
    oauth_config = OpenIdConfig()


@users_router.get(
    "/users/all",
    status_code=status.HTTP_200_OK,
    dependencies=[Depends(current_superuser)],
)
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
