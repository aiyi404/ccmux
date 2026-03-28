package cmd

import (
	"fmt"
	"os"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/farion1231/ccmux/internal/cli/tui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Interactive TUI menu",
	RunE: func(cmd *cobra.Command, args []string) error {
		state, err := getState()
		if err != nil {
			return err
		}
		m := tui.NewModel(state)
		p := tea.NewProgram(m, tea.WithAltScreen())
		finalModel, err := p.Run()
		if err != nil {
			return err
		}
		if fm, ok := finalModel.(tui.Model); ok {
			if sig := fm.GetExecSignal(); sig != nil && sig.Result != nil {
				fmt.Printf("▸ launching claude with session profile\n")
				return syscall.Exec(sig.Result.Binary, sig.Result.Args, os.Environ())
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
