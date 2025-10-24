package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const ConfigFileName = ".gosync.json"

// Config represents the gosync configuration
type Config struct {
	LocalDir       string `json:"localDir"`
	Host           string `json:"host"`
	User           string `json:"user"`
	AuthType       string `json:"authType"` // "password" or "key"
	Password       string `json:"password,omitempty"`
	PrivateKeyPath string `json:"privateKeyPath,omitempty"`
	RemoteDir      string `json:"remoteDir"`
}

// LoadConfig loads configuration from .gosync.json in the current directory
func LoadConfig() (*Config, error) {
	configPath := filepath.Join(".", ConfigFileName)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("invalid config file: %w", err)
	}

	return &config, nil
}

// SaveConfig saves configuration to .gosync.json in the current directory
func SaveConfig(config *Config) error {
	configPath := filepath.Join(".", ConfigFileName)

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600) // Restricted permissions for security
}

// ConfigExists checks if .gosync.json exists in the current directory
func ConfigExists() bool {
	configPath := filepath.Join(".", ConfigFileName)
	_, err := os.Stat(configPath)
	return err == nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.LocalDir == "" {
		return fmt.Errorf("localDir is required")
	}
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
