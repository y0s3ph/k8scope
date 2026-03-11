package preflight

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

type StorageClassCheck struct {
	kubeconfig string
	mode       string
}

func (s *StorageClassCheck) Name() string { return "StorageClass" }

func (s *StorageClassCheck) Run(ctx context.Context) CheckResult {
	if s.mode == "dev" {
		return CheckResult{
			Name:    s.Name(),
			Passed:  true,
			Message: "not required for dev mode",
		}
	}

	config, err := clientcmd.BuildConfigFromFlags("", s.kubeconfig)
	if err != nil {
		return CheckResult{
			Name:    s.Name(),
			Passed:  false,
			Message: fmt.Sprintf("cannot load kubeconfig: %v", err),
		}
	}

	client, err := buildClientsetFromConfig(config)
	if err != nil {
		return CheckResult{
			Name:    s.Name(),
			Passed:  false,
			Message: fmt.Sprintf("cannot create client: %v", err),
		}
	}

	scList, err := client.StorageV1().StorageClasses().List(ctx, metav1.ListOptions{})
	if err != nil {
		return CheckResult{
			Name:    s.Name(),
			Passed:  false,
			Message: fmt.Sprintf("cannot list storage classes: %v", err),
		}
	}

	if len(scList.Items) == 0 {
		return CheckResult{
			Name:    s.Name(),
			Passed:  false,
			Message: "no StorageClass found (required for persistent volumes)",
		}
	}

	defaultSC := ""
	for _, sc := range scList.Items {
		if sc.Annotations["storageclass.kubernetes.io/is-default-class"] == "true" {
			defaultSC = sc.Name
			break
		}
	}

	if defaultSC != "" {
		return CheckResult{
			Name:    s.Name(),
			Passed:  true,
			Message: fmt.Sprintf("%q (default)", defaultSC),
		}
	}

	return CheckResult{
		Name:    s.Name(),
		Passed:  true,
		Warning: true,
		Message: fmt.Sprintf("%d available, but no default set", len(scList.Items)),
	}
}
