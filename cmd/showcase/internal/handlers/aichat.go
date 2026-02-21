package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/plaenen/webx/ds"
	"github.com/plaenen/webx/ui/aichat"
	"github.com/starfederation/datastar-go/datastar"
)

type aichatHandlers struct{}

func newAIChatHandlers() *aichatHandlers {
	return &aichatHandlers{}
}

func (h *aichatHandlers) register(r chi.Router) {
	r.Post("/api/aichat/send", h.send())
	r.Post("/api/aichat/action", h.action())
}

const demoChatID = "demo-aichat"

func (h *aichatHandlers) send() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var signals aichat.AIChatSignals
		if err := aichat.ReadSignals(demoChatID, r, &signals); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		input := strings.TrimSpace(signals.Input)
		if input == "" {
			return
		}

		sse := datastar.NewSSE(w, r)
		messagesSelector := "#" + aichat.MessagesID(demoChatID)

		// 1. Append user message
		_ = sse.PatchElementTempl(
			aichat.UserMessage(aichat.UserMessageProps{Text: input}),
			datastar.WithSelector(messagesSelector),
			datastar.WithModeAppend(),
		)

		// 2. Show typing indicator
		_ = sse.PatchElementTempl(
			aichat.Typing(aichat.TypingProps{ID: "typing-indicator"}),
			datastar.WithSelector(messagesSelector),
			datastar.WithModeAppend(),
		)

		// 3. Simulate AI thinking
		time.Sleep(800 * time.Millisecond)

		// 4. Remove typing indicator and append AI response
		_ = sse.RemoveElementByID("typing-indicator")

		submitURL := "/showcase/api/aichat/send"
		actionURL := "/showcase/api/aichat/action"

		lower := strings.ToLower(input)
		switch {
		case contains(lower, "cancel", "subscription"):
			_ = sse.PatchElementTempl(
				subscriptionResponse(demoChatID, submitURL),
				datastar.WithSelector(messagesSelector),
				datastar.WithModeAppend(),
			)

		case contains(lower, "plan", "date"):
			_ = sse.PatchElementTempl(
				dateNightResponse(demoChatID, submitURL),
				datastar.WithSelector(messagesSelector),
				datastar.WithModeAppend(),
			)

		case contains(lower, "buy", "needs", "football", "shoes", "boots"):
			_ = sse.PatchElementTempl(
				shoppingResponse(input, actionURL),
				datastar.WithSelector(messagesSelector),
				datastar.WithModeAppend(),
			)

		case contains(lower, "netflix", "spotify", "youtube", "show all"):
			_ = sse.PatchElementTempl(
				subscriptionDetailResponse(lower, demoChatID, submitURL),
				datastar.WithSelector(messagesSelector),
				datastar.WithModeAppend(),
			)

		case contains(lower, "friday", "saturday", "week"):
			_ = sse.PatchElementTempl(
				dateConfirmResponse(lower),
				datastar.WithSelector(messagesSelector),
				datastar.WithModeAppend(),
			)

		default:
			_ = sse.PatchElementTempl(
				defaultResponse(input, actionURL),
				datastar.WithSelector(messagesSelector),
				datastar.WithModeAppend(),
			)
		}
	}
}

func (h *aichatHandlers) action() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)
		ds.Send.Toast(sse, ds.ToastSuccess, "Added to inbox!")
	}
}

func contains(s string, keywords ...string) bool {
	for _, kw := range keywords {
		if strings.Contains(s, kw) {
			return true
		}
	}
	return false
}
