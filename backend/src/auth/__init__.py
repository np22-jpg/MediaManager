import logging
from datetime import datetime, timedelta, timezone

import jwt
from fastapi import APIRouter, Depends, HTTPException, status
from fastapi.security import OAuth2PasswordBearer
from jwt.exceptions import InvalidTokenError
from pydantic import BaseModel

from auth.config import AuthConfig
from database import DbSessionDependency
from database.users import User


# TODO: evaluate FASTAPI-Users package

class Token(BaseModel):
    access_token: str
    token_type: str


class TokenData(BaseModel):
    uid: str | None = None


log = logging.getLogger(__name__)

oauth2_scheme = OAuth2PasswordBearer(tokenUrl="api/v1/token")

router = APIRouter()


async def get_current_user(db: DbSessionDependency, token: str = Depends(oauth2_scheme)) -> User:
    credentials_exception = HTTPException(
        status_code=status.HTTP_401_UNAUTHORIZED,
        detail="Could not validate credentials",
        headers={"WWW-Authenticate": "Bearer"},
    )
    auth_config = AuthConfig
    log.debug("token: " + token)

    try:
        payload = jwt.decode(token, auth_config.jwt_signing_key, algorithms=[auth_config.jwt_signing_algorithm])
        log.debug("jwt payload: " + payload.__str__())
        user_uid: str = payload.get("sub")
        log.debug("jwt payload sub (USER uid): " + user_uid)
        if user_uid is None:
            raise credentials_exception
        token_data = TokenData(uid=user_uid)
    except InvalidTokenError:
        log.warning("received invalid token: " + token)
        raise credentials_exception

    user: User | None = db.get(User, token_data.uid)

    if user is None:
        log.debug("USER not found")
        raise credentials_exception

    log.debug("received USER: " + user.__str__())
    return user


def create_access_token(data: dict, expires_delta: timedelta | None = None):
    to_encode = data.copy()
    auth_config = AuthConfig
    if expires_delta:
        expire = datetime.now(timezone.utc) + expires_delta
    else:
        expire = datetime.now(timezone.utc) + timedelta(minutes=auth_config.jwt_access_token_lifetime)
    to_encode.update({"exp": expire})
    encoded_jwt = jwt.encode(to_encode, auth_config.jwt_signing_key, algorithm=auth_config.jwt_signing_algorithm)
    return encoded_jwt
