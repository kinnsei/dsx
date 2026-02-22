package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/plaenen/webx/ds"
	"github.com/plaenen/webx/ui/aichat"
	"github.com/plaenen/webx/ui/commandbar"
	"github.com/starfederation/datastar-go/datastar"
)

type aichatHandlers struct{}

func newAIChatHandlers() *aichatHandlers {
	return &aichatHandlers{}
}

func (h *aichatHandlers) register(r chi.Router) {
	r.Post("/api/aichat/send", h.send(demoChatID, h.readAIChatInput))
	r.Post("/api/aichat/send-combined", h.send(combinedChatID, h.readCommandBarInput))
	r.Post("/api/aichat/action", h.action())
	r.Post("/api/aichat/upload", h.action())
	r.Post("/api/aichat/voice", h.action())
}

const (
	demoChatID     = "demo-aichat"
	combinedChatID = "demo-combined"
	combinedBarID  = "combined-bar"
)

// inputReader extracts the user's text from the request signals.
type inputReader func(r *http.Request) (string, error)

func (h *aichatHandlers) readAIChatInput(r *http.Request) (string, error) {
	var signals aichat.AIChatSignals
	if err := aichat.ReadSignals(demoChatID, r, &signals); err != nil {
		return "", fmt.Errorf("read aichat signals: %w", err)
	}
	return signals.Input, nil
}

func (h *aichatHandlers) readCommandBarInput(r *http.Request) (string, error) {
	var signals commandbar.CommandBarSignals
	if err := ds.ReadSignals(combinedBarID, r, &signals); err != nil {
		return "", fmt.Errorf("read commandbar signals: %w", err)
	}
	return signals.Text, nil
}

func (h *aichatHandlers) send(chatID string, readInput inputReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		text, err := readInput(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		input := strings.TrimSpace(text)
		if input == "" {
			return
		}

		sse := datastar.NewSSE(w, r)
		chat := aichat.Chat(sse, chatID)

		_ = chat.UserMessage(input)
		_ = chat.ShowTyping()
		time.Sleep(800 * time.Millisecond)
		_ = chat.HideTyping()

		submitURL := "/showcase/api/aichat/send"
		if chatID == combinedChatID {
			submitURL = "/showcase/api/aichat/send-combined"
		}
		actionURL := "/showcase/api/aichat/action"

		lower := strings.ToLower(input)
		switch {
		case contains(lower, "cancel", "subscription"):
			_ = chat.Append(subscriptionResponse(chatID, submitURL))

		case contains(lower, "plan", "date"):
			_ = chat.Append(dateNightResponse(chatID, submitURL))

		case contains(lower, "buy", "needs", "football", "shoes", "boots"):
			_ = chat.Append(shoppingResponse(input, actionURL))

		case contains(lower, "netflix", "spotify", "youtube", "show all"):
			_ = chat.Append(subscriptionDetailResponse(lower, chatID, submitURL))

		case contains(lower, "friday", "saturday", "week"):
			_ = chat.Append(dateConfirmResponse(lower))

		default:
			_ = chat.Append(defaultResponse(input, actionURL))
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
