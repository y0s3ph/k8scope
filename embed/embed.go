package embed

import "embed"

//go:embed all:charts all:values all:dashboards all:alerts
var Assets embed.FS
