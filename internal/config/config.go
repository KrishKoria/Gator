package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
    DBURL          string `json:"db_url"`
    CurrentUserName string `json:"current_user_name"`
}

func Read() (*Config, error) {
    configPath, err := getConfigFilePath()
    if err != nil {
        return nil, fmt.Errorf("failed to get config file path: %w", err)
    }

    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config Config
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("failed to unmarshal config: %w", err)
    }

    return &config, nil
}

func (cfg *Config) SetUser(username string) error {
    cfg.CurrentUserName = username
    return write(*cfg)
}

func getConfigFilePath() (string, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "", fmt.Errorf("failed to get home directory: %w", err)
    }
    return filepath.Join(homeDir, configFileName), nil
}

func write(cfg Config) error {
    configPath, err := getConfigFilePath()
    if err != nil {
        return fmt.Errorf("failed to get config file path: %w", err)
    }

    data, err := json.MarshalIndent(cfg, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal config: %w", err)
    }

    if err := os.WriteFile(configPath, data, 0644); err != nil {
        return fmt.Errorf("failed to write config file: %w", err)
    }

    return nil
}