package toast

import (
	"bytes"
	"context"
	"fmt"

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

// Show appends a toast notification to the toast container via SSE.
// The toast auto-dismisses after the specified duration (in milliseconds).
// If duration is 0, it defaults to 3000ms.
func Show(sse *datastar.ServerSentEventGenerator, level Level, message string, durationMs int) error {
	if durationMs == 0 {
		durationMs = 3000
	}

	variant := levelToVariant(level)
	id := "toast-" + utils.RandomID()

	html := fmt.Sprintf(
		`<div id="%s" class="alert %s alert-soft animate-bounce-in" data-on:load="setTimeout(() => document.getElementById('%s')?.remove(), %d)">
	<span>%s</span>
</div>`,
		id, variant, id, durationMs, templ.EscapeString(message),
	)

	return sse.PatchElements(html,
		datastar.WithSelector("#"+ContainerID),
		datastar.WithModeAppend(),
	)
}

// ShowComponent appends a custom templ component as a toast via SSE.
// Useful when you need more than a simple text message.
func ShowComponent(sse *datastar.ServerSentEventGenerator, component templ.Component) error {
	var buf bytes.Buffer
	if err := component.Render(context.Background(), &buf); err != nil {
		return fmt.Errorf("rendering toast component: %w", err)
	}

	return sse.PatchElements(buf.String(),
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
