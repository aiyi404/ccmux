package commands

import (
	"fmt"
	"strings"

	"github.com/aiyi404/ccmux/internal/cli/i18n"
	"github.com/aiyi404/ccmux/internal/store"
)

func RunList(state *store.AppState) error {
	providers, err := state.Service.List()
	if err != nil {
		return err
	}
	if len(providers) == 0 {
		fmt.Println(i18n.T("err_no_providers"))
		return nil
	}
	current, _ := state.Service.GetCurrent()
	currentName := ""
	if current != nil {
		currentName = current.Name
	}
	fmt.Printf("  %-18s %-36s %s\n", "NAME", "BASE_URL", "MODEL")
	for _, p := range providers {
		marker := "  "
		if p.Name == currentName {
			marker = "→ "
		}
		url := p.Env["ANTHROPIC_BASE_URL"]
		url = strings.TrimPrefix(url, "http://")
		url = strings.TrimPrefix(url, "https://")
		model := p.Env["ANTHROPIC_MODEL"]
		fmt.Printf("%s%-18s %-36s %s\n", marker, p.Name, url, model)
	}
	return nil
}
