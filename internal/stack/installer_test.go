package stack

import (
	"testing"
	"testing/fstest"

	"gopkg.in/yaml.v3"
)

func TestLoadValues(t *testing.T) {
	valuesContent := `prometheus:
  prometheusSpec:
    replicas: 1
    retention: 7d
grafana:
  enabled: false
`

	mockFS := fstest.MapFS{
		"values/startup/prometheus.yaml": &fstest.MapFile{
			Data: []byte(valuesContent),
		},
	}

	o := &Orchestrator{assets: mockFS}
	values, err := o.loadValues("values/startup/prometheus.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out, _ := yaml.Marshal(values)
	if len(out) == 0 {
		t.Fatal("expected non-empty values")
	}

	prom, ok := values["prometheus"].(map[string]interface{})
	if !ok {
		t.Fatal("expected prometheus key in values")
	}
	spec, ok := prom["prometheusSpec"].(map[string]interface{})
	if !ok {
		t.Fatal("expected prometheusSpec in prometheus")
	}
	if spec["replicas"] != 1 {
		t.Errorf("expected replicas=1, got %v", spec["replicas"])
	}
}

func TestLoadValuesInvalidYAML(t *testing.T) {
	mockFS := fstest.MapFS{
		"values/bad.yaml": &fstest.MapFile{
			Data: []byte("{{invalid yaml"),
		},
	}

	o := &Orchestrator{assets: mockFS}
	_, err := o.loadValues("values/bad.yaml")
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestLoadValuesNotFound(t *testing.T) {
	mockFS := fstest.MapFS{}
	o := &Orchestrator{assets: mockFS}
	_, err := o.loadValues("values/missing.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
