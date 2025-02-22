import logging
from datetime import datetime, timedelta, timezone
from typing import Annotated

import jwt
from fastapi import Depends, HTTPException, status, APIRouter
from fastapi.security import OAuth2, OAuth2AuthorizationCodeBearer
from jwt.exceptions import InvalidTokenError
from pydantic import BaseModel
import database
from database import UserInternal


class Token(BaseModel):
    access_token: str
    token_type: str


class TokenData(BaseModel):
    uid: str | None = None


# to get a string like this run:
# openssl rand -hex 32
# TODO: remove secrets from files
SECRET_KEY = "09d25e094faa6ca2556c818166b7a9563b93f7099f6f0f4caa6cf63b88e8d3e7"
ALGORITHM = "HS256"
ACCESS_TOKEN_EXPIRE_MINUTES = 30

log = logging.getLogger(__name__)
log.level = logging.DEBUG
log.addHandler(logging.StreamHandler())

router = APIRouter()

async def get_current_user(token: str) -> UserInternal:
    credentials_exception = HTTPException(
        status_code=status.HTTP_401_UNAUTHORIZED,
        detail="Could not validate credentials",
        headers={"WWW-Authenticate": "Bearer"},
    )
    log.debug("token: "+ token)
    try:
        payload = jwt.decode(token, SECRET_KEY, algorithms=[ALGORITHM])
        log.debug("jwt payload: "+payload.__str__())
        user_uid: str = payload.get("sub")
        log.debug("jwt payload sub (user uid): "+user_uid)
        if user_uid is None:
            raise credentials_exception
        token_data = TokenData(uid=user_uid)
    except InvalidTokenError:
        log.warning("received invalid token: "+token)
        raise credentials_exception
    user = database.get_user(uid=token_data.uid)
    if user is None:
        log.debug("user not found")
        raise credentials_exception
    log.debug("received user: "+user.__str__())
    return user

def create_access_token(data: dict, expires_delta: timedelta | None = None):
    to_encode = data.copy()
    if expires_delta:
        expire = datetime.now(timezone.utc) + expires_delta
    else:
        expire = datetime.now(timezone.utc) + timedelta(minutes=15)
    to_encode.update({"exp": expire})
    encoded_jwt = jwt.encode(to_encode, SECRET_KEY, algorithm=ALGORITHM)
    return encoded_jwt
