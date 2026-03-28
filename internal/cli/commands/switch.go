package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aiyi404/ccmux/internal/config"
	"github.com/aiyi404/ccmux/internal/store"
)

func backupSettings() error {
	if _, err := os.Stat(config.ClaudeSettings); os.IsNotExist(err) {
		return nil
	}
	os.MkdirAll(config.BackupDir, 0755)
	ts := time.Now().UnixMilli()
	dst := filepath.Join(config.BackupDir, fmt.Sprintf(".claude.json.backup.%d", ts))
	data, err := os.ReadFile(config.ClaudeSettings)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

func RunSwitch(state *store.AppState, name string) error {
	overlay, err := state.Service.BuildOverlay(name)
	if err != nil {
		return err
	}
	if err := backupSettings(); err != nil {
		return fmt.Errorf("backup failed: %w", err)
	}
	p, _ := state.Service.GetByName(name)
	providerName := name
	if p != nil {
		providerName = p.Name
	}
	existing := make(map[string]interface{})
	if data, err := os.ReadFile(config.ClaudeSettings); err == nil {
		json.Unmarshal(data, &existing)
	}
	existingEnv, _ := existing["env"].(map[string]interface{})
	if existingEnv == nil {
		existingEnv = make(map[string]interface{})
	}
	for k, v := range overlay.Env {
		existingEnv[k] = v
	}
	existing["env"] = existingEnv
	if overlay.Model != "" {
		existing["model"] = overlay.Model
	}
	data, err := json.MarshalIndent(existing, "", "  ")
	if err != nil {
		return err
	}
	os.MkdirAll(filepath.Dir(config.ClaudeSettings), 0755)
	if err := os.WriteFile(config.ClaudeSettings, data, 0644); err != nil {
		return err
	}
	if err := state.Service.SetCurrent(providerName); err != nil {
		return err
	}
	fmt.Printf("✓ switched to '%s'\n", providerName)
	return nil
}
