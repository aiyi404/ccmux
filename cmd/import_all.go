package cmd

import (
	"github.com/aiyi404/ccmux/internal/cli/commands"
	"github.com/spf13/cobra"
)

var importAllCmd = &cobra.Command{
	Use:   "import-all",
	Short: "Import all providers from cc-switch database",
	RunE: func(cmd *cobra.Command, args []string) error {
		return commands.RunImportAllCLI()
	},
}

func init() {
	rootCmd.AddCommand(importAllCmd)
}
