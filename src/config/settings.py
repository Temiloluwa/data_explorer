from enum import Enum
from typing import Optional

from pydantic import BaseModel, Field
from pydantic_settings import BaseSettings


class AgentMode(Enum):
    BASIC = "basic"  # Fast but less accurate
    INTERMEDIATE = "intermediate"  # Balanced between speed and accuracy
    ADVANCED = "advanced"  # Higher accuracy but slower
    EXPERT = "expert"  # Maximum accuracy with the most computational cost


class LLMConfig(BaseModel):
    temperature: float = 0.7
    max_tokens: int = 2000


class EmbeddingConfig(BaseModel):
    chunk_size: int = 500
    chunk_overlap: int = 100


class VectorStoreConfig(BaseModel):
    mode: str = "local"
    uri: Optional[str] = "../../local_vectorstore/chroma"  # Path for local or URL for cloud
    port: Optional[int] = None
    api_key: Optional[str] = None
    collection_name: Optional[str] = "default"


class Settings(BaseSettings):
    openai_api_key: str = Field(..., env="OPENAI_API_KEY")
    embedding_model: str = "sentence-transformers/all-MiniLM-L6-v2"
    llm_model: str = "gpt-3.5-turbo-0125"
    vector_store: str = "chroma"
    agent_mode: AgentMode = AgentMode.BASIC  # Default to BASIC

    llm_config: LLMConfig = LLMConfig()
    embedding_config: EmbeddingConfig = EmbeddingConfig()
    vector_store_config: VectorStoreConfig = VectorStoreConfig()

    class Config:
        env_file = ".env"
