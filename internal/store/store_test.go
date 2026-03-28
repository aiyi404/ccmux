package store

import "testing"

func TestNew_ReturnsAppState(t *testing.T) {
	state, err := New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state.Service == nil {
		t.Error("expected Service to be non-nil")
	}
	if state.Config == nil {
		t.Error("expected Config to be non-nil")
	}
}
