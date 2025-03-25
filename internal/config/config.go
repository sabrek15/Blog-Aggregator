package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL	string `json:"db_url"`
	CurrentUserName		string `json:"current_user_name"`
}

func Read() (Config, error) {
	filePath, err := getConfigFilePath();
	if err != nil {
		return Config{}, err
	}

	resp, err := os.Open(filePath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to open config file: %w", err)
	}
	defer resp.Close()

	var cfg Config

	err = json.NewDecoder(resp).Decode(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return cfg, nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir();
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(homeDir, configFileName), nil
}

func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName
	return Write(*cfg)
}

func Write(cfg Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("failed to get file path: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder:= json.NewEncoder(file)

	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("failed to encode json: %w", err)
	}

	return nil
}