from fastapi import APIRouter, Depends
from fastapi import status
from sqlalchemy import select

from media_manager.config import AllEncompassingConfig
from media_manager.auth.db import User
from media_manager.auth.schemas import UserRead
from media_manager.auth.users import current_superuser
from media_manager.database import DbSessionDependency

users_router = APIRouter()
auth_metadata_router = APIRouter()
oauth_config = AllEncompassingConfig().auth.openid_connect


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
    if oauth_config.enabled:
        return {
            "oauth_name": oauth_config.name,
        }
    else:
        return {"oauth_name": None}
