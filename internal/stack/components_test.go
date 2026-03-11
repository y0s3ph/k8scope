package stack

import (
	"testing"
)

func TestForModeStartup(t *testing.T) {
	components := ForMode("startup", nil)
	if len(components) < 3 {
		t.Fatalf("expected at least 3 components for startup mode, got %d", len(components))
	}

	expected := []struct {
		name        string
		releaseName string
	}{
		{"Prometheus", "k8scope-prometheus"},
		{"Loki", "k8scope-loki"},
		{"Grafana", "k8scope-grafana"},
	}
	for i, e := range expected {
		if components[i].Name != e.name {
			t.Errorf("component[%d]: expected %s, got %s", i, e.name, components[i].Name)
		}
		if components[i].ReleaseName != e.releaseName {
			t.Errorf("component[%d]: expected release %s, got %s", i, e.releaseName, components[i].ReleaseName)
		}
		if !components[i].Enabled {
			t.Errorf("component[%d] %s: expected enabled", i, e.name)
		}
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
