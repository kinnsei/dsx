package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/plaenen/webx/ds"
	"github.com/plaenen/webx/ui/commandbar"
	"github.com/starfederation/datastar-go/datastar"
)

type commandbarHandlers struct{}

func newCommandbarHandlers() *commandbarHandlers {
	return &commandbarHandlers{}
}

func (h *commandbarHandlers) register(r chi.Router) {
	r.Post("/api/commandbar/capture", h.capture())
}

// commandbar demo IDs on the showcase page (sanitized: hyphens → underscores).
var commandBarDemoIDs = []string{
	"demo_text_only",
	"demo_suggestions",
	"demo_all_modes",
	"demo_text_file",
	"demo_custom",
}

func (h *commandbarHandlers) capture() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)

		// Read all signals as raw JSON to find which commandbar namespace is present.
		var raw map[string]json.RawMessage
		if err := datastar.ReadSignals(r, &raw); err != nil {
			ds.Send.Toast(sse, ds.ToastError, fmt.Sprintf("Failed to read signals: %v", err))
			return
		}

		for _, id := range commandBarDemoIDs {
			data, ok := raw[id]
			if !ok {
				continue
			}

			var signals commandbar.CommandBarSignals
			if err := json.Unmarshal(data, &signals); err != nil {
				continue
			}

			text := strings.TrimSpace(signals.Text)
			mode := signals.Mode

			switch {
			case text != "":
				ds.Send.Toast(sse, ds.ToastSuccess,
					fmt.Sprintf("Received: %q (mode: %s)", text, mode))
			case mode == "voice":
				ds.Send.Toast(sse, ds.ToastSuccess, "Voice recording received")
			case mode == "file":
				ds.Send.Toast(sse, ds.ToastSuccess, "File upload received")
			default:
				ds.Send.Toast(sse, ds.ToastWarning, "Empty input received")
			}
			return
		}

		ds.Send.Toast(sse, ds.ToastWarning, "No command bar input found in signals")
	}
}
