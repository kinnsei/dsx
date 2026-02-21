package aichat

import (
	"github.com/a-h/templ"
	"github.com/starfederation/datastar-go/datastar"
)

const typingIndicatorID = "typing-indicator"

// ChatSender provides SSE operations scoped to a specific AI chat instance.
type ChatSender struct {
	sse      *datastar.ServerSentEventGenerator
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
