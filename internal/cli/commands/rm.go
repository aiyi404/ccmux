package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/aiyi404/ccmux/internal/store"
)

func RunRm(state *store.AppState, name string) error {
	p, err := state.Service.GetByName(name)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Remove profile '%s'? [y/N] ", p.Name)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))
	if confirm != "y" {
		fmt.Println("cancelled")
		return nil
	}
	if err := state.Service.Remove(p.Name); err != nil {
		return err
	}
	fmt.Printf("✓ profile '%s' removed\n", p.Name)
	return nil
}
