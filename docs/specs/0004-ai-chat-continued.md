# Spec: dsx — AIChat Full-Page Mode

> **Library**: `github.com/laenen-partners/dsx`
> **Package**: `ui/aichat`
> **Extends**: Existing `aichat.AIChat` component

---

## Context

The `aichat.AIChat` component is a collapsible widget: sparkle icon → click to expand → `max-h-96` message area → click to collapse. This works for floating chat widgets on dashboards.

But apps like steward need a **full-page chat** that fills a container (e.g. a drawer content area). The message list should grow to fill available height, with the compose bar pinned at the bottom. No collapse/expand behaviour.

The `ChatSender` server-side API is already mode-agnostic — it works with any layout. Only the template needs a variant.

## Proposed change

Add a `Mode` prop to `aichat.Props`:

```go
type Mode string

const (
    ModeWidget   Mode = ""         // default: collapsible widget with max-h-96
    ModeFullPage Mode = "fullpage" // fills container, no collapse, compose pinned at bottom
)

type Props struct {
    ID               string
    Class            string
    Attributes       templ.Attributes
    Mode             Mode            // NEW: widget (default) or fullpage
    Placeholder      string
    ReplyPlaceholder string
    SubmitURL        string
    Shortcut         string
    InputSlot        templ.Component
}
```

### Widget mode (default, unchanged)

```
┌─────────────────────────────────┐
│ ✦ Tell your butler something... │  ← collapsed
└─────────────────────────────────┘

       ↓ click ↓

┌─────────────────────────────────┐
│                            [×]  │
│  messages (max-h-96, scroll)    │
│  ...                            │
│  ...                            │
├─────────────────────────────────┤
│  [input] Type a reply...        │
└─────────────────────────────────┘
```

### Full-page mode (new)

```
┌─────────────────────────────────┐
│                                 │
│  messages (flex-1, fills space) │
│  ...                            │
│  ...                            │
│  ...                            │
│  ...                            │
│                                 │
│                                 │
├─────────────────────────────────┤
│  [input] Type a message...  [➤] │
└─────────────────────────────────┘
```

No collapsed state. No sparkle icon. No close button. No `max-h-96`. The component fills its parent container using flexbox.

## Template changes

```templ
templ AIChat(props Props) {
    if props.Mode == ModeFullPage {
        @fullPageChat(props)
    } else {
        @widgetChat(props)  // existing implementation
    }
}

templ fullPageChat(props Props) {
    {{ /* same signal setup as widget */ }}
    <div
        id={ id }
        data-signals={ signals.DataSignals }
        class={ utils.TwMerge("flex flex-col h-full", props.Class) }
        { props.Attributes... }
    >
        <!-- Messages area (grows to fill) -->
        <div
            id={ messagesID }
            class="flex-1 overflow-y-auto flex flex-col gap-3 px-4 py-4"
        >
            { children... }
        </div>

        <!-- AI response target (between messages and input) -->
        <div id={ id + "-ai-response" } class="px-4"></div>

        <!-- Compose bar (pinned at bottom) -->
        if props.InputSlot != nil {
            <div class="border-t border-base-200">
                @props.InputSlot
            </div>
        } else {
            <div class="flex items-center gap-2 border-t border-base-200 px-4 py-3">
                <input
                    type="text"
                    class="input input-sm input-bordered flex-1"
                    placeholder={ replyPlaceholder }
                    { ds.Bind(id, "input")... }
                    data-on:keydown={ /* Enter to submit */ }
                />
                <button class="btn btn-sm btn-primary btn-circle"
                    { ds.OnClick(ds.Post(props.SubmitURL))... }
                >
                    <!-- send icon -->
                </button>
            </div>
        }
    </div>
}
```

## ChatSender changes

None. `ChatSender` targets `#project-chat-messages` (or whatever MessagesID returns). It doesn't care about the layout mode — it just appends/patches fragments.

The only new addition:

```go
// AIResponseID returns the ID of the AI response target for SSE patching.
func AIResponseID(chatID string) string {
    return chatID + "-ai-response"
}
```

## Usage in steward

```templ
// Project detail page — chat fills the drawer content area
@aichat.AIChat(aichat.Props{
    ID:        "project-chat",
    Mode:      aichat.ModeFullPage,
    SubmitURL: "/projects/" + props.Project.ID + "/messages",
    ReplyPlaceholder: "Type a message...",
}) {
    // Initial messages rendered server-side
    for _, m := range props.Messages {
        @messageBubble(m)
    }
}

// SSE connection for AI streaming (separate, on page load)
<div style="display:none"
    { ds.Init(ds.GetOnce("/projects/" + props.Project.ID + "/ai-stream"))... }
></div>
```

Steward replaces ~40 lines of custom layout with one `@aichat.AIChat` call.

## Acceptance criteria

- [ ] `ModeFullPage` fills parent container (no `max-h-96`, no collapse)
- [ ] Messages area uses `flex-1` to grow
- [ ] Compose bar pinned at bottom with `border-t`
- [ ] `#ai-response` div between messages and compose for SSE streaming
- [ ] Enter to submit, input clears after send
- [ ] `ChatSender` works identically in both modes
- [ ] Widget mode (default) unchanged
- [ ] `InputSlot` works in full-page mode (custom compose bar)
- [ ] Showcase demonstrates both modes
