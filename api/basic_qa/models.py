from pydantic import BaseModel
from typing import List

class Document(BaseModel):
    document_id: str #index of document
    filename: str
    number_of_chunks: int
    metadata: dict | None = None

class DocumentChunk(BaseModel):
    document_chunk_id: str
    chunk_id: str
    text: str
    metadata: dict | None = None

class Message(BaseModel):
    message_id: str
    role: str
    text: str
    timestamp: str

class QuestionWithDocumentChunks(Message):
    document_chunks: List[DocumentChunk]

class Answer(Message):
    document_chunk_ids: List[str]