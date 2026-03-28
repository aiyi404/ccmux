package cmd

import (
	"github.com/aiyi404/ccmux/internal/cli/commands"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Create a new profile interactively (standalone only)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		state, err := getState()
		if err != nil {
			return err
		}
		return commands.RunAdd(state, args[0])
	},
}

func init() { rootCmd.AddCommand(addCmd) }
