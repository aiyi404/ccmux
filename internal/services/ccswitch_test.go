package services

import (
	"database/sql"
	"testing"

	"github.com/farion1231/ccmux/internal/config"
	_ "modernc.org/sqlite"
)

func setupCCSwitchTest(t *testing.T) *CCSwitchService {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`CREATE TABLE providers (
		id INTEGER PRIMARY KEY, name TEXT, app_type TEXT,
		settings_config TEXT, is_current INTEGER DEFAULT 0, sort_index INTEGER DEFAULT 0
	)`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`INSERT INTO providers (name, app_type, settings_config, is_current, sort_index) VALUES
		('proxy1', 'claude', '{"env":{"ANTHROPIC_BASE_URL":"https://p1.com","ANTHROPIC_MODEL":"claude-sonnet-4-6"}}', 1, 0),
		('proxy2', 'claude', '{"env":{"ANTHROPIC_BASE_URL":"https://p2.com","ANTHROPIC_MODEL":"claude-opus-4-6"}}', 0, 1)`)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	cfg := &config.AppConfig{Mode: "ccswitch"}
	return NewCCSwitchService(cfg, db)
}

func TestCCSwitch_List(t *testing.T) {
	svc := setupCCSwitchTest(t)
	providers, err := svc.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(providers) != 2 {
		t.Errorf("expected 2, got %d", len(providers))
	}
}

func TestCCSwitch_GetCurrent(t *testing.T) {
	svc := setupCCSwitchTest(t)
	p, err := svc.GetCurrent()
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != "proxy1" {
		t.Errorf("expected 'proxy1', got %q", p.Name)
	}
}

func TestCCSwitch_SetCurrent(t *testing.T) {
	svc := setupCCSwitchTest(t)
	err := svc.SetCurrent("proxy2")
	if err != nil {
		t.Fatal(err)
	}
	p, _ := svc.GetCurrent()
	if p.Name != "proxy2" {
		t.Errorf("expected 'proxy2', got %q", p.Name)
	}
}

func TestCCSwitch_AddNotSupported(t *testing.T) {
	svc := setupCCSwitchTest(t)
	err := svc.Add(Provider{Name: "test"})
	if err != ErrNotSupported {
		t.Errorf("expected ErrNotSupported, got %v", err)
	}
}
