// Package dsx provides request-scoped context for DSX templates.
//
// This package re-exports from internal sub-packages for backward compatibility.
// New code can import the sub-packages directly:
//
//	import "github.com/kinnsei/dsx/internal/uxcontext"
//	import "github.com/kinnsei/dsx/internal/middleware"
package dsx

import (
	"context"

	"github.com/kinnsei/dsx/internal/uxcontext"
)

// Context carries request-scoped state through the middleware chain and into
// templ components. Retrieve it with FromContext.
type Context = uxcontext.Context

// NewContext returns an empty Context.
func NewContext() *Context {
	return uxcontext.NewContext()
}

// FromContext extracts the dsx Context from a context.Context.
// If none is present, a new empty Context is returned.
func FromContext(ctx context.Context) *Context {
	return uxcontext.FromContext(ctx)
}
