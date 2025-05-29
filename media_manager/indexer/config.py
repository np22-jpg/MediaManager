from pydantic_settings import BaseSettings, SettingsConfigDict


class ProwlarrConfig(BaseSettings):
    model_config = SettingsConfigDict(env_prefix="PROWLARR_")
    enabled: bool = True
    api_key: str
    url: str = "http://localhost:9696"
