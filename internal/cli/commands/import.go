package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aiyi404/ccmux/internal/config"
	"github.com/aiyi404/ccmux/internal/store"
)

func RunImport(state *store.AppState, name string) error {
	if name == "" {
		env, err := readCurrentSettingsEnv()
		if err != nil {
			return fmt.Errorf("cannot read settings.json: %w", err)
		}
		url := env["ANTHROPIC_BASE_URL"]
		if url != "" {
			name = strings.TrimPrefix(url, "http://")
			name = strings.TrimPrefix(name, "https://")
			if idx := strings.IndexAny(name, ":/"); idx >= 0 {
				name = name[:idx]
			}
			name = strings.Map(func(r rune) rune {
				if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '.' || r == '-' || r == '_' {
					return r
				}
				return '-'
			}, name)
		}
		if name == "" {
			return fmt.Errorf("cannot derive name from settings, provide one: ccc import <name>")
		}
	}
	if err := state.Service.Import(name); err != nil {
		return err
	}
	fmt.Printf("✓ imported current settings as profile '%s'\n", name)
	return nil
}

func readCurrentSettingsEnv() (map[string]string, error) {
	data, err := os.ReadFile(config.ClaudeSettings)
	if err != nil {
		return nil, err
	}
	var settings struct {
		Env map[string]string `json:"env"`
	}
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, err
	}
	return settings.Env, nil
}
