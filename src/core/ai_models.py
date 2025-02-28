from langchain_huggingface.embeddings import HuggingFaceEmbeddings
from langchain_openai import ChatOpenAI, OpenAIEmbeddings


class AIModels:
    @classmethod
    def get_chat_model(cls, settings):
        match settings.llm_model:
            case "gpt-3.5-turbo-0125":
                return cls.openai_chat_model(settings)
            case _:
                raise ValueError(f"Unknown model: {settings.llm_model}")

    @classmethod
    def get_embeddings_model(cls, settings):
        match settings.embedding_model:
            case "sentence-transformers/all-MiniLM-L6-v2":
                return cls.huggingface_embeddings_model(settings)
            case _:
                raise ValueError(f"Unknown model: {settings.embedding_model}")


    @classmethod
    def openai_chat_model(cls, settings):
        return ChatOpenAI(
            model=settings.llm_model,
            temperature=settings.llm_config.temperature,
            max_tokens=settings.llm_config.max_tokens,
        )

    @classmethod
    def openai_embeddings_model(cls, settings):
        return OpenAIEmbeddings(model=settings.embedding_model)

    @classmethod
    def huggingface_chat_model(cls, settings):
        pass
        
    @classmethod
    def huggingface_embeddings_model(cls, settings):
        return HuggingFaceEmbeddings(model_name=settings.embedding_model)

