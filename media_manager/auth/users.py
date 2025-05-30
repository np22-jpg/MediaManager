import os
import uuid
from typing import Optional

import httpx
from fastapi import Depends, Request
from fastapi_users import BaseUserManager, FastAPIUsers, UUIDIDMixin, models
from fastapi_users.authentication import (
    AuthenticationBackend,
    BearerTransport,
    CookieTransport,
    JWTStrategy,
)
from fastapi_users.db import SQLAlchemyUserDatabase
from httpx_oauth.clients.openid import OpenID
from fastapi.responses import RedirectResponse, Response
from starlette import status

from media_manager.auth.config import AuthConfig, OpenIdConfig
from media_manager.auth.db import User, get_user_db
from media_manager.auth.schemas import UserUpdate
from media_manager.config import BasicConfig

config = AuthConfig()
SECRET = config.token_secret
LIFETIME = config.session_lifetime

if (
        os.getenv("OAUTH_ENABLED") is not None
        and os.getenv("OAUTH_ENABLED").upper() == "TRUE"
):
    openid_config = OpenIdConfig()
    openid_client = OpenID(
        client_id=openid_config.client_id,
        client_secret=openid_config.client_secret,
        name=openid_config.name,
        openid_configuration_endpoint=openid_config.configuration_endpoint,
        base_scopes=["openid", "email", "profile"],
    )
else:
    openid_client = None


# TODO: implement on_xxx methods
class UserManager(UUIDIDMixin, BaseUserManager[User, uuid.UUID]):
    reset_password_token_secret = SECRET
    verification_token_secret = SECRET

    async def on_after_register(self, user: User, request: Optional[Request] = None):
        print(f"User {user.id} has registered.")
        if user.email in config.admin_email:
            updated_user = UserUpdate(is_superuser=True, is_verified=True)
            await self.update(user=user, user_update=updated_user)

    async def on_after_forgot_password(
            self, user: User, token: str, request: Optional[Request] = None
    ):
        print(f"User {user.id} has forgot their password. Reset token: {token}")

    async def on_after_reset_password(
            self, user: User, request: Optional[Request] = None
    ):
        print(f"User {user.id} has reset their password.")

    async def on_after_request_verify(
            self, user: User, token: str, request: Optional[Request] = None
    ):
        print(f"Verification requested for user {user.id}. Verification token: {token}")

    async def on_after_verify(self, user: User, request: Optional[Request] = None):
        print(f"User {user.id} has been verified")


async def get_user_manager(user_db: SQLAlchemyUserDatabase = Depends(get_user_db)):
    yield UserManager(user_db)


def get_jwt_strategy() -> JWTStrategy[models.UP, models.ID]:
    return JWTStrategy(secret=SECRET, lifetime_seconds=LIFETIME)


# needed because the default CookieTransport does not redirect after login,
# thus the user would be stuck on the OAuth Providers "redirecting" page
class RedirectingCookieTransport(CookieTransport):
    async def get_login_response(self, token: str) -> Response:
        response = RedirectResponse(
            str(BasicConfig().FRONTEND_URL) + "dashboard",
            status_code=status.HTTP_302_FOUND,
        )
        return self._set_login_cookie(response, token)


bearer_transport = BearerTransport(tokenUrl="auth/jwt/login")
cookie_transport = CookieTransport(cookie_max_age=LIFETIME)
openid_cookie_transport = RedirectingCookieTransport(cookie_max_age=LIFETIME)

bearer_auth_backend = AuthenticationBackend(
    name="jwt",
    transport=bearer_transport,
    get_strategy=get_jwt_strategy,
)
cookie_auth_backend = AuthenticationBackend(
    name="cookie",
    transport=cookie_transport,
    get_strategy=get_jwt_strategy,
)
openid_cookie_auth_backend = AuthenticationBackend(
    name="cookie",
    transport=openid_cookie_transport,
    get_strategy=get_jwt_strategy,
)

fastapi_users = FastAPIUsers[User, uuid.UUID](
    get_user_manager, [bearer_auth_backend, cookie_auth_backend]
)

current_active_user = fastapi_users.current_user(active=True, verified=True)
current_superuser = fastapi_users.current_user(
    active=True, verified=True, superuser=True
)
