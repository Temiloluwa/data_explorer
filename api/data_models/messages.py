from pydantic import BaseModel
from typing import List
from data_models.documents import DocumentChunk

class Message(BaseModel):
    message_id: str
    role: str
    text: str
    timestamp: str

class QuestionWithDocumentChunks(Message):
    document_chunks: List[DocumentChunk]

class Answer(Message):
    document_chunk_ids: List[str]