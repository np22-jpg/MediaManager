from typing import Annotated

import hashlib

import bcrypt
from fastapi import Depends, HTTPException, status
from fastapi.security import OAuth2PasswordRequestForm
from sqlmodel import Session, select

import database
from auth import create_access_token, Token, router
from database import users, SessionDependency
from database.users import UserInternal


def verify_password(plain_password, hashed_password):
    return bcrypt.checkpw(
        bytes(plain_password, encoding="utf-8"),
        bytes(hashed_password, encoding="utf-8"),
    )


def get_password_hash(password: str) -> str:
    return bcrypt.hashpw(password.encode("utf-8"), bcrypt.gensalt()).decode("utf-8")


def authenticate_user(db: SessionDependency, email: str, password: str) -> bool | UserInternal:
    """

    :param email: email of the user
    :param password:  password of the user
    :return:  if authentication succeeds, returns the user object with added name and lastname, otherwise  or if the user doesn't exist returns False
    """
    user: UserInternal | None = db.exec(select(UserInternal).where(UserInternal.email == email)).first()
    if not user:
        return False
    if not verify_password(password, user.hashed_password):
        return False
    return user


@router.post("/token")
async def login_for_access_token(
        form_data: Annotated[OAuth2PasswordRequestForm, Depends()],
        db: SessionDependency,
) -> Token:
    user = authenticate_user(db,form_data.username, form_data.password)
    if not user:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Incorrect email or password",
            headers={"WWW-Authenticate": "Bearer"},
        )
    # id needs to be converted because a UUID object isn't json serializable
    access_token = create_access_token(data={"sub": user.id.__str__()})
    return Token(access_token=access_token, token_type="bearer")
