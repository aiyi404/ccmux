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
	if cfg.Mode != "auto" {
		t.Errorf("expected mode 'auto', got %q", cfg.Mode)
	}
}

func TestLoadConfig_FromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	os.WriteFile(path, []byte(`{"mode":"standalone","lang":"zh","current":"myproxy"}`), 0644)
	cfg, err := LoadConfigFrom(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Mode != "standalone" {
		t.Errorf("expected mode 'standalone', got %q", cfg.Mode)
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
	cfg := &AppConfig{Mode: "ccswitch", Lang: "en", Current: "test"}
	err := SaveConfigTo(cfg, path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	loaded, err := LoadConfigFrom(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loaded.Mode != "ccswitch" {
		t.Errorf("expected mode 'ccswitch', got %q", loaded.Mode)
	}
	if loaded.Current != "test" {
		t.Errorf("expected current 'test', got %q", loaded.Current)
	}
}
