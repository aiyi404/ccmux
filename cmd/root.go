package cmd

import (
	"github.com/aiyi404/ccmux/internal/store"
	"github.com/aiyi404/ccmux/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "ccc",
	Short:   "Claude Code Provider Multiplexer",
	Version: version.Version,
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
