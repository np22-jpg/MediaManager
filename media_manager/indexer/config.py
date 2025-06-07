from pydantic_settings import BaseSettings, SettingsConfigDict


class ProwlarrConfig(BaseSettings):
    model_config = SettingsConfigDict(env_prefix="PROWLARR_")
    enabled: bool | None = False
    api_key: str | None = None
    url: str = "http://localhost:9696"


class JackettConfig(BaseSettings):
    model_config = SettingsConfigDict(env_prefix="JACKETT_")
    enabled: bool | None = False
    api_key: str | None = None
    url: str = "http://localhost:9696"
    indexers: list[str] = ["all"]
