# aichat — AI Conversation Component

> **Import**: `"github.com/laenen-partners/dsx/ui/aichat"`

A conversational AI chat widget with collapsed/expanded states, multi-turn messaging, typing indicators, quick replies, rich response cards, and member assignment. Messages are delivered via SSE.

---

## Component Overview

| Component | Description |
|---|---|
| `AIChat` | Main container — collapsed bar, expandable messages area + input |
| `UserMessage` | Dark pill bubble, right-aligned (user's message) |
| `AssistantMessage` | Light bubble with sparkle avatar, left-aligned (takes children) |
| `AssistantMessageWithContent` | Same as above, takes `templ.Component` instead of children |
| `AssistantText` | Text bubble for use inside an assistant message |
| `Typing` | Sparkle avatar + animated dots (AI is thinking) |
| `QuickReplies` | Row of clickable reply buttons below an AI message |
| `ResponseCard` | Rich card with title, tags, description (takes children) |
| `AssignRow` | "Assign:" label + avatar chips |
| `ActionBar` | Primary action button + discard button |

---

## AIChat — Main Container

```go
type Props struct {
    ID               string
    Class            string
    Attributes       templ.Attributes
    Placeholder      string          // collapsed placeholder (default "Tell your butler something...")
    ReplyPlaceholder string          // reply input placeholder (default "Type a reply...")
    SubmitURL        string          // POST endpoint for user messages
    Shortcut         string          // keyboard shortcut hint (default "⌘K")
    InputSlot        templ.Component // optional: replaces default input row
}
```

**Basic usage:**

```templ
@aichat.AIChat(aichat.Props{
    ID:        "my-chat",
    SubmitURL: "/api/chat/send",
})
```

**With embedded command bar** (text + file + voice input):

```templ
@aichat.AIChat(aichat.Props{
    ID: "my-chat",
    InputSlot: commandbar.CommandBar(commandbar.Props{
        ID:       "my-bar",
        Embedded: true,
        SubmitURL: "/api/chat/send",
        UploadURL: "/api/chat/upload",
        VoiceURL:  "/api/chat/voice",
    }),
})
```

**Signals:** `{ open: false, input: "" }` — namespaced by component ID.

**Behavior:**
- **Collapsed**: sparkle icon + placeholder + keyboard shortcut. Click to expand.
- **Expanded**: scrollable messages area + input row + × close button.
- Enter submits, Escape closes. Input is cleared after submit.
- When `InputSlot` is set, it replaces the default text input row.
- The × close button is always rendered at the AI chat level.

---

## Message Components

### UserMessage

```templ
@aichat.UserMessage(aichat.UserMessageProps{
    Text: "max needs football boots size 38",
})
```

Dark pill, right-aligned. No avatar.

### AssistantMessage

Takes children — use in templ files:

```templ
@aichat.AssistantMessage() {
    @aichat.AssistantText("Love it. When are you thinking?")
}
```

Light bubble with sparkle avatar, left-aligned. Can contain text, cards, or any content.

### AssistantMessageWithContent

Takes a `templ.Component` — use from Go handler code where `{ children... }` isn't available:

```go
chat.AssistantMessage(aichat.AssistantText("Hello"))
```

### AssistantText

Simple text bubble. Use inside `AssistantMessage`:

```templ
@aichat.AssistantText("Here are your options.")
```

### Typing

Animated dots indicator (AI is thinking):

```templ
@aichat.Typing(aichat.TypingProps{ID: "typing-indicator"})
```

Give it an ID so you can remove it later via `sse.RemoveElementByID()` or `chat.HideTyping()`.

### QuickReplies

Row of bordered buttons. Clicking sets the input signal and POSTs:

```templ
@aichat.QuickReplies(aichat.QuickRepliesProps{
    ChatID:    "my-chat",
    SubmitURL: "/api/chat/send",
    Options: []aichat.QuickReply{
        {Label: "Netflix €12.99", Value: "Netflix"},
        {Label: "Spotify €17.99", Value: "Spotify"},
    },
})
```

**Required props for posting:** `ChatID` (matches the parent AIChat ID) and `SubmitURL`.

---

## Card Components

### ResponseCard

Rich card with title, tags, description. Takes children for action rows:

```templ
@aichat.ResponseCard(aichat.ResponseCardProps{
    Title:       "Buy football boots for Max (size 38)",
    Description: "Decathlon has options €35–85.",
    Tags: []aichat.CardTag{
        {Label: "Shopping"},
        {Label: "Needs input", Variant: "warning"},
    },
}) {
    @aichat.AssignRow(aichat.AssignRowProps{Members: members})
    @aichat.ActionBar(aichat.ActionBarProps{
        PrimaryLabel: "Add to inbox ✦",
        PrimaryURL:   "/api/chat/action",
        DiscardURL:   "/api/chat/action",
    })
}
```

**Tag variants:** `""` (default outline), `"warning"`, `"success"`, `"info"`, `"error"`

### AssignRow

"Assign:" label with avatar chips:

```go
type Member struct {
    Initials string // "ME", "SA"
    Name     string // display name
    Color    string // Tailwind bg class: "bg-success", "bg-secondary"
    Active   bool   // shows ring highlight
}
```

### ActionBar

Primary button + discard button. Both POST on click:

```go
type ActionBarProps struct {
    PrimaryLabel string // "Add to inbox ✦"
    PrimaryURL   string // POST endpoint
    DiscardURL   string // POST endpoint
}
```

---

## Server-Side Chat Helper

`aichat.Chat(sse, chatID)` returns a `*ChatSender` with operations scoped to a specific chat instance. All methods target the chat's `#<chatID>-messages` container.

```go
chat := aichat.Chat(sse, "my-chat")
```

### Methods

| Method | Description |
|---|---|
| `chat.UserMessage(text)` | Append a user message bubble |
| `chat.AssistantText(text)` | Append a simple text response from the assistant |
| `chat.AssistantMessage(component)` | Append an assistant message wrapping any templ component |
| `chat.ShowTyping()` | Append typing indicator |
| `chat.HideTyping()` | Remove typing indicator |
| `chat.Append(component)` | Append any templ component (cards, quick replies, etc.) |

### Handler example

```go
func (h *handler) send() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var signals aichat.AIChatSignals
        if err := aichat.ReadSignals("my-chat", r, &signals); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        input := strings.TrimSpace(signals.Input)
        if input == "" {
            return
        }

        sse := datastar.NewSSE(w, r)
        chat := aichat.Chat(sse, "my-chat")

        _ = chat.UserMessage(input)
        _ = chat.ShowTyping()

        // ... call AI, process input ...

        _ = chat.HideTyping()
        _ = chat.AssistantText("Here's what I found.")
    }
}
```

### Multi-turn with quick replies and cards

```go
chat := aichat.Chat(sse, chatID)

_ = chat.UserMessage(input)
_ = chat.ShowTyping()
time.Sleep(800 * time.Millisecond)
_ = chat.HideTyping()

// Simple text
_ = chat.AssistantText("You have 8 subscriptions.")

// Quick replies
_ = chat.Append(aichat.QuickReplies(aichat.QuickRepliesProps{
    ChatID:    chatID,
    SubmitURL: "/api/chat/send",
    Options: []aichat.QuickReply{
        {Label: "Netflix", Value: "Netflix"},
        {Label: "Show all", Value: "Show all"},
    },
}))

// Response card (use AssistantMessage to wrap with avatar)
_ = chat.AssistantMessage(aichat.ResponseCard(aichat.ResponseCardProps{
    Title:       "Task created",
    Description: "I'll handle this for you.",
}))
```

---

## Reading Signals

```go
func ReadSignals(chatID string, r *http.Request, dest *AIChatSignals) error
```

Reads the chat's namespaced signals from the request. The `chatID` is sanitized (hyphens → underscores) to match the JavaScript signal namespace.

```go
type AIChatSignals struct {
    Open  bool   `json:"open"`
    Input string `json:"input"`
}
```

---

## SSE Targeting

Messages are appended to `#<chatID>-messages`. Use `aichat.MessagesID(chatID)` to get this ID:

```go
messagesID := aichat.MessagesID("my-chat") // "my-chat-messages"
```

The `ChatSender` handles this automatically — you rarely need `MessagesID` directly.

---

## CommandBar Embedded Mode

When using a `commandbar.CommandBar` as the `InputSlot`, set `Embedded: true`:

```go
commandbar.Props{
    Embedded: true, // no collapsed state, no close, no outer chrome
    // ...
}
```

**Embedded mode changes:**
- No collapsed state — starts in text mode
- No × close button (the AI chat provides its own)
- No outer border/shadow (parent provides chrome)
- On submit: clears text but stays in current mode (doesn't collapse)

**Signal reading:** When the command bar is embedded, input comes from the command bar's `text` signal (not the AI chat's `input` signal). Use `datastar.ReadSignals` with the command bar's ID namespace.

---

## Package Structure

```
ui/aichat/
  aichat.templ    — AIChat main container (collapsed/expanded, input, messages area)
  messages.templ  — UserMessage, AssistantMessage, AssistantText, Typing, QuickReplies
  card.templ      — ResponseCard, AssignRow, ActionBar
  chat.go         — ChatSender (server-side SSE helper)
  helpers.go      — quickReplyExpr, ReadSignals
```
