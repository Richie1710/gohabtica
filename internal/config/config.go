package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	// DefaultBaseURL is the default base URL of the Habitica API.
	DefaultBaseURL = "https://habitica.com/api/v3"

	envUserID   = "HABITICA_USER_ID"
	envAPIToken = "HABITICA_API_TOKEN"
)

// Config holds the settings required to talk to the Habitica API.
type Config struct {
	BaseURL  string
	UserID   string
	APIToken string
}

// Options control how configuration is loaded.
type Options struct {
	// ConfigPath allows specifying an explicit path to a YAML configuration file.
	ConfigPath string
	// BaseURLOverride allows overriding the API base URL via a flag.
	BaseURLOverride string
}

// Load reads configuration from environment variables and optionally from a YAML file.
// Order:
//   1. Environment variables HABITICA_USER_ID and HABITICA_API_TOKEN
//   2. If not set: config file (default path or explicitly provided)
func Load(opts Options) (*Config, error) {
	baseURL := DefaultBaseURL
	if opts.BaseURLOverride != "" {
		baseURL = opts.BaseURLOverride
	}

	userID := os.Getenv(envUserID)
	apiToken := os.Getenv(envAPIToken)

	// If both values are present in the environment, use them directly.
	if userID != "" && apiToken != "" {
		return &Config{
			BaseURL:  baseURL,
			UserID:   userID,
			APIToken: apiToken,
		}, nil
	}

	// Otherwise, try to read a configuration file.
	configPath := opts.ConfigPath
	if configPath == "" {
		var err error
		configPath, err = defaultConfigPath()
		if err != nil {
			return nil, err
		}
	}

	cfgFromFile, err := loadFromFile(configPath)
	if err != nil {
		// If the file does not exist, return a clear error about missing credentials.
		if errors.Is(err, os.ErrNotExist) {
			return nil, newMissingCredentialsError()
		}
		return nil, err
	}

	// Base URL from file, if present, otherwise from the value computed above.
	if cfgFromFile.BaseURL == "" {
		cfgFromFile.BaseURL = baseURL
	}

	if cfgFromFile.UserID == "" || cfgFromFile.APIToken == "" {
		return nil, newMissingCredentialsError()
	}

	return cfgFromFile, nil
}

// defaultConfigPath determines the default path of the configuration file.
func defaultConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "gohabitica", "config.yaml"), nil
}

// loadFromFile loads configuration from the given YAML file.
func loadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var raw struct {
		BaseURL  string `yaml:"base_url"`
		UserID   string `yaml:"user_id"`
		APIToken string `yaml:"api_token"`
	}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("configuration file %s could not be read: %w", path, err)
	}

	return &Config{
		BaseURL:  raw.BaseURL,
		UserID:   raw.UserID,
		APIToken: raw.APIToken,
	}, nil
}

// MissingCredentialsError describes missing Habitica credentials.
type MissingCredentialsError struct{}

func (e *MissingCredentialsError) Error() string {
	return "Habitica credentials are missing – please set HABITICA_USER_ID and HABITICA_API_TOKEN or create a config file"
}

func newMissingCredentialsError() error {
	return &MissingCredentialsError{}
}

