package cli

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/sanspareilsmyn/historai/internal/config"
	"github.com/sanspareilsmyn/historai/internal/history"
	"github.com/sanspareilsmyn/historai/internal/llm"
)

const (
	defaultSuggestHistoryContextLimit = 100
)

// suggestCmd represents the suggest command
var suggestCmd = &cobra.Command{
	Use:   "suggest \"<natural language description of task>\"",
	Short: "Suggest shell commands based on a task description using AI",
	Long: `Asks an LLM to suggest shell commands for the task you describe.
It can optionally use your recent shell history (currently Zsh) as context
to potentially provide more relevant suggestions based on tools you typically use.

You can control the history context using flags:
  --limit / -n        : How many recent history entries to provide as context (default: 100).
  --no-history-context: Disable using shell history as context for the suggestion.

Example:
  historai suggest "how to convert a video file to an animated gif"
  historai suggest --limit 200 "command to find all python files modified today"
  historai suggest --no-history-context "recursively remove all .DS_Store files"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		if query == "" {
			return errors.New("task description cannot be empty")
		}

		// 1. Parse and validate flags (limit, no-history-context)
		limit, noHistoryContext, err := parseSuggestFlags(cmd)
		if err != nil {
			return err
		}

		// 2. Execute the core suggestion logic
		suggestions, err := runSuggestCore(logger, query, limit, noHistoryContext)
		if err != nil {
			return err
		}

		// 3. Print Header to Stderr & Result to Stdout
		err = printCommandOutput(
			logger,
			suggestions,
			"--- Suggested Commands ---",
			"No suggestions generated or suggestions indicate failure.",
		)
		if err != nil {
			return err
		}
		return nil
	},
}

// parseSuggestFlags extracts and validates flags specific to the suggest command.
func parseSuggestFlags(cmd *cobra.Command) (limit int, noHistoryContext bool, err error) {
	limit, err = cmd.Flags().GetInt("limit")
	if err != nil {
		logger.Error("Failed to get 'limit' flag value", zap.Error(err))
		err = fmt.Errorf("internal error getting limit flag: %w", err)
		return
	}

	noHistoryContext, err = cmd.Flags().GetBool("no-history-context")
	if err != nil {
		logger.Error("Failed to get 'no-history-context' flag value", zap.Error(err))
		err = fmt.Errorf("internal error getting no-history-context flag: %w", err)
		return
	}

	return limit, noHistoryContext, nil
}

// runSuggestCore executes the main logic: config, optional history, LLM interaction.
func runSuggestCore(logger *zap.Logger, query string, limit int, noHistoryContext bool) (string, error) {
	// 1. Load Configuration
	cfg, err := config.LoadConfig(logger)
	if err != nil {
		return "", fmt.Errorf("failed to load configuration: %w", err)
	}

	// 2. Read Shell History (Optional, for Context)
	var historyEntries []history.HistoryEntry
	if !noHistoryContext {
		// TODO: Replace with multi-shell logic later
		historyReader, err := history.NewZshHistoryReader(logger)
		if err != nil {
			logger.Error("Failed to initialize Zsh history reader for context", zap.Error(err))
			return "", fmt.Errorf("failed to initialize history reader for context: %w", err)
		}

		// Call ReadHistory with only limit, matching the updated interface/implementation
		historyEntries, err = historyReader.ReadHistory(limit)
		if err != nil {
			logger.Error("Failed to read history for context", zap.Error(err))
			return "", fmt.Errorf("failed to read history for context: %w", err)
		}
		if len(historyEntries) == 0 {
			logger.Warn("No history entries found matching the criteria (limit) to provide as context.")
		}
	} else {
		logger.Debug("Skipping history reading as --no-history-context flag was provided.")
	}

	// 3. Initialize LLM Client
	ctx := context.Background()
	logger.Debug("Initializing LLM client (Gemini)...")
	llmClient, err := llm.NewGeminiClient(ctx, logger, cfg.GoogleAPIKey)
	if err != nil {
		return "", fmt.Errorf("failed to initialize LLM client: %w", err)
	}
	defer func() {
		if closeErr := llmClient.Close(); closeErr != nil {
			logger.Error("Failed to close LLM client", zap.Error(closeErr))
		}
	}()

	// 4. Call LLM API to suggest commands
	suggestions, err := llmClient.SuggestCommands(query, historyEntries)
	if err != nil {
		return "", fmt.Errorf("failed to get suggestions from LLM: %w", err)
	}

	return suggestions, nil
}

// init adds the suggestCmd and its flags to the rootCmd.
func init() {
	rootCmd.AddCommand(suggestCmd)

	suggestCmd.Flags().IntP("limit", "n", defaultSuggestHistoryContextLimit, "Limit the number of most recent history entries to provide as context")
	suggestCmd.Flags().Bool("no-history-context", false, "Do not use shell history as context for suggestions")
}
