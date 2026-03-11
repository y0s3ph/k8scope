package cli

import (
	"github.com/spf13/cobra"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "k8scope",
	Short: "Opinionated observability stack for Kubernetes",
	Long: `k8scope deploys a production-ready observability stack on Kubernetes
with battle-tested defaults. Choose a mode that fits your stage:

  dev          Docker Compose stack for local testing
  startup      Single-replica, lightweight, low resource usage
  production   HA, persistent storage, alerting, retention policies
  enterprise   OIDC, external storage (S3/GCS), multi-tenant`,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: $HOME/.k8scope.yaml)")
	rootCmd.PersistentFlags().String("kubeconfig", "", "path to kubeconfig file")
	rootCmd.PersistentFlags().StringP("namespace", "n", "k8scope", "target namespace")
}

func Execute() error {
	return rootCmd.Execute()
}
