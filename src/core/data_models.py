from pydantic import BaseModel


class Document(BaseModel):
    content: str
    metadata: dict = {}


class QueryResult(BaseModel):
    answer: str
    sources: list[Document] = []
    confidence: float = 0.0
