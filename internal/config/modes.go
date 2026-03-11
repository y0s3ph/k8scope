package config

import "fmt"

type ComponentSpec struct {
	Name     string
	Enabled  bool
	Replicas int
}

type Mode struct {
	Name        string
	Description string
	Components  []ComponentSpec
	Features    []string
}

func GetMode(name string) (Mode, error) {
	if err := ValidateMode(name); err != nil {
		return Mode{}, err
	}

	mode, ok := modes[name]
	if !ok {
		return Mode{}, fmt.Errorf("mode %q not found", name)
	}
	return mode, nil
}

var modes = map[string]Mode{
	"dev": {
		Name:        "dev",
		Description: "Local Docker Compose stack for development and testing",
		Components: []ComponentSpec{
			{Name: "Prometheus", Enabled: true, Replicas: 1},
			{Name: "Grafana", Enabled: true, Replicas: 1},
			{Name: "Loki", Enabled: true, Replicas: 1},
			{Name: "Alertmanager", Enabled: false, Replicas: 0},
		},
		Features: []string{
			"Ephemeral storage",
			"Pre-loaded dashboards",
			"No authentication",
		},
	},
	"startup": {
		Name:        "startup",
		Description: "Lightweight single-replica stack for small clusters",
		Components: []ComponentSpec{
			{Name: "Prometheus", Enabled: true, Replicas: 1},
			{Name: "Grafana", Enabled: true, Replicas: 1},
			{Name: "Loki", Enabled: true, Replicas: 1},
			{Name: "Alertmanager", Enabled: true, Replicas: 1},
		},
		Features: []string{
			"Persistent storage (10Gi default)",
			"7-day retention",
			"Basic alerting rules",
			"Pre-loaded dashboards",
			"Ingress-ready",
		},
	},
	"production": {
		Name:        "production",
		Description: "High-availability stack with full alerting and retention",
		Components: []ComponentSpec{
			{Name: "Prometheus", Enabled: true, Replicas: 2},
			{Name: "Grafana", Enabled: true, Replicas: 2},
			{Name: "Loki", Enabled: true, Replicas: 3},
			{Name: "Alertmanager", Enabled: true, Replicas: 3},
		},
		Features: []string{
			"High availability",
			"Persistent storage (50Gi default)",
			"30-day retention",
			"Full alerting rules (critical, warning, info)",
			"Pod anti-affinity",
			"Resource limits enforced",
			"PodDisruptionBudgets",
		},
	},
	"enterprise": {
		Name:        "enterprise",
		Description: "Multi-tenant stack with SSO, external storage, and compliance",
		Components: []ComponentSpec{
			{Name: "Prometheus", Enabled: true, Replicas: 2},
			{Name: "Grafana", Enabled: true, Replicas: 3},
			{Name: "Loki", Enabled: true, Replicas: 3},
			{Name: "Alertmanager", Enabled: true, Replicas: 3},
		},
		Features: []string{
			"Everything in production mode",
			"OIDC / SSO authentication",
			"External object storage (S3/GCS/Azure Blob)",
			"Multi-tenant isolation",
			"Audit logging",
			"90-day retention",
			"Horizontal Pod Autoscaling",
		},
	},
}
