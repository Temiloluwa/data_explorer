package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/data-explorer/internal/config" // Adjust
	dsc "github.com/yourusername/data-explorer/internal/datasources/common" // Adjust
	"github.com/yourusername/data-explorer/internal/datasources/local" // Adjust
	"github.com/yourusername/data-explorer/internal/datasources/web"  // Adjust
	"github.com/yourusername/data-explorer/internal/domain"       // Adjust
	"github.com/yourusername/data-explorer/internal/platform/logger"    // Adjust
	"github.com/yourusername/data-explorer/internal/qengine"      // Adjust
	// Import gRPC client packages if connecting to the server instead of direct engine use
	// pb "github.com/yourusername/data-explorer/api/proto/dataexplorer/v1"
	// "google.golang.org/grpc"
	// "google.golang.org/grpc/credentials/insecure"
)

var (
	// These flags will be bound by Cobra
	configPath string
	logLevel   string
	// Add other global CLI flags if needed
)

func main() {
	// Cobra handles errors and prints them, so we just exit if Execute fails.
	if err := rootCmd().Execute(); err != nil {
		// Cobra typically prints the error, but we can add extra logging if needed.
		// fmt.Fprintf(os.Stderr, "CLI execution failed: %v\n", err)
		os.Exit(1)
	}
}

// rootCmd represents the base command when called without any subcommands
func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "data-explorer-cli",
		Short: "Data Explorer CLI - Ask questions from various sources.",
		Long: `A command-line interface to interact with the Data Explorer engine.
You can ask questions and get answers synthesized from configured data sources.

Example:
  data-explorer-cli ask "What is the weather in Lagos?" --source web
  data-explorer-cli ask "Summarize chapter 3" --source local -o local.filename=mydoc.pdf`,
		// PersistentPreRunE runs before *any* command's RunE.
		// Good for initializing things needed by all commands, like logging/config,
		// but be aware flags might not be fully parsed depending on where they are defined.
		// Often better to initialize within the specific command's RunE.
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Example: You could initialize a basic logger here if needed early.
			return nil
		},
	}

	// Define persistent flags - available to this command and all subcommands
	cmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "Path to configuration file (default: ./configs/config.yaml or ./config.yaml)")
	cmd.PersistentFlags().StringVar(&logLevel, "log-level", "", "Override log level (e.g., debug, info, warn, error)")

	// Add subcommands
	cmd.AddCommand(askCmd())
	// Add other commands if needed (e.g., 'config view', 'source list', 'index local')

	return cmd
}

// askCmd represents the ask command
func askCmd() *cobra.Command {
	// Flags specific to the 'ask' command
	var sourcesFlag []string // Flag to specify sources, e.g., --source web --source local
	var optionsFlag []string // Flag for key=value options, e.g., --option local.max_files=10

	cmd := &cobra.Command{
		Use:   "ask \"Your question here?\"",
		Short: "Ask a question to the Data Explorer",
		Long: `Sends a question to the Data Explorer engine and prints the synthesized answer.
You can specify data sources to use and provide source-specific options.`,
		Example: `  data-explorer-cli ask "What is Go?"
  data-explorer-cli ask "Latest news about AI" -s web
  data-explorer-cli ask "Find section on 'widgets'" -s local -o local.path=./docs/manual.txt`,
		Args: cobra.ExactArgs(1), // Require exactly one argument (the question)
		// RunE is executed when the 'ask' command is called.
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get the question from the arguments
			questionQuery := args[0]

			// --- Setup Logger and Config (inside command RunE to ensure flags are parsed) ---
			// Use a temporary logger for config loading itself
			tempLogForConfig, _ := logger.New(logger.Config{Level: "info", Format: "console"})
			cfg, err := config.LoadConfig(tempLogForConfig, configPath)
			if err != nil {
				// Use fmt.Errorf for errors returned from RunE
				return fmt.Errorf("failed to load configuration: %w", err)
			}

			// Override log level from flag if provided *after* loading config
			if logLevel != "" {
				cfg.Log.Level = logLevel
			}
			log, err := logger.New(cfg.Log)
			if err != nil {
				return fmt.Errorf("failed to initialize logger: %w", err)
			}
            // Defer logger sync for this command execution
            defer func() {
                if loggerImpl, ok := log.(*logger.zapLogger); ok { loggerImpl.sugar.Sync() }
            }()


			// --- Initialize Engine (Directly for now) ---
			// This duplicates setup logic from cmd/server.
			// TODO: Consider refactoring setup into internal/app or connect to the server via gRPC.
			// Direct engine use is simpler for a basic CLI.
			var availableSources []dsc.DataSource
			if cfg.DataSources.WebSearch.Enable {
				webSource, err := web.NewWebSearchSource(log, cfg.DataSources.WebSearch.APIKey)
                if err != nil { log.Warnf("Failed to init web source for CLI: %v", err) }
				if webSource != nil { availableSources = append(availableSources, webSource) }
			}
			if cfg.DataSources.LocalFiles.Enable {
				localSource, err := local.NewLocalFileSource(log, cfg.DataSources.LocalFiles.WatchPaths)
                if err != nil { log.Warnf("Failed to init local source for CLI: %v", err) }
                if localSource != nil { availableSources = append(availableSources, localSource) }
			}
			// Add other sources...

			if len(availableSources) == 0 {
				log.Warnf("CLI mode: No data sources enabled or initialized! Results may be empty.")
			}

			engine, err := qengine.NewEngine(log, cfg, availableSources)
			if err != nil {
				return fmt.Errorf("failed to initialize QA engine: %w", err)
			}

			// --- Prepare Question ---
			domainQuestion := domain.Question{
				Query:   questionQuery,
				Options: parseOptions(optionsFlag), // Use helper to parse options
				Sources: make([]domain.DataSourceSpec, len(sourcesFlag)),
			}
            // Map source type strings from flags to DataSourceSpec
            for i, sType := range sourcesFlag {
                domainQuestion.Sources[i] = domain.DataSourceSpec{Type: sType}
            }


			log.Infof("Asking question via CLI: Query length=%d", len(domainQuestion.Query))
            log.Debugf("Full CLI Query: %s, Sources: %+v, Options: %+v",
                domainQuestion.Query, domainQuestion.Sources, domainQuestion.Options)


			// --- Call Engine ---
			// Use command's context which handles SIGINT etc. automatically with Cobra v1.7+
			ctx := cmd.Context()
			answer, err := engine.ProcessQuestion(ctx, domainQuestion)
			if err != nil {
				// Engine errors are logged internally, return a user-friendly message via error
				log.Errorf("Engine processing failed: %v", err) // Log the detailed error
				return fmt.Errorf("failed to get answer (see logs for details)") // User-facing error
			}
            if answer == nil {
                 // Should not happen if engine returns error on failure
                 log.Errorf("Engine returned nil answer without error")
                 return fmt.Errorf("internal error: received nil answer from engine")
            }

			// --- Print Answer to Stdout ---
			fmt.Println("Answer:")
			fmt.Println(answer.Text)
			fmt.Println("\nSources Used:")
			if len(answer.SourcesUsed) > 0 {
				for i, src := range answer.SourcesUsed {
					fmt.Printf("  %d. [%s] %s (Score: %.2f)\n", i+1, src.SourceType, src.Identifier, src.RelevanceScore)
					// Use a helper function for consistent truncation
					fmt.Printf("     Snippet: %s...\n", truncateString(src.ContentSnippet, 150))
				}
			} else {
				fmt.Println("  (No specific sources cited)")
			}
			// Print non-fatal errors/warnings from the answer object
			if answer.Error != nil {
				fmt.Printf("\nWarning/Error during processing: %v\n", answer.Error)
			}

			// Return nil to indicate success to Cobra
			return nil
		},
	}

	// Add flags specific to the 'ask' command
	// Use StringSlice for sources as a type can be specified multiple times (though our current engine logic dedupes)
	cmd.Flags().StringSliceVarP(&sourcesFlag, "source", "s", []string{}, "Specify data source types to query (e.g., 'web', 'local') (default: all enabled)")
	// Use StringArray for options as keys might be repeated (though mapstructure usually takes the last one)
	cmd.Flags().StringArrayVarP(&optionsFlag, "option", "o", []string{}, "Pass key=value options to data sources (e.g., 'local.max_files=10')")

	return cmd
}

// Helper function to parse --option flags (key=value format) into a map.
func parseOptions(opts []string) map[string]string {
	parsed := make(map[string]string)
	for _, o := range opts {
		parts := strings.SplitN(o, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key != "" { // Avoid empty keys
				parsed[key] = value
			}
		}
		// Silently ignore malformed options or log a warning?
	}
	return parsed
}

// Helper function to truncate strings safely (considers runes).
func truncateString(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen])
}
