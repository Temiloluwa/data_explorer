from fastapi import APIRouter, HTTPException
from typing import List, Dict
from data_models import Document, DocumentChunk, QuestionWithDocumentChunks, Answer
from .handler import check_if_document_has_been_ingested
import uuid

router = APIRouter()

@router.put("/ingest_document", response_model=str)
def ingest_document(document: Document):
    
    if check_if_document_has_been_ingested(document):
        pass
    
    return f"Document {document.document_id} ingested successfully"
