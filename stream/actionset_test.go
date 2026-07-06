package stream_test

import (
	"testing"
	"time"

	"github.com/kinnsei/dsx/stream"
)

func TestAction_ID(t *testing.T) {
	// Direct ActionSet.ID().Get() chain.
	reaction := stream.Updated.ID(42).Get("/api/items/42")
	if reaction == (stream.Reaction{}) {
		t.Fatal("reaction should not be zero value")
	}
}

func TestActionOr_ID(t *testing.T) {
	reaction := stream.Created.Or(stream.Updated).ID(42).Get("/api/items/42")
	if reaction == (stream.Reaction{}) {
		t.Fatal("reaction should not be zero value")
	}
}

func TestAction_Debounce(t *testing.T) {
	reaction := stream.Structural.Debounce(500 * time.Millisecond).Get("/api/items/list")
	if reaction == (stream.Reaction{}) {
		t.Fatal("reaction should not be zero value")
	}
}

func TestAction_IDAndDebounce(t *testing.T) {
	reaction := stream.Updated.ID(42).Debounce(200*time.Millisecond).Get("/api/items/42")
	if reaction == (stream.Reaction{}) {
		t.Fatal("reaction should not be zero value")
	}
}

func TestAction_StructuralDefinition(t *testing.T) {
	reaction := stream.Structural.Get("/api/items")
	if reaction == (stream.Reaction{}) {
		t.Fatal("reaction should not be zero value")
	}
}

func TestAction_AnyDefinition(t *testing.T) {
	reaction := stream.Any.Get("/api/count")
	if reaction == (stream.Reaction{}) {
		t.Fatal("reaction should not be zero value")
	}
}

func TestActionOr_MultipleChained(t *testing.T) {
	// Chain three action sets.
	reaction := stream.Created.Or(stream.Updated).Or(stream.Action("archived")).Get("/api/items")
	if reaction == (stream.Reaction{}) {
		t.Fatal("reaction should not be zero value")
	}
}
