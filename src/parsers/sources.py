import os
from abc import ABC, abstractmethod
from typing import List

from document import Document


class DocumentSource(ABC):
    @abstractmethod
    async def load(self, source_identifier: str) -> List[Document]:
        pass


class FileSource(DocumentSource):
    async def load(self, source_identifier: str) -> List[Document]:
        documents = []
        for dirpath, _, filenames in os.walk(source_identifier):
            for filename in filenames:
                file_path = os.path.join(dirpath, filename)
                stat = os.stat(file_path)
                metadata = {"size": stat.st_size, "modified_time": stat.st_mtime}
                documents.append(
                    Document(
                        identifier=file_path, source_type="file", metadata=metadata
                    )
                )
        return documents


class URLSource(DocumentSource):
    async def load(self, source_identifier: str) -> List[Document]:
        return [Document(identifier=source_identifier, source_type="url", metadata={})]
