package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ProjectName string `yaml:"project_name"`
	Author      string `yaml:"author"`
}

func DefaultConfig() Config {
	return Config{}
}

func Load(batonDir string) (Config, error) {
	path := filepath.Join(batonDir, "config.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return DefaultConfig(), nil
		}
		return Config{}, fmt.Errorf("config load: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("config load: %w", err)
	}
	return cfg, nil
}

func Save(batonDir string, cfg Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("config save: %w", err)
	}
	path := filepath.Join(batonDir, "config.yaml")
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("config save: %w", err)
	}
	return nil
}
