package stack

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

type DashboardManifest struct {
	Name      string
	Filename  string
	JSON      string
}

func LoadDashboards(assets fs.FS) ([]DashboardManifest, error) {
	var dashboards []DashboardManifest

	entries, err := fs.ReadDir(assets, "dashboards")
	if err != nil {
		return nil, fmt.Errorf("reading dashboards directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		data, err := fs.ReadFile(assets, filepath.Join("dashboards", entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("reading dashboard %s: %w", entry.Name(), err)
		}

		if !json.Valid(data) {
			return nil, fmt.Errorf("invalid JSON in dashboard %s", entry.Name())
		}

		dashboards = append(dashboards, DashboardManifest{
			Name:     strings.TrimSuffix(entry.Name(), ".json"),
			Filename: entry.Name(),
			JSON:     string(data),
		})
	}

	return dashboards, nil
}
