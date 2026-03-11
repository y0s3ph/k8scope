package stack

import (
	"testing"

	"github.com/y0s3ph/k8scope/embed"
)

func TestLoadDashboards(t *testing.T) {
	dashboards, err := LoadDashboards(embed.Assets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(dashboards) != 5 {
		t.Fatalf("expected 5 dashboards, got %d", len(dashboards))
	}

	expected := map[string]bool{
		"cluster-overview": false,
		"node-resources":   false,
		"pod-resources":    false,
		"networking":       false,
		"logs-overview":    false,
	}

	for _, d := range dashboards {
		if _, ok := expected[d.Name]; !ok {
			t.Errorf("unexpected dashboard: %s", d.Name)
		}
		expected[d.Name] = true

		if len(d.JSON) == 0 {
			t.Errorf("dashboard %s has empty JSON", d.Name)
		}
	}

	for name, found := range expected {
		if !found {
			t.Errorf("missing dashboard: %s", name)
		}
	}
}
