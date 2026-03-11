package stack

import (
	"testing"
)

func TestForModeStartup(t *testing.T) {
	components := ForMode("startup", nil)
	if len(components) == 0 {
		t.Fatal("expected at least one component for startup mode")
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
	if prom.ChartPath != "charts/kube-prometheus-stack" {
		t.Errorf("unexpected chart path: %s", prom.ChartPath)
	}
	if prom.ValuesPath != "values/startup/prometheus.yaml" {
		t.Errorf("unexpected values path: %s", prom.ValuesPath)
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
