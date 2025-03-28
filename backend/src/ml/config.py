from pydantic_settings import BaseSettings


class MachineLearningConfig(BaseSettings):
    ollama_model_name: str = "qwen2.5:0.5b"
