package toast

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/a-h/templ"
	"github.com/plaenen/webx/ui/alert"
	"github.com/plaenen/webx/utils"
	"github.com/starfederation/datastar-go/datastar"
)

// ContainerID is the fixed ID of the toast container in the base template.
const ContainerID = "toast-container"

// Level represents the severity of a toast notification.
type Level string

const (
	LevelInfo    Level = "info"
	LevelSuccess Level = "success"
	LevelWarning Level = "warning"
	LevelError   Level = "error"
)

// Show appends an auto-dismissing toast notification via SSE.
// The toast disappears after durationMs (default 3000ms) and has a close button.
func Show(sse *datastar.ServerSentEventGenerator, level Level, message string, durationMs int) error {
	if durationMs == 0 {
		durationMs = 3000
	}
	id := "toast-" + utils.RandomID()
	variant := levelToVariant(level)

	body := fmt.Sprintf(`<span>%s</span>`, templ.EscapeString(message))
	html := buildToast(id, variant, body, durationMs, true)

	return patchToast(sse, html)
}

// ShowPersistent appends a toast that stays until the user clicks the close button.
func ShowPersistent(sse *datastar.ServerSentEventGenerator, level Level, message string) error {
	id := "toast-" + utils.RandomID()
	variant := levelToVariant(level)

	body := fmt.Sprintf(`<span>%s</span>`, templ.EscapeString(message))
	html := buildToast(id, variant, body, 0, true)

	return patchToast(sse, html)
}

// ShowAction appends a toast with an action button that triggers a Datastar GET.
// The toast stays until the action or close button is clicked.
func ShowAction(sse *datastar.ServerSentEventGenerator, level Level, message string, actionLabel string, actionURL string) error {
	id := "toast-" + utils.RandomID()
	variant := levelToVariant(level)

	body := fmt.Sprintf(
		`<span>%s</span><button class="btn btn-sm" data-on:click="@get('%s'); document.getElementById('%s')?.remove()">%s</button>`,
		templ.EscapeString(message),
		templ.EscapeString(actionURL),
		id,
		templ.EscapeString(actionLabel),
	)
	html := buildToast(id, variant, body, 0, true)

	return patchToast(sse, html)
}

// ShowLink appends a toast with a clickable link.
// The toast auto-dismisses after durationMs (default 5000ms).
func ShowLink(sse *datastar.ServerSentEventGenerator, level Level, message string, linkText string, linkURL string, durationMs int) error {
	if durationMs == 0 {
		durationMs = 5000
	}
	id := "toast-" + utils.RandomID()
	variant := levelToVariant(level)

	body := fmt.Sprintf(
		`<span>%s <a href="%s" class="underline font-semibold" data-on:click="document.getElementById('%s')?.remove()">%s</a></span>`,
		templ.EscapeString(message),
		templ.EscapeString(linkURL),
		id,
		templ.EscapeString(linkText),
	)
	html := buildToast(id, variant, body, durationMs, true)

	return patchToast(sse, html)
}

// ShowComponent appends a custom templ component as a toast via SSE.
// Useful when you need more than a simple text message.
func ShowComponent(sse *datastar.ServerSentEventGenerator, component templ.Component) error {
	var buf bytes.Buffer
	if err := component.Render(context.Background(), &buf); err != nil {
		return fmt.Errorf("rendering toast component: %w", err)
	}

	return patchToast(sse, buf.String())
}

// buildToast constructs the HTML for a single toast alert.
// durationMs=0 means no auto-dismiss. closable adds a close button.
func buildToast(id, variant, body string, durationMs int, closable bool) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf(`<div id="%s" class="alert %s shadow-lg"`, id, variant))

	if durationMs > 0 {
		b.WriteString(fmt.Sprintf(` data-init="setTimeout(() => document.getElementById('%s')?.remove(), %d)"`, id, durationMs))
	}

	b.WriteString(`>`)
	b.WriteString(body)

	if closable {
		b.WriteString(fmt.Sprintf(
			`<button class="btn btn-sm btn-ghost btn-circle" data-on:click="document.getElementById('%s')?.remove()">✕</button>`,
			id,
		))
	}

	b.WriteString(`</div>`)
	return b.String()
}

func patchToast(sse *datastar.ServerSentEventGenerator, html string) error {
	return sse.PatchElements(html,
		datastar.WithSelector("#"+ContainerID),
		datastar.WithModeAppend(),
	)
}

func levelToVariant(level Level) string {
	switch level {
	case LevelInfo:
		return string(alert.VariantInfo)
	case LevelSuccess:
		return string(alert.VariantSuccess)
	case LevelWarning:
		return string(alert.VariantWarning)
	case LevelError:
		return string(alert.VariantError)
	default:
		return string(alert.VariantInfo)
	}
}
