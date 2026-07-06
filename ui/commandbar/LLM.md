# commandbar — Multi-Mode Command Input

> **Import**: `"github.com/kinnsei/dsx/ui/commandbar"`

An expandable multi-mode command input supporting text, file upload, and voice recording. Collapsed by default as a compact bar, expanding to a tabbed interface on interaction.

---

## Component Overview

| Component | Description |
|---|---|
| `CommandBar` | Main component — collapsed bar, expandable tabbed input |

---

## Props

```go
type Props struct {
    ID          string
    Class       string
    Attributes  templ.Attributes
    Placeholder string   // collapsed placeholder (default "Type something...")
    SubmitURL   string   // POST endpoint for text submit
    UploadURL   string   // POST endpoint for file upload (enables file tab)
    VoiceURL    string   // POST endpoint for voice (enables voice tab)
    Suggestions []string // quick-fill suggestion chips in text mode
    FileHint    string   // subtitle in file drop zone (default "Drop a file or click to browse")
    Embedded    bool     // when true: no collapsed state, no close behavior, no outer chrome
}
```

---

## Modes

Modes are enabled by providing the corresponding URL prop:

| Mode | Enabled by | Behavior |
|---|---|---|
| Text | `SubmitURL` (always present) | Textarea + Send button, Enter submits |
| File | `UploadURL` | Drag-and-drop zone, camera, browse buttons |
| Voice | `VoiceURL` | Mic button toggles recording UI (visual mock) |

---

## Basic Usage

### Text only

```templ
@commandbar.CommandBar(commandbar.Props{
    ID:        "my-bar",
    SubmitURL: "/api/command/send",
})
```

### With suggestions

```templ
@commandbar.CommandBar(commandbar.Props{
    ID:          "my-bar",
    Placeholder: "How can I help?",
    SubmitURL:   "/api/command/send",
    Suggestions: []string{"Check balance", "Transfer funds", "Report issue"},
})
```

### All modes (text + file + voice)

```templ
@commandbar.CommandBar(commandbar.Props{
    ID:          "my-bar",
    Placeholder: "Type, upload, or record...",
    SubmitURL:   "/api/command/send",
    UploadURL:   "/api/command/upload",
    VoiceURL:    "/api/command/voice",
    Suggestions: []string{"Open a case", "Upload document"},
})
```

---

## Embedded Mode

When `Embedded: true`, the command bar is designed to be composed inside another component (e.g. `aichat.AIChat`):

- No collapsed state — starts in text mode
- No close button (parent provides its own)
- No outer border/shadow (parent provides chrome)
- On submit: clears text but stays in current mode (doesn't collapse)

```templ
@aichat.AIChat(aichat.Props{
    ID: "my-chat",
    InputSlot: commandbar.CommandBar(commandbar.Props{
        ID:        "my-bar",
        Embedded:  true,
        SubmitURL: "/api/chat/send",
        UploadURL: "/api/chat/upload",
        VoiceURL:  "/api/chat/voice",
    }),
})
```

---

## Signals

```go
type CommandBarSignals struct {
    Mode      string `json:"mode"`      // "" | "text" | "file" | "voice"
    Text      string `json:"text"`      // text input value
    Recording bool   `json:"recording"` // voice recording active
}
```

Signals are namespaced by component ID (hyphens → underscores). A command bar with `ID: "my-bar"` creates signals at `my_bar.mode`, `my_bar.text`, `my_bar.recording`.

---

## Reading Signals in Handlers

Use `ds.ReadSignals()` to read namespaced signals from the request. Pass the component ID (hyphens are automatically converted to underscores).

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Read signals BEFORE creating SSE (SSE consumes the request body).
    var signals commandbar.CommandBarSignals
    if err := ds.ReadSignals("my-bar", r, &signals); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // signals.Text contains the user's text input
    // signals.Mode is the active mode ("text", "file", "voice")

    sse := datastar.NewSSE(w, r)
    ds.Send.Toast(sse, ds.ToastSuccess, fmt.Sprintf("Received: %q", signals.Text))
}
```

**Important**: Call `ds.ReadSignals()` **before** `datastar.NewSSE()` — creating the SSE consumes the request body.

---

## Multiple Instances on One Page

When multiple command bars exist on the same page, **all instances send their signals** in every POST (Datastar sends all signals). To identify which instance triggered the request:

1. Each instance has a unique ID → unique signal namespace
2. Check which namespace has active content (non-empty text or non-idle mode)

```go
// Read raw JSON to inspect all namespaces
var raw map[string]json.RawMessage
if err := datastar.ReadSignals(r, &raw); err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
}

for _, id := range knownIDs {
    data, ok := raw[id]
    if !ok {
        continue
    }

    var signals commandbar.CommandBarSignals
    if err := json.Unmarshal(data, &signals); err != nil {
        continue
    }

    // Skip instances in idle state
    if signals.Text == "" && (signals.Mode == "" || signals.Mode == "text") {
        continue
    }

    // This is the active instance
    fmt.Println("Input from", id, ":", signals.Text)
    break
}
```

---

## Behavior Details

### Standalone mode (default)
- **Collapsed**: Placeholder bar with shortcut buttons (file, voice). Click to expand.
- **Expanded**: Tabbed interface (Type / File / Voice) with mode-specific content.
- **Submit**: Enter or Send button POSTs to `SubmitURL`, then collapses (resets mode + text).
- **Escape**: Collapses (same as close button).

### Embedded mode (`Embedded: true`)
- **No collapsed state**: Starts directly in text mode.
- **Submit**: POSTs to `SubmitURL`, clears text but stays in current mode.
- **No close button**: Parent component handles closing.

### File mode
- Drag-and-drop zone or browse button triggers `data-on:change` → POSTs to `UploadURL`.
- Camera button available for mobile (uses `capture="environment"`).
- File data is **not** sent via signals — handle file upload server-side via the POST body or use a separate upload mechanism.

### Voice mode
- **Visual mock only** — no actual audio recording.
- Mic button toggles the `recording` signal (true/false).
- Stop button POSTs to `VoiceURL` with current signals.
- Real voice recording requires browser `MediaRecorder` API integration.

---

## Package Structure

```
ui/commandbar/
  commandbar.templ  — CommandBar component (collapsed/expanded, tabs, modes)
```
