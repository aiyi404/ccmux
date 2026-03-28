package store

import (
	"os"

	"github.com/farion1231/ccmux/internal/cli/i18n"
	"github.com/farion1231/ccmux/internal/config"
	"github.com/farion1231/ccmux/internal/database"
	"github.com/farion1231/ccmux/internal/services"
)

type AppState struct {
	Mode    string
	Config  *config.AppConfig
	Service services.ProviderService
	Lang    string
}

// DetectMode determines the operating mode by priority:
// 1. CLI flag, 2. env var, 3. config file, 4. auto-detect (DB exists -> ccswitch)
func DetectMode(flagMode, envMode, configMode, dbPath string) string {
	if flagMode == "standalone" || flagMode == "ccswitch" {
		return flagMode
	}
	if envMode == "standalone" || envMode == "ccswitch" {
		return envMode
	}
	if configMode != "" && configMode != "auto" {
		return configMode
	}
	if _, err := os.Stat(dbPath); err == nil {
		return "ccswitch"
	}
	return "standalone"
}

// New creates an AppState with the appropriate service based on detected mode
func New(flagMode string) (*AppState, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		cfg = &config.AppConfig{Mode: "auto"}
	}

	envMode := os.Getenv("CCC_MODE")
	mode := DetectMode(flagMode, envMode, cfg.Mode, config.CCswitchDB)

	lang := i18n.DetectLang(cfg.Lang, os.Getenv("CCC_LANG"))
	i18n.SetLang(lang)

	var svc services.ProviderService
	if mode == "ccswitch" {
		db, err := database.Open(config.CCswitchDB)
		if err != nil {
			return nil, err
		}
		svc = services.NewCCSwitchService(cfg, db)
	} else {
		svc = services.NewStandaloneService(cfg, config.CCCConfig, config.CCCProfiles)
	}

	return &AppState{
		Mode:    mode,
		Config:  cfg,
		Service: svc,
		Lang:    lang,
	}, nil
}
