package llm

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/sanspareilsmyn/historai/internal/history"

	"github.com/google/generative-ai-go/genai"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

const (
	defaultModelName = "gemini-1.5-flash-latest"

	findHistoryContextLimit    = 150
	suggestHistoryContextLimit = 50
)

// GeminiClient implements the LLMClient interface using the Google AI (Gemini) API.
type GeminiClient struct {
	logger *zap.Logger
	client *genai.Client
	model  *genai.GenerativeModel
}

// NewGeminiClient creates a new client specifically for the Google Gemini models.
func NewGeminiClient(ctx context.Context, logger *zap.Logger, apiKey string) (*GeminiClient, error) {
	if apiKey == "" {
		return nil, errors.New("google AI (Gemini) API key is required")
	}

	modelName := defaultModelName
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		logger.Error("Failed to create Google AI (genai) client", zap.Error(err))
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	model := client.GenerativeModel(modelName)
	model.SafetySettings = defaultSafetySettings()

	return &GeminiClient{
		logger: logger,
		client: client,
		model:  model,
	}, nil
}

// defaultSafetySettings returns a default set of safety settings.
func defaultSafetySettings() []*genai.SafetySetting {
	return []*genai.SafetySetting{
		{Category: genai.HarmCategoryHarassment, Threshold: genai.HarmBlockMediumAndAbove},
		{Category: genai.HarmCategoryHateSpeech, Threshold: genai.HarmBlockMediumAndAbove},
		{Category: genai.HarmCategorySexuallyExplicit, Threshold: genai.HarmBlockMediumAndAbove},
		{Category: genai.HarmCategoryDangerousContent, Threshold: genai.HarmBlockMediumAndAbove},
	}
}

// Close closes the underlying Google AI (genai) client.
func (c *GeminiClient) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// FindHistoryEntries implements the LLMClient interface method.
func (c *GeminiClient) FindHistoryEntries(query string, historyContext []history.HistoryEntry) (string, error) {
	prompt := c.buildFindPrompt(query, historyContext)
	if prompt == "" {
		return "(No relevant commands found or AI response was empty)", nil
	}

	result, err := c.generateGeminiContent(context.Background(), prompt)
	if err != nil {
		c.logger.Error("Gemini content generation failed for FindHistoryEntries", zap.Error(err))
		return "", fmt.Errorf("gemini API call failed (Find): %w", err)
	}

	if result == "" || result == "No relevant commands found." {
		c.logger.Info("Gemini indicated no relevant commands found for the query.")
		return "(No relevant commands found or AI response was empty)", nil
	}

	return result, nil
}

// SuggestCommands implements the LLMClient interface method.
func (c *GeminiClient) SuggestCommands(taskDescription string, historyContext []history.HistoryEntry) (string, error) {
	prompt := c.buildSuggestPrompt(taskDescription, historyContext)

	result, err := c.generateGeminiContent(context.Background(), prompt)
	if err != nil {
		c.logger.Error("Gemini content generation failed for SuggestCommands", zap.Error(err))
		if strings.Contains(err.Error(), "blocked due to safety settings") {
			return "", errors.New("suggestion blocked due to safety settings") // Return specific user-friendly error
		}
		return "", fmt.Errorf("gemini API call failed (Suggest): %w", err)
	}

	if result == "" || result == "Cannot suggest a command for this task." {
		c.logger.Info("Gemini indicated it cannot suggest a command for the task.")
		return "(AI could not suggest a command for this task or the response was empty)", nil
	}

	return result, nil
}

// generateGeminiContent calls the Gemini API and handles common error/safety checks.
func (c *GeminiClient) generateGeminiContent(ctx context.Context, prompt string) (string, error) {
	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))

	// 1. Check for API call error (network, auth, etc.)
	if err != nil {
		if resp != nil && resp.PromptFeedback != nil && resp.PromptFeedback.BlockReason == genai.BlockReasonSafety {
			c.logger.Warn("Prompt blocked by safety settings during API call", zap.Any("feedback", resp.PromptFeedback))
			return "", fmt.Errorf("prompt blocked due to safety settings (Reason: %s)", resp.PromptFeedback.BlockReason.String())
		}
		return "", fmt.Errorf("API call error: %w", err)
	}

	// 2. Check for safety block in prompt feedback (even if err is nil)
	if resp.PromptFeedback != nil && resp.PromptFeedback.BlockReason == genai.BlockReasonSafety {
		c.logger.Warn("Prompt blocked by safety settings", zap.Any("feedback", resp.PromptFeedback))
		return "", fmt.Errorf("prompt blocked due to safety settings (Reason: %s)", resp.PromptFeedback.BlockReason.String())
	}

	// 3. Check for safety block in candidate finish reason
	if len(resp.Candidates) > 0 && resp.Candidates[0].FinishReason == genai.FinishReasonSafety {
		c.logger.Warn("Response candidate blocked by safety settings", zap.Any("candidate", resp.Candidates[0]))
		return "", errors.New("response blocked due to safety settings")
	}

	// 4. Extract text response
	aiResponseText := extractTextFromResponse(resp)
	if aiResponseText == "" {
		c.logger.Warn("Received empty text response from Gemini (or response was blocked)")
		return "", nil
	}

	return aiResponseText, nil
}

// buildFindPrompt constructs the prompt string for finding history entries.
func (c *GeminiClient) buildFindPrompt(query string, historyContext []history.HistoryEntry) string {
	if len(historyContext) == 0 {
		c.logger.Warn("Cannot build find prompt: history context is empty")
		return ""
	}

	var promptBuilder strings.Builder
	promptBuilder.WriteString("You are an expert shell history analyzer.\n")
	promptBuilder.WriteString("The user is searching their shell history for commands based on a description.\n")
	promptBuilder.WriteString(fmt.Sprintf("User's search query: \"%s\"\n\n", query))
	promptBuilder.WriteString("Please analyze the following shell history entries. Return ONLY the command text of the entry or entries that BEST match the user's query. If multiple commands are good matches, list each matching command on a new line.\n")
	promptBuilder.WriteString("If NO history entries strongly match the query, return the exact phrase: 'No relevant commands found.'\n\n")

	promptBuilder.WriteString(formatHistoryContext("Shell History Entries Provided", historyContext, findHistoryContextLimit))

	promptBuilder.WriteString("Matching command(s) from the history above:\n")

	return promptBuilder.String()
}

// buildSuggestPrompt constructs the prompt for generating command suggestions.
func (c *GeminiClient) buildSuggestPrompt(taskDescription string, historyContext []history.HistoryEntry) string {
	var promptBuilder strings.Builder

	promptBuilder.WriteString("You are an AI assistant expert in generating safe and useful POSIX-compliant shell commands (like for Linux or macOS).\n")
	promptBuilder.WriteString("The user wants a shell command to accomplish the following task:\n")
	promptBuilder.WriteString(fmt.Sprintf("Task: \"%s\"\n\n", taskDescription))

	promptBuilder.WriteString(formatHistoryContext("Recent History Context (Optional)", historyContext, suggestHistoryContextLimit))

	promptBuilder.WriteString("Instructions for generating the command:\n")
	promptBuilder.WriteString("1. Generate one or more shell commands that directly address the user's task.\n")
	promptBuilder.WriteString("2. **Prioritize Safety:** Avoid suggesting potentially destructive commands (like `rm -rf /`, `dd`, etc.) unless absolutely necessary for the task AND explicitly confirmed by the user's request phrasing. If suggesting a command with potential side effects (e.g., modifying files, deleting data), add a brief `# Warning: This command modifies/deletes...` comment before it.\n")
	promptBuilder.WriteString("3. Provide ONLY the raw command(s), each on a new line.\n")
	promptBuilder.WriteString("4. If multiple steps or commands are needed, list them sequentially.\n")
	promptBuilder.WriteString("5. If the task is ambiguous, too complex for a simple command, or cannot be safely achieved, respond with the exact phrase: 'Cannot suggest a command for this task.'\n\n")

	promptBuilder.WriteString("Suggested Command(s):\n")

	return promptBuilder.String()
}

// formatHistoryContext formats the history entries for inclusion in a prompt.
func formatHistoryContext(header string, historyContext []history.HistoryEntry, maxEntries int) string {
	if len(historyContext) == 0 {
		return "No specific user history context provided.\n\n"
	}

	var builder strings.Builder
	builder.WriteString(header + ":\n")
	builder.WriteString(strings.Repeat("-", len(header)+1) + "\n")

	startIdx := 0
	if maxEntries > 0 && len(historyContext) > maxEntries {
		startIdx = len(historyContext) - maxEntries
	}

	for i := startIdx; i < len(historyContext); i++ {
		builder.WriteString(historyContext[i].Command)
		builder.WriteString("\n")
	}
	builder.WriteString(strings.Repeat("-", len(header)+1) + "\n\n") // Dynamic underline

	return builder.String()
}

// extractTextFromResponse safely extracts the text content from the Gemini API response candidates.
func extractTextFromResponse(resp *genai.GenerateContentResponse) string {
	if resp == nil || len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		return ""
	}

	var result strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			result.WriteString(string(text))
		}
	}
	return result.String()
}
