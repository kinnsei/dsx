# AI Chat Component Guide

The `ui/aichat` package provides a complete conversational AI chat component with two modes: a collapsible widget and a full-page layout.

## Quick start

### Frontend (templ)

```go
import "github.com/laenen-partners/dsx/ui/aichat"

// Widget mode — collapsible sparkle bar (dashboards, sidebars)
@aichat.AIChat(aichat.Props{
    ID:        "my-chat",
    SubmitURL: "/api/chat/send",
})

// Full-page mode — fills parent container (drawers, dedicated pages)
@aichat.AIChat(aichat.Props{
    ID:        "my-chat",
    Mode:      aichat.ModeFullPage,
    SubmitURL: "/api/chat/send",
})
```

### Backend (handler)

```go
func sendMessage(w http.ResponseWriter, r *http.Request) {
    // 1. Read the user's input from signals.
    var signals aichat.AIChatSignals
    aichat.ReadSignals("my-chat", r, &signals)
    input := strings.TrimSpace(signals.Input)

    // 2. Create SSE and the chat sender.
    sse := datastar.NewSSE(w, r)
    chat := aichat.Chat(sse, "my-chat")

    // 3. Clear input and show user message.
    chat.ClearInput()
    chat.UserMessage(input)

    // 4. Show AI response.
    chat.ShowTyping()
    response := callYourAI(input)
    chat.HideTyping()
    chat.AssistantText(response)
}
```

## Modes

### Widget mode (default)

```go
@aichat.AIChat(aichat.Props{
    ID:          "assistant",
    SubmitURL:   "/api/chat/send",
    Placeholder: "Ask anything...",   // collapsed state text
    Shortcut:    "⌘K",               // keyboard hint
})
```

Renders a compact sparkle bar. Click to expand into a chat with `max-h-96` scrollable messages and a text input. Click ✕ or press Escape to collapse. Ideal for floating widgets on dashboards.

### Full-page mode

```go
@aichat.AIChat(aichat.Props{
    ID:               "project-chat",
    Mode:             aichat.ModeFullPage,
    SubmitURL:        "/api/chat/send",
    ReplyPlaceholder: "Type a message...",
})
```

Fills the parent container using `flex flex-col h-full`. Messages area grows with `flex-1`, compose bar pinned at bottom. No collapse/expand. Ideal for drawers, dedicated pages, or panels.

**Important:** The parent must have a defined height (e.g. `h-full`, `h-screen`, or explicit `height`).

## File attachments

Add `UploadURL` and `RemoveURL` to enable the `+` button for file attachments:

```go
@aichat.AIChat(aichat.Props{
    ID:        "project-chat",
    Mode:      aichat.ModeFullPage,
    SubmitURL: "/api/chat/send",
    UploadURL: "/api/upload/files",    // POST endpoint for file uploads
    RemoveURL: "/api/upload/remove",   // POST endpoint for file removal
})
```

This reuses the `ui/fileupload` handlers:

```go
import "github.com/laenen-partners/dsx/ui/fileupload"

store := fileupload.NewStore()
r.Post("/api/upload/files", fileupload.UploadHandler(store))
r.Post("/api/upload/remove", fileupload.RemoveHandler(store))
```

In the send handler, read and clear attached files:

```go
func sendMessage(store fileupload.Store) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var signals aichat.AIChatSignals
        aichat.ReadSignals("project-chat", r, &signals)

        // Get attached files.
        wxctx := dsx.FromContext(r.Context())
        key := fileupload.StoreKey(wxctx.SessionID, aichat.AttachmentsID("project-chat"))
        files := store.List(key)

        sse := datastar.NewSSE(w, r)
        chat := aichat.Chat(sse, "project-chat")

        chat.ClearInput()
        chat.ClearAttachments()

        // Build message with file info.
        msg := signals.Input
        if len(files) > 0 {
            names := make([]string, len(files))
            for i, f := range files {
                names[i] = f.Name
            }
            msg = fmt.Sprintf("%s\n📎 %d file(s): %s", msg, len(files), strings.Join(names, ", "))
            store.Clear(key)
        }

        chat.UserMessage(msg)

        // AI response...
        chat.ShowTyping()
        response := callYourAI(msg, files)
        chat.HideTyping()
        chat.AssistantText(response)
    }
}
```

## Streaming responses

For token-by-token streaming (like ChatGPT/Claude):

```go
chat := aichat.Chat(sse, "project-chat")

chat.ClearInput()
chat.UserMessage(input)

// Start streaming — appends an empty AI bubble with loading dots.
chat.StreamStart()

// Stream tokens — updates the bubble in place (merge mode).
accumulated := ""
for token := range tokenStream {
    accumulated += token
    chat.StreamText(accumulated)
}

// Done — removes loading dots, finalizes the bubble.
chat.StreamDone(accumulated)
```

`StreamStart` appends a bubble with a stable ID. `StreamText` merges into it (Datastar morph mode). `StreamDone` removes the loading indicator. The bubble stays in the message list as a regular message.

This works with any token source — Claude API, OpenAI, gRPC stream, local LLM, etc.

## Custom compose bar

Replace the default compose bar with `InputSlot`:

```go
@aichat.AIChat(aichat.Props{
    ID:        "project-chat",
    Mode:      aichat.ModeFullPage,
    InputSlot: myCustomComposeBar(),
})
```

Your component controls the entire compose area. It must handle submitting via `@post(submitURL)` and clearing input itself.

## Footer content

Add content below the compose bar with `FooterSlot`:

```go
@aichat.AIChat(aichat.Props{
    ID:         "project-chat",
    Mode:       aichat.ModeFullPage,
    SubmitURL:  "/api/chat/send",
    FooterSlot: disclaimer(),
})

templ disclaimer() {
    <p class="text-xs text-base-content/40 px-4 py-2 text-center">
        AI responses may be inaccurate. Always verify important information.
    </p>
}
```

## ChatSender API reference

All methods target the chat's `#{id}-messages` container.

| Method | What it does |
|--------|-------------|
| `chat.UserMessage(text)` | Appends a user bubble (right-aligned, dark) |
| `chat.AssistantText(text)` | Appends an assistant bubble (left-aligned, with sparkle avatar) |
| `chat.AssistantMessage(component)` | Appends an assistant bubble wrapping a templ component |
| `chat.Append(component)` | Appends any templ component to the messages area |
| `chat.ShowTyping()` | Appends a typing indicator (animated dots) |
| `chat.HideTyping()` | Removes the typing indicator |
| `chat.StreamStart()` | Appends an empty streaming bubble with loading dots |
| `chat.StreamText(accumulated)` | Updates the streaming bubble in place |
| `chat.StreamDone(fullText)` | Finalizes the streaming bubble (removes dots) |
| `chat.ClearInput()` | Clears the compose textarea via signal patch |
| `chat.ClearAttachments()` | Clears file attachment cards from the compose bar |

## Message types

Use the message templates directly for initial content or custom layouts:

```go
// In templ — initial messages rendered server-side
@aichat.AIChat(aichat.Props{...}) {
    @aichat.AssistantMessage() {
        @aichat.AssistantText("Welcome! How can I help?")
    }
    @aichat.UserMessage(aichat.UserMessageProps{Text: "Hello!"})
    @aichat.AssistantMessage() {
        @aichat.AssistantText("Hi there!")
        @aichat.QuickReplies(aichat.QuickRepliesProps{
            ChatID:    "my-chat",
            SubmitURL: "/api/chat/send",
            Options: []aichat.QuickReply{
                {Label: "Tell me more", Value: "Tell me more"},
                {Label: "Thanks!", Value: "Thanks!"},
            },
        })
    }
}
```

## ID conventions

| Function | Returns |
|----------|---------|
| `aichat.MessagesID(chatID)` | `{chatID}-messages` — message list container |
| `aichat.AttachmentsID(chatID)` | `{chatID}-attachments` — file attachments container |

These are the SSE patch targets. The `ChatSender` handles this automatically — you rarely need to use them directly.
