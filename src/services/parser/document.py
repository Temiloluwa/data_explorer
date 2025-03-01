from dataclasses import dataclass, field
from typing import Any, Dict, List, Optional


@dataclass
class SubDocument:
    identifier: Optional[str] = None  # e.g. "page=3" or a sublink URL
    content: Optional[Any] = None
    metadata: Dict[str, Any] = field(default_factory=dict)

    def __repr__(self):
        content_str = ""
        if isinstance(self.content, str):
            content_str = self.content[: min(100, len(self.content))]
        return (
            f"SubDocument(identifier={self.identifier}, "
            f"content={content_str!r}..., metadata={self.metadata})"
        )


@dataclass
class Document:
    identifier: str  # File path or URL
    source_type: str  # "file" or "url"
    content: Optional[Any] = None  # Raw content (loaded later)
    metadata: Dict[str, Any] = field(default_factory=dict)
    subdocuments: List[SubDocument] = field(default_factory=list)

    def __repr__(self):
        content_str = ""
        if isinstance(self.content, str):
            content_str = self.content[: min(100, len(self.content))]
        return (
            f"Document(identifier={self.identifier}, source_type={self.source_type}, "
            f"content={content_str!r}..., metadata={self.metadata}, "
            f"subdocuments={len(self.subdocuments)} subdocuments)"
        )
