package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Mode       string           `yaml:"mode"`
	Namespace  string           `yaml:"namespace"`
	Kubeconfig string           `yaml:"kubeconfig"`
	Components ComponentsConfig `yaml:"components"`
	Ingress    IngressConfig    `yaml:"ingress"`
}

type ComponentsConfig struct {
	Prometheus   ComponentOverride `yaml:"prometheus"`
	Grafana      GrafanaOverride   `yaml:"grafana"`
	Loki         ComponentOverride `yaml:"loki"`
	Alertmanager ComponentOverride `yaml:"alertmanager"`
	OtelCollector OtelOverride     `yaml:"otelCollector"`
}

type ComponentOverride struct {
	Enabled   *bool  `yaml:"enabled,omitempty"`
	Replicas  *int   `yaml:"replicas,omitempty"`
	Storage   string `yaml:"storage,omitempty"`
	Retention string `yaml:"retention,omitempty"`
}

type GrafanaOverride struct {
	ComponentOverride `yaml:",inline"`
	AdminPassword     string `yaml:"adminPassword,omitempty"`
}

type OtelOverride struct {
	ComponentOverride `yaml:",inline"`
	DeployMode        string `yaml:"mode,omitempty"` // daemonset | gateway | daemonset+gateway
}

type IngressConfig struct {
	Enabled   bool   `yaml:"enabled"`
	ClassName string `yaml:"className,omitempty"`
	Domain    string `yaml:"domain,omitempty"`
	TLS      bool   `yaml:"tls"`
}

func DefaultConfig() *Config {
	return &Config{
		Namespace: "k8scope",
	}
}

func Load(path string) (*Config, error) {
	if path != "" {
		return loadFromFile(path)
	}

	for _, candidate := range configSearchPaths() {
		if _, err := os.Stat(candidate); err == nil {
			return loadFromFile(candidate)
		}
	}

	return DefaultConfig(), nil
}

func loadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	return cfg, nil
}

func configSearchPaths() []string {
	paths := []string{".k8scope.yaml", ".k8scope.yml"}

	home, err := os.UserHomeDir()
	if err == nil {
		paths = append(paths,
			filepath.Join(home, ".k8scope.yaml"),
			filepath.Join(home, ".k8scope.yml"),
		)
	}

	return paths
}

// ApplyFlags merges CLI flag values into the config.
// Flags take highest priority and override config file values.
func (c *Config) ApplyFlags(mode, namespace, kubeconfig string) {
	if mode != "" {
		c.Mode = mode
	}
	if namespace != "" {
		c.Namespace = namespace
	}
	if kubeconfig != "" {
		c.Kubeconfig = kubeconfig
	}
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
