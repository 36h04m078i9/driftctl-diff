package cmd

import (
	""

	"github.com/spf13/cobra"
	"github.com/user/driftctl-diff/internal/drift"
	"github.com/user/driftctl-diff/internal/output"
	"github.com/user/driftctl-diff/internal/state"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Compare Terraform state against live cloud resources",
	Example: `  driftctl-diff diff --state terraform.tfstate --region us-west-2
  driftctl-diff diff -s ./infra/terraform.tfstate -o json`,
	RunE: runDiff,
}

func runDiff(cmd *cobra.Command, args []string) error {
	parser := state.NewParser()
	tfState, err := parser.ParseFile(statePath)
	if err != nil {
		return fmt.Errorf("failed to parse state file %q: %w", statePath, err)
	}

	detector := drift.NewDetector(region)
	results, err := detector.Detect(tfState)
	if err != nil {
		return fmt.Errorf("drift detection failed: %w", err)
	}

	printer := output.NewPrinter(outputFmt, os.Stdout)
	if err := printer.Print(results); err != nil {
		return fmt.Errorf("failed to render output: %w", err)
	}

	if results.HasDrift() {
		os.Exit(2)
	}
	return nil
}
