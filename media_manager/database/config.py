from pydantic_settings import BaseSettings


class DbConfig(BaseSettings):
    host: str
    port: int = 5432
    user: str = "MediaManager"
    password: str = "MediaManager"
    dbname: str = "MediaManager"
