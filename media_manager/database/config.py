from pydantic_settings import BaseSettings


class DbConfig(BaseSettings):
    HOST: str = "localhost"
    PORT: int = 5432
    USER: str = "MediaManager"
    PASSWORD: str = "MediaManager"
    DBNAME: str = "MediaManager"
