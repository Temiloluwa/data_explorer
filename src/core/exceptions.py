class RAGException(Exception):
    """Base exception for RAG application"""


class IngestionError(RAGException):
    """Error during document ingestion"""


class RetrievalError(RAGException):
    """Error during document retrieval"""
