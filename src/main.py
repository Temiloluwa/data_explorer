import asyncio

from retrievers import get_vector_store
from services.ingestion_service import IngestionService
from services.query_service import QueryService


async def main():
    vector_store = get_vector_store()
    ingestion_service = IngestionService(vector_store)
    query_service = QueryService(vector_store)

    await ingestion_service.ingest("../../research/data/pdfs/well-intervention")
    result = await query_service.query("What is RAG?")
    print(f"Answer: {result.answer}")


if __name__ == "__main__":
    asyncio.run(main())
