package cmd

import (
	"fmt"
	"os"

	"github.com/snyk/driftctl-diff/internal/baseline"
	"github.com/spf13/cobra"
)

var baselineCmd = &cobra.Command{
	Use:   "baseline",
	Short: "Manage the drift baseline",
}

var baselineAddCmd = &cobra.Command{
	Use:   "add <resource-type> <resource-id> <attribute>",
	Short: "Acknowledge a drift item and add it to the baseline",
	Args:  cobra.ExactArgs(3),
	RunE:  runBaselineAdd,
}

var baselinePathFlag string

func init() {
	baselineCmd.PersistentFlags().StringVar(&baselinePathFlag, "baseline-file", ".driftctl-baseline.json", "path to baseline file")
	baselineCmd.AddCommand(baselineAddCmd)
	rootCmd.AddCommand(baselineCmd)
}

func runBaselineAdd(cmd *cobra.Command, args []string) error {
	resourceType, resourceID, attribute := args[0], args[1], args[2]

	var b *baseline.Baseline
	loaded, err := baseline.LoadFrom(baselinePathFlag)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("loading baseline: %w", err)
		}
		b = baseline.New()
	} else {
		b = loaded
	}

	if b.Contains(resourceType, resourceID, attribute) {
		fmt.Fprintf(cmd.OutOrStdout(), "already acknowledged: %s.%s[%s]\n", resourceType, resourceID, attribute)
		return nil
	}

	b.Add(resourceType, resourceID, attribute)
	if err := b.SaveTo(baselinePathFlag); err != nil {
		return fmt.Errorf("saving baseline: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "added to baseline: %s.%s[%s]\n", resourceType, resourceID, attribute)
	return nil
}
