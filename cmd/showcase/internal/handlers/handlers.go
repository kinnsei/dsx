package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/plaenen/webx/stream"
)

// Handlers wires all showcase SSE/API handlers.
type Handlers struct {
	form       *formHandlers
	upload     *uploadHandlers
	toast      *toastHandlers
	stream     *streamHandlers
	drawer     *drawerHandlers
	modal      *modalHandlers
	commandbar *commandbarHandlers
	aichat     *aichatHandlers
	yamltree   *yamltreeHandlers
}

func New(broker *stream.Broker) *Handlers {
	fileStore := newFileStore()
	return &Handlers{
		form:       newFormHandlers(),
		upload:     newUploadHandlers(fileStore),
		toast:      newToastHandlers(),
		stream:     newStreamHandlers(broker),
		drawer:     newDrawerHandlers(),
		modal:      newModalHandlers(),
		commandbar: newCommandbarHandlers(),
		aichat:     newAIChatHandlers(),
		yamltree:   newYamlTreeHandlers(),
	}
}

// RegisterRoutes mounts all API handlers onto the given router.
func (h *Handlers) RegisterRoutes(r chi.Router) {
	h.form.register(r)
	h.upload.register(r)
	h.toast.register(r)
	h.stream.register(r)
	h.drawer.register(r)
	h.modal.register(r)
	h.commandbar.register(r)
	h.aichat.register(r)
	h.yamltree.register(r)
}
