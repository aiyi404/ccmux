package cmd

import (
	"github.com/farion1231/ccmux/internal/store"
	"github.com/spf13/cobra"
)

var (
	flagStandalone bool
	flagCCSwitch   bool
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
	rootCmd.PersistentFlags().BoolVar(&flagStandalone, "standalone", false, "Force standalone mode")
	rootCmd.PersistentFlags().BoolVar(&flagCCSwitch, "cc-switch", false, "Force cc-switch mode")
}

func getState() (*store.AppState, error) {
	flagMode := ""
	if flagStandalone {
		flagMode = "standalone"
	} else if flagCCSwitch {
		flagMode = "ccswitch"
	}
	return store.New(flagMode)
}
