package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/farion1231/ccmux/internal/config"
)

type profileFile struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
	Model       string            `json:"model,omitempty"`
}

type StandaloneService struct {
	cfg         *config.AppConfig
	configPath  string
	profilesDir string
}

func NewStandaloneService(cfg *config.AppConfig, configPath, profilesDir string) *StandaloneService {
	return &StandaloneService{
		cfg:         cfg,
		configPath:  configPath,
		profilesDir: profilesDir,
	}
}

// ProfilePath returns the absolute path to a profile JSON file.
func (s *StandaloneService) ProfilePath(name string) string {
	return filepath.Join(s.profilesDir, name+".json")
}

func (s *StandaloneService) loadProfile(path string) (*profileFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var pf profileFile
	if err := json.Unmarshal(data, &pf); err != nil {
		return nil, err
	}
	return &pf, nil
}

func (s *StandaloneService) saveProfile(pf *profileFile) error {
	if err := os.MkdirAll(s.profilesDir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(pf, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.ProfilePath(pf.Name), data, 0644)
}

func (s *StandaloneService) toProvider(pf *profileFile) Provider {
	return Provider{
		Name:        pf.Name,
		Description: pf.Description,
		Env:         pf.Env,
		ModelAlias:  pf.Model,
	}
}

func (s *StandaloneService) List() ([]Provider, error) {
	entries, err := os.ReadDir(s.profilesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var providers []Provider
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		pf, err := s.loadProfile(filepath.Join(s.profilesDir, e.Name()))
		if err != nil {
			continue
		}
		providers = append(providers, s.toProvider(pf))
	}
	return providers, nil
}


func (s *StandaloneService) GetCurrent() (*Provider, error) {
	if s.cfg.Current == "" {
		return nil, fmt.Errorf("%w: no current provider set", ErrNotFound)
	}
	return s.GetByName(s.cfg.Current)
}

func (s *StandaloneService) GetByName(name string) (*Provider, error) {
	providers, err := s.List()
	if err != nil {
		return nil, err
	}
	resolved, err := ResolveName(name, providers)
	if err != nil {
		return nil, err
	}
	pf, err := s.loadProfile(s.ProfilePath(resolved))
	if err != nil {
		return nil, err
	}
	p := s.toProvider(pf)
	return &p, nil
}

func (s *StandaloneService) SetCurrent(name string) error {
	providers, err := s.List()
	if err != nil {
		return err
	}
	resolved, err := ResolveName(name, providers)
	if err != nil {
		return err
	}
	s.cfg.Current = resolved
	return config.SaveConfigTo(s.cfg, s.configPath)
}

func (s *StandaloneService) BuildOverlay(name string) (*SettingsOverlay, error) {
	p, err := s.GetByName(name)
	if err != nil {
		return nil, err
	}
	overlay := &SettingsOverlay{
		Env:   p.Env,
		Model: p.ModelAlias,
	}
	return overlay, nil
}


func (s *StandaloneService) Add(p Provider) error {
	path := s.ProfilePath(p.Name)
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("profile %q already exists", p.Name)
	}
	pf := &profileFile{
		Name:        p.Name,
		Description: p.Description,
		Env:         p.Env,
		Model:       p.ModelAlias,
	}
	return s.saveProfile(pf)
}

func (s *StandaloneService) Edit(name string) error {
	path := s.ProfilePath(name)
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w: '%s'", ErrNotFound, name)
		}
		return err
	}
	return nil
}

func (s *StandaloneService) Remove(name string) error {
	path := s.ProfilePath(name)
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w: '%s'", ErrNotFound, name)
		}
		return err
	}
	if err := os.Remove(path); err != nil {
		return err
	}
	if s.cfg.Current == name {
		s.cfg.Current = ""
		return config.SaveConfigTo(s.cfg, s.configPath)
	}
	return nil
}

func (s *StandaloneService) Import(name string) error {
	data, err := os.ReadFile(config.ClaudeSettings)
	if err != nil {
		return fmt.Errorf("cannot read claude settings: %w", err)
	}
	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return fmt.Errorf("invalid claude settings: %w", err)
	}
	env := make(map[string]string)
	if envRaw, ok := settings["env"].(map[string]interface{}); ok {
		for k, v := range envRaw {
			if vs, ok := v.(string); ok {
				env[k] = vs
			}
		}
	}
	var model string
	if m, ok := settings["model"].(string); ok {
		model = m
	}
	p := Provider{
		Name:       name,
		Env:        env,
		ModelAlias: model,
	}
	return s.Add(p)
}

// Compile-time interface check.
var _ ProviderService = (*StandaloneService)(nil)
