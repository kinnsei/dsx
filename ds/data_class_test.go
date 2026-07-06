package ds_test

import (
	"testing"

	"github.com/laenen-partners/dsx/ds"
)

func TestNewDataClass_Empty(t *testing.T) {
	d := ds.NewDataClass()
	got := d.Build()
	want := "{}"
	if got != want {
		t.Errorf("empty DataClass = %q, want %q", got, want)
	}
}

func TestDataClass_Single(t *testing.T) {
	d := ds.NewDataClass().Add("active", "$item.selected")
	got := d.Build()
	want := "{'active': $item.selected}"
	if got != want {
		t.Errorf("DataClass = %q, want %q", got, want)
	}
}

func TestDataClass_Multiple(t *testing.T) {
	d := ds.NewDataClass().
		Add("collapse-open", "$accordion.open").
		Add("collapse-close", "!$accordion.open")
	got := d.Build()
	want := "{'collapse-open': $accordion.open, 'collapse-close': !$accordion.open}"
	if got != want {
		t.Errorf("DataClass = %q, want %q", got, want)
	}
}

func TestDataClass_Chaining(t *testing.T) {
	d := ds.NewDataClass().Add("a", "$a").Add("b", "$b").Add("c", "$c")
	got := d.Build()
	want := "{'a': $a, 'b': $b, 'c': $c}"
	if got != want {
		t.Errorf("DataClass = %q, want %q", got, want)
	}
}

func TestDataClass_HyphenatedNames(t *testing.T) {
	d := ds.NewDataClass().Add("drawer-open", "$drawer.open")
	got := d.Build()
	// Hyphens in class names should be preserved (CSS class names can have hyphens).
	if got != "{'drawer-open': $drawer.open}" {
		t.Errorf("DataClass with hyphens = %q, want %q", got, "{'drawer-open': $drawer.open}")
	}
}
