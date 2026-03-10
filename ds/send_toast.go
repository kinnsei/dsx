package ds

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/a-h/templ"
	"github.com/laenen-partners/dsx/utils"
	"github.com/starfederation/datastar-go/datastar"
)

// ToastContainerID is the fixed ID of the toast container in the base template.
const ToastContainerID = "toast-container"

// ToastLevel represents the severity of a toast notification.
type ToastLevel string

const (
	ToastInfo    ToastLevel = "info"
	ToastSuccess ToastLevel = "success"
	ToastWarning ToastLevel = "warning"
	ToastError   ToastLevel = "error"
)

// ToastOption customizes a toast notification.
type ToastOption func(*toastConfig)

type toastConfig struct {
	duration    int    // auto-dismiss in ms; 0 = persistent
	persistent  bool   // stays until user closes
	actionLabel string // action button text
	actionURL   string // action button @get URL
	linkText    string // link text
	linkURL     string // link href
}

// WithToastDuration sets the auto-dismiss duration in milliseconds.
// Default is 3000ms. Use WithToastPersistent() to disable auto-dismiss.
func WithToastDuration(ms int) ToastOption {
	return func(c *toastConfig) { c.duration = ms }
}

// WithToastPersistent makes the toast stay until the user clicks the close button.
func WithToastPersistent() ToastOption {
	return func(c *toastConfig) { c.persistent = true }
}

// WithToastAction adds an action button that triggers a Datastar GET.
// The toast stays until the action or close button is clicked.
func WithToastAction(label, url string) ToastOption {
	return func(c *toastConfig) {
		c.actionLabel = label
		c.actionURL = url
		c.persistent = true
	}
}

// WithToastLink adds a clickable link to the toast message.
func WithToastLink(text, url string) ToastOption {
	return func(c *toastConfig) {
		c.linkText = text
		c.linkURL = url
	}
}

// Toast appends a toast notification via SSE.
// Default behavior: auto-dismiss after 3 seconds with a close button.
func (s *Sender) Toast(sse *datastar.ServerSentEventGenerator, level ToastLevel, message string, opts ...ToastOption) error {
	cfg := &toastConfig{duration: 3000}
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.persistent {
		cfg.duration = 0
	}

	id := "toast-" + utils.RandomID()
	variant := toastLevelToVariant(level)

	var body strings.Builder
	if cfg.linkText != "" && cfg.linkURL != "" {
		fmt.Fprintf(&body,
			`<span>%s <a href="%s" class="underline font-semibold" data-on:click="document.getElementById('%s')?.remove()">%s</a></span>`,
			templ.EscapeString(message),
			templ.EscapeString(cfg.linkURL),
			id,
			templ.EscapeString(cfg.linkText),
		)
	} else {
		fmt.Fprintf(&body, `<span>%s</span>`, templ.EscapeString(message))
	}

	if cfg.actionLabel != "" && cfg.actionURL != "" {
		fmt.Fprintf(&body,
			`<button class="btn btn-sm" data-on:click="@get('%s'); document.getElementById('%s')?.remove()">%s</button>`,
			templ.EscapeString(cfg.actionURL),
			id,
			templ.EscapeString(cfg.actionLabel),
		)
	}

	html := buildToastHTML(id, variant, body.String(), cfg.duration)
	return patchToast(sse, html)
}

// ToastComponent appends a custom templ component as a toast via SSE.
func (s *Sender) ToastComponent(sse *datastar.ServerSentEventGenerator, component templ.Component) error {
	var buf bytes.Buffer
	if err := component.Render(context.Background(), &buf); err != nil {
		return fmt.Errorf("rendering toast component: %w", err)
	}
	return patchToast(sse, buf.String())
}

func buildToastHTML(id, variant, body string, durationMs int) string {
	var b strings.Builder

	fmt.Fprintf(&b, `<div id="%s" class="alert %s shadow-lg"`, id, variant)

	if durationMs > 0 {
		fmt.Fprintf(&b, ` data-init="setTimeout(() => document.getElementById('%s')?.remove(), %d)"`, id, durationMs)
	}

	b.WriteString(`>`)
	b.WriteString(body)

	// Close button
	fmt.Fprintf(&b,
		`<button class="btn btn-sm btn-ghost btn-circle" data-on:click="document.getElementById('%s')?.remove()">✕</button>`,
		id,
	)

	b.WriteString(`</div>`)
	return b.String()
}

func patchToast(sse *datastar.ServerSentEventGenerator, html string) error {
	return sse.PatchElements(html,
		datastar.WithSelector("#"+ToastContainerID),
		datastar.WithModeAppend(),
	)
}

func toastLevelToVariant(level ToastLevel) string {
	switch level {
	case ToastInfo:
		return "alert-info"
	case ToastSuccess:
		return "alert-success"
	case ToastWarning:
		return "alert-warning"
	case ToastError:
		return "alert-error"
	default:
		return "alert-info"
	}
}
