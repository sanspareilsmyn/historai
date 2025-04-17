package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger

	// Flag variable to store the value of the --debug flag.
	debugMode bool

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "historai",
		Short: "An AI-powered CLI tool to find/suggest commands based on shell history.",
		Long: `historai helps you search your shell command history
or get command suggestions using natural language queries powered by LLM APIs.
Find or discover commands based on what they do, not just keywords.`,
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var zapConfig zap.Config
			if debugMode {
				zapConfig = zap.NewDevelopmentConfig()
			} else {
				zapConfig = zap.NewProductionConfig()
				zapConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
			}

			logger, err = zapConfig.Build()
			if err != nil {
				return fmt.Errorf("critical: Failed to initialize logger: %w", err)
			}

			logger.Debug("Debug logging enabled.")
			logger.Debug("Logger initialized successfully.")

			return nil
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	err := rootCmd.Execute()
	return err
}

// init is called when the package is imported.
func init() {
	rootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", false, "Enable debug logging")
}
