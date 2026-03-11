package engine

import "context"

type ReleaseOptions struct {
	Name      string
	Namespace string
	ChartPath string
	Values    map[string]interface{}
	DryRun    bool
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
