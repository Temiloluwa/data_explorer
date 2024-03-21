import uvicorn
from fastapi import FastAPI
from auth import auth_router
from basic_qa import basic_qa_router
from utils import load_api_kwargs
from fastapi.middleware.cors import CORSMiddleware

origins = ["*"]

app = FastAPI()
app =  FastAPI(**load_api_kwargs())

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins
)

app.include_router(basic_qa_router, prefix="/api/v1", tags=["Basic QA"])

if __name__ == "__main__":
    uvicorn.run(app)
