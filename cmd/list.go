package cmd

import (
	"github.com/farion1231/ccmux/internal/cli/commands"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all providers/profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		state, err := getState()
		if err != nil {
			return err
		}
		return commands.RunList(state)
	},
}

func init() { rootCmd.AddCommand(listCmd) }
