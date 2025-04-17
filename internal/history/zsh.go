package history

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"go.uber.org/zap"
)

// ZshHistoryReader implements the HistoryReader interface for Zsh.
type ZshHistoryReader struct {
	logger      *zap.Logger
	historyFile string
}

// NewZshHistoryReader remains the same.
func NewZshHistoryReader(logger *zap.Logger) (*ZshHistoryReader, error) {
	histFilePath, err := getDefaultZshHistoryPath()
	if err != nil {
		logger.Error("Failed to get default Zsh history path", zap.Error(err))
		return nil, fmt.Errorf("could not determine Zsh history file path: %w", err)
	}
	if _, err := os.Stat(histFilePath); os.IsNotExist(err) {
		logger.Error("Zsh history file does not exist", zap.String("path", histFilePath))
		return nil, fmt.Errorf("zsh history file not found at %s", histFilePath)
	}
	logger.Debug("Using Zsh history file", zap.String("path", histFilePath))
	return &ZshHistoryReader{
		logger:      logger,
		historyFile: histFilePath,
	}, nil
}

// getDefaultZshHistoryPath remains the same.
func getDefaultZshHistoryPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, ".zsh_history"), nil
}

// ReadHistory opens the history file and delegates parsing and filtering.
func (r *ZshHistoryReader) ReadHistory(limit int) ([]HistoryEntry, error) {
	file, err := os.Open(r.historyFile)
	if err != nil {
		r.logger.Error("Failed to open Zsh history file", zap.String("path", r.historyFile), zap.Error(err))
		return nil, fmt.Errorf("failed to open history file %s: %w", r.historyFile, err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// Delegate parsing to a separate method
	allEntries, err := r.parseHistory(file)
	if err != nil {
		return nil, err
	}

	// Apply the limit filter
	filteredEntries := applyLimitFilter(r.logger, allEntries, limit)
	return filteredEntries, nil
}

// parseHistory reads from the reader, parses entries, handles multi-line and UTF-8.
func (r *ZshHistoryReader) parseHistory(reader io.Reader) ([]HistoryEntry, error) {
	var allEntries []HistoryEntry
	scanner := bufio.NewScanner(reader)
	// Regex still needs to capture timestamp
	re := regexp.MustCompile(`^: (\d{10,}):\d+;(.+)`)
	var currentCommand strings.Builder
	var currentTimestamp int64
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		originalLineBytes := scanner.Bytes()
		// Process line ensuring valid UTF-8
		line := r.ensureValidUTF8(originalLineBytes)

		match := re.FindStringSubmatch(line)
		if len(match) == 3 {
			// Finalize the previous command if one was being built
			if currentCommand.Len() > 0 {
				commandStr := strings.ToValidUTF8(strings.TrimSpace(currentCommand.String()), "\uFFFD")
				allEntries = append(allEntries, HistoryEntry{
					Timestamp: currentTimestamp,
					Command:   commandStr,
				})
			}
			// Start the new command
			currentTimestamp, _ = strconv.ParseInt(match[1], 10, 64)
			currentCommand.Reset()
			currentCommand.WriteString(strings.ToValidUTF8(match[2], "\uFFFD"))
		} else if currentCommand.Len() > 0 {
			currentStr := currentCommand.String()
			nextLineStr := line

			if strings.HasSuffix(currentStr, "\\") {
				currentCommand.Reset()
				currentCommand.WriteString(currentStr[:len(currentStr)-1])
				currentCommand.WriteString(nextLineStr)
			} else {
				currentCommand.WriteString("\n")
				currentCommand.WriteString(nextLineStr)
			}
		}
	}

	// Add the very last command entry if it exists
	if currentCommand.Len() > 0 {
		commandStr := strings.ToValidUTF8(strings.TrimSpace(currentCommand.String()), "\uFFFD")
		allEntries = append(allEntries, HistoryEntry{
			Timestamp: currentTimestamp,
			Command:   commandStr,
		})
	}

	// Check for scanner errors after the loop finishes
	if err := scanner.Err(); err != nil {
		r.logger.Error("Error reading Zsh history data", zap.Error(err))
		return nil, fmt.Errorf("error reading history data: %w", err)
	}

	return allEntries, nil
}

// ensureValidUTF8 checks for valid UTF-8 and replaces invalid sequences.
func (r *ZshHistoryReader) ensureValidUTF8(lineBytes []byte) string {
	if utf8.Valid(lineBytes) {
		return string(lineBytes)
	}

	return strings.ToValidUTF8(string(lineBytes), "\uFFFD")
}

// applyLimitFilter applies the limit to the list of entries.
func applyLimitFilter(logger *zap.Logger, entries []HistoryEntry, limit int) []HistoryEntry {
	if limit <= 0 || len(entries) <= limit {
		return entries
	}

	logger.Debug("Applying 'limit'", zap.Int("limit", limit), zap.Int("initial_count", len(entries)))
	filtered := entries[len(entries)-limit:]
	return filtered
}
