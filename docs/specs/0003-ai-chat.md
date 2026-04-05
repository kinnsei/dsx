# Spec: dsx — AI Chat Component

> **Status:** Closed — already implemented
>
> **Library**: `github.com/laenen-partners/dsx`
> **Package**: `ui/aichat`
> **Showcase**: `showcase/aichat`
> **Depends on**: `ui/chat`, `ui/form`, `ui/badge`, `ds` (Datastar helpers)

---

## Resolution

The core of this spec is **already implemented** in `ui/aichat/`. The existing implementation covers the full chat loop with a cleaner API than what this spec proposes.

### What exists today

| Spec item | Status | Implementation |
|-----------|--------|----------------|
| Chat container (messages + compose) | Done | `aichat.AiChat()` in `aichat.templ` |
| Message bubbles (user, ai, system, file) | Done | `aichat.MessageBubble()` in `messages.templ` |
| Compose bar with Enter/Shift+Enter | Done | Built into `aichat.AiChat()` |
| Server-side append user message | Done | `aichat.Chat(sse, id).UserMessage(text)` |
| Show/hide typing indicator | Done | `chat.ShowTyping()` / `chat.HideTyping()` |
| Stream AI response (per-token) | Done | `chat.StreamStart()` → `chat.StreamText(accumulated)` → `chat.StreamDone(fullText)` — merge-mode patching into a stable bubble |
| Stream AI response (full) | Done | `chat.AssistantText(text)` — appends a complete bubble |
| Clear compose after submit | Done | `chat.HideTyping()` resets input; or patch signal directly |
| Auto-scroll | Done | Uses Datastar SSE append mode |
| Showcase demo | Done | `cmd/showcase/internal/handlers/aichat.go` with mock streaming |
| Quick reply buttons | Done | `aichat.AICard` with `QuickReplies` in `card.templ` |
| ReadSignals helper | Done | `aichat.ReadSignals(chatID, r, &signals)` |

### What this spec proposes that we won't implement

**`StreamBridge[T]` (Connect-go stream → SSE)** — This couples the library to a specific RPC framework. The bridge belongs in the consuming application (e.g. steward), not in dsx. The existing `ChatSender` API is framework-agnostic: any handler can call `chat.AssistantText()` in a loop regardless of where the tokens come from (Connect-go, gRPC, HTTP, WebSocket, local LLM).

**`PatchAIToken` / `PatchAIDone` as separate functions** — Implemented as `chat.StreamStart()`, `chat.StreamText()`, and `chat.StreamDone()` on `ChatSender`. Uses a `StreamBubble` with a stable ID and Datastar merge-mode patching to update the bubble in place as tokens arrive.

**Standalone showcase binary** — Low value since the main showcase already demonstrates the component fully, including mock AI streaming.

### Existing API (for reference)

```go
// Handler side
sse := datastar.NewSSE(w, r)
chat := aichat.Chat(sse, "my-chat")

chat.UserMessage("What is Go?")       // append user bubble
chat.StreamStart()                     // append empty AI bubble with loading dots

// Stream tokens (from any source — Connect-go, gRPC, HTTP, local LLM)
accumulated := ""
for token := range tokenStream {
    accumulated += token
    chat.StreamText(accumulated)       // merge into AI bubble in place
}
chat.StreamDone(accumulated)           // finalize — remove loading dots

// Or for a complete response (no streaming):
chat.AssistantText("Go is a language.") // append full bubble at once
```

---

## Original spec (for historical context)

## Context

Every conversational AI app needs the same pattern: message list, compose bar, AI streaming response. Steward, and any future app using dsx, shouldn't rebuild this from scratch. The existing `ui/chat` component provides message bubbles but not the full chat experience (compose, submit, stream, scroll).

This spec defines a complete AI chat component that handles the full loop: compose → submit → append user bubble → stream AI response → save → scroll.

---

## Component design

### `aichat.Chat` — the container

```go
type Props struct {
    ID          string           // unique component ID (required)
    Class       string           // additional CSS classes
    Attributes  templ.Attributes
    SubmitURL   string           // POST endpoint for sending messages
    Messages    []Message        // initial message list (server-rendered)
    Placeholder string           // compose textarea placeholder (default: "Type a message...")
    ShowAvatar  bool             // show avatars on messages (default: true)
}

type Message struct {
    ID            string
    Type          string // "user", "other", "ai", "system", "file"
    SenderName    string
    SenderInitial string
    Content       string
    SentAt        string
    IsCurrentUser bool
    // File messages
    DocumentID string
    Filename   string
    FileURL    string
}
```

**Rendering**:
```
┌────────────────────────────────────┐
│  #messages (scrollable)            │
│                                    │
│  [system] Project created          │
│                                    │
│  Alice                             │
│  ┌──────────────────────┐          │
│  │ Can you help with... │          │
│  └──────────────────────┘          │
│                                    │
│               ┌──────────────────┐ │
│               │ Sure, I can...   │ │ You
│               └──────────────────┘ │
│                                    │
│  AI                                │
│  ┌──────────────────────┐          │
│  │ Here's what I found: │          │
│  │ ...streaming...       │          │
│  └──────────────────────┘          │
│                                    │
├────────────────────────────────────┤
│ [textarea] Type a message... [➤]   │ compose bar
└────────────────────────────────────┘
```

**IDs for Datastar patching**:
- `#{id}-messages` — message list container (append target)
- `#{id}-ai-response` — AI streaming bubble (merge target)

### `aichat.MessageBubble` — renders a single message

Uses `chat.Chat` + `chat.Bubble` internally. Handles all message types:

| Type | Position | Variant | Avatar |
|---|---|---|---|
| `user` (current) | End (right) | Primary | Initials |
| `other` | Start (left) | Default | Initials |
| `ai` | Start (left) | Secondary | "AI" badge |
| `system` | Center | — | No avatar, muted pill |
| `file` | Center | — | No avatar, link pill |

### `aichat.ComposeBar` — the input form

Uses `form.Form` internally with `ds.Post`:
- Textarea bound to `${id}.text` signal
- Submit button (send icon)
- Enter to submit (Shift+Enter for newline) — via `data-on:keydown`
- Auto-clear text signal after successful submit

```go
type ComposeBarProps struct {
    ID          string // must match parent Chat ID
    SubmitURL   string // POST endpoint
    Placeholder string
    Class       string
}
```

### `aichat.AIBubble` — streaming AI response

```go
// AIBubble renders the streaming AI message.
// id="{chatID}-ai-response" for Datastar merge patching.
templ AIBubble(chatID string, text string, streaming bool)

// AIBubbleAppend wraps AIBubble for SSE append into the message list.
templ AIBubbleAppend(chatID string, text string, streaming bool)
```

When `streaming=true`: shows loading dots indicator.
When `streaming=false`: shows final text, no indicator.

---

## Server-side helpers

### `aichat.Handler` — SSE response helpers

```go
package aichat

// AppendUserMessage sends an SSE patch that appends a user message bubble
// to the message list.
func AppendUserMessage(sse *datastar.SSE, chatID string, msg Message)

// AppendAIPlaceholder sends an SSE patch that appends an empty AI bubble
// with loading indicator to the message list.
func AppendAIPlaceholder(sse *datastar.SSE, chatID string)

// PatchAIToken sends an SSE patch that updates the AI bubble with
// accumulated text. Call for each token.
func PatchAIToken(sse *datastar.SSE, chatID string, accumulated string)

// PatchAIDone sends a final SSE patch that marks the AI bubble as complete
// (removes loading indicator).
func PatchAIDone(sse *datastar.SSE, chatID string, fullText string)

// ClearCompose sends an SSE signal patch to reset the compose textarea.
func ClearCompose(sse *datastar.SSE, chatID string)
```

### `aichat.StreamBridge` — Connect-go stream → SSE

Generic bridge that forwards a Connect-go server stream to Datastar SSE:

```go
// StreamBridge reads from a Connect-go server stream and forwards
// each chunk as a Datastar SSE patch into the AI bubble.
// Returns the accumulated full response text.
func StreamBridge[T interface{ GetText() string; GetDone() bool }](
    sse *datastar.SSE,
    chatID string,
    stream interface{ Receive() bool; Msg() T; Err() error },
) (string, error)
```

Usage in steward's web handler:

```go
func (w *Web) projectSendMessage(rw http.ResponseWriter, r *http.Request) {
    // ... send user message via API ...

    sse := datastar.NewSSE(rw, r)
    aichat.AppendUserMessage(sse, "project-chat", msg)
    aichat.ClearCompose(sse, "project-chat")
    aichat.AppendAIPlaceholder(sse, "project-chat")

    // Bridge the Connect-go stream to SSE.
    aiStream, _ := client.StreamAIResponse(r.Context(), req)
    fullText, _ := aichat.StreamBridge(sse, "project-chat", aiStream)
    _ = fullText // already saved by API
}
```

---

## Showcase: `showcase/aichat`

Self-contained demo app with no external dependencies:

```go
// showcase/aichat/main.go
func main() {
    mux := http.NewServeMux()

    // Serve the demo page.
    mux.HandleFunc("GET /", handlePage)

    // Handle message submission + mock AI streaming.
    mux.HandleFunc("POST /messages", handleSendMessage)

    http.ListenAndServe(":8080", mux)
}

func handleSendMessage(w http.ResponseWriter, r *http.Request) {
    var signals struct{ Text string `json:"text"` }
    ds.ReadSignals("demo-chat", r, &signals)

    sse := datastar.NewSSE(w, r)

    // Append user message.
    aichat.AppendUserMessage(sse, "demo-chat", aichat.Message{
        Type:    "user", Content: signals.Text,
        SentAt:  time.Now().Format("15:04"), IsCurrentUser: true,
    })
    aichat.ClearCompose(sse, "demo-chat")

    // Stream mock AI response (word by word with delay).
    aichat.AppendAIPlaceholder(sse, "demo-chat")
    response := "This is a simulated AI response that streams word by word."
    words := strings.Fields(response)
    var accumulated strings.Builder
    for _, word := range words {
        time.Sleep(100 * time.Millisecond)
        if accumulated.Len() > 0 {
            accumulated.WriteString(" ")
        }
        accumulated.WriteString(word)
        aichat.PatchAIToken(sse, "demo-chat", accumulated.String())
    }
    aichat.PatchAIDone(sse, "demo-chat", accumulated.String())
}
```

**Demo page** (`showcase/aichat/page.templ`):

```go
@layouts.Base(layouts.BaseProps{Title: "AI Chat Showcase"}) {
    <div class="max-w-2xl mx-auto h-screen">
        @aichat.Chat(aichat.Props{
            ID:        "demo-chat",
            SubmitURL: "/messages",
            Messages: []aichat.Message{
                {Type: "system", Content: "Welcome to the AI chat demo."},
            },
        })
    </div>
}
```

### Running the showcase

```bash
cd dsx/showcase/aichat
go run .
# Open http://localhost:8080
```

No Ollama, no database, no auth — just a working chat with simulated streaming.

---

## Integration with steward

After the component is built, steward replaces its custom chat rendering:

```go
// Before (custom templates)
for _, m := range props.Messages {
    @messageBubble(m)
}

// After (dsx component)
@aichat.Chat(aichat.Props{
    ID:        "project-chat",
    SubmitURL: "/projects/" + props.Project.ID + "/messages",
    Messages:  toAIChatMessages(props.Messages),
})
```

The web handler uses `aichat.StreamBridge` instead of manual SSE patching.

---

## Files to create

| File | Purpose |
|---|---|
| `ui/aichat/aichat.templ` | Chat container + message list + compose bar |
| `ui/aichat/message.templ` | Message bubble variants (user, ai, system, file) |
| `ui/aichat/compose.templ` | Compose bar using form.Form |
| `ui/aichat/ai_bubble.templ` | Streaming AI bubble with loading state |
| `ui/aichat/handler.go` | Server-side SSE helpers (append, patch, clear) |
| `ui/aichat/bridge.go` | Generic Connect-go stream → Datastar SSE bridge |
| `showcase/aichat/main.go` | Demo server with mock AI streaming |
| `showcase/aichat/page.templ` | Demo page |

## Acceptance criteria

- [ ] Showcase runs standalone: `go run ./showcase/aichat` → working chat at localhost:8080
- [ ] User sends message → bubble appears immediately (no page reload)
- [ ] AI response streams word by word with loading indicator
- [ ] Loading dots disappear when stream completes
- [ ] Enter submits, Shift+Enter adds newline
- [ ] Textarea clears after successful submit
- [ ] Messages auto-scroll to bottom on new message
- [ ] All message types render correctly (user, other, ai, system, file)
- [ ] Component works with any backend — just needs a POST endpoint that returns SSE
- [ ] StreamBridge works with any Connect-go server stream that has `GetText()/GetDone()` methods
