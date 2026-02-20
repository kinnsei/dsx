package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/plaenen/webx/stream"
)

// Handlers wires all showcase SSE/API handlers.
type Handlers struct {
	validate *validateHandlers
	parse    *parseHandlers
	form     *formHandlers
	upload   *uploadHandlers
	preview  *previewHandlers
	toast    *toastHandlers
	stream   *streamHandlers
	drawer   *drawerHandlers
	modal    *modalHandlers
}

func New(broker *stream.Broker) *Handlers {
	fileStore := newFileStore()
	return &Handlers{
		validate: newValidateHandlers(),
		parse:    newParseHandlers(),
		form:     newFormHandlers(),
		upload:   newUploadHandlers(fileStore),
		preview:  newPreviewHandlers(),
		toast:    newToastHandlers(),
		stream:   newStreamHandlers(broker),
		drawer:   newDrawerHandlers(),
		modal:    newModalHandlers(),
	}
}

// RegisterRoutes mounts all API handlers onto the given router.
func (h *Handlers) RegisterRoutes(r chi.Router) {
	h.validate.register(r)
	h.parse.register(r)
	h.form.register(r)
	h.upload.register(r)
	h.preview.register(r)
	h.toast.register(r)
	h.stream.register(r)
	h.drawer.register(r)
	h.modal.register(r)
}
