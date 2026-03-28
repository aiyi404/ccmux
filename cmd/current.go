package cmd

import (
	"github.com/farion1231/ccmux/internal/cli/commands"
	"github.com/spf13/cobra"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show the active provider",
	RunE: func(cmd *cobra.Command, args []string) error {
		state, err := getState()
		if err != nil {
			return err
		}
		return commands.RunCurrent(state)
	},
}

func init() { rootCmd.AddCommand(currentCmd) }
