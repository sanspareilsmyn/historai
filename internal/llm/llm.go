package llm

import "github.com/sanspareilsmyn/historai/internal/history"

// LLMClient defines the interface for interacting with an LLM API.
type LLMClient interface {
	FindHistoryEntries(query string, historyContext []history.HistoryEntry) (string, error)

	SuggestCommands(taskDescription string, historyContext []history.HistoryEntry) (string, error)

	Close() error
}
