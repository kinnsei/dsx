package ds

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/a-h/templ"
	"github.com/starfederation/datastar-go/datastar"
)

// ModalContainerID is the fixed ID of the modal container in the base template.
const ModalContainerID = "modal-panel"

const modalCloseExpr = "document.getElementById('" + ModalContainerID + "').innerHTML=''"

// ModalOption customizes modal appearance.
type ModalOption func(*modalConfig)

type modalConfig struct {
	maxWidth string
}

// WithModalMaxWidth sets the modal max-width class (default "max-w-lg").
func WithModalMaxWidth(class string) ModalOption {
	return func(c *modalConfig) { c.maxWidth = class }
}

// Modal renders a templ component inside a centered modal dialog
// and patches it into #modal-panel via SSE.
func (s *Sender) Modal(sse *datastar.ServerSentEventGenerator, content templ.Component, opts ...ModalOption) error {
	cfg := &modalConfig{maxWidth: "max-w-lg"}
	for _, opt := range opts {
		opt(cfg)
	}

	var contentBuf bytes.Buffer
	if err := content.Render(context.Background(), &contentBuf); err != nil {
		return fmt.Errorf("rendering modal content: %w", err)
	}

	var b strings.Builder
	fmt.Fprintf(&b, `<div id="%s">`, ModalContainerID)

	// Overlay
	fmt.Fprintf(&b,
		`<div class="fixed inset-0 bg-black/40 z-40" data-on:click="%s"></div>`,
		modalCloseExpr,
	)

	// Dialog
	fmt.Fprintf(&b,
		`<div class="fixed inset-0 z-50 flex items-center justify-center p-4">`,
	)
	fmt.Fprintf(&b,
		`<div class="bg-base-100 rounded-box shadow-xl w-full %s overflow-y-auto max-h-[90vh]">`,
		cfg.maxWidth,
	)
	b.WriteString(`<div class="p-6 relative">`)

	// Close button
	fmt.Fprintf(&b,
		`<button class="btn btn-sm btn-circle btn-ghost absolute right-4 top-4" data-on:click="%s">✕</button>`,
		modalCloseExpr,
	)

	// Content
	b.WriteString(contentBuf.String())

	b.WriteString(`</div></div></div></div>`)

	return sse.PatchElements(b.String())
}

// HideModal patches #modal-panel back to its empty placeholder via SSE.
// Useful for server-initiated close (e.g. after a form submit inside the modal).
func (s *Sender) HideModal(sse *datastar.ServerSentEventGenerator) error {
	return sse.PatchElements(fmt.Sprintf(`<div id="%s"></div>`, ModalContainerID))
}
