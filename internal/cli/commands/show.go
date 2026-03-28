package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/farion1231/ccmux/internal/store"
)

func maskToken(val string) string {
	if len(val) <= 8 {
		return "****"
	}
	return val[:4] + "****" + val[len(val)-4:]
}

func RunShow(state *store.AppState, name string) error {
	p, err := state.Service.GetByName(name)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", p.Name)
	if p.Description != "" {
		fmt.Printf("  description: %s\n", p.Description)
	}
	keys := make([]string, 0, len(p.Env))
	for k := range p.Env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := p.Env[k]
		if strings.Contains(strings.ToUpper(k), "TOKEN") || strings.Contains(strings.ToUpper(k), "KEY") {
			v = maskToken(v)
		}
		fmt.Printf("  %s: %s\n", k, v)
	}
	if p.ModelAlias != "" {
		fmt.Printf("  model_alias: %s\n", p.ModelAlias)
	}
	return nil
}
