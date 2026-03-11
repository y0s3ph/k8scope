package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Mode       string `yaml:"mode"`
	Namespace  string `yaml:"namespace"`
	Kubeconfig string `yaml:"kubeconfig"`
}

func Load(path string) (*Config, error) {
	cfg := &Config{
		Namespace: "k8scope",
	}

	if path == "" {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	return cfg, nil
}

var validModes = map[string]bool{
	"dev":        true,
	"startup":    true,
	"production": true,
	"enterprise": true,
}

func ValidateMode(name string) error {
	if !validModes[name] {
		return fmt.Errorf("invalid mode %q: must be one of dev, startup, production, enterprise", name)
	}
	return nil
}
