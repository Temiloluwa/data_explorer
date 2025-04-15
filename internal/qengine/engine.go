package qengine

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync" // For concurrent data fetching

	"github.com/yourusername/data-explorer/internal/config"       // Adjust path
	"github.com/yourusername/data-explorer/internal/datasources/common" // Adjust path
	"github.com/yourusername/data-explorer/internal/domain"       // Adjust path
	"github.com/yourusername/data-explorer/internal/platform/logger"    // Adjust path
)

// QEngine implements the Engine interface.
type QEngine struct {
	log         logger.Logger // Use your logger interface
	cfg         *config.Config
	dataSources map[string]common.DataSource // Map type string to implementation
}

// NewEngine creates a new QA Engine instance.
func NewEngine(log logger.Logger, cfg *config.Config, sources []common.DataSource) (Engine, error) {
	if len(sources) == 0 {
		log.Warnf("QA Engine initialized with zero data sources.")
	}
	dsMap := make(map[string]common.DataSource)
	for _, s := range sources {
		if _, exists := dsMap[s.Type()]; exists {
			// Depending on config, this might be an error or just a warning
			log.Warnf("Duplicate data source type registered: %s. Overwriting.", s.Type())
		}
		dsMap[s.Type()] = s
		log.Infof("Registered data source: %s", s.Type())
	}

	return &QEngine{
		log:         log.With("component", "qengine"),
		cfg:         cfg,
		dataSources: dsMap,
	}, nil
}

// ProcessQuestion is the core logic (currently with placeholder synthesis).
func (e *QEngine) ProcessQuestion(ctx context.Context, question domain.Question) (*domain.Answer, error) {
	e.log.Infof("Processing question: Query='%s', Sources=%+v, Options=%+v", question.Query, question.Sources, question.Options)

	// 1. Determine which data sources to query
	sourcesToQuery := e.selectSources(question)
	if len(sourcesToQuery) == 0 {
		e.log.Warnf("No applicable data sources found or selected for query.")
		// Return a specific answer indicating no sources were queried
		return &domain.Answer{
			Text: "I cannot answer this question as no relevant data sources are available or selected.",
		}, nil
	}

	// 2. Fetch data from selected sources concurrently
	fetchedDataChan := make(chan []domain.FetchedData, len(sourcesToQuery))
	errorChan := make(chan error, len(sourcesToQuery))
	var wg sync.WaitGroup

	for _, ds := range sourcesToQuery {
		wg.Add(1)
		go func(dataSource common.DataSource) {
			defer wg.Done()
			// TODO: Determine the appropriate identifier for this source based on the question/config.
			// For web search, it's likely question.Query. For local files, it might be question.Query or a path.
			sourceIdentifier := question.Query // Placeholder - needs refinement

			e.log.Debugf("Querying source '%s' with identifier: %s", dataSource.Type(), sourceIdentifier)
			// Pass source-specific options if available in question.Options
			// TODO: Filter question.Options relevant to this dataSource.Type()
			opts := question.Options

			fetched, err := dataSource.FetchData(ctx, question, sourceIdentifier, opts)
			if err != nil {
				// Don't fail the whole request for one source error, log it and send to error chan
				e.log.Errorf("Failed to fetch data from source '%s': %v", dataSource.Type(), err)
				errorChan <- fmt.Errorf("source '%s': %w", dataSource.Type(), err)
				return // Don't send to fetchedDataChan
			}
			if len(fetched) > 0 {
				fetchedDataChan <- fetched
			}
			e.log.Debugf("Fetched %d items from source '%s'", len(fetched), dataSource.Type())
		}(ds)
	}

	// Wait for all fetches to complete
	wg.Wait()
	close(fetchedDataChan)
	close(errorChan)

	// Collect results and errors
	var allFetchedData []domain.FetchedData
	var fetchErrors []error
	for data := range fetchedDataChan {
		allFetchedData = append(allFetchedData, data...)
	}
	for err := range errorChan {
		fetchErrors = append(fetchErrors, err)
	}

	// Combine non-fatal errors to potentially include in the answer
	var combinedFetchError error
	if len(fetchErrors) > 0 {
		errorStrings := make([]string, len(fetchErrors))
		for i, err := range fetchErrors {
			errorStrings[i] = err.Error()
		}
		combinedFetchError = errors.New(strings.Join(errorStrings, "; "))
		e.log.Warnf("Encountered errors fetching data: %v", combinedFetchError)
	}

	if len(allFetchedData) == 0 && len(fetchErrors) == len(sourcesToQuery) {
		e.log.Errorf("All data sources failed for query: %s", question.Query)
		return nil, fmt.Errorf("failed to fetch data from all sources: %w", combinedFetchError) // Return fatal error
	}
    if len(allFetchedData) == 0 {
         e.log.Warnf("No data fetched for query: %s (some sources might have failed)", question.Query)
         return &domain.Answer{
             Text: "Sorry, I couldn't find any information for your query from the available sources.",
             Error: combinedFetchError, // Include errors if any occurred
         }, nil
    }

	// 3. Process/Analyze Fetched Data (This is the complex part - Placeholder)
	// TODO: Implement actual NLP/LLM logic here.
	// - Filter/rank allFetchedData based on relevance to question.Query.
	// - Extract key information.
	// - Synthesize an answer using the extracted info (e.g., call an LLM with context).
	// - Populate sourcesUsed based *only* on the data actually used for the final answer.

	// --- Placeholder Synthesis Logic ---
	rankedData := allFetchedData // Placeholder: Assume all fetched data is relevant for now
	var sourcesUsed []domain.SourceDocument
	var combinedContent strings.Builder
	combinedContent.WriteString(fmt.Sprintf("Based on %d potentially relevant piece(s) of information:\n", len(rankedData)))

	for i, data := range rankedData {
        // Limit the number of sources cited in the placeholder for brevity
        if i >= 5 {
             combinedContent.WriteString(fmt.Sprintf("- ...and %d more sources.\n", len(rankedData)-i))
             break
        }
		snippet := truncate(string(data.Content), 100) // Use helper
		combinedContent.WriteString(fmt.Sprintf("- From %s (%s): %s...\n", data.SourceType, data.Identifier, snippet))
		sourcesUsed = append(sourcesUsed, domain.SourceDocument{
			SourceType:     data.SourceType,
			Identifier:     data.Identifier,
			ContentSnippet: snippet, // Store the truncated snippet or a more relevant one
			RelevanceScore: 0.5, // Placeholder score - should come from ranking step
		})
	}
	answerText := fmt.Sprintf("Placeholder answer for '%s'.\nDetails:\n%s", question.Query, combinedContent.String())
	// --- End Placeholder Synthesis ---


	// 4. Format the answer
	answer := &domain.Answer{
		Text:        answerText,
		SourcesUsed: sourcesUsed, // Use the sources actually contributing to the answer
		Error:       combinedFetchError, // Include non-fatal fetch errors
	}

	e.log.Infof("Successfully processed question: %s", question.Query)
	return answer, nil
}

// selectSources determines which data sources to query based on the question and config.
func (e *QEngine) selectSources(question domain.Question) []common.DataSource {
	// If specific sources are requested in the question, use those.
	if len(question.Sources) > 0 {
		selected := make([]common.DataSource, 0, len(question.Sources))
		requestedTypes := make(map[string]bool) // Track types to avoid duplicates if specified multiple times
		for _, spec := range question.Sources {
			if ds, exists := e.dataSources[spec.Type]; exists && !requestedTypes[spec.Type] {
				// TODO: Add logic here if spec.Identifier needs to match a specific configured instance
				selected = append(selected, ds)
                requestedTypes[spec.Type] = true
				e.log.Debugf("Selecting source explicitly requested: %s", spec.Type)
			} else if !exists {
                e.log.Warnf("Requested data source type '%s' is not available/configured.", spec.Type)
            }
		}
		return selected
	}

	// Otherwise, use all registered/configured sources (default behavior).
	// TODO: Add more sophisticated selection logic if needed (e.g., based on query analysis).
	allSources := make([]common.DataSource, 0, len(e.dataSources))
	for _, ds := range e.dataSources {
		allSources = append(allSources, ds)
	}
	e.log.Debugf("No specific sources requested, selecting all available: %d", len(allSources))
	return allSources
}


// Helper to truncate strings (example)
func truncate(s string, maxLen int) string {
	// Consider rune count for multi-byte characters if necessary
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
