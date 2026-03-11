package stack

import (
	"testing"
)

func TestForModeStartup(t *testing.T) {
	components := ForMode("startup", nil)
	if len(components) < 5 {
		t.Fatalf("expected at least 5 components for startup mode, got %d", len(components))
	}

	expected := []struct {
		name        string
		releaseName string
		enabled     bool
	}{
		{"Prometheus", "k8scope-prometheus", true},
		{"Loki", "k8scope-loki", true},
		{"Alertmanager", "k8scope-alertmanager", false}, // bundled in kube-prometheus-stack
		{"OTel Collector", "k8scope-otel", true},
		{"Grafana", "k8scope-grafana", true},
	}
	for i, e := range expected {
		if components[i].Name != e.name {
			t.Errorf("component[%d]: expected %s, got %s", i, e.name, components[i].Name)
		}
		if components[i].ReleaseName != e.releaseName {
			t.Errorf("component[%d]: expected release %s, got %s", i, e.releaseName, components[i].ReleaseName)
		}
		if components[i].Enabled != e.enabled {
			t.Errorf("component[%d] %s: expected enabled=%v, got %v", i, e.name, e.enabled, components[i].Enabled)
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
