package history

// HistoryEntry represents a single command from the shell history.
type HistoryEntry struct {
	Timestamp int64
	Command   string // The command itself
}

// HistoryReader defines the interface for reading shell history.
type HistoryReader interface {
	ReadHistory(limit int) ([]HistoryEntry, error)
}
