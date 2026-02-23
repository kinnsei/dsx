package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/plaenen/webx"
	"github.com/plaenen/webx/ds"
	"github.com/starfederation/datastar-go/datastar"
)

type toastHandlers struct{}

func newToastHandlers() *toastHandlers {
	return &toastHandlers{}
}

func (t *toastHandlers) register(r chi.Router) {
	// Auto-dismiss toasts
	r.Get("/api/toast/info", t.showToast(ds.ToastInfo, "This is an informational message."))
	r.Get("/api/toast/success", t.showToast(ds.ToastSuccess, "Operation completed successfully!"))
	r.Get("/api/toast/warning", t.showToast(ds.ToastWarning, "Warning: please check your input."))
	r.Get("/api/toast/error", t.showToast(ds.ToastError, "Something went wrong. Please try again."))

	// Persistent toast
	r.Get("/api/toast/persistent", t.showPersistent())

	// Action toast
	r.Get("/api/toast/action", t.showAction())
	r.Get("/api/toast/action-callback", t.actionCallback())

	// Link toast
	r.Get("/api/toast/link", t.showLink())
}

func (t *toastHandlers) showToast(level ds.ToastLevel, message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)
		ds.Send.Toast(sse, level, message)
	}
}

func (t *toastHandlers) showPersistent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)
		ds.Send.Toast(sse, ds.ToastWarning, "This toast will stay until you close it.", ds.WithToastPersistent())
	}
}

func (t *toastHandlers) showAction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wxctx := webx.FromContext(r.Context())
		sse := datastar.NewSSE(w, r)
		ds.Send.Toast(sse, ds.ToastError, "Item deleted.", ds.WithToastAction("Undo", wxctx.APIPath("/api/toast/action-callback")))
	}
}

func (t *toastHandlers) actionCallback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)
		ds.Send.Toast(sse, ds.ToastSuccess, "Action undone!", ds.WithToastDuration(2000))
	}
}

func (t *toastHandlers) showLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wxctx := webx.FromContext(r.Context())
		sse := datastar.NewSSE(w, r)
		ds.Send.Toast(sse, ds.ToastInfo, "New alert component available.",
			ds.WithToastLink("View Alert", wxctx.BasePath+"/components/alert"),
			ds.WithToastDuration(5000),
		)
	}
}
