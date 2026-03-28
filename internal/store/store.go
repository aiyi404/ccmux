package store

import (
	"os"

	"github.com/aiyi404/ccmux/internal/cli/i18n"
	"github.com/aiyi404/ccmux/internal/config"
	"github.com/aiyi404/ccmux/internal/services"
)

type AppState struct {
	Config  *config.AppConfig
	Service services.ProviderService
	Lang    string
}

func New() (*AppState, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		cfg = &config.AppConfig{}
	}

	lang := i18n.DetectLang(cfg.Lang, os.Getenv("CCC_LANG"))
	i18n.SetLang(lang)

	svc := services.NewStandaloneService(cfg, config.CCCConfig, config.CCCProfiles)

	return &AppState{
		Config:  cfg,
		Service: svc,
		Lang:    lang,
	}, nil
}
