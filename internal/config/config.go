package config

import (
	"errors"
	"os"

	"go.uber.org/zap"
)

const (
	EnvGoogleAPIKey = "GOOGLE_API_KEY"
)

// Config holds the application configuration.
type Config struct {
	GoogleAPIKey string
}

// LoadConfig loads the configuration, currently only from environment variables.
func LoadConfig(logger *zap.Logger) (*Config, error) {
	apiKey := os.Getenv(EnvGoogleAPIKey)
	if apiKey == "" {
		logger.Error("Google AI API Key not found in environment variable",
			zap.String("variable_name", EnvGoogleAPIKey))
		return nil, errors.New("required environment variable " + EnvGoogleAPIKey + " is not set")
	}
	logger.Debug("Successfully loaded Google AI API Key from environment variable")

	cfg := &Config{
		GoogleAPIKey: apiKey,
	}

	return cfg, nil
}
