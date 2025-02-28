import uvicorn
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from rag import rag_router
from ingestion import ingestion_router
from backend.utils import *

app = FastAPI(**load_api_kwargs())

# CORS settings
origins = ["*"]
app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Include routers
app.include_router(ingestion_router, prefix="/api/v1/ingest", tags=["ingest"])
app.include_router(rag_router, prefix="/api/v1/rag", tags=["rag"])

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)
