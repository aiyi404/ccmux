package cmd

import (
	"github.com/aiyi404/ccmux/internal/cli/commands"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm <name>",
	Short: "Remove a profile (standalone only)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		state, err := getState()
		if err != nil {
			return err
		}
		return commands.RunRm(state, args[0])
	},
}

func init() { rootCmd.AddCommand(rmCmd) }
