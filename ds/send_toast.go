package ds

import (
	"bytes"
	"context"
	"fmt"

	"github.com/a-h/templ"
	"github.com/laenen-partners/dsx/ui/toastcontainer"
	"github.com/laenen-partners/dsx/utils"
	"github.com/starfederation/datastar-go/datastar"
)

// ToastContainerID is the fixed ID of the toast container in the base template.
const ToastContainerID = toastcontainer.ContainerID

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

	tc := toastcontainer.Config{
		ID:          "toast-" + utils.RandomID(),
		Message:     message,
		Variant:     toastLevelToVariant(level),
		DurationMs:  cfg.duration,
		LinkText:    cfg.linkText,
		LinkURL:     cfg.linkURL,
		ActionLabel: cfg.actionLabel,
		ActionURL:   cfg.actionURL,
	}

	var buf bytes.Buffer
	if err := toastcontainer.Toast(tc).Render(context.Background(), &buf); err != nil {
		return fmt.Errorf("rendering toast: %w", err)
	}
	return patchToast(sse, buf.String())
}

// ToastComponent appends a custom templ component as a toast via SSE.
func (s *Sender) ToastComponent(ctx context.Context, sse *datastar.ServerSentEventGenerator, component templ.Component) error {
	var buf bytes.Buffer
	if err := component.Render(ctx, &buf); err != nil {
		return fmt.Errorf("rendering toast component: %w", err)
	}
	return patchToast(sse, buf.String())
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
