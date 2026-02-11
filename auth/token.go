package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	configDirName  = "fiken"
	tokenFileName  = "token"
	configFileName = "config.json"
)

// Config holds the CLI configuration.
type Config struct {
	DefaultCompany string `json:"default_company,omitempty"`
}

// configDir returns the configuration directory path.
func configDir() (string, error) {
	home, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("getting config dir: %w", err)
	}
	return filepath.Join(home, configDirName), nil
}

// ensureConfigDir creates the config directory if it doesn't exist.
func ensureConfigDir() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("creating config dir: %w", err)
	}
	return dir, nil
}

// SaveToken stores the API token to ~/.config/fiken/token.
func SaveToken(token string) error {
	dir, err := ensureConfigDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, tokenFileName)
	if err := os.WriteFile(path, []byte(token), 0600); err != nil {
		return fmt.Errorf("writing token: %w", err)
	}
	return nil
}

// LoadToken reads the API token from ~/.config/fiken/token.
func LoadToken() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, tokenFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("no token found. Run 'fiken auth token <token>' to set up authentication")
		}
		return "", fmt.Errorf("reading token: %w", err)
	}
	return string(data), nil
}

// TokenExists checks whether a token is stored.
func TokenExists() bool {
	dir, err := configDir()
	if err != nil {
		return false
	}
	path := filepath.Join(dir, tokenFileName)
	_, err = os.Stat(path)
	return err == nil
}

// RemoveToken deletes the stored token.
func RemoveToken() error {
	dir, err := configDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, tokenFileName)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing token: %w", err)
	}
	return nil
}

// SaveConfig saves the CLI configuration.
func SaveConfig(cfg *Config) error {
	dir, err := ensureConfigDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, configFileName)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding config: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}
	return nil
}

// LoadConfig loads the CLI configuration.
func LoadConfig() (*Config, error) {
	dir, err := configDir()
	if err != nil {
		return &Config{}, nil
	}
	path := filepath.Join(dir, configFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("decoding config: %w", err)
	}
	return &cfg, nil
}
