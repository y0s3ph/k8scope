package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove the observability stack",
	Long:  `Uninstall all k8scope components from the cluster.`,
	RunE:  runUninstall,
}

func init() {
	uninstallCmd.Flags().Bool("confirm", false, "skip confirmation prompt")
	rootCmd.AddCommand(uninstallCmd)
}

func runUninstall(cmd *cobra.Command, args []string) error {
	// TODO: implement uninstall logic (Phase 2)
	fmt.Println("⚠️  Uninstall not yet implemented. Coming soon!")
	return nil
}
