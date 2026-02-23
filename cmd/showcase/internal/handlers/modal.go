package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/plaenen/webx"
	"github.com/plaenen/webx/cmd/showcase/internal/pages"
	"github.com/plaenen/webx/ds"
	"github.com/starfederation/datastar-go/datastar"
)

type modalHandlers struct{}

func newModalHandlers() *modalHandlers {
	return &modalHandlers{}
}

func (h *modalHandlers) register(r chi.Router) {
	r.Get("/api/modal/show", h.showModal())
	r.Get("/api/modal/show-wide", h.showWideModal())
	r.Get("/api/modal/confirm", h.showConfirm())
	r.Get("/api/modal/confirm-danger", h.showDangerConfirm())
	r.Post("/api/modal/confirmed", h.confirmed())
	r.Get("/api/modal/patch", h.patchContent())
	r.Get("/api/modal/redirect", h.redirect())
	r.Get("/api/modal/download", h.download())
}

func (h *modalHandlers) showModal() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)
		ds.Send.Modal(sse, pages.ModalContent())
	}
}

func (h *modalHandlers) showWideModal() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)
		ds.Send.Modal(sse, pages.ModalContentWide(), ds.WithModalMaxWidth("max-w-2xl"))
	}
}

func (h *modalHandlers) showConfirm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wxctx := webx.FromContext(r.Context())
		sse := datastar.NewSSE(w, r)
		ds.Send.Confirm(sse, "Are you sure you want to proceed with this action?",
			wxctx.APIPath("/api/modal/confirmed"),
		)
	}
}

func (h *modalHandlers) showDangerConfirm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wxctx := webx.FromContext(r.Context())
		sse := datastar.NewSSE(w, r)
		ds.Send.Confirm(sse, "This will permanently delete all data. This action cannot be undone.",
			wxctx.APIPath("/api/modal/confirmed"),
			ds.WithConfirmTitle("Danger Zone"),
			ds.WithConfirmLabel("Delete Everything"),
			ds.WithConfirmClass("btn btn-error"),
		)
	}
}

func (h *modalHandlers) confirmed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)
		ds.Send.Toast(sse, ds.ToastSuccess, "Action confirmed!")
	}
}

func (h *modalHandlers) patchContent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)
		timestamp := time.Now().Format("15:04:05")
		ds.Send.Patch(sse, pages.PatchedContent(timestamp))
	}
}

func (h *modalHandlers) redirect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wxctx := webx.FromContext(r.Context())
		sse := datastar.NewSSE(w, r)
		ds.Send.Redirect(sse, wxctx.BasePath+"/")
	}
}

func (h *modalHandlers) download() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wxctx := webx.FromContext(r.Context())
		sse := datastar.NewSSE(w, r)
		ds.Send.Download(sse, wxctx.APIPath("/api/modal/export.csv"), "export.csv")
	}
}
