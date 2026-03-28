package cmd

import (
	"github.com/farion1231/ccmux/internal/cli/commands"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit <name>",
	Short: "Edit a profile with $EDITOR (standalone only)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		state, err := getState()
		if err != nil {
			return err
		}
		return commands.RunEdit(state, args[0])
	},
}

func init() { rootCmd.AddCommand(editCmd) }
