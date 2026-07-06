package ds

import (
	"bytes"
	"context"
	"fmt"

	"github.com/kinnsei/dsx/ui/modalpanel"
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

	var actionExpr string
	if cfg.method == "get" {
		actionExpr = GetOnce(confirmURL)
	} else {
		actionExpr = PostOnce(confirmURL)
	}

	cc := modalpanel.ConfirmConfig{
		MaxWidth:     cfg.maxWidth,
		Title:        cfg.title,
		Message:      message,
		CancelLabel:  cfg.cancelLabel,
		ConfirmLabel: cfg.confirmLabel,
		ConfirmClass: cfg.confirmClass,
		ActionExpr:   actionExpr,
	}

	var buf bytes.Buffer
	if err := modalpanel.Confirm(cc).Render(context.Background(), &buf); err != nil {
		return fmt.Errorf("rendering confirm dialog: %w", err)
	}
	return sse.PatchElements(buf.String())
}
