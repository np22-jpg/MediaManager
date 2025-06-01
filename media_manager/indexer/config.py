from pydantic_settings import BaseSettings, SettingsConfigDict


class ProwlarrConfig(BaseSettings):
    model_config = SettingsConfigDict(env_prefix="PROWLARR_")
    enabled: bool | None = False
    api_key: str | None
    url: str = "http://localhost:9696"


# TODO: add this to docs
class JackettConfig(BaseSettings):
    model_config = SettingsConfigDict(env_prefix="JACKETT_")
    enabled: bool | None = False
    api_key: str | None
    url: str = "http://localhost:9696"
    indexers: list[str] = [
        "all"
    ]  # needs to be formatted like this ["indexer1", "indexer2", "indexer3"] in env file (note double quotes!)
