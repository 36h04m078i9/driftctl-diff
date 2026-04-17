package cmd

import (
	"fmt"
	"os"

	"github.com/owner/driftctl-diff/internal/remediate"
	"github.com/owner/driftctl-diff/internal/runner"
	"github.com/spf13/cobra"
)

var remediateCmd = &cobra.Command{
	Use:   "remediate",
	Short: "Print remediation suggestions for detected drift",
	RunE:  runRemediate,
}

func init() {
	remediateCmd.Flags().StringP("state", "s", "terraform.tfstate", "path to Terraform state file")
	rootCmd.AddCommand(remediateCmd)
}

func runRemediate(cmd *cobra.Command, _ []string) error {
	statePath, err := cmd.Flags().GetString("state")
	if err != nil {
		return fmt.Errorf("reading flag: %w", err)
	}

	r := runner.New(statePath, nil)
	results, err := r.Run(cmd.Context())
	if err != nil {
		return fmt.Errorf("running drift detection: %w", err)
	}

	suggestions := remediate.Suggest(results)
	p := remediate.New(os.Stdout)
	p.Print(suggestions)
	return nil
}
