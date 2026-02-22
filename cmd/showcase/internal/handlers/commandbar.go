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
		// Read signals BEFORE creating SSE (SSE consumes the request body).
		var raw map[string]json.RawMessage
		if err := datastar.ReadSignals(r, &raw); err != nil {
			http.Error(w, fmt.Sprintf("read signals: %v", err), http.StatusBadRequest)
			return
		}

		sse := datastar.NewSSE(w, r)

		// All commandbar instances on the page send their signals.
		// Find the one with active content (non-empty text or non-idle mode).
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

			// Skip instances in their default/idle state.
			if text == "" && (mode == "" || mode == "text") {
				continue
			}

			switch {
			case text != "":
				ds.Send.Toast(sse, ds.ToastSuccess,
					fmt.Sprintf("Received: %q (mode: %s)", text, mode))
			case mode == "voice":
				ds.Send.Toast(sse, ds.ToastSuccess, "Voice recording received")
			case mode == "file":
				ds.Send.Toast(sse, ds.ToastSuccess, "File upload received")
			}
			return
		}

		// No active commandbar found — fallback to any with text content.
		ds.Send.Toast(sse, ds.ToastSuccess, "Action received")
	}
}
