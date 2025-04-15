package domain

// SourceDocument represents evidence used for the answer, derived from FetchedData.
type SourceDocument struct {
	SourceType     string  // e.g., "web", "local_file", "db_row"
	Identifier     string  // e.g., URL, file path, table/row id
	ContentSnippet string  // Relevant snippet of the source content
	RelevanceScore float32 // How relevant this source was (optional)
}

// FetchedData represents raw data retrieved from a source before processing.
type FetchedData struct {
	SourceType string                 // e.g., "web", "local_file"
	Identifier string                 // e.g., URL, file path
	Content    []byte                 // Raw content
	Metadata   map[string]interface{} // Additional info (e.g., title, last modified)
}
