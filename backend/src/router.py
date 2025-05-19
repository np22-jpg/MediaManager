from fastapi import APIRouter, Depends
from fastapi import status
from sqlalchemy import select

from auth.db import User
from auth.schemas import UserRead
from auth.users import current_superuser
from database import DbSessionDependency

router = APIRouter()


@router.get("/users/all", status_code=status.HTTP_200_OK, dependencies=[Depends(current_superuser)])
def get_all_users(db: DbSessionDependency) -> list[UserRead]:
    stmt = select(User)
    result = db.execute(stmt).scalars().unique()
    return [UserRead.model_validate(user) for user in result]
