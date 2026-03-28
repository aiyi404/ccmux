package services

import (
	"database/sql"

	"github.com/farion1231/ccmux/internal/config"
	"github.com/farion1231/ccmux/internal/database/dao"
)

type CCSwitchService struct {
	cfg *config.AppConfig
	dao *dao.ProviderDAO
	db  *sql.DB
}

func NewCCSwitchService(cfg *config.AppConfig, db *sql.DB) *CCSwitchService {
	return &CCSwitchService{
		cfg: cfg,
		dao: dao.NewProviderDAO(db, config.AppType),
		db:  db,
	}
}

func (s *CCSwitchService) rowToProvider(r *dao.ProviderRow) (*Provider, error) {
	env, model, err := dao.ParseSettingsConfig(r.SettingsConfig)
	if err != nil {
		return nil, err
	}
	return &Provider{Name: r.Name, Env: env, ModelAlias: model}, nil
}

func (s *CCSwitchService) List() ([]Provider, error) {
	rows, err := s.dao.ListAll()
	if err != nil {
		return nil, err
	}
	var providers []Provider
	for _, r := range rows {
		p, err := s.rowToProvider(&r)
		if err != nil {
			continue
		}
		providers = append(providers, *p)
	}
	return providers, nil
}

func (s *CCSwitchService) GetCurrent() (*Provider, error) {
	r, err := s.dao.GetCurrent()
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, ErrNotFound
	}
	return s.rowToProvider(r)
}

func (s *CCSwitchService) GetByName(name string) (*Provider, error) {
	providers, err := s.List()
	if err != nil {
		return nil, err
	}
	resolved, err := ResolveName(name, providers)
	if err != nil {
		return nil, err
	}
	for _, p := range providers {
		if p.Name == resolved {
			return &p, nil
		}
	}
	return nil, ErrNotFound
}

func (s *CCSwitchService) SetCurrent(name string) error {
	p, err := s.GetByName(name)
	if err != nil {
		return err
	}
	return s.dao.SetCurrent(p.Name)
}

func (s *CCSwitchService) BuildOverlay(name string) (*SettingsOverlay, error) {
	p, err := s.GetByName(name)
	if err != nil {
		return nil, err
	}
	overlay := &SettingsOverlay{Env: p.Env}
	if p.ModelAlias != "" {
		overlay.Model = p.ModelAlias
	}
	return overlay, nil
}

func (s *CCSwitchService) Add(_ Provider) error { return ErrNotSupported }
func (s *CCSwitchService) Edit(_ string) error   { return ErrNotSupported }
func (s *CCSwitchService) Remove(_ string) error { return ErrNotSupported }
func (s *CCSwitchService) Import(_ string) error { return ErrNotSupported }
