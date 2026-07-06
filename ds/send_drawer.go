package ds

import (
	"bytes"
	"context"
	"fmt"

	"github.com/a-h/templ"
	"github.com/laenen-partners/dsx/ui/drawerpanel"
	"github.com/starfederation/datastar-go/datastar"
)

// DrawerContainerID is the fixed ID of the drawer container in the base template.
const DrawerContainerID = drawerpanel.ContainerID

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
func (s *Sender) Drawer(ctx context.Context, sse *datastar.ServerSentEventGenerator, content templ.Component, opts ...DrawerOption) error {
	cfg := &drawerConfig{maxWidth: "max-w-lg"}
	for _, opt := range opts {
		opt(cfg)
	}

	dc := drawerpanel.Config{
		MaxWidth:   cfg.maxWidth,
		Expandable: cfg.expandable,
	}
	drawer := drawerpanel.Drawer(dc, content)

	var buf bytes.Buffer
	if err := drawer.Render(ctx, &buf); err != nil {
		return fmt.Errorf("rendering drawer: %w", err)
	}
	return sse.PatchElements(buf.String())
}

// HideDrawer patches #drawer-panel back to its empty placeholder via SSE.
// Useful for server-initiated close (e.g. after a form submit inside the drawer).
func (s *Sender) HideDrawer(sse *datastar.ServerSentEventGenerator) error {
	return sse.PatchElements(fmt.Sprintf(`<div id="%s"></div>`, DrawerContainerID))
}
