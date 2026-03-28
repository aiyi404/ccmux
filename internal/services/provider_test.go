package services

import "testing"

func TestResolveName_ExactMatch(t *testing.T) {
	providers := []Provider{
		{Name: "MyProxy"},
		{Name: "OpenRouter"},
	}
	got, err := ResolveName("myproxy", providers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "MyProxy" {
		t.Errorf("expected 'MyProxy', got %q", got)
	}
}

func TestResolveName_PrefixMatch(t *testing.T) {
	providers := []Provider{
		{Name: "MyProxy"},
		{Name: "OpenRouter"},
	}
	got, err := ResolveName("my", providers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "MyProxy" {
		t.Errorf("expected 'MyProxy', got %q", got)
	}
}

func TestResolveName_Ambiguous(t *testing.T) {
	providers := []Provider{
		{Name: "MyProxy1"},
		{Name: "MyProxy2"},
	}
	_, err := ResolveName("my", providers)
	if err == nil {
		t.Fatal("expected error for ambiguous match")
	}
}

func TestResolveName_NotFound(t *testing.T) {
	providers := []Provider{
		{Name: "MyProxy"},
	}
	_, err := ResolveName("zzz", providers)
	if err == nil {
		t.Fatal("expected error for no match")
	}
}
