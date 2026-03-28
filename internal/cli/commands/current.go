package commands

import (
	"fmt"

	"github.com/aiyi404/ccmux/internal/store"
)

func RunCurrent(state *store.AppState) error {
	p, err := state.Service.GetCurrent()
	if err != nil {
		return fmt.Errorf("no active provider")
	}
	fmt.Printf("→ %s\n", p.Name)
	fmt.Printf("  base_url: %s\n", p.Env["ANTHROPIC_BASE_URL"])
	fmt.Printf("  model:    %s\n", p.Env["ANTHROPIC_MODEL"])
	return nil
}
