package ds

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/a-h/templ"
	"github.com/starfederation/datastar-go/datastar"
)

// DrawerContainerID is the fixed ID of the drawer container in the base template.
const DrawerContainerID = "drawer-panel"

// closeExpr is the inline Datastar expression used by close button and overlay.
const closeExpr = "document.getElementById('" + DrawerContainerID + "').innerHTML=''"

// DrawerOption customizes drawer appearance.
type DrawerOption func(*drawerConfig)

type drawerConfig struct {
	maxWidth string
}

// WithDrawerMaxWidth sets the panel max-width class (default "max-w-lg").
func WithDrawerMaxWidth(class string) DrawerOption {
	return func(c *drawerConfig) { c.maxWidth = class }
}

// Drawer renders a templ component inside a slide-in drawer panel
// and patches it into #drawer-panel via SSE.
func (s *Sender) Drawer(sse *datastar.ServerSentEventGenerator, content templ.Component, opts ...DrawerOption) error {
	cfg := &drawerConfig{maxWidth: "max-w-lg"}
	for _, opt := range opts {
		opt(cfg)
	}

	var contentBuf bytes.Buffer
	if err := content.Render(context.Background(), &contentBuf); err != nil {
		return fmt.Errorf("rendering drawer content: %w", err)
	}

	var b strings.Builder
	fmt.Fprintf(&b, `<div id="%s">`, DrawerContainerID)

	// Overlay
	fmt.Fprintf(&b,
		`<div class="fixed inset-0 bg-black/40 z-40" data-on:click="%s"></div>`,
		closeExpr,
	)

	// Panel
	fmt.Fprintf(&b,
		`<div class="fixed inset-y-0 right-0 w-full %s bg-base-100 shadow-xl z-50 overflow-y-auto">`,
		cfg.maxWidth,
	)
	b.WriteString(`<div class="p-6 relative">`)

	// Close button
	fmt.Fprintf(&b,
		`<button class="btn btn-sm btn-circle btn-ghost absolute right-4 top-4" data-on:click="%s">✕</button>`,
		closeExpr,
	)

	// Content
	b.WriteString(contentBuf.String())

	b.WriteString(`</div></div></div>`)

	return sse.PatchElements(b.String())
}

// HideDrawer patches #drawer-panel back to its empty placeholder via SSE.
// Useful for server-initiated close (e.g. after a form submit inside the drawer).
func (s *Sender) HideDrawer(sse *datastar.ServerSentEventGenerator) error {
	return sse.PatchElements(fmt.Sprintf(`<div id="%s"></div>`, DrawerContainerID))
}
