package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/farion1231/ccmux/internal/services"
	"github.com/farion1231/ccmux/internal/store"
)

func RunEdit(state *store.AppState, name string) error {
	svc, ok := state.Service.(*services.StandaloneService)
	if !ok {
		return services.ErrNotSupported
	}
	p, err := state.Service.GetByName(name)
	if err != nil {
		return err
	}
	path := svc.ProfilePath(p.Name)
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("editor failed: %w", err)
	}
	fmt.Printf("✓ profile '%s' updated\n", p.Name)
	return nil
}
