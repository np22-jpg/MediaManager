from pydantic_settings import BaseSettings


class BasicConfig(BaseSettings):
    storage_directory: str = "."
