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
	method       string // HTTP method for confirm action ("get" or "post", default "post")
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

// WithConfirmGet uses GET instead of POST for the confirm action.
// By default, confirm uses POST with CSRF protection since confirmations
// typically trigger destructive/mutating operations.
func WithConfirmGet() ConfirmOption {
	return func(c *confirmConfig) { c.method = "get" }
}

// Confirm shows a confirmation dialog via SSE using the modal container.
// When the user clicks confirm, it triggers a POST request (with CSRF) to confirmURL.
// Use WithConfirmGet() to use GET instead.
// Cancel closes the dialog without any request.
func (s *Sender) Confirm(sse *datastar.ServerSentEventGenerator, message string, confirmURL string, opts ...ConfirmOption) error {
	cfg := &confirmConfig{
		title:        "Confirm",
		confirmLabel: "Confirm",
		cancelLabel:  "Cancel",
		confirmClass: "btn btn-primary",
		maxWidth:     "max-w-sm",
		method:       "post",
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

	// Confirm button — uses PostOnce (with CSRF) by default, GetOnce if WithConfirmGet.
	var actionExpr string
	if cfg.method == "get" {
		actionExpr = GetOnce(confirmURL)
	} else {
		actionExpr = PostOnce(confirmURL)
	}
	fmt.Fprintf(&b,
		`<button class="%s" data-on:click="%s; %s">%s</button>`,
		cfg.confirmClass,
		actionExpr,
		modalCloseExpr,
		templ.EscapeString(cfg.confirmLabel),
	)

	b.WriteString(`</div></div></div></div></div>`)

	return sse.PatchElements(b.String())
}
