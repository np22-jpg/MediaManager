from pydantic_settings import BaseSettings, SettingsConfigDict
from pydantic import Field
import secrets


class AuthConfig(BaseSettings):
    # to get a signing key run:
    # openssl rand -hex 32
    model_config = SettingsConfigDict(env_prefix="AUTH_")
    token_secret: str = Field(default_factory=secrets.token_hex)
    session_lifetime: int = 60 * 60 * 24
    admin_email: list[str] = []
    email_password_resets: bool = False

    @property
    def jwt_signing_key(self):
        return self._jwt_signing_key


class EmailConfig(BaseSettings):
    model_config = SettingsConfigDict(env_prefix="EMAIL_")
    smtp_host: str
    smtp_port: int
    smtp_user: str
    smtp_password: str
    from_email: str
    use_tls: bool = False


class OpenIdConfig(BaseSettings):
    model_config = SettingsConfigDict(env_prefix="OPENID_")
    client_id: str
    client_secret: str
    configuration_endpoint: str
    name: str = "OpenID"
