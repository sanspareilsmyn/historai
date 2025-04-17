package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"go.uber.org/zap"
)

var knownFailureMessages = map[string]struct{}{
	"": {},
	"(No relevant commands found or AI response was empty)":                    {},
	"No relevant commands found.":                                              {},
	"(AI could not suggest a command for this task or the response was empty)": {},
	"Cannot suggest a command for this task.":                                  {},
	"suggestion blocked due to safety settings":                                {},
}

func printCommandOutput(logger *zap.Logger, output string, header string, logOnFailure string) (err error) {
	// Trim whitespace just in case
	trimmedOutput := strings.TrimSpace(output)

	// Check if the output is effectively empty or a known failure message
	_, isKnownFailure := knownFailureMessages[trimmedOutput]

	if !isKnownFailure {
		infoColor := color.New(color.FgYellow)
		_, err = infoColor.Fprintln(os.Stderr, "\n"+header)
		if err != nil {
			return err
		}

		resultColor := color.New(color.FgGreen)
		_, err = resultColor.Println(trimmedOutput)
		if err != nil {
			return err
		}

	} else {
		logger.Warn(logOnFailure, zap.String("response", trimmedOutput))
		_, err = fmt.Fprintln(os.Stderr, trimmedOutput)
		if err != nil {
			return err
		}
	}
	return nil
}
