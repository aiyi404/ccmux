package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_Default(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	cfg, err := LoadConfigFrom(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Lang != "" {
		t.Errorf("expected empty lang, got %q", cfg.Lang)
	}
}

func TestLoadConfig_FromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	os.WriteFile(path, []byte(`{"lang":"zh","current":"myproxy"}`), 0644)
	cfg, err := LoadConfigFrom(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Lang != "zh" {
		t.Errorf("expected lang 'zh', got %q", cfg.Lang)
	}
	if cfg.Current != "myproxy" {
		t.Errorf("expected current 'myproxy', got %q", cfg.Current)
	}
}

func TestSaveConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	cfg := &AppConfig{Lang: "en", Current: "test"}
	err := SaveConfigTo(cfg, path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	loaded, err := LoadConfigFrom(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loaded.Current != "test" {
		t.Errorf("expected current 'test', got %q", loaded.Current)
	}
}
