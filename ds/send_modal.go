package ds

import (
	"bytes"
	"context"
	"fmt"

	"github.com/a-h/templ"
	"github.com/kinnsei/dsx/ui/modalpanel"
	"github.com/starfederation/datastar-go/datastar"
)

// ModalContainerID is the fixed ID of the modal container in the base template.
const ModalContainerID = modalpanel.ContainerID

// modalCloseExpr is the inline Datastar expression for closing the modal.
// Used internally by Confirm and other modal-based operations.
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
func (s *Sender) Modal(ctx context.Context, sse *datastar.ServerSentEventGenerator, content templ.Component, opts ...ModalOption) error {
	cfg := &modalConfig{maxWidth: "max-w-lg"}
	for _, opt := range opts {
		opt(cfg)
	}

	mc := modalpanel.Config{MaxWidth: cfg.maxWidth}
	modal := modalpanel.Modal(mc, content)

	var buf bytes.Buffer
	if err := modal.Render(ctx, &buf); err != nil {
		return fmt.Errorf("rendering modal: %w", err)
	}
	return sse.PatchElements(buf.String())
}

// HideModal patches #modal-panel back to its empty placeholder via SSE.
// Useful for server-initiated close (e.g. after a form submit inside the modal).
func (s *Sender) HideModal(sse *datastar.ServerSentEventGenerator) error {
	return sse.PatchElements(fmt.Sprintf(`<div id="%s"></div>`, ModalContainerID))
}
