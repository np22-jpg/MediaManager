from pydantic_settings import BaseSettings


class ProwlarrConfig(BaseSettings):
    enabled: bool = False
    api_key: str = ""
    url: str = "http://localhost:9696"


class JackettConfig(BaseSettings):
    enabled: bool = False
    api_key: str = ""
    url: str = "http://localhost:9696"
    indexers: list[str] = ["all"]


class IndexerConfig(BaseSettings):
    prowlarr: ProwlarrConfig
    jackett: JackettConfig
