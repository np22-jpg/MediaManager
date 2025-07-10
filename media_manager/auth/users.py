import logging
import uuid
from typing import Optional

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

import media_manager.notification.utils
from media_manager.auth.db import User, get_user_db
from media_manager.auth.schemas import UserUpdate
from media_manager.config import AllEncompassingConfig

log = logging.getLogger(__name__)

config = AllEncompassingConfig().auth
SECRET = config.token_secret
LIFETIME = config.session_lifetime

if config.openid_connect.enabled:
    openid_config = AllEncompassingConfig().auth.openid_connect
    openid_client = OpenID(
        base_scopes=["openid", "email", "profile"],
        client_id=openid_config.client_id,
        client_secret=openid_config.client_secret,
        name=openid_config.name,
        openid_configuration_endpoint=openid_config.configuration_endpoint,
    )
    openid_client.base_scopes = ["openid", "email", "profile"]
else:
    openid_client = None


class UserManager(UUIDIDMixin, BaseUserManager[User, uuid.UUID]):
    reset_password_token_secret = SECRET
    verification_token_secret = SECRET

    async def on_after_register(self, user: User, request: Optional[Request] = None):
        log.info(f"User {user.id} has registered.")
        if user.email in config.admin_emails:
            updated_user = UserUpdate(is_superuser=True, is_verified=True)
            await self.update(user=user, user_update=updated_user)

    async def on_after_forgot_password(
        self, user: User, token: str, request: Optional[Request] = None
    ):
        link = f"{AllEncompassingConfig().misc.frontend_url}login/reset-password?token={token}"
        log.info(f"User {user.id} has forgot their password. Reset Link: {link}")

        if not config.email_password_resets:
            log.info("Email password resets are disabled, not sending email.")
            return

        subject = "MediaManager - Password Reset Request"
        html = f"""\
        <html>
          <body>
            <p>Hi {user.email},
            <br>
            <br>
            if you forgot your password, <a href="{link}">reset you password here</a>.<br>
            If you did not request a password reset, you can ignore this email.</p>
            <br>
            <br>
            If the link does not work, copy the following link into your browser: {link}<br>
          </body>
        </html>
        """
        media_manager.notification.utils.send_email(
            subject=subject, html=html, addressee=user.email
        )
        log.info(f"Sent password reset email to {user.email}")

    async def on_after_reset_password(
        self, user: User, request: Optional[Request] = None
    ):
        log.info(f"User {user.id} has reset their password.")

    async def on_after_request_verify(
        self, user: User, token: str, request: Optional[Request] = None
    ):
        log.info(
            f"Verification requested for user {user.id}. Verification token: {token}"
        )

    async def on_after_verify(self, user: User, request: Optional[Request] = None):
        log.info(f"User {user.id} has been verified")


async def get_user_manager(user_db: SQLAlchemyUserDatabase = Depends(get_user_db)):
    yield UserManager(user_db)


def get_jwt_strategy() -> JWTStrategy[models.UP, models.ID]:
    return JWTStrategy(secret=SECRET, lifetime_seconds=LIFETIME)


# needed because the default CookieTransport does not redirect after login,
# thus the user would be stuck on the OAuth Providers "redirecting" page
class RedirectingCookieTransport(CookieTransport):
    async def get_login_response(self, token: str) -> Response:
        response = RedirectResponse(
            str(AllEncompassingConfig().misc.frontend_url) + "dashboard",
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
