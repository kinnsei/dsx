package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/plaenen/webx"
	"github.com/plaenen/webx/ui/toast"
	"github.com/starfederation/datastar-go/datastar"
)

type toastHandlers struct{}

func newToastHandlers() *toastHandlers {
	return &toastHandlers{}
}

func (t *toastHandlers) register(r chi.Router) {
	// Auto-dismiss toasts
	r.Get("/api/toast/info", t.showToast(toast.LevelInfo, "This is an informational message."))
	r.Get("/api/toast/success", t.showToast(toast.LevelSuccess, "Operation completed successfully!"))
	r.Get("/api/toast/warning", t.showToast(toast.LevelWarning, "Warning: please check your input."))
	r.Get("/api/toast/error", t.showToast(toast.LevelError, "Something went wrong. Please try again."))

	// Persistent toast
	r.Get("/api/toast/persistent", t.showPersistent())

	// Action toast
	r.Get("/api/toast/action", t.showAction())
	r.Get("/api/toast/action-callback", t.actionCallback())

	// Link toast
	r.Get("/api/toast/link", t.showLink())
}

func (t *toastHandlers) showToast(level toast.Level, message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)
		toast.Show(sse, level, message, 3000)
	}
}

func (t *toastHandlers) showPersistent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)
		toast.ShowPersistent(sse, toast.LevelWarning, "This toast will stay until you close it.")
	}
}

func (t *toastHandlers) showAction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wctx := webx.FromContext(r.Context())
		sse := datastar.NewSSE(w, r)
		toast.ShowAction(sse, toast.LevelError, "Item deleted.", "Undo", wctx.APIPath("/api/toast/action-callback"))
	}
}

func (t *toastHandlers) actionCallback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)
		toast.Show(sse, toast.LevelSuccess, "Action undone!", 2000)
	}
}

func (t *toastHandlers) showLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wctx := webx.FromContext(r.Context())
		sse := datastar.NewSSE(w, r)
		toast.ShowLink(sse, toast.LevelInfo, "New alert component available.", "View Alert", wctx.BasePath+"/components/alert", 5000)
	}
}
