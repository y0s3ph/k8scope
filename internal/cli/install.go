package cli

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/y0s3ph/k8scope/internal/config"
	"github.com/y0s3ph/k8scope/internal/preflight"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the observability stack",
	Long:  `Install Prometheus, Grafana, Loki, and Alertmanager with curated defaults for the selected mode.`,
	RunE:  runInstall,
}

func init() {
	installCmd.Flags().StringP("mode", "m", "", "deployment mode: dev, startup, production, enterprise (required)")
	installCmd.Flags().Bool("dry-run", false, "show what would be installed without applying")
	installCmd.Flags().Bool("skip-preflight", false, "skip preflight checks")
	_ = installCmd.MarkFlagRequired("mode")
	rootCmd.AddCommand(installCmd)
}

func runInstall(cmd *cobra.Command, args []string) error {
	modeName, _ := cmd.Flags().GetString("mode")
	namespace, _ := cmd.Flags().GetString("namespace")
	kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	skipPreflight, _ := cmd.Flags().GetBool("skip-preflight")

	mode, err := config.GetMode(modeName)
	if err != nil {
		return err
	}

	if modeName != "dev" && !skipPreflight {
		fmt.Printf("Preflight checks for mode %q:\n", mode.Name)
		checkers := preflight.ForMode(modeName, kubeconfig)
		results, ok := preflight.RunAll(context.Background(), checkers)
		preflight.PrintResults(results)
		fmt.Println()

		if !ok {
			return fmt.Errorf("preflight checks failed — fix the issues above or use --skip-preflight to bypass")
		}
	}

	if dryRun {
		fmt.Printf("Dry run: showing installation plan for mode %q\n\n", mode.Name)
	} else {
		fmt.Printf("Installing k8scope in mode %q\n\n", mode.Name)
	}

	printInstallPlan(mode, namespace)

	if dryRun {
		fmt.Println("\nNo changes applied (dry-run mode).")
		return nil
	}

	// TODO: wire up HelmInstaller for actual deployments (Phase 2 issues #4-#8)
	fmt.Println("\nInstallation engine ready. Component charts not yet embedded (see Phase 2 issues).")
	return nil
}

func printInstallPlan(mode config.Mode, namespace string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	_, _ = fmt.Fprintf(w, "Namespace:\t%s\n", namespace)
	_, _ = fmt.Fprintf(w, "Mode:\t%s\n", mode.Name)
	_, _ = fmt.Fprintf(w, "Description:\t%s\n", mode.Description)
	_, _ = fmt.Fprintln(w, "\nComponents:")
	_, _ = fmt.Fprintf(w, "  COMPONENT\tENABLED\tREPLICAS\n")
	for _, c := range mode.Components {
		_, _ = fmt.Fprintf(w, "  %s\t%v\t%d\n", c.Name, c.Enabled, c.Replicas)
	}

	if len(mode.Features) > 0 {
		_, _ = fmt.Fprintln(w, "\nFeatures:")
		for _, f := range mode.Features {
			_, _ = fmt.Fprintf(w, "  ✓ %s\n", f)
		}
	}
	_ = w.Flush()
}
