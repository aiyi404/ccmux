package cmd

import (
	"github.com/aiyi404/ccmux/internal/cli/commands"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import [name]",
	Short: "Import current settings.json as a profile",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		state, err := getState()
		if err != nil {
			return err
		}
		name := ""
		if len(args) > 0 {
			name = args[0]
		}
		return commands.RunImport(state, name)
	},
}

func init() { rootCmd.AddCommand(importCmd) }
