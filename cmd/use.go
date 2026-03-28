package cmd

import (
	"github.com/farion1231/ccmux/internal/cli/commands"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use <name> [-- claude-args...]",
	Short: "Launch claude with the given profile (session-level, no global change)",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		state, err := getState()
		if err != nil {
			return err
		}
		name := args[0]
		var extraArgs []string
		if cmd.ArgsLenAtDash() >= 0 {
			extraArgs = args[cmd.ArgsLenAtDash():]
		} else if len(args) > 1 {
			extraArgs = args[1:]
		}
		return commands.RunUse(state, name, extraArgs)
	},
}

func init() { rootCmd.AddCommand(useCmd) }
