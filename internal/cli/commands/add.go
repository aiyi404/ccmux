package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/farion1231/ccmux/internal/services"
	"github.com/farion1231/ccmux/internal/store"
)

func RunAdd(state *store.AppState, name string) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Creating profile: %s\n\n", name)
	fmt.Print("Description (optional): ")
	desc, _ := reader.ReadString('\n')
	desc = strings.TrimSpace(desc)
	fmt.Print("ANTHROPIC_BASE_URL: ")
	baseURL, _ := reader.ReadString('\n')
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" {
		return fmt.Errorf("base_url is required")
	}
	fmt.Print("ANTHROPIC_AUTH_TOKEN: ")
	token, _ := reader.ReadString('\n')
	token = strings.TrimSpace(token)
	if token == "" {
		return fmt.Errorf("auth_token is required")
	}
	fmt.Print("ANTHROPIC_MODEL: ")
	model, _ := reader.ReadString('\n')
	model = strings.TrimSpace(model)
	if model == "" {
		return fmt.Errorf("model is required")
	}
	fmt.Printf("ANTHROPIC_DEFAULT_HAIKU_MODEL (Enter = %s): ", model)
	haiku, _ := reader.ReadString('\n')
	haiku = strings.TrimSpace(haiku)
	if haiku == "" {
		haiku = model
	}
	fmt.Printf("ANTHROPIC_DEFAULT_OPUS_MODEL (Enter = %s): ", model)
	opus, _ := reader.ReadString('\n')
	opus = strings.TrimSpace(opus)
	if opus == "" {
		opus = model
	}
	fmt.Printf("ANTHROPIC_DEFAULT_SONNET_MODEL (Enter = %s): ", model)
	sonnet, _ := reader.ReadString('\n')
	sonnet = strings.TrimSpace(sonnet)
	if sonnet == "" {
		sonnet = model
	}
	fmt.Printf("ANTHROPIC_REASONING_MODEL (Enter = %s): ", model)
	reasoning, _ := reader.ReadString('\n')
	reasoning = strings.TrimSpace(reasoning)
	if reasoning == "" {
		reasoning = model
	}
	fmt.Print("Model alias (e.g. opus[1m], Enter to skip): ")
	modelAlias, _ := reader.ReadString('\n')
	modelAlias = strings.TrimSpace(modelAlias)
	p := services.Provider{
		Name:        name,
		Description: desc,
		Env: map[string]string{
			"ANTHROPIC_BASE_URL":              baseURL,
			"ANTHROPIC_AUTH_TOKEN":            token,
			"ANTHROPIC_MODEL":                 model,
			"ANTHROPIC_DEFAULT_HAIKU_MODEL":   haiku,
			"ANTHROPIC_DEFAULT_OPUS_MODEL":    opus,
			"ANTHROPIC_DEFAULT_SONNET_MODEL":  sonnet,
			"ANTHROPIC_REASONING_MODEL":       reasoning,
		},
		ModelAlias: modelAlias,
	}
	if err := state.Service.Add(p); err != nil {
		return err
	}
	fmt.Printf("✓ profile '%s' created\n", name)
	return nil
}
