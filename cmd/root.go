package cmd

import (
	"github.com/farion1231/ccmux/internal/store"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "ccc",
	Short:   "Claude Code Provider Multiplexer",
	Version: "0.2.0",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return tuiCmd.RunE(cmd, args)
	}
}

func getState() (*store.AppState, error) {
	return store.New()
}
