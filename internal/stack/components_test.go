package stack

import (
	"testing"
)

func TestForModeStartup(t *testing.T) {
	components := ForMode("startup", nil)
	if len(components) < 2 {
		t.Fatalf("expected at least 2 components for startup mode, got %d", len(components))
	}

	prom := components[0]
	if prom.Name != "Prometheus" {
		t.Errorf("expected first component to be Prometheus, got %s", prom.Name)
	}
	if prom.ReleaseName != "k8scope-prometheus" {
		t.Errorf("unexpected release name: %s", prom.ReleaseName)
	}
	if !prom.Enabled {
		t.Error("expected Prometheus to be enabled")
	}

	grafana := components[1]
	if grafana.Name != "Grafana" {
		t.Errorf("expected second component to be Grafana, got %s", grafana.Name)
	}
	if grafana.ReleaseName != "k8scope-grafana" {
		t.Errorf("unexpected release name: %s", grafana.ReleaseName)
	}
	if grafana.ValuesPath != "values/startup/grafana.yaml" {
		t.Errorf("unexpected values path: %s", grafana.ValuesPath)
	}
}

func TestForModeUnknown(t *testing.T) {
	components := ForMode("unknown", nil)
	if components != nil {
		t.Errorf("expected nil for unknown mode, got %d components", len(components))
	}
}

func TestForModeDev(t *testing.T) {
	components := ForMode("dev", nil)
	if components != nil {
		t.Errorf("expected nil for dev mode (not yet implemented), got %d components", len(components))
	}
}
