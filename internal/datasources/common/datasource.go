package common

import (
	"context"
	"github.com/yourusername/data-explorer/internal/domain" // Adjust path if needed
)

// DataSource is the interface implemented by all data source connectors.
type DataSource interface {
	// Type returns a unique identifier for the data source type (e.g., "web", "local").
	// This should match configuration keys.
	Type() string

	// FetchData retrieves data relevant to the question from the specific source.
	// The 'identifier' parameter here is typically derived from the question or config,
	// representing what to fetch (e.g., search query for web, path for local).
	// Options can provide source-specific parameters.
	FetchData(ctx context.Context, question domain.Question, identifier string, options map[string]string) ([]domain.FetchedData, error)
}
