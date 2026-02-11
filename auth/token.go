package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/99designs/keyring"
)

const (
	configDirName = "fiken"
	serviceName   = "fiken-cli"

	keyAPIToken       = "api-token"
	keyDefaultCompany = "default-company"

	// Legacy file names for migration.
	legacyTokenFileName  = "token"
	legacyConfigFileName = "config.json"
)

// KeyringBackend holds the user-selected backend override ("auto" = default).
var KeyringBackend = "auto"

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

// openKeyring opens the keyring with the configured backend.
func openKeyring() (keyring.Keyring, error) {
	dir, err := ensureConfigDir()
	if err != nil {
		return nil, err
	}

	cfg := keyring.Config{
		ServiceName:      serviceName,
		FileDir:          filepath.Join(dir, "keyring"),
		FilePasswordFunc: keyring.TerminalPrompt,
		AllowedBackends:  resolveBackends(KeyringBackend),
	}

	ring, err := keyring.Open(cfg)
	if err != nil {
		return nil, fmt.Errorf("opening keyring: %w", err)
	}
	return ring, nil
}

// resolveBackends returns the allowed backends based on the user choice.
func resolveBackends(backend string) []keyring.BackendType {
	switch strings.ToLower(backend) {
	case "secret-service":
		return []keyring.BackendType{keyring.SecretServiceBackend}
	case "keychain":
		return []keyring.BackendType{keyring.KeychainBackend}
	case "wincred":
		return []keyring.BackendType{keyring.WinCredBackend}
	case "pass":
		return []keyring.BackendType{keyring.PassBackend}
	case "file":
		return []keyring.BackendType{keyring.FileBackend}
	default: // "auto"
		return []keyring.BackendType{
			keyring.SecretServiceBackend,
			keyring.KeychainBackend,
			keyring.WinCredBackend,
			keyring.PassBackend,
			keyring.FileBackend,
		}
	}
}

// migrateLegacyToken checks for a plaintext token file and migrates it to the keyring.
// Returns true if a migration was performed.
func migrateLegacyToken(ring keyring.Keyring) bool {
	dir, err := configDir()
	if err != nil {
		return false
	}
	path := filepath.Join(dir, legacyTokenFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		return false // no legacy file
	}
	token := strings.TrimSpace(string(data))
	if token == "" {
		return false
	}

	// Store in keyring.
	err = ring.Set(keyring.Item{
		Key:  keyAPIToken,
		Data: []byte(token),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not migrate token to keyring: %v\n", err)
		return false
	}

	// Remove the plaintext file.
	if err := os.Remove(path); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not remove legacy token file: %v\n", err)
	}

	fmt.Fprintf(os.Stderr, "✓ Migrated API token from plaintext file to secure keyring storage.\n")
	return true
}

// migrateLegacyConfig checks for a legacy config.json and migrates the default company.
// Returns true if a migration was performed.
func migrateLegacyConfig(ring keyring.Keyring) bool {
	dir, err := configDir()
	if err != nil {
		return false
	}
	path := filepath.Join(dir, legacyConfigFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return false
	}
	if cfg.DefaultCompany == "" {
		// Nothing to migrate, just clean up.
		os.Remove(path)
		return false
	}

	err = ring.Set(keyring.Item{
		Key:  keyDefaultCompany,
		Data: []byte(cfg.DefaultCompany),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not migrate default company to keyring: %v\n", err)
		return false
	}

	if err := os.Remove(path); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not remove legacy config file: %v\n", err)
	}

	fmt.Fprintf(os.Stderr, "✓ Migrated default company from config.json to secure keyring storage.\n")
	return true
}

// SaveToken stores the API token in the keyring.
func SaveToken(token string) error {
	ring, err := openKeyring()
	if err != nil {
		return err
	}
	return ring.Set(keyring.Item{
		Key:  keyAPIToken,
		Data: []byte(token),
	})
}

// LoadToken reads the API token from the keyring, migrating from plaintext if needed.
func LoadToken() (string, error) {
	ring, err := openKeyring()
	if err != nil {
		return "", err
	}

	// Attempt migration of legacy files on first access.
	migrateLegacyToken(ring)
	migrateLegacyConfig(ring)

	item, err := ring.Get(keyAPIToken)
	if err != nil {
		if err == keyring.ErrKeyNotFound {
			return "", fmt.Errorf("no token found. Run 'fiken auth token <token>' to set up authentication")
		}
		return "", fmt.Errorf("reading token from keyring: %w", err)
	}
	return string(item.Data), nil
}

// TokenExists checks whether a token is stored in the keyring.
func TokenExists() bool {
	ring, err := openKeyring()
	if err != nil {
		return false
	}

	// Check for legacy plaintext file too.
	dir, dirErr := configDir()
	if dirErr == nil {
		path := filepath.Join(dir, legacyTokenFileName)
		if _, statErr := os.Stat(path); statErr == nil {
			return true
		}
	}

	_, err = ring.Get(keyAPIToken)
	return err == nil
}

// RemoveToken deletes the stored token from the keyring.
func RemoveToken() error {
	ring, err := openKeyring()
	if err != nil {
		return err
	}
	err = ring.Remove(keyAPIToken)
	if err != nil && err != keyring.ErrKeyNotFound {
		return fmt.Errorf("removing token from keyring: %w", err)
	}

	// Also remove legacy file if it exists.
	dir, dirErr := configDir()
	if dirErr == nil {
		path := filepath.Join(dir, legacyTokenFileName)
		os.Remove(path) // ignore error
	}

	return nil
}

// SaveConfig saves the CLI configuration to the keyring.
func SaveConfig(cfg *Config) error {
	ring, err := openKeyring()
	if err != nil {
		return err
	}
	if cfg.DefaultCompany != "" {
		return ring.Set(keyring.Item{
			Key:  keyDefaultCompany,
			Data: []byte(cfg.DefaultCompany),
		})
	}
	// If empty, remove the key.
	err = ring.Remove(keyDefaultCompany)
	if err != nil && err != keyring.ErrKeyNotFound {
		return fmt.Errorf("removing default company from keyring: %w", err)
	}
	return nil
}

// LoadConfig loads the CLI configuration from the keyring.
func LoadConfig() (*Config, error) {
	ring, err := openKeyring()
	if err != nil {
		return &Config{}, nil
	}

	// Attempt migration of legacy config.
	migrateLegacyConfig(ring)

	item, err := ring.Get(keyDefaultCompany)
	if err != nil {
		if err == keyring.ErrKeyNotFound {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("reading config from keyring: %w", err)
	}
	return &Config{DefaultCompany: string(item.Data)}, nil
}
