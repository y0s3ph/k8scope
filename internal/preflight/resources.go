package preflight

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ResourceCheck struct {
	kubeconfig   string
	requirements ResourceRequirements
}

func (r *ResourceCheck) Name() string { return "Allocatable resources" }

func (r *ResourceCheck) Run(ctx context.Context) CheckResult {
	config, err := clientcmd.BuildConfigFromFlags("", r.kubeconfig)
	if err != nil {
		return CheckResult{
			Name:    r.Name(),
			Passed:  false,
			Message: fmt.Sprintf("cannot load kubeconfig: %v", err),
		}
	}

	client, err := buildClientsetFromConfig(config)
	if err != nil {
		return CheckResult{
			Name:    r.Name(),
			Passed:  false,
			Message: fmt.Sprintf("cannot create client: %v", err),
		}
	}

	nodes, err := client.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return CheckResult{
			Name:    r.Name(),
			Passed:  false,
			Message: fmt.Sprintf("cannot list nodes: %v", err),
		}
	}

	var totalCPU, totalMem int64
	for _, node := range nodes.Items {
		cpu := node.Status.Allocatable.Cpu()
		mem := node.Status.Allocatable.Memory()
		if cpu != nil {
			totalCPU += cpu.MilliValue()
		}
		if mem != nil {
			totalMem += mem.Value()
		}
	}

	reqCPU := parseCPURequirement(r.requirements.CPU)
	reqMem := parseMemoryRequirement(r.requirements.Memory)

	cpuOK := totalCPU >= reqCPU
	memOK := totalMem >= reqMem

	msg := fmt.Sprintf("%s CPU, %s RAM available (minimum: %s CPU, %s RAM)",
		formatQuantity(totalCPU, false),
		formatQuantity(totalMem, true),
		r.requirements.CPU,
		r.requirements.Memory,
	)

	if !cpuOK || !memOK {
		return CheckResult{
			Name:    r.Name(),
			Passed:  false,
			Message: msg + " — insufficient",
		}
	}

	return CheckResult{
		Name:    r.Name(),
		Passed:  true,
		Message: msg,
	}
}

func buildClientsetFromConfig(config *rest.Config) (*kubernetes.Clientset, error) {
	return kubernetes.NewForConfig(config)
}
