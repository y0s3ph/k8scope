package engine

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type ReleaseOptions struct {
	Name      string
	Namespace string
	ChartPath string
	ChartFS   fs.FS
	Values    map[string]interface{}
	DryRun    bool
	Wait      bool
	Timeout   time.Duration
}

type ReleaseStatus struct {
	Name      string
	Namespace string
	Status    string
	Version   int
}

type Installer interface {
	Install(ctx context.Context, opts ReleaseOptions) (*ReleaseStatus, error)
	Upgrade(ctx context.Context, opts ReleaseOptions) (*ReleaseStatus, error)
	Uninstall(ctx context.Context, name, namespace string) error
	Status(ctx context.Context, name, namespace string) (*ReleaseStatus, error)
	IsInstalled(ctx context.Context, name, namespace string) (bool, error)
}

type HelmInstaller struct {
	kubeconfig string
}

func NewHelmInstaller(kubeconfig string) *HelmInstaller {
	return &HelmInstaller{kubeconfig: kubeconfig}
}

func (h *HelmInstaller) Install(ctx context.Context, opts ReleaseOptions) (*ReleaseStatus, error) {
	cfg, err := h.actionConfig(opts.Namespace)
	if err != nil {
		return nil, fmt.Errorf("initializing helm config: %w", err)
	}

	chartPath, cleanup, err := h.resolveChart(opts)
	if err != nil {
		return nil, fmt.Errorf("resolving chart: %w", err)
	}
	if cleanup != nil {
		defer cleanup()
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		return nil, fmt.Errorf("loading chart %q: %w", chartPath, err)
	}

	install := action.NewInstall(cfg)
	install.ReleaseName = opts.Name
	install.Namespace = opts.Namespace
	install.CreateNamespace = true
	install.DryRun = opts.DryRun
	install.Wait = opts.Wait
	install.Timeout = opts.timeout()

	rel, err := install.RunWithContext(ctx, chart, opts.Values)
	if err != nil {
		return nil, fmt.Errorf("installing release %q: %w", opts.Name, err)
	}

	return releaseToStatus(rel), nil
}

func (h *HelmInstaller) Upgrade(ctx context.Context, opts ReleaseOptions) (*ReleaseStatus, error) {
	cfg, err := h.actionConfig(opts.Namespace)
	if err != nil {
		return nil, fmt.Errorf("initializing helm config: %w", err)
	}

	chartPath, cleanup, err := h.resolveChart(opts)
	if err != nil {
		return nil, fmt.Errorf("resolving chart: %w", err)
	}
	if cleanup != nil {
		defer cleanup()
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		return nil, fmt.Errorf("loading chart %q: %w", chartPath, err)
	}

	upgrade := action.NewUpgrade(cfg)
	upgrade.Namespace = opts.Namespace
	upgrade.DryRun = opts.DryRun
	upgrade.Wait = opts.Wait
	upgrade.Timeout = opts.timeout()
	upgrade.ReuseValues = false
	upgrade.Atomic = true

	rel, err := upgrade.RunWithContext(ctx, opts.Name, chart, opts.Values)
	if err != nil {
		return nil, fmt.Errorf("upgrading release %q: %w", opts.Name, err)
	}

	return releaseToStatus(rel), nil
}

func (h *HelmInstaller) Uninstall(_ context.Context, name, namespace string) error {
	cfg, err := h.actionConfig(namespace)
	if err != nil {
		return fmt.Errorf("initializing helm config: %w", err)
	}

	uninstall := action.NewUninstall(cfg)
	uninstall.KeepHistory = false

	_, err = uninstall.Run(name)
	if err != nil {
		return fmt.Errorf("uninstalling release %q: %w", name, err)
	}

	return nil
}

func (h *HelmInstaller) Status(_ context.Context, name, namespace string) (*ReleaseStatus, error) {
	cfg, err := h.actionConfig(namespace)
	if err != nil {
		return nil, fmt.Errorf("initializing helm config: %w", err)
	}

	status := action.NewStatus(cfg)
	rel, err := status.Run(name)
	if err != nil {
		return nil, fmt.Errorf("getting status for %q: %w", name, err)
	}

	return releaseToStatus(rel), nil
}

func (h *HelmInstaller) IsInstalled(_ context.Context, name, namespace string) (bool, error) {
	cfg, err := h.actionConfig(namespace)
	if err != nil {
		return false, fmt.Errorf("initializing helm config: %w", err)
	}

	list := action.NewList(cfg)
	list.Filter = fmt.Sprintf("^%s$", name)
	list.StateMask = action.ListDeployed | action.ListFailed | action.ListPendingInstall

	results, err := list.Run()
	if err != nil {
		return false, fmt.Errorf("listing releases: %w", err)
	}

	return len(results) > 0, nil
}

func (h *HelmInstaller) InstallOrUpgrade(ctx context.Context, opts ReleaseOptions) (*ReleaseStatus, error) {
	installed, err := h.IsInstalled(ctx, opts.Name, opts.Namespace)
	if err != nil {
		return nil, err
	}

	if installed {
		return h.Upgrade(ctx, opts)
	}
	return h.Install(ctx, opts)
}

func (h *HelmInstaller) actionConfig(namespace string) (*action.Configuration, error) {
	settings := cli.New()
	if h.kubeconfig != "" {
		settings.KubeConfig = h.kubeconfig
	}

	flags := &genericclioptions.ConfigFlags{
		KubeConfig: &settings.KubeConfig,
		Namespace:  &namespace,
	}

	cfg := new(action.Configuration)
	if err := cfg.Init(flags, namespace, "secret", func(format string, v ...interface{}) {}); err != nil {
		return nil, err
	}

	return cfg, nil
}

// resolveChart extracts an embedded chart to a temp directory or returns
// the chart path directly if ChartPath is set.
func (h *HelmInstaller) resolveChart(opts ReleaseOptions) (string, func(), error) {
	if opts.ChartPath != "" {
		return opts.ChartPath, nil, nil
	}

	if opts.ChartFS == nil {
		return "", nil, fmt.Errorf("either ChartPath or ChartFS must be set")
	}

	tmpDir, err := os.MkdirTemp("", "k8scope-chart-*")
	if err != nil {
		return "", nil, fmt.Errorf("creating temp dir: %w", err)
	}

	cleanup := func() { os.RemoveAll(tmpDir) }

	err = fs.WalkDir(opts.ChartFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		target := filepath.Join(tmpDir, path)

		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		data, err := fs.ReadFile(opts.ChartFS, path)
		if err != nil {
			return fmt.Errorf("reading embedded file %q: %w", path, err)
		}

		return os.WriteFile(target, data, 0o644)
	})
	if err != nil {
		cleanup()
		return "", nil, fmt.Errorf("extracting chart: %w", err)
	}

	return tmpDir, cleanup, nil
}

func (opts ReleaseOptions) timeout() time.Duration {
	if opts.Timeout > 0 {
		return opts.Timeout
	}
	return 5 * time.Minute
}

func releaseToStatus(rel *release.Release) *ReleaseStatus {
	if rel == nil {
		return nil
	}
	return &ReleaseStatus{
		Name:      rel.Name,
		Namespace: rel.Namespace,
		Status:    strings.ToLower(string(rel.Info.Status)),
		Version:   rel.Version,
	}
}

// IsNotFound returns true if the error indicates a release was not found.
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	return err == driver.ErrReleaseNotFound || strings.Contains(err.Error(), "not found")
}
