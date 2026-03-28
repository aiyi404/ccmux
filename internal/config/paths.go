package config

import (
	"os"
	"path/filepath"
)

var (
	HomeDir          = os.Getenv("HOME")
	CCswitchDB       = filepath.Join(HomeDir, ".cc-switch", "cc-switch.db")
	CCswitchSettings = filepath.Join(HomeDir, ".cc-switch", "settings.json")
	ClaudeSettings   = filepath.Join(HomeDir, ".claude", "settings.json")
	BackupDir        = filepath.Join(HomeDir, ".claude", "backups")
	CCCDir           = filepath.Join(HomeDir, ".config", "ccc")
	CCCConfig        = filepath.Join(HomeDir, ".config", "ccc", "config.json")
	CCCProfiles      = filepath.Join(HomeDir, ".config", "ccc", "profiles")
)

const AppType = "claude"
