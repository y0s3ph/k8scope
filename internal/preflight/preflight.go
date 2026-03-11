package preflight

import (
	"context"
	"fmt"
	"strings"
)

type CheckResult struct {
	Name    string
	Passed  bool
	Message string
	Warning bool
}

type Checker interface {
	Name() string
	Run(ctx context.Context) CheckResult
}

func RunAll(ctx context.Context, checkers []Checker) ([]CheckResult, bool) {
	results := make([]CheckResult, 0, len(checkers))
	allPassed := true

	for _, c := range checkers {
		result := c.Run(ctx)
		results = append(results, result)
		if !result.Passed && !result.Warning {
			allPassed = false
		}
	}

	return results, allPassed
}

func PrintResults(results []CheckResult) {
	for _, r := range results {
		icon := "✓"
		if !r.Passed {
			if r.Warning {
				icon = "⚠"
			} else {
				icon = "✗"
			}
		}
		fmt.Printf("  %s %s: %s\n", icon, r.Name, r.Message)
	}
}

func ForMode(mode, kubeconfig string) []Checker {
	reqs := modeResources[mode]
	if reqs.CPU == "" {
		reqs = modeResources["startup"]
	}

	return []Checker{
		&ClusterCheck{kubeconfig: kubeconfig},
		&VersionCheck{kubeconfig: kubeconfig},
		&ResourceCheck{kubeconfig: kubeconfig, requirements: reqs},
		&StorageClassCheck{kubeconfig: kubeconfig, mode: mode},
	}
}

type ResourceRequirements struct {
	CPU    string
	Memory string
}

var modeResources = map[string]ResourceRequirements{
	"startup":    {CPU: "2", Memory: "4Gi"},
	"production": {CPU: "6", Memory: "12Gi"},
	"enterprise": {CPU: "8", Memory: "16Gi"},
}

func formatQuantity(value int64, isMemory bool) string {
	if isMemory {
		gi := value / (1024 * 1024 * 1024)
		if gi > 0 {
			return fmt.Sprintf("%dGi", gi)
		}
		mi := value / (1024 * 1024)
		return fmt.Sprintf("%dMi", mi)
	}
	if value >= 1000 {
		return fmt.Sprintf("%d", value/1000)
	}
	return fmt.Sprintf("%dm", value)
}

func parseMemoryRequirement(s string) int64 {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "Gi") {
		var val int64
		_, _ = fmt.Sscanf(s, "%dGi", &val)
		return val * 1024 * 1024 * 1024
	}
	if strings.HasSuffix(s, "Mi") {
		var val int64
		_, _ = fmt.Sscanf(s, "%dMi", &val)
		return val * 1024 * 1024
	}
	return 0
}

func parseCPURequirement(s string) int64 {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "m") {
		var val int64
		_, _ = fmt.Sscanf(s, "%dm", &val)
		return val
	}
	var val int64
	_, _ = fmt.Sscanf(s, "%d", &val)
	return val * 1000
}
