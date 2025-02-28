from config import settings
from langchain_core.output_parsers import StrOutputParser
from langchain_core.prompts import ChatPromptTemplate
from langchain_openai import ChatOpenAI
from retrievers.vector_store import VectorStore


def get_qa_chain(vector_store: VectorStore):
    prompt = ChatPromptTemplate.from_template(
        "Answer the question based only on the context:\n\n{context}\n\nQuestion: {question}"
    )

    llm = ChatOpenAI(model=settings.llm_model)

    return (
        {
            "context": lambda x: vector_store.similarity_search(x["question"]),
            "question": lambda x: x["question"],
        }
        | prompt
        | llm
        | StrOutputParser()
    )
