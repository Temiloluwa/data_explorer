import asyncio
from sources import FileSource, URLSource, DocumentSource
from content_processors import get_processor_for_document
from document import Document

async def process_documents(source_identifier: str, level: str):
    """
    Loads documents from the given source and processes each one using
    the appropriate content processor based on the document type and level.
    """
    # Select the source based on whether the identifier is a URL or file path.
    if source_identifier.startswith("http"):
        source: DocumentSource = URLSource()
    else:
        source: DocumentSource = FileSource()
    
    # Load lightweight Document objects (without content loaded)
    documents = await source.load(source_identifier)
    
    processed_documents = []
    for doc in documents:
        # Choose the processor based on document type and desired level.
        processor = get_processor_for_document(doc, level)
        processed_doc = await processor.process(doc)
        processed_documents.append(processed_doc)
    
    return processed_documents

async def main():
    # Set source_identifier to a folder path (for file sources) or a URL (for URL sources)
    source_identifier = "/Users/temmie/codes-and-scripts/projects/data_explorer/data"  # Change to your local folder or URL.
    level = "basic"  # Options: "basic", "intermediate", "advanced"
    
    processed_docs = await process_documents(source_identifier, level)
    
    for doc in processed_docs:
        print("Document:", doc.identifier)
        print("Metadata:", doc.metadata)
        content_preview = doc.content[:100] if isinstance(doc.content, str) else str(doc.content)[:100]
        print("Content Preview:", content_preview)
        print("Subdocuments:")
        for sub in doc.subdocuments:
            sub_preview = sub.content[:50] if isinstance(sub.content, str) else str(sub.content)[:50]
            print(f"   {sub.identifier}: {sub_preview}")
        print("----")

if __name__ == "__main__":
    asyncio.run(main())
