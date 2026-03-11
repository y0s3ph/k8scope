package preflight

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"k8s.io/client-go/tools/clientcmd"
)

const (
	minMajor = 1
	minMinor = 25
	maxTested = 32
)

type VersionCheck struct {
	kubeconfig string
}

func (v *VersionCheck) Name() string { return "Kubernetes version" }

func (v *VersionCheck) Run(ctx context.Context) CheckResult {
	config, err := clientcmd.BuildConfigFromFlags("", v.kubeconfig)
	if err != nil {
		return CheckResult{
			Name:    v.Name(),
			Passed:  false,
			Message: fmt.Sprintf("cannot load kubeconfig: %v", err),
		}
	}

	client, err := buildClientsetFromConfig(config)
	if err != nil {
		return CheckResult{
			Name:    v.Name(),
			Passed:  false,
			Message: fmt.Sprintf("cannot create client: %v", err),
		}
	}

	serverVersion, err := client.Discovery().ServerVersion()
	if err != nil {
		return CheckResult{
			Name:    v.Name(),
			Passed:  false,
			Message: fmt.Sprintf("cannot get server version: %v", err),
		}
	}

	major, _ := strconv.Atoi(serverVersion.Major)
	minorStr := strings.TrimRight(serverVersion.Minor, "+")
	minor, _ := strconv.Atoi(minorStr)

	version := fmt.Sprintf("v%d.%d", major, minor)

	if major < minMajor || (major == minMajor && minor < minMinor) {
		return CheckResult{
			Name:    v.Name(),
			Passed:  false,
			Message: fmt.Sprintf("%s (minimum required: v%d.%d)", version, minMajor, minMinor),
		}
	}

	if minor > maxTested {
		return CheckResult{
			Name:    v.Name(),
			Passed:  true,
			Warning: true,
			Message: fmt.Sprintf("%s (untested, latest verified: v%d.%d)", version, minMajor, maxTested),
		}
	}

	return CheckResult{
		Name:    v.Name(),
		Passed:  true,
		Message: fmt.Sprintf("%s (supported)", version),
	}
}
