package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	statePath  string
	region     string
	outputFmt  string
)

var rootCmd = &cobra.Command{
	Use:   "driftctl-diff",
	Short: "Surface infrastructure drift between Terraform state and live cloud resources",
	Long: `driftctl-diff compares your Terraform state file against live cloud resources
and presents the differences in a human-readable diff format.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&statePath, "state", "s", "terraform.tfstate", "path to Terraform state file")
	rootCmd.PersistentFlags().StringVarP(&region, "region", "r", "us-east-1", "AWS region to scan")
	rootCmd.PersistentFlags().StringVarP(&outputFmt, "output", "o", "text", "output format: text, json")

	rootCmd.AddCommand(diffCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
