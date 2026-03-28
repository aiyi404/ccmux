package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type AppConfig struct {
	Mode    string `json:"mode"`
	Lang    string `json:"lang,omitempty"`
	Current string `json:"current,omitempty"`
}

func LoadConfigFrom(path string) (*AppConfig, error) {
	cfg := &AppConfig{Mode: "auto"}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func SaveConfigTo(cfg *AppConfig, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func LoadConfig() (*AppConfig, error) {
	return LoadConfigFrom(CCCConfig)
}

func SaveConfig(cfg *AppConfig) error {
	return SaveConfigTo(cfg, CCCConfig)
}
