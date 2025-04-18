from pydantic_settings import BaseSettings


class DownloadClientConfig(BaseSettings):
    download_client: str = "qbit"
