package stack

import (
	"io/fs"
)

type Component struct {
	Name       string
	ReleaseName string
	ChartPath  string
	ValuesPath string
	Order      int
	Enabled    bool
}

func ForMode(mode string, assets fs.FS) []Component {
	switch mode {
	case "startup":
		return startupComponents()
	case "production":
		return productionComponents()
	case "enterprise":
		return enterpriseComponents()
	default:
		return nil
	}
}

func startupComponents() []Component {
	return []Component{
		{
			Name:        "Prometheus",
			ReleaseName: "k8scope-prometheus",
			ChartPath:   "charts/kube-prometheus-stack",
			ValuesPath:  "values/startup/prometheus.yaml",
			Order:       1,
			Enabled:     true,
		},
	}
}

func productionComponents() []Component {
	return startupComponents()
}

func enterpriseComponents() []Component {
	return startupComponents()
}
