import os

from pydantic import BaseModel


class BasicConfig(BaseModel):
    storage_directory: str = os.getenv("STORAGE_FILE_PATH") or "."









class MachineLearningConfig(BaseModel):
    model_name: str = os.getenv("OLLAMA_MODEL_NAME") or "qwen2.5:0.5b"