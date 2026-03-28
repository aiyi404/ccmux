package commands

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func setupImportAllTest(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	profilesDir := filepath.Join(dir, "profiles")
	os.MkdirAll(profilesDir, 0755)
	dbPath := filepath.Join(dir, "cc-switch.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE providers (
		id INTEGER PRIMARY KEY, name TEXT, app_type TEXT,
		settings_config TEXT, is_current INTEGER DEFAULT 0, sort_index INTEGER DEFAULT 0
	)`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`INSERT INTO providers (name, app_type, settings_config, sort_index) VALUES
		('proxy1', 'claude', '{"env":{"ANTHROPIC_BASE_URL":"https://p1.com","ANTHROPIC_MODEL":"claude-sonnet-4-6","ANTHROPIC_AUTH_TOKEN":"sk-111"}}', 0),
		('proxy2', 'claude', '{"env":{"ANTHROPIC_BASE_URL":"https://p2.com","ANTHROPIC_MODEL":"claude-opus-4-6","ANTHROPIC_AUTH_TOKEN":"sk-222"}}', 1)`)
	if err != nil {
		t.Fatal(err)
	}

	return dbPath, profilesDir
}

func TestImportAll_Basic(t *testing.T) {
	dbPath, profilesDir := setupImportAllTest(t)
	imported, skipped, err := RunImportAll(dbPath, profilesDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if imported != 2 {
		t.Errorf("expected 2 imported, got %d", imported)
	}
	if skipped != 0 {
		t.Errorf("expected 0 skipped, got %d", skipped)
	}
	files, _ := filepath.Glob(filepath.Join(profilesDir, "*.json"))
	if len(files) != 2 {
		t.Errorf("expected 2 profile files, got %d", len(files))
	}
}

func TestImportAll_Dedup(t *testing.T) {
	dbPath, profilesDir := setupImportAllTest(t)
	// pre-create a profile with same base_url+model as proxy1
	existing := map[string]interface{}{
		"name": "my-existing",
		"env": map[string]string{
			"ANTHROPIC_BASE_URL": "https://p1.com",
			"ANTHROPIC_MODEL":    "claude-sonnet-4-6",
		},
	}
	data, _ := json.Marshal(existing)
	os.WriteFile(filepath.Join(profilesDir, "my-existing.json"), data, 0644)

	imported, skipped, err := RunImportAll(dbPath, profilesDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if imported != 1 {
		t.Errorf("expected 1 imported, got %d", imported)
	}
	if skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", skipped)
	}
}
