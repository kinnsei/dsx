package ds_test

import (
	"testing"

	"github.com/kinnsei/dsx/ds"
)

// TestNewSignals verifies the signal key and signal references.
func TestNewSignals(t *testing.T) {
	s := ds.NewSignals("tabs", struct {
		Active string `json:"active"`
	}{Active: "overview"})

	if got := s.DataSignals; got == "" {
		t.Fatal("DataSignals should not be empty")
	}
	if !containsJSON(t, s.DataSignals, "active") {
		t.Errorf("DataSignals should contain field 'active', got: %s", s.DataSignals)
	}
}

// TestSignal verifies the $namespace.field reference.
func TestSignal(t *testing.T) {
	s := ds.NewSignals("tabs", nil)
	got := s.Signal("active")
	want := "$tabs.active"
	if got != want {
		t.Errorf("Signal() = %q, want %q", got, want)
	}
}

// TestSignal_SanitizedID verifies hyphens are converted to underscores.
func TestSignal_SanitizedID(t *testing.T) {
	s := ds.NewSignals("my-tabs", nil)
	got := s.Signal("active")
	want := "$my_tabs.active"
	if got != want {
		t.Errorf("Signal() = %q, want %q", got, want)
	}
}

func TestSet(t *testing.T) {
	s := ds.NewSignals("tabs", nil)
	got := s.Set("active", "'settings'")
	want := "$tabs.active = 'settings'"
	if got != want {
		t.Errorf("Set() = %q, want %q", got, want)
	}
}

func TestSetString(t *testing.T) {
	s := ds.NewSignals("tabs", nil)
	got := s.SetString("active", "settings")
	want := "$tabs.active = 'settings'"
	if got != want {
		t.Errorf("SetString() = %q, want %q", got, want)
	}
}

func TestToggle(t *testing.T) {
	s := ds.NewSignals("tabs", nil)
	got := s.Toggle("open")
	want := "$tabs.open = !$tabs.open"
	if got != want {
		t.Errorf("Toggle() = %q, want %q", got, want)
	}
}

func TestEquals(t *testing.T) {
	s := ds.NewSignals("tabs", nil)
	got := s.Equals("active", "'overview'")
	want := "$tabs.active === 'overview'"
	if got != want {
		t.Errorf("Equals() = %q, want %q", got, want)
	}
}

func TestNotEquals(t *testing.T) {
	s := ds.NewSignals("tabs", nil)
	got := s.NotEquals("active", "'overview'")
	want := "$tabs.active !== 'overview'"
	if got != want {
		t.Errorf("NotEquals() = %q, want %q", got, want)
	}
}

func TestEquals_stringLiteral(t *testing.T) {
	s := ds.NewSignals("tabs", nil)
	// Using Equals with a raw string — the value is used as-is, no extra quotes.
	got := s.Equals("active", "''")
	want := "$tabs.active === ''"
	if got != want {
		t.Errorf("Equals(empty) = %q, want %q", got, want)
	}
}

func TestConditional(t *testing.T) {
	s := ds.NewSignals("tabs", nil)
	got := s.Conditional("admin", "'show'", "'hide'")
	want := "$tabs.admin ? 'show' : 'hide'"
	if got != want {
		t.Errorf("Conditional() = %q, want %q", got, want)
	}
}

// Helper.
func containsJSON(t *testing.T, s, substr string) bool {
	t.Helper()
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i] == '"' || s[i] == '\'' {
			continue
		}
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
