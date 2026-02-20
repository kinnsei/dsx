package ds

import (
	"fmt"
	"strings"

	"github.com/a-h/templ"
	"github.com/starfederation/datastar-go/datastar"
)

// ConfirmOption customizes the confirm dialog.
type ConfirmOption func(*confirmConfig)

type confirmConfig struct {
	title        string
	confirmLabel string
	cancelLabel  string
	confirmClass string
	maxWidth     string
}

// WithConfirmTitle sets the dialog title (default: "Confirm").
func WithConfirmTitle(title string) ConfirmOption {
	return func(c *confirmConfig) { c.title = title }
}

// WithConfirmLabel sets the confirm button text (default: "Confirm").
func WithConfirmLabel(label string) ConfirmOption {
	return func(c *confirmConfig) { c.confirmLabel = label }
}

// WithCancelLabel sets the cancel button text (default: "Cancel").
func WithCancelLabel(label string) ConfirmOption {
	return func(c *confirmConfig) { c.cancelLabel = label }
}

// WithConfirmClass sets the confirm button class (default: "btn btn-primary").
func WithConfirmClass(class string) ConfirmOption {
	return func(c *confirmConfig) { c.confirmClass = class }
}

// WithConfirmMaxWidth sets the dialog max-width class (default: "max-w-sm").
func WithConfirmMaxWidth(class string) ConfirmOption {
	return func(c *confirmConfig) { c.maxWidth = class }
}

// Confirm shows a confirmation dialog via SSE using the modal container.
// When the user clicks confirm, it triggers a GET request to confirmURL.
// Cancel closes the dialog without any request.
func (s *Sender) Confirm(sse *datastar.ServerSentEventGenerator, message string, confirmURL string, opts ...ConfirmOption) error {
	cfg := &confirmConfig{
		title:        "Confirm",
		confirmLabel: "Confirm",
		cancelLabel:  "Cancel",
		confirmClass: "btn btn-primary",
		maxWidth:     "max-w-sm",
	}
	for _, opt := range opts {
		opt(cfg)
	}

	var b strings.Builder
	fmt.Fprintf(&b, `<div id="%s">`, ModalContainerID)

	// Overlay
	fmt.Fprintf(&b,
		`<div class="fixed inset-0 bg-black/40 z-40" data-on:click="%s"></div>`,
		modalCloseExpr,
	)

	// Dialog
	b.WriteString(`<div class="fixed inset-0 z-50 flex items-center justify-center p-4">`)
	fmt.Fprintf(&b,
		`<div class="bg-base-100 rounded-box shadow-xl w-full %s">`,
		cfg.maxWidth,
	)
	b.WriteString(`<div class="p-6">`)

	// Title
	fmt.Fprintf(&b,
		`<h3 class="text-lg font-bold">%s</h3>`,
		templ.EscapeString(cfg.title),
	)

	// Message
	fmt.Fprintf(&b,
		`<p class="py-4">%s</p>`,
		templ.EscapeString(message),
	)

	// Actions
	b.WriteString(`<div class="flex justify-end gap-2">`)

	// Cancel button
	fmt.Fprintf(&b,
		`<button class="btn" data-on:click="%s">%s</button>`,
		modalCloseExpr,
		templ.EscapeString(cfg.cancelLabel),
	)

	// Confirm button
	fmt.Fprintf(&b,
		`<button class="%s" data-on:click="@get('%s'); %s">%s</button>`,
		cfg.confirmClass,
		templ.EscapeString(confirmURL),
		modalCloseExpr,
		templ.EscapeString(cfg.confirmLabel),
	)

	b.WriteString(`</div></div></div></div></div>`)

	return sse.PatchElements(b.String())
}
