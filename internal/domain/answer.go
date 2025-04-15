package domain

// Answer holds the result of a query processed by the engine.
type Answer struct {
	// Text is the synthesized answer string.
	Text string
	// SourcesUsed lists the evidence documents that contributed to the answer.
	SourcesUsed []SourceDocument
	// Error holds any non-fatal error or warning encountered during processing
	// (e.g., one data source failed but others succeeded). Fatal errors should
	// be returned directly by the engine method.
	Error error
}
