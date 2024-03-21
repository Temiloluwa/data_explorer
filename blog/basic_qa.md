# Basic Question and Answering with LLMs

## Introduction
With the proliferation of the Internet, came big data generation. Consequent applications
like Web Search, Recommendation Systems, Data Analytics are an attempt to distill information and provide answers, as accurate as possible to users. The most recent approach employs Retrieval-Augmented-Generation (RAG) deploys large language models (LLM) like ChatGPT to interprete documents retrieved from a database. This system can be broken down into three simple steps

1. Retrieval Step: A user provides a Question or a query in technical terms, then a number of relevant documents are retrieved from a database. This is analogous to the user experience from Google Search.
2. Augmentation Step: Further processing can be employed to the retrieved documents for example, reranking, summarization etc.
3. Generation Step: The augmented information is combined with the user query in a prompt and the LLM provides an answer

LLMs are fraut with challenges like hallucination and bias so the answers are not bound to be accurate or aligned to the user's interest. The goal of this project is to shed light on various parts of the Question Answering System, and explore improvements in further Projects.

## Retrieval System