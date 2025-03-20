import logging
import os
from typing import Literal

from pydantic import BaseModel


class DbConfig(BaseModel):
    host: str = os.getenv("DB_HOST") or "localhost"
    port: int = int(os.getenv("DB_PORT")) or 5432
    user: str = os.getenv("DB_USERNAME") or "MediaManager"
    _password: str = os.getenv("DB_PASSWORD") or "MediaManager"
    dbname: str = os.getenv("DB_NAME") or "MediaManager"

    @property
    def password(self):
        return self._password

class TvConfig(BaseModel):
    api_key: str = os.getenv("TMDB_API_KEY")


class IndexerConfig(BaseModel):
    default_indexer: Literal["tmdb"] = os.getenv("INDEXER") or "tmdb"
    _default_indexer_api_key: str = os.getenv("INDEXER_API_KEY")


class AuthConfig(BaseModel):
    # to get a signing key run:
    # openssl rand -hex 32
    _jwt_signing_key: str = os.getenv("JWT_SIGNING_KEY")
    jwt_signing_algorithm: str = "HS256"
    jwt_access_token_lifetime: int = int(os.getenv("JWT_ACCESS_TOKEN_LIFETIME")) or 60 * 24 * 30

    @property
    def jwt_signing_key(self):
        return self._jwt_signing_key


class MachineLearningConfig(BaseModel):
    model_name: str = os.getenv("OLLAMA_MODEL_NAME") or "qwen2.5:0.5b"


def get_db_config() -> DbConfig:
    return DbConfig()



log = logging.getLogger(__name__)

def load_config():
    log.info(f"loaded config: DbConfig: {DbConfig().__str__()}")
    log.info(f"loaded config: IndexerConfig: {IndexerConfig().__str__()}")
    log.info(f"loaded config: AuthConfig: {AuthConfig().__str__()}")
    log.info(f"loaded config: TvConfig: {TvConfig().__str__()}")


if __name__ == "__main__":
    db: DbConfig = DbConfig()
    indexer: IndexerConfig = IndexerConfig()
    auth: AuthConfig = AuthConfig()

    print(db.__str__())
    print(indexer.__str__())
    print(auth.__str__())
