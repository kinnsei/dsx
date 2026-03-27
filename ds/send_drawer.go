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

// Inline SVG icon paths (Lucide maximize-2 and minimize-2, 16x16).
const (
	iconMaximize2 = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M15 3h6v6"/><path d="m21 3-7 7"/><path d="m3 21 7-7"/><path d="M9 21H3v-6"/></svg>`
	iconMinimize2 = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m14 10 7-7"/><path d="M20 10h-6V4"/><path d="m3 21 7-7"/><path d="M4 14h6v6"/></svg>`
)

// DrawerOption customizes drawer appearance.
type DrawerOption func(*drawerConfig)

type drawerConfig struct {
	maxWidth   string
	expandable bool
}

// WithDrawerMaxWidth sets the panel max-width class (default "max-w-lg").
func WithDrawerMaxWidth(class string) DrawerOption {
	return func(c *drawerConfig) { c.maxWidth = class }
}

// WithDrawerExpandable adds an expand/collapse toggle button to the drawer
// header. When expanded, the drawer fills the full viewport width. When
// collapsed, it returns to the configured max-width.
func WithDrawerExpandable() DrawerOption {
	return func(c *drawerConfig) { c.expandable = true }
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

	if cfg.expandable {
		s.buildExpandableDrawer(&b, cfg, contentBuf.String())
	} else {
		s.buildStandardDrawer(&b, cfg, contentBuf.String())
	}

	b.WriteString(`</div>`)

	return sse.PatchElements(b.String())
}

// buildStandardDrawer renders the original drawer layout: absolutely positioned
// close button inside a padded content area.
func (s *Sender) buildStandardDrawer(b *strings.Builder, cfg *drawerConfig, content string) {
	fmt.Fprintf(b,
		`<div class="fixed inset-y-0 right-0 w-full %s bg-base-100 shadow-xl z-50 overflow-y-auto">`,
		cfg.maxWidth,
	)
	b.WriteString(`<div class="p-6 relative">`)

	// Close button — absolutely positioned top-right
	fmt.Fprintf(b,
		`<button class="btn btn-sm btn-circle btn-ghost absolute right-4 top-4" data-on:click="%s">✕</button>`,
		closeExpr,
	)

	b.WriteString(content)
	b.WriteString(`</div></div>`)
}

// buildExpandableDrawer renders a drawer with a sticky header bar containing
// expand/collapse and close buttons, keeping them out of the content flow.
func (s *Sender) buildExpandableDrawer(b *strings.Builder, cfg *drawerConfig, content string) {
	// Panel — always w-full, constrained by max-width via data-class when not expanded.
	// The class name must be quoted in data-class because it contains hyphens.
	fmt.Fprintf(b,
		`<div data-signals="{_drawerExpanded: false}" class="fixed inset-y-0 right-0 w-full %s bg-base-100 shadow-xl z-50 overflow-y-auto" data-class="{'%s': !$_drawerExpanded}">`,
		cfg.maxWidth, cfg.maxWidth,
	)

	// Sticky header bar — expand on the left, close on the right.
	b.WriteString(`<div class="sticky top-0 z-10 flex justify-between p-2 bg-base-100 border-b border-base-200">`)

	// Expand/collapse toggle — left side.
	fmt.Fprintf(b,
		`<button class="btn btn-sm btn-circle btn-ghost" data-on:click="$_drawerExpanded = !$_drawerExpanded" title="Toggle full width"><span data-show="!$_drawerExpanded">%s</span><span data-show="$_drawerExpanded">%s</span></button>`,
		iconMaximize2, iconMinimize2,
	)

	// Close button — right side.
	fmt.Fprintf(b,
		`<button class="btn btn-sm btn-circle btn-ghost" data-on:click="%s" title="Close">✕</button>`,
		closeExpr,
	)

	b.WriteString(`</div>`)

	// Content area — padded, below the header bar.
	b.WriteString(`<div class="p-6">`)
	b.WriteString(content)
	b.WriteString(`</div>`)

	b.WriteString(`</div>`)
}

// HideDrawer patches #drawer-panel back to its empty placeholder via SSE.
// Useful for server-initiated close (e.g. after a form submit inside the drawer).
func (s *Sender) HideDrawer(sse *datastar.ServerSentEventGenerator) error {
	return sse.PatchElements(fmt.Sprintf(`<div id="%s"></div>`, DrawerContainerID))
}
