package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/plaenen/webx/stream"
)

// Handlers wires all showcase SSE/API handlers.
type Handlers struct {
	validate *validateHandlers
	form     *formHandlers
	upload   *uploadHandlers
	toast    *toastHandlers
	stream   *streamHandlers
	drawer   *drawerHandlers
	modal    *modalHandlers
	butler   *butlerHandlers
	aichat   *aichatHandlers
}

func New(broker *stream.Broker) *Handlers {
	fileStore := newFileStore()
	return &Handlers{
		validate: newValidateHandlers(),
		form:     newFormHandlers(),
		upload:   newUploadHandlers(fileStore),
		toast:    newToastHandlers(),
		stream:   newStreamHandlers(broker),
		drawer:   newDrawerHandlers(),
		modal:    newModalHandlers(),
		butler:   newButlerHandlers(),
		aichat:   newAIChatHandlers(),
	}
}

// RegisterRoutes mounts all API handlers onto the given router.
func (h *Handlers) RegisterRoutes(r chi.Router) {
	h.validate.register(r)
	h.form.register(r)
	h.upload.register(r)
	h.toast.register(r)
	h.stream.register(r)
	h.drawer.register(r)
	h.modal.register(r)
	h.butler.register(r)
	h.aichat.register(r)
}
