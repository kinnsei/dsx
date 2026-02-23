package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
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

func (h *commandbarHandlers) capture() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// With filterSignals, only the triggering commandbar's namespace arrives.
		// Read raw JSON and find the single namespace present.
		var raw map[string]json.RawMessage
		if err := ds.ReadRaw(r, &raw); err != nil {
			slog.Error("commandbar: failed to read signals", "error", err)
			http.Error(w, "failed to read input", http.StatusBadRequest)
			return
		}

		sse := datastar.NewSSE(w, r)

		for id, data := range raw {
			var signals commandbar.CommandBarSignals
			if err := json.Unmarshal(data, &signals); err != nil {
				continue
			}

			text := strings.TrimSpace(signals.Text)

			switch {
			case text != "":
				ds.Send.Toast(sse, ds.ToastSuccess,
					fmt.Sprintf("[%s] Received: %q", id, text))
			case signals.Mode == "voice":
				ds.Send.Toast(sse, ds.ToastSuccess,
					fmt.Sprintf("[%s] Voice recording received", id))
			case signals.Mode == "file":
				ds.Send.Toast(sse, ds.ToastSuccess,
					fmt.Sprintf("[%s] File upload received", id))
			default:
				ds.Send.Toast(sse, ds.ToastSuccess,
					fmt.Sprintf("[%s] Action received", id))
			}
			return
		}

		ds.Send.Toast(sse, ds.ToastWarning, "No signals received")
	}
}
