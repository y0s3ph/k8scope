package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the current state of the observability stack",
	Long:  `Check the health and status of all k8scope components in the cluster.`,
	RunE:  runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	// TODO: implement status checks (Phase 2)
	fmt.Println("⚠️  Status command not yet implemented. Coming soon!")
	return nil
}
