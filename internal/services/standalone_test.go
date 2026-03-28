package services

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/farion1231/ccmux/internal/config"
)

func setupStandaloneTest(t *testing.T) (*StandaloneService, string) {
	t.Helper()
	dir := t.TempDir()
	profilesDir := filepath.Join(dir, "profiles")
	os.MkdirAll(profilesDir, 0755)
	configPath := filepath.Join(dir, "config.json")
	cfg := &config.AppConfig{Mode: "standalone"}
	svc := NewStandaloneService(cfg, configPath, profilesDir)
	return svc, profilesDir
}

func writeProfile(t *testing.T, dir, name string) {
	t.Helper()
	p := map[string]interface{}{
		"name": name,
		"env": map[string]string{
			"ANTHROPIC_BASE_URL":   "https://api.example.com",
			"ANTHROPIC_AUTH_TOKEN": "sk-test1234",
			"ANTHROPIC_MODEL":     "claude-sonnet-4-6",
		},
	}
	data, _ := json.Marshal(p)
	os.WriteFile(filepath.Join(dir, name+".json"), data, 0644)
}

func TestStandalone_List(t *testing.T) {
	svc, dir := setupStandaloneTest(t)
	writeProfile(t, dir, "proxy1")
	writeProfile(t, dir, "proxy2")
	providers, err := svc.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(providers) != 2 {
		t.Errorf("expected 2 providers, got %d", len(providers))
	}
}

func TestStandalone_GetByName(t *testing.T) {
	svc, dir := setupStandaloneTest(t)
	writeProfile(t, dir, "myproxy")
	p, err := svc.GetByName("my")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "myproxy" {
		t.Errorf("expected 'myproxy', got %q", p.Name)
	}
}

func TestStandalone_Add(t *testing.T) {
	svc, dir := setupStandaloneTest(t)
	p := Provider{
		Name: "newproxy",
		Env: map[string]string{
			"ANTHROPIC_BASE_URL":   "https://new.example.com",
			"ANTHROPIC_AUTH_TOKEN": "sk-new",
			"ANTHROPIC_MODEL":     "claude-sonnet-4-6",
		},
	}
	err := svc.Add(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "newproxy.json")); err != nil {
		t.Errorf("profile file not created: %v", err)
	}
}

func TestStandalone_Remove(t *testing.T) {
	svc, dir := setupStandaloneTest(t)
	writeProfile(t, dir, "todelete")
	err := svc.Remove("todelete")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "todelete.json")); !os.IsNotExist(err) {
		t.Error("profile file should be deleted")
	}
}

func TestStandalone_BuildOverlay(t *testing.T) {
	svc, dir := setupStandaloneTest(t)
	writeProfile(t, dir, "myproxy")
	overlay, err := svc.BuildOverlay("myproxy")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if overlay.Env["ANTHROPIC_BASE_URL"] != "https://api.example.com" {
		t.Errorf("unexpected base_url: %s", overlay.Env["ANTHROPIC_BASE_URL"])
	}
}
