from pydantic_settings import BaseSettings, SettingsConfigDict


class AuthConfig(BaseSettings):
    # to get a signing key run:
    # openssl rand -hex 32
    model_config = SettingsConfigDict(env_prefix="AUTH_")
    token_secret: str
    session_lifetime: int = 60 * 60 * 24
    admin_email: str | list[str]

    @property
    def jwt_signing_key(self):
        return self._jwt_signing_key


class OAuth2Config(BaseSettings):
    model_config = SettingsConfigDict(env_prefix="OAUTH_")
    client_id: str
    client_secret: str
    authorize_endpoint: str
    access_token_endpoint: str
    user_info_endpoint: str
    name: str = "OAuth2"
