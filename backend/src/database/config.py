from pydantic_settings import BaseSettings, SettingsConfigDict


class DbConfig(BaseSettings):
    model_config = SettingsConfigDict(env_prefix="DB_")
    HOST: str = "localhost"
    PORT: int = 5432
    USER: str = "MediaManager"
    PASSWORD: str = "MediaManager"
    DBNAME: str = "MediaManager"
