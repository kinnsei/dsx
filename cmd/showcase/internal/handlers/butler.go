package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/plaenen/webx"
	"github.com/plaenen/webx/ds"
	"github.com/starfederation/datastar-go/datastar"
)

type butlerHandlers struct{}

func newButlerHandlers() *butlerHandlers {
	return &butlerHandlers{}
}

func (h *butlerHandlers) register(r chi.Router) {
	r.Post("/api/butler/capture", h.capture())
}

// captureSignals matches the CommandBar signal structure.
type captureSignals struct {
	Open bool   `json:"open"`
	Text string `json:"text"`
}

func (h *butlerHandlers) capture() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = webx.FromContext(r.Context())
		sse := datastar.NewSSE(w, r)

		ds.Send.Toast(sse, ds.ToastSuccess, "Message received! AI is processing your request.")
	}
}
