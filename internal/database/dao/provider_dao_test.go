package dao

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`CREATE TABLE providers (
		id INTEGER PRIMARY KEY,
		name TEXT,
		app_type TEXT,
		settings_config TEXT,
		is_current INTEGER DEFAULT 0,
		sort_index INTEGER DEFAULT 0
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
	return db
}

func TestDAO_ListAll(t *testing.T) {
	db := setupTestDB(t)
	dao := NewProviderDAO(db, "claude")
	rows, err := dao.ListAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Errorf("expected 2 rows, got %d", len(rows))
	}
	if rows[0].Name != "proxy1" {
		t.Errorf("expected 'proxy1', got %q", rows[0].Name)
	}
}

func TestDAO_GetCurrent(t *testing.T) {
	db := setupTestDB(t)
	dao := NewProviderDAO(db, "claude")
	cur, err := dao.GetCurrent()
	if err != nil {
		t.Fatal(err)
	}
	if cur == nil || cur.Name != "proxy1" {
		t.Errorf("expected current 'proxy1'")
	}
}

func TestDAO_SetCurrent(t *testing.T) {
	db := setupTestDB(t)
	dao := NewProviderDAO(db, "claude")
	err := dao.SetCurrent("proxy2")
	if err != nil {
		t.Fatal(err)
	}
	cur, _ := dao.GetCurrent()
	if cur == nil || cur.Name != "proxy2" {
		t.Errorf("expected current 'proxy2'")
	}
}

func TestParseSettingsConfig(t *testing.T) {
	env, model, err := ParseSettingsConfig(`{"env":{"ANTHROPIC_BASE_URL":"https://test.com"},"model":"opus"}`)
	if err != nil {
		t.Fatal(err)
	}
	if env["ANTHROPIC_BASE_URL"] != "https://test.com" {
		t.Errorf("unexpected base_url")
	}
	if model != "opus" {
		t.Errorf("expected model 'opus', got %q", model)
	}
}
