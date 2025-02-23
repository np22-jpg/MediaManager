import os
from typing import Literal

from pydantic import BaseModel

class DbConfig(BaseModel):
    host: str = os.getenv("DB_HOST") or "localhost"
    port: int = int(os.getenv("DB_PORT")) or 5432
    user: str = os.getenv("DB_USERNAME") or "MediaManager"
    password: str = os.getenv("DB_PASSWORD") or "MediaManager"
    dbname: str = os.getenv("DB_NAME") or "MediaManager"


class IndexerConfig(BaseModel):
    default_indexer: Literal["tmdb"] = os.getenv("INDEXER") or "tmdb"
    default_indexer_api_key: str = os.getenv("INDEXER_API_KEY")

class AuthConfig(BaseModel):
    # to get a signing key run:
    # openssl rand -hex 32
    jwt_signing_key: str = os.getenv("JWT_SIGNING_KEY")
    jwt_signing_algorithm: str = "HS256"
    jwt_access_token_lifetime: int = int(os.getenv("JWT_ACCESS_TOKEN_LIFETIME")) or 60*24*30

db: DbConfig = DbConfig()
indexer: IndexerConfig = IndexerConfig()
auth: AuthConfig = AuthConfig()

if __name__ == "__main__":
    print(db.__str__())
    print(indexer.__str__())
    print(auth.__str__())