package aichat

import (
	"fmt"
	"strings"

	"github.com/a-h/templ"
	"github.com/starfederation/datastar-go/datastar"
)

const typingIndicatorID = "typing-indicator"

// ChatSender provides SSE operations scoped to a specific AI chat instance.
type ChatSender struct {
	sse      *datastar.ServerSentEventGenerator
	chatID   string
	selector string
}

// Chat returns a ChatSender scoped to the given chat ID.
// All operations target the chat's messages container.
//
//	chat := aichat.Chat(sse, "my-chat")
//	chat.UserMessage("Hello")
//	chat.ShowTyping()
//	chat.HideTyping()
//	chat.AssistantText("Hi there!")
func Chat(sse *datastar.ServerSentEventGenerator, chatID string) *ChatSender {
	return &ChatSender{
		sse:      sse,
		chatID:   chatID,
		selector: "#" + MessagesID(chatID),
	}
}

// UserMessage appends a user message bubble.
func (c *ChatSender) UserMessage(text string) error {
	return c.sse.PatchElementTempl(
		UserMessage(UserMessageProps{Text: text}),
		datastar.WithSelector(c.selector),
		datastar.WithModeAppend(),
	)
}

// AssistantMessage appends an assistant message wrapping a templ component.
func (c *ChatSender) AssistantMessage(content templ.Component) error {
	return c.sse.PatchElementTempl(
		AssistantMessageWithContent(content),
		datastar.WithSelector(c.selector),
		datastar.WithModeAppend(),
	)
}

// AssistantText appends a simple text response from the assistant.
func (c *ChatSender) AssistantText(text string) error {
	return c.AssistantMessage(AssistantText(text))
}

// ShowTyping appends a typing indicator.
func (c *ChatSender) ShowTyping() error {
	return c.sse.PatchElementTempl(
		Typing(TypingProps{ID: typingIndicatorID}),
		datastar.WithSelector(c.selector),
		datastar.WithModeAppend(),
	)
}

// HideTyping removes the typing indicator.
func (c *ChatSender) HideTyping() error {
	return c.sse.RemoveElementByID(typingIndicatorID)
}

// Append appends any templ component to the messages area.
func (c *ChatSender) Append(component templ.Component) error {
	return c.sse.PatchElementTempl(
		component,
		datastar.WithSelector(c.selector),
		datastar.WithModeAppend(),
	)
}

// ClearAttachments empties the file attachments list via SSE patch.
func (c *ChatSender) ClearAttachments() error {
	return c.sse.PatchElements(
		fmt.Sprintf(`<div id="%s-list"></div>`, AttachmentsID(c.chatID)),
	)
}

// ClearInput resets the compose input signal to empty via SSE patch.
func (c *ChatSender) ClearInput() error {
	sanitizedID := strings.ReplaceAll(c.chatID, "-", "_")
	return c.sse.MarshalAndPatchSignals(map[string]any{
		sanitizedID: map[string]any{
			"input": "",
		},
	})
}

// StreamStart appends an empty streaming AI bubble to the messages area.
// The bubble has a stable ID so subsequent StreamText calls merge into it.
// Call this once before streaming tokens.
//
//	chat.StreamStart()
//	for token := range tokens {
//	    accumulated += token
//	    chat.StreamText(accumulated)
//	}
//	chat.StreamDone(fullText)
func (c *ChatSender) StreamStart() error {
	return c.sse.PatchElementTempl(
		StreamBubble(c.chatID, "", true),
		datastar.WithSelector(c.selector),
		datastar.WithModeAppend(),
	)
}

// StreamText updates the streaming AI bubble with accumulated text.
// Uses merge mode (Datastar's default morph) to update in place.
func (c *ChatSender) StreamText(accumulated string) error {
	return c.sse.PatchElementTempl(
		StreamBubble(c.chatID, accumulated, true),
	)
}

// StreamDone finalizes the streaming AI bubble, removing the loading
// indicator. The bubble remains in the message list as a regular message.
func (c *ChatSender) StreamDone(fullText string) error {
	return c.sse.PatchElementTempl(
		StreamBubble(c.chatID, fullText, false),
	)
}
