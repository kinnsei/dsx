// Provides context for dsx templates
package dsx

import (
	"context"
	"fmt"
	"strconv"
)

// Watcher pairs a scope with a unique signal key for that watcher.
// Multiple components can watch the same scope, each with its own key.
type Watcher struct {
	Scope string // e.g. "customers:*"
	Key   string // unique signal key, e.g. "customers_WILD", "customers_WILD_2"
}

// Context carries request-scoped state through the middleware chain and into
// templ components. Retrieve it with FromContext.
type Context struct {
	SessionID string
	CSRFToken string
	Theme     string
	BasePath  string    // prefix for all SSE handler routes (e.g. "/showcase")
	StreamURL string    // URL for the reactive SSE stream endpoint (e.g. "/showcase/stream")
	Watchers  []Watcher // reactive scope watchers accumulated during render

	keyCounts map[string]int // tracks how many watchers use each base key
}

// WatchScope registers a scope watcher and returns its unique signal key.
// Multiple components can watch the same scope — each gets a distinct key
// so their data-effects don't interfere with each other.
func (ctx *Context) WatchScope(scope string, baseKey string) string {
	if ctx.keyCounts == nil {
		ctx.keyCounts = make(map[string]int)
	}
	ctx.keyCounts[baseKey]++
	key := baseKey
	if ctx.keyCounts[baseKey] > 1 {
		key = baseKey + "_" + strconv.Itoa(ctx.keyCounts[baseKey])
	}
	ctx.Watchers = append(ctx.Watchers, Watcher{Scope: scope, Key: key})
	return key
}

// NewContext returns an empty Context.
func NewContext() *Context {
	return &Context{}
}

type ctxKey struct{}

// FromContext extracts the dsx Context from a context.Context.
// If none is present, a new empty Context is returned.
func FromContext(ctx context.Context) *Context {
	if wctx, ok := ctx.Value(ctxKey{}).(*Context); ok {
		return wctx
	}
	return NewContext()
}

// WithContext stores the dsx Context into a context.Context.
func (c *Context) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey{}, c)
}

// APIPath returns the full path for a component SSE handler by prepending
// the BasePath. For example, with BasePath "/showcase" and path "/calendar/navigate",
// it returns "/showcase/calendar/navigate".
func (c *Context) APIPath(path string) string {
	return c.BasePath + path
}

// Post returns a Datastar expression that performs a POST request to the given URL.
func Post(url string) string {
	return fmt.Sprintf("@post('%s')", url)
}
