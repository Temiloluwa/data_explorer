from abc import ABC, abstractmethod
from typing import List, Optional
import os
import chromadb

from core import AIModels
from langchain_chroma import Chroma
from langchain_core.documents import Document
from config import settings


class VectorStore(ABC):
    @abstractmethod
    async def add_documents(self, documents: list[Document]):
        pass

    @abstractmethod
    async def similarity_search(self, query: str, k: int = 4) -> list[Document]:
        pass


class ChromaStore(VectorStore):
    def __init__(self):
        self.config = settings.vector_store_config
        self.embeddings = AIModels.get_embeddings_model(settings)
        self.client_settings = self._get_client_settings()
        self.vectorstore: Optional[Chroma] = None
        self.collection_name = self.config.collection_name or "qa_agent_collection"

    def _get_client_settings(self) -> chromadb.config.Settings:
        """Configure client based on deployment mode"""
        if self.config.mode == "local":
            if not os.path.isdir(self.config.uri):
                raise ValueError(f"Directory not found: {self.config.uri}")
            return chromadb.PersistentClient(path=self.config.uri)
        elif self.config.mode == "online":
            return chromadb.HttpClient(host=self.config.uri, port=self.config.port) 
        else:
            raise ValueError(f"Unsupported Chroma mode: {self.config.mode}")

    async def _get_or_create_vectorstore(self) -> Chroma:
        """Initialize or reuse existing Chroma instance"""
        if not self.vectorstore:
            self.vectorstore = Chroma(
                collection_name=self.collection_name,
                embedding_function=self.embeddings,
                client=self._get_client_settings()
            )
        return self.vectorstore

    async def add_documents(self, documents: list[Document]):
        vectorstore = await self._get_or_create_vectorstore()
        await vectorstore.aadd_documents(documents)

    async def similarity_search(self, query: str, k: int = 4) -> list[Document]:
        vectorstore = await self._get_or_create_vectorstore()
        return await vectorstore.asimilarity_search(query, k=k)

    async def delete_collection_async(self):
        """Optional cleanup method"""
        if self.vectorstore:
            client = self._get_client_settings()
            client.delete_collection(self.collection_name)


def get_vector_store() -> VectorStore:
    """Factory function to return the appropriate VectorStore instance"""
    match settings.vector_store:
        case "chroma":
            return ChromaStore()
        case _:
            raise ValueError(f"Unsupported vector store: {settings.vector_store}")
