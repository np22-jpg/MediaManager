import os

from pydantic import BaseModel


class BasicConfig(BaseModel):
    storage_directory: str = os.getenv("STORAGE_FILE_PATH") or "."

class ProwlarrConfig(BaseModel):
    enabled: bool = bool(os.getenv("PROWLARR_ENABLED") or True)
    api_key: str = os.getenv("PROWLARR_API_KEY")
    url: str = os.getenv("PROWLARR_URL")








class MachineLearningConfig(BaseModel):
    model_name: str = os.getenv("OLLAMA_MODEL_NAME") or "qwen2.5:0.5b"