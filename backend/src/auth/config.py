from pydantic_settings import BaseSettings


class AuthConfig(BaseSettings):
    # to get a signing key run:
    # openssl rand -hex 32
    jwt_signing_key: str
    jwt_signing_algorithm: str = "HS256"
    jwt_access_token_lifetime: int = 60 * 24 * 30

    @property
    def jwt_signing_key(self):
        return self._jwt_signing_key
