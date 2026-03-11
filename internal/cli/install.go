package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/y0s3ph/k8scope/internal/config"
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
	_ = installCmd.MarkFlagRequired("mode")
	rootCmd.AddCommand(installCmd)
}

func runInstall(cmd *cobra.Command, args []string) error {
	modeName, _ := cmd.Flags().GetString("mode")
	namespace, _ := cmd.Flags().GetString("namespace")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	mode, err := config.GetMode(modeName)
	if err != nil {
		return err
	}

	if dryRun {
		fmt.Printf("🔍 Dry run: showing installation plan for mode %q\n\n", mode.Name)
	} else {
		fmt.Printf("🚀 Installing k8scope in mode %q\n\n", mode.Name)
	}

	printInstallPlan(mode, namespace)

	if dryRun {
		fmt.Println("\nNo changes applied (dry-run mode).")
		return nil
	}

	// TODO: actual Helm-based installation (Phase 2)
	fmt.Println("\n⚠️  Installation engine not yet implemented. Coming soon!")
	return nil
}

func printInstallPlan(mode config.Mode, namespace string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintf(w, "Namespace:\t%s\n", namespace)
	fmt.Fprintf(w, "Mode:\t%s\n", mode.Name)
	fmt.Fprintf(w, "Description:\t%s\n", mode.Description)
	fmt.Fprintln(w, "\nComponents:")
	fmt.Fprintf(w, "  COMPONENT\tENABLED\tREPLICAS\n")
	for _, c := range mode.Components {
		fmt.Fprintf(w, "  %s\t%v\t%d\n", c.Name, c.Enabled, c.Replicas)
	}

	if len(mode.Features) > 0 {
		fmt.Fprintln(w, "\nFeatures:")
		for _, f := range mode.Features {
			fmt.Fprintf(w, "  ✓ %s\n", f)
		}
	}
	w.Flush()
}
