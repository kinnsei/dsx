package ds_test

import (
	"testing"

	"github.com/kinnsei/dsx/ds"
)

func TestNewExpression_Empty(t *testing.T) {
	e := ds.NewExpression()
	got := e.Build()
	want := ""
	if got != want {
		t.Errorf("empty expression = %q, want %q", got, want)
	}
}

func TestExpression_SingleStatement(t *testing.T) {
	e := ds.NewExpression().SetSignal("drawer.open", "true")
	got := e.Build()
	want := "$drawer.open = true"
	if got != want {
		t.Errorf("expression = %q, want %q", got, want)
	}
}

func TestExpression_MultiStatement(t *testing.T) {
	e := ds.NewExpression().
		SetSignal("drawer.open", "false").
		Statement("@post('/api/save')")
	got := e.Build()
	want := "$drawer.open = false; @post('/api/save')"
	if got != want {
		t.Errorf("expression = %q, want %q", got, want)
	}
}

func TestExpression_EmptyStatementIgnored(t *testing.T) {
	e := ds.NewExpression().Statement("").Statement("@get('/api')")
	got := e.Build()
	want := "@get('/api')"
	if got != want {
		t.Errorf("expression with empty stmt = %q, want %q", got, want)
	}
}

func TestExpression_Conditional(t *testing.T) {
	e := ds.NewExpression().
		Conditional("$visible", "@get('/api')", "")
	got := e.Build()
	want := "$visible ? @get('/api') : null"
	if got != want {
		t.Errorf("expression conditional = %q, want %q", got, want)
	}
}

func TestExpression_Conditional_WithFalse(t *testing.T) {
	e := ds.NewExpression().
		Conditional("$visible", "@get('/api')", "@get('/fallback')")
	got := e.Build()
	want := "$visible ? @get('/api') : @get('/fallback')"
	if got != want {
		t.Errorf("expression = %q, want %q", got, want)
	}
}

func TestBuildConditional(t *testing.T) {
	got := ds.BuildConditional("$visible", "'show'", "'hide'")
	want := "$visible ? 'show' : 'hide'"
	if got != want {
		t.Errorf("BuildConditional = %q, want %q", got, want)
	}
}

func TestBuildConditional_EmptyFalse(t *testing.T) {
	got := ds.BuildConditional("$visible", "'show'", "")
	want := "$visible ? 'show' : null"
	if got != want {
		t.Errorf("BuildConditional = %q, want %q", got, want)
	}
}
