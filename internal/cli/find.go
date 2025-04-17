package cli

import (
	"context"
	"errors"
	"fmt"
	"github.com/sanspareilsmyn/historai/internal/config"
	"github.com/sanspareilsmyn/historai/internal/history"
	"github.com/sanspareilsmyn/historai/internal/llm"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	defaultFindHistoryLimit = 300
)

// findCmd represents the find command
var findCmd = &cobra.Command{
	Use:   "find \"<natural language query>\"",
	Short: "Find commands in shell history using a natural language query",
	Long: `Searches your shell history file (currently Zsh: ~/.zsh_history) using an LLM
to find commands that match the provided natural language description.

You can limit the scope of the history search using the flag:
  --limit / -n : How many recent entries to consider (default: 300).

Example:
  historai find "how I listed files sorted by size last month"
  historai find --limit 500 "the ssh command to connect to the webserver"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Debug("Executing find command")

		query := args[0]
		if query == "" {
			return errors.New("query cannot be empty")
		}
		logger.Debug("Received query", zap.String("query", query))

		// 1. Parse and validate flags
		limit, err := parseFindFlags(cmd)
		if err != nil {
			return err
		}

		// 2. Execute the core finding logic
		result, err := runFind(logger, query, limit)
		if err != nil {
			return err
		}

		// 3. Print Header to Stderr & Result to Stdout
		err = printCommandOutput(
			logger,
			result,
			"--- Found Commands ---",
			"No relevant commands found or response indicates failure.",
		)
		if err != nil {
			return err
		}

		return nil
	},
}

// parseFindFlags extracts and validates flags specific to the find command.
func parseFindFlags(cmd *cobra.Command) (limit int, err error) {
	limit, err = cmd.Flags().GetInt("limit")
	if err != nil {
		logger.Error("Failed to get 'limit' flag value", zap.Error(err))
		err = fmt.Errorf("internal error getting limit flag: %w", err)
		return
	}
	return limit, nil
}

// runFind executes the main logic: config, history, LLM interaction.
func runFind(logger *zap.Logger, query string, limit int) (string, error) {
	// 1. Load Configuration
	logger.Debug("Loading configuration...")
	cfg, err := config.LoadConfig(logger)
	if err != nil {
		return "", fmt.Errorf("failed to load configuration: %w", err)
	}
	logger.Debug("Configuration loaded successfully")

	// 2. Read Shell History
	// TODO: Replace with a factory when supporting multiple shells!!
	historyReader, err := history.NewZshHistoryReader(logger)
	if err != nil {
		return "", fmt.Errorf("failed to initialize history reader: %w", err)
	}

	historyEntries, err := historyReader.ReadHistory(limit)
	if err != nil {
		return "", fmt.Errorf("failed to read history: %w", err)
	}
	logger.Debug("History read successfully", zap.Int("entries_count", len(historyEntries)))

	// 3. Initialize LLM Client
	ctx := context.Background()
	logger.Debug("Initializing LLM client...")
	// TODO: Add support for other LLMs!!
	llmClient, err := llm.NewGeminiClient(ctx, logger, cfg.GoogleAPIKey)
	if err != nil {
		return "", fmt.Errorf("failed to initialize LLM client: %w", err)
	}
	defer func() {
		logger.Debug("Closing LLM client...")
		if closeErr := llmClient.Close(); closeErr != nil {
			logger.Error("Failed to close LLM client", zap.Error(closeErr))
		}
	}()
	logger.Debug("LLM client initialized successfully")

	// 4. Call LLM API
	logger.Debug("Sending query and history context to LLM...", zap.Int("history_context_size", len(historyEntries)))
	result, err := llmClient.FindHistoryEntries(query, historyEntries)
	if err != nil {
		return "", fmt.Errorf("failed to get results from LLM: %w", err)
	}
	logger.Debug("Received response from LLM")

	return result, nil
}

// init adds the findCmd and its flags to the rootCmd.
func init() {
	rootCmd.AddCommand(findCmd)

	findCmd.Flags().IntP("limit", "n", defaultFindHistoryLimit, "Limit the number of most recent history entries to analyze")
}
