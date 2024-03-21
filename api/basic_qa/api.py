from fastapi import APIRouter, HTTPException
from typing import List, Dict
from .models import Document, DocumentChunk, QuestionWithDocumentChunks, Answer
from .service import check_if_document_has_been_ingested
import uuid

router = APIRouter()

@router.put("/ingest_document", response_model=str)
def ingest_document(document: Document):
    
    if check_if_document_has_been_ingested(document):
        pass
    
    return f"Document {document.document_id} ingested successfully"


@router.get("/retrieve_documents/{document_id}", response_model=Document)
def retrieve_documents(document_id: str):
    if document_id not in documents:
        raise HTTPException(status_code=404, detail="Document not found")
    return documents[document_id]

@router.post("/ask_question", response_model=str)
def ask_question(question: Question):
    # Handle the question and associated document IDs
    # For demonstration purposes, just return a success message
    return "Question received and processed successfully"


@router.delete("/cleanup_document/{document_id}", response_model=str)
def cleanup_document(document_id: str):
    if document_id not in documents:
        raise HTTPException(status_code=404, detail="Document not found")
    del documents[document_id]
    return f"Document {document_id} cleaned up successfully"
