// Provides context for webx templates
package webx

import (
	"context"
	"fmt"
	"slices"
)

// Stylesheet represents a <link rel="stylesheet"> tag to inject in <head>.
type Stylesheet struct {
	Href string
}

// Script represents a <script> tag to inject in <head>.
type Script struct {
	Src  string // script src attribute
	Type string // script type attribute (defaults to "module")
}

// BodyTag represents a custom element to inject at the end of <body>.
type BodyTag struct {
	Tag string // e.g. "<datastar-inspector></datastar-inspector>"
}

type WebXContext struct {
	CSRFToken   string
	DevMode     bool
	SessionID   string
	BasePath    string       // prefix for all SSE handler routes (e.g. "/showcase")
	Theme       string       // DaisyUI theme name applied via data-theme on <html>
	Store       SessionStore // session store for handlers that need to persist data
	StreamURL   string       // URL for the reactive SSE stream endpoint (e.g. "/showcase/stream")
	Stylesheets []Stylesheet
	Scripts     []Script
	BodyTags    []BodyTag
	Scopes      []string // reactive scopes accumulated during render by components
}

// WatchScope registers a scope to be watched by the stream SSE connection.
// Components call this during render to declare what data they depend on.
// Duplicates are ignored.
func (wctx *WebXContext) WatchScope(scope string) {
	if slices.Contains(wctx.Scopes, scope) {
		return
	}
	wctx.Scopes = append(wctx.Scopes, scope)
}

func NewContext(ctx context.Context) *WebXContext {
	return &WebXContext{}
}

type ctxKey struct{}

func FromContext(ctx context.Context) *WebXContext {
	if wctx, ok := ctx.Value(ctxKey{}).(*WebXContext); ok {
		return wctx
	}
	return NewContext(ctx)
}

func (wctx *WebXContext) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey{}, wctx)
}

// APIPath returns the full path for a component SSE handler by prepending
// the BasePath. For example, with BasePath "/showcase" and path "/api/calendar/navigate",
// it returns "/showcase/api/calendar/navigate".
func (wctx *WebXContext) APIPath(path string) string {
	return wctx.BasePath + path
}

// Post returns a Datastar expression that performs a POST request to the given URL.
func Post(url string) string {
	return fmt.Sprintf("@post('%s')", url)
}
