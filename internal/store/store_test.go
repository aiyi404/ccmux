package store

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectMode_FlagOverride(t *testing.T) {
	mode := DetectMode("standalone", "", "", "")
	if mode != "standalone" {
		t.Errorf("expected 'standalone', got %q", mode)
	}
}

func TestDetectMode_EnvOverride(t *testing.T) {
	mode := DetectMode("", "ccswitch", "", "")
	if mode != "ccswitch" {
		t.Errorf("expected 'ccswitch', got %q", mode)
	}
}

func TestDetectMode_ConfigOverride(t *testing.T) {
	mode := DetectMode("", "", "standalone", "")
	if mode != "standalone" {
		t.Errorf("expected 'standalone', got %q", mode)
	}
}

func TestDetectMode_AutoWithDB(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "cc-switch.db")
	os.WriteFile(dbPath, []byte("fake"), 0644)
	mode := DetectMode("", "", "auto", dbPath)
	if mode != "ccswitch" {
		t.Errorf("expected 'ccswitch', got %q", mode)
	}
}

func TestDetectMode_AutoWithoutDB(t *testing.T) {
	mode := DetectMode("", "", "auto", "/nonexistent/path.db")
	if mode != "standalone" {
		t.Errorf("expected 'standalone', got %q", mode)
	}
}
