package local

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io/fs" // Use io/fs for walking directories
	"os"
	"path/filepath"
	"strings"
	"sync" // Added for potentially parallel file processing

	"github.com/yourusername/data-explorer/internal/datasources/common"
	"github.com/yourusername/data-explorer/internal/domain"
	"github.com/yourusername/data-explorer/internal/platform/logger"
	// Add imports for file parsing libraries if needed (e.g., PDF, DOCX)
	// "github.com/ledongthuc/pdf" // Example PDF library
)

const DataSourceTypeLocal = "local"

// LocalFileSource implements common.DataSource for reading local files.
type LocalFileSource struct {
	log        logger.Logger
	watchPaths []string
	// Potentially add an index here (e.g., Bleve, map[string][]string) for faster searching
	// index      map[string]string // Example simple index: path -> content (inefficient for large files)
	// indexMutex sync.RWMutex      // Mutex for concurrent index access if used
}

// NewLocalFileSource creates a new instance for local file reading.
func NewLocalFileSource(log logger.Logger, watchPaths []string) (common.DataSource, error) {
	if len(watchPaths) == 0 {
		log.Warnf("Local file source initialized with no watch paths. It will not find any files.")
	} else {
		log.Infof("Initializing local file source with watch paths: %v", watchPaths)
		// TODO: Validate paths exist and are directories?
	}

	source := &LocalFileSource{
		log:        log.With("datasource", DataSourceTypeLocal),
		watchPaths: watchPaths,
		// Initialize index if using one
	}

	// TODO: Implement background indexing if needed.
	// go source.startIndexing()

	return source, nil
}

// Type returns the identifier for this data source.
func (s *LocalFileSource) Type() string {
	return DataSourceTypeLocal
}

// FetchData searches configured local files for relevant content.
// The 'identifier' could be used to specify a sub-path or pattern, but here we assume
// it's related to the search query itself (similar to question.Query).
// Options could specify things like "max_files_to_scan", "file_extensions", etc.
func (s *LocalFileSource) FetchData(ctx context.Context, question domain.Question, identifier string, options map[string]string) ([]domain.FetchedData, error) {
	// Use question.Query as the primary search term within files.
	searchQuery := question.Query
	if searchQuery == "" {
		return nil, errors.New("local search requires a non-empty query")
	}
	s.log.Debugf("Searching local files in paths %v for query: '%s'", s.watchPaths, searchQuery)

	// TODO: Use options map to get parameters like max results, allowed extensions etc.
	// maxResults := parseOptionInt(options, "local.max_results", 10)
	// allowedExts := parseOptionStringSlice(options, "local.extensions", []string{".txt", ".md"})

	var results []domain.FetchedData
	var resultMutex sync.Mutex // Mutex to protect concurrent appends to results slice
	var wg sync.WaitGroup      // WaitGroup for concurrent file processing

	processedFileCount := 0 // Counter for logging/limits

	for _, dir := range s.watchPaths {
		// Check context before starting walk
		if ctx.Err() != nil {
			s.log.Infof("Context cancelled before walking directory: %s", dir)
			break // Stop processing further directories
		}

		// Walk the directory using io/fs.WalkDir for better error handling
		err := fs.WalkDir(os.DirFS(dir), ".", func(path string, d fs.DirEntry, err error) error {
			// Check context periodically within the walk function
			select {
			case <-ctx.Done():
				s.log.Infof("Context cancelled during file walk.")
				return ctx.Err() // Stop walking immediately
			default:
				// Continue processing
			}

			// Handle errors accessing path (permissions etc.)
			if err != nil {
				s.log.Warnf("Error accessing path %q during walk: %v", filepath.Join(dir, path), err)
				// Decide whether to skip this entry or stop the walk
				// Returning the error will stop the walk for this directory.
				// Returning fs.SkipDir will skip the contents of a directory on error.
				// Returning nil will log the error and continue.
				return nil // Log and continue for now
			}

			// Skip directories
			if d.IsDir() {
				return nil
			}

			// Skip files that don't match allowed extensions (example)
			// if !s.isAllowedExtension(path, allowedExts) {
			// 	return nil
			// }

            // --- Process File Concurrently ---
            // Increment WaitGroup counter *before* launching goroutine
            wg.Add(1)
            go func(fullPath string) {
                defer wg.Done() // Decrement counter when goroutine finishes

                // Check context again within the goroutine
                if ctx.Err() != nil { return }

                // TODO: Implement actual file reading and searching logic
                // 1. Read file content (os.ReadFile or streaming for large files).
                // 2. Handle different file types (plain text, PDF, DOCX...).
                // 3. Search content for relevance to 'searchQuery' (case-insensitive string contains is basic).
                // 4. If relevant, create domain.FetchedData.

                // --- Placeholder: Simple text file search ---
                if strings.HasSuffix(strings.ToLower(fullPath), ".txt") || strings.HasSuffix(strings.ToLower(fullPath), ".md") {
                    contentBytes, readErr := os.ReadFile(fullPath)
                    if readErr != nil {
                        s.log.Warnf("Failed to read file %s: %v", fullPath, readErr)
                        return // Skip this file on read error
                    }
                    content := string(contentBytes)
                    // Basic case-insensitive search
                    if strings.Contains(strings.ToLower(content), strings.ToLower(searchQuery)) {
                        s.log.Debugf("Found match in file: %s", fullPath)
                        // Extract a snippet around the match (more advanced logic needed)
                        snippet := s.extractSnippet(content, searchQuery, 200)

                        fetched := domain.FetchedData{
                            SourceType: DataSourceTypeLocal,
                            Identifier: fullPath, // Use full path as identifier
                            Content:    contentBytes, // Could store snippet instead to save memory
                            Metadata:   map[string]interface{}{"filename": filepath.Base(fullPath)},
                        }
                        // Append result safely using mutex
                        resultMutex.Lock()
                        results = append(results, fetched)
                        processedFileCount++
                        resultMutex.Unlock()
                    }
                }
                // --- End Placeholder ---

            }(filepath.Join(dir, path)) // Pass the full path to the goroutine

			return nil // Continue walking
		}) // End of WalkDir

		// Handle error returned by WalkDir itself (e.g., context cancelled)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				s.log.Infof("File walk interrupted for directory %s: %v", dir, err)
				break // Stop processing further directories if context is done
			}
			s.log.Errorf("Error walking directory %s: %v", dir, err)
			// Decide if this error is critical for the whole operation
		}
	} // End of loop over watchPaths

    // Wait for all file processing goroutines to complete
    wg.Wait()

	// Check context one last time after waiting
	if ctx.Err() != nil {
		s.log.Infof("Context cancelled after file processing completed or during wait.")
		// Return partial results along with context error
		return results, ctx.Err()
	}


	s.log.Infof("Found %d potentially relevant local file(s) for query: '%s'", len(results), searchQuery)
	// TODO: Apply ranking/filtering to 'results' before returning if needed.
	return results, nil
}

// isAllowedExtension checks if the file path has one of the allowed extensions.
// func (s *LocalFileSource) isAllowedExtension(path string, allowedExts []string) bool {
// 	ext := strings.ToLower(filepath.Ext(path))
// 	for _, allowed := range allowedExts {
// 		if ext == allowed {
// 			return true
// 		}
// 	}
// 	return false
// }

// extractSnippet finds the search query and returns surrounding text. Very basic example.
func (s *LocalFileSource) extractSnippet(content, query string, length int) string {
	lowerContent := strings.ToLower(content)
	lowerQuery := strings.ToLower(query)

	index := strings.Index(lowerContent, lowerQuery)
	if index == -1 {
		// Query not found, return beginning of content
		if len(content) > length {
			return content[:length]
		}
		return content
	}

	start := index - (length / 2) + (len(query) / 2)
	if start < 0 {
		start = 0
	}
	end := start + length
	if end > len(content) {
		end = len(content)
		// Adjust start if possible to still get desired length
		start = end - length
		if start < 0 {
			start = 0
		}
	}

    // Ensure start/end are within bounds after adjustments
    if start < 0 { start = 0 }
    if end > len(content) { end = len(content) }
    if start >= end { // Should not happen with logic above, but safeguard
        if len(content) > length { return content[:length] }
        return content
    }

	return content[start:end]
}


// TODO: Implement helper functions for parsing options if needed
// func parseOptionInt(options map[string]string, key string, defaultValue int) int { ... }
// func parseOptionStringSlice(options map[string]string, key string, defaultValue []string) []string { ... }
