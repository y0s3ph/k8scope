package preflight

import (
	"context"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type ClusterCheck struct {
	kubeconfig string
}

func (c *ClusterCheck) Name() string { return "Cluster reachable" }

func (c *ClusterCheck) Run(ctx context.Context) CheckResult {
	config, err := clientcmd.BuildConfigFromFlags("", c.kubeconfig)
	if err != nil {
		return CheckResult{
			Name:    c.Name(),
			Passed:  false,
			Message: fmt.Sprintf("cannot load kubeconfig: %v", err),
		}
	}

	config.Timeout = 5_000_000_000 // 5 seconds

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return CheckResult{
			Name:    c.Name(),
			Passed:  false,
			Message: fmt.Sprintf("cannot create client: %v", err),
		}
	}

	_, err = client.Discovery().ServerVersion()
	if err != nil {
		return CheckResult{
			Name:    c.Name(),
			Passed:  false,
			Message: fmt.Sprintf("cluster unreachable: %v", err),
		}
	}

	raw, _ := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	contextName := ""
	if raw != nil {
		contextName = raw.CurrentContext
	}

	return CheckResult{
		Name:    c.Name(),
		Passed:  true,
		Message: fmt.Sprintf("context %q", contextName),
	}
}

