// internal/cli/i18n/i18n_test.go
package i18n

import "testing"

func TestT_English(t *testing.T) {
	SetLang("en")
	if got := T("provider"); got != "Provider" {
		t.Errorf("expected 'Provider', got %q", got)
	}
}

func TestT_Chinese(t *testing.T) {
	SetLang("zh")
	if got := T("provider"); got != "服务商" {
		t.Errorf("expected '服务商', got %q", got)
	}
}

func TestT_Fallback(t *testing.T) {
	SetLang("en")
	if got := T("nonexistent_key"); got != "nonexistent_key" {
		t.Errorf("expected key itself as fallback, got %q", got)
	}
}

func TestDetectLang(t *testing.T) {
	lang := DetectLang("", "")
	if lang != "en" {
		t.Errorf("expected 'en', got %q", lang)
	}
	lang = DetectLang("zh", "")
	if lang != "zh" {
		t.Errorf("expected 'zh', got %q", lang)
	}
}
