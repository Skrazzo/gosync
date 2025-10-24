package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const FileName = ".gosync.json"

// Config represents the gosync configuration
type Config struct {
	Host           string `json:"host"`
	User           string `json:"user"`
	AuthType       string `json:"authType"` // "password" or "key"
	Password       string `json:"password,omitempty"`
	PrivateKeyPath string `json:"privateKeyPath,omitempty"`
	RemoteDir      string `json:"remoteDir"`
}

// Load loads configuration from .gosync.json in the current directory
func (c *Config) Load() error {
	configPath := filepath.Join(".", FileName)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, c); err != nil {
		return fmt.Errorf("invalid config file: %w", err)
	}

	return nil
}

// Save saves configuration to .gosync.json in the current directory
func (c *Config) Save() error {
	configPath := filepath.Join(".", FileName)

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600) // Restricted permissions for security
}

// Exists checks if .gosync.json exists in the current directory
func (c *Config) Exists() bool {
	configPath := filepath.Join(".", FileName)
	_, err := os.Stat(configPath)
	return err == nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("host is required")
	}
	if c.User == "" {
		return fmt.Errorf("user is required")
	}
	if c.RemoteDir == "" {
		return fmt.Errorf("remoteDir is required")
	}
	if c.AuthType != "password" && c.AuthType != "key" {
		return fmt.Errorf("authType must be 'password' or 'key'")
	}
	if c.AuthType == "password" && c.Password == "" {
		return fmt.Errorf("password is required when authType is 'password'")
	}
	if c.AuthType == "key" && c.PrivateKeyPath == "" {
		return fmt.Errorf("privateKeyPath is required when authType is 'key'")
	}

	return nil
}

// New creates a new Config instance
func New() *Config {
	return &Config{}
}

// NewFromFile creates a new Config instance and loads it from disk
func NewFromFile() (*Config, error) {
	c := New()
	if err := c.Load(); err != nil {
		return nil, err
	}
	return c, nil
}
