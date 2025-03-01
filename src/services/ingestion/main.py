import asyncio
import os

from config import AgentMode, settings
from core.data_models import Document
from core.exceptions import IngestionError
from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain_community.document_loaders import PyPDFLoader
from services.core.vector_store import VectorStore
from tqdm.asyncio import tqdm


class IngestionService:
    splitter = None

    def __init__(self, vector_store):
        self.vector_store = vector_store
        self.setup()

    def setup(self):
        match settings.agent_mode:
            case "expert":
                raise NotImplementedError("Expert mode setup is not implemented.")
            case "advanced":
                raise NotImplementedError("Advanced mode setup is not implemented.")
            case "intermediate":
                raise NotImplementedError("Intermediate mode setup is not implemented.")
            case _:
                self._basic_setup()

    def _basic_setup(self):
        self.splitter = RecursiveCharacterTextSplitter(
            chunk_size=settings.embedding_config.chunk_size,
            chunk_overlap=settings.embedding_config.chunk_overlap,
        )

    async def _collect_pdf_files(self, data_path: str):
        if not os.path.isdir(data_path):
            raise IngestionError(f"Invalid directory path: {data_path}")

        pdf_files = []
        for root, _, files in os.walk(data_path):
            for file in files:
                if file.endswith(".pdf"):
                    pdf_files.append(os.path.join(root, file))
        return pdf_files

    async def _process_file(self, file: str, documents: list):
        try:
            docs = PyPDFLoader(file).load()
            split_docs = self.splitter.split_documents(docs)
            documents.extend(split_docs)
            print(f"Processed file: {file}")
        except Exception as e:
            print(f"Error processing {file}: {e}")

    async def _basic_ingestion(self, data_path: str):
        if not self.splitter:
            raise IngestionError("Splitter is not initialized. Call 'setup()' first.")

        documents = []
        try:
            pdf_files = await self._collect_pdf_files(data_path)
            tasks = [
                self._process_file(file, documents)
                for file in tqdm(pdf_files, desc="Commenced Processing File", unit="file")
            ]
            await asyncio.gather(*tasks)

            if documents:
                print(f"Commened Ingesting : {len(documents)} documents")
                await self.vector_store.add_documents(documents)

        except Exception as e:
            raise IngestionError(f"Ingestion failed: {str(e)}") from e

    async def ingest(self, data_path: str):
        if settings.agent_mode == AgentMode.BASIC:
            await self._basic_ingestion(data_path)
        else:
            raise IngestionError(f"Unsupported agent mode: {settings.agent_mode}")
