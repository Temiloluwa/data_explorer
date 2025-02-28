from chains.qa_chain import get_qa_chain
from config import settings
from core.data_models import QueryResult
from core.exceptions import RetrievalError
from langchain_core.runnables import RunnableLambda
from retrievers.vector_store import VectorStore


class QueryService:
    def __init__(self, vector_store: VectorStore):
        self.vector_store = vector_store
        self.qa_chain = get_qa_chain(vector_store)

    async def query(self, question: str) -> QueryResult:
        try:
            result = await self.qa_chain.ainvoke({"question": question})
            return result
            #return QueryResult(**result)
        except Exception as e:
            raise RetrievalError(f"Query failed: {str(e)}") from e
