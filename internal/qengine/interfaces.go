package qengine

import (
	"context"
	"github.com/yourusername/data-explorer/internal/domain" // Adjust path if needed
)

// Engine defines the interface for the core Question Answering logic.
type Engine interface {
	// ProcessQuestion takes a question, queries relevant data sources,
	// synthesizes an answer, and returns it.
	// Errors returned should indicate a failure to process the question overall.
	// Non-fatal issues (like one source failing) might be included in the Answer.Error field.
	ProcessQuestion(ctx context.Context, question domain.Question) (*domain.Answer, error)
}
