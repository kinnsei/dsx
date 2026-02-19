package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/plaenen/webx/ui/toast"
	"github.com/starfederation/datastar-go/datastar"
)

type toastHandlers struct{}

func newToastHandlers() *toastHandlers {
	return &toastHandlers{}
}

func (t *toastHandlers) register(r chi.Router) {
	r.Get("/api/toast/info", t.showToast(toast.LevelInfo, "This is an informational message."))
	r.Get("/api/toast/success", t.showToast(toast.LevelSuccess, "Operation completed successfully!"))
	r.Get("/api/toast/warning", t.showToast(toast.LevelWarning, "Warning: please check your input."))
	r.Get("/api/toast/error", t.showToast(toast.LevelError, "Something went wrong. Please try again."))
}

func (t *toastHandlers) showToast(level toast.Level, message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)
		toast.Show(sse, level, message, 3000)
	}
}
