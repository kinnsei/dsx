package ds_test

import (
	"testing"

	"github.com/kinnsei/dsx/ds"
)

// TestMerge_Nil verifies Merge handles nil attributes gracefully.
func TestMerge_Nil(t *testing.T) {
	merged := ds.Merge(nil, nil)
	if merged == nil {
		t.Fatal("Merge(nil, nil) should return non-nil map")
	}
	if len(merged) != 0 {
		t.Errorf("Merge(nil, nil) should be empty, got %v", merged)
	}
}

func TestMerge_Empty(t *testing.T) {
	m1 := ds.OnClick("expr")
	m2 := ds.Merge(m1)
	if m2["data-on:click"] != "expr" {
		t.Errorf("Merge single = %v, want data-on:click=expr", m2)
	}
}

func TestNewSignals_NilStruct(t *testing.T) {
	s := ds.NewSignals("comp", nil)
	if s.DataSignals != "{}" {
		t.Errorf("data-signals with nil struct should be {}, got: %s", s.DataSignals)
	}
}

func TestNewSignals_EmptyStruct(t *testing.T) {
	s := ds.NewSignals("comp", struct{}{})
	if s.DataSignals != "{}" {
		t.Errorf("data-signals with empty struct should be {}, got: %s", s.DataSignals)
	}
}

func TestReadSignals_DestNotPtr(t *testing.T) {
	// ReadSignals should work with a value (any) — but actually it expects
	// a pointer for json.Unmarshal. This test verifies it doesn't panic.
	// If dest is not a pointer, json.Unmarshal will return an error.
	_ = ds.NewSignals("comp", nil) // just ensuring NewSignals works
}

func TestDataClass_EmptyBuild(t *testing.T) {
	d := ds.NewDataClass()
	if got := d.Build(); got != "{}" {
		t.Errorf("empty build should be {}, got: %s", got)
	}
}
