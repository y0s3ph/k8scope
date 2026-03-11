package stack

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/y0s3ph/k8scope/internal/engine"
	"gopkg.in/yaml.v3"
)

type Orchestrator struct {
	helm      *engine.HelmInstaller
	assets    fs.FS
	namespace string
}

func NewOrchestrator(kubeconfig, namespace string, assets fs.FS) *Orchestrator {
	return &Orchestrator{
		helm:      engine.NewHelmInstaller(kubeconfig),
		assets:    assets,
		namespace: namespace,
	}
}

func (o *Orchestrator) Install(ctx context.Context, mode string, dryRun bool) error {
	components := ForMode(mode, o.assets)
	if len(components) == 0 {
		return fmt.Errorf("no components defined for mode %q", mode)
	}

	for _, comp := range components {
		if !comp.Enabled {
			continue
		}

		fmt.Printf("  → Installing %s...\n", comp.Name)

		values, err := o.loadValues(comp.ValuesPath)
		if err != nil {
			return fmt.Errorf("loading values for %s: %w", comp.Name, err)
		}

		chartFS, err := fs.Sub(o.assets, comp.ChartPath)
		if err != nil {
			return fmt.Errorf("accessing chart for %s: %w", comp.Name, err)
		}

		opts := engine.ReleaseOptions{
			Name:      comp.ReleaseName,
			Namespace: o.namespace,
			ChartFS:   chartFS,
			Values:    values,
			DryRun:    dryRun,
			Wait:      !dryRun,
		}

		status, err := o.helm.InstallOrUpgrade(ctx, opts)
		if err != nil {
			return fmt.Errorf("installing %s: %w", comp.Name, err)
		}

		if dryRun {
			fmt.Printf("  ✓ %s (dry-run)\n", comp.Name)
		} else {
			fmt.Printf("  ✓ %s installed (status: %s, revision: %d)\n",
				comp.Name, status.Status, status.Version)
		}
	}

	return nil
}

func (o *Orchestrator) Uninstall(ctx context.Context, mode string) error {
	components := ForMode(mode, o.assets)

	for i := len(components) - 1; i >= 0; i-- {
		comp := components[i]
		if !comp.Enabled {
			continue
		}

		fmt.Printf("  → Uninstalling %s...\n", comp.Name)

		err := o.helm.Uninstall(ctx, comp.ReleaseName, o.namespace)
		if err != nil {
			if engine.IsNotFound(err) {
				fmt.Printf("  - %s not found, skipping\n", comp.Name)
				continue
			}
			return fmt.Errorf("uninstalling %s: %w", comp.Name, err)
		}

		fmt.Printf("  ✓ %s uninstalled\n", comp.Name)
	}

	return nil
}

func (o *Orchestrator) loadValues(path string) (map[string]interface{}, error) {
	data, err := fs.ReadFile(o.assets, path)
	if err != nil {
		return nil, err
	}

	var values map[string]interface{}
	if err := yaml.Unmarshal(data, &values); err != nil {
		return nil, fmt.Errorf("parsing values %q: %w", path, err)
	}

	return values, nil
}
