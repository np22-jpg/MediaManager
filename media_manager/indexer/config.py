from pydantic_settings import BaseSettings


class ProwlarrConfig(BaseSettings):
    enabled: bool | None = False
    api_key: str | None = None
    url: str = "http://localhost:9696"


class JackettConfig(BaseSettings):
    enabled: bool | None = False
    api_key: str | None = None
    url: str = "http://localhost:9696"
    indexers: list[str] = ["all"]


class IndexerConfig(BaseSettings):
    prowlarr: ProwlarrConfig = ProwlarrConfig()
    jackett: JackettConfig = JackettConfig()
