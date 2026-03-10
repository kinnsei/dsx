# ADR-0002: ds — Type-Safe Datastar Helpers

## Status

Accepted

## Context

WebX uses [Datastar](https://data-star.dev) for all frontend interactivity. Datastar is powerful but its HTML attribute API has several sharp edges that cause silent failures, security gaps, and repetitive boilerplate in Go/templ code.

### Problem 1: Silent Attribute Syntax Failures

Datastar uses **colon-separated** plugin syntax for parameterized attributes:

```html
<!-- Correct — works -->
<button data-on:click="@get('/api/items')">Load</button>

<!-- Wrong — silently ignored, no error in console -->
<button data-on-click="@get('/api/items')">Load</button>
```

The hyphen variant looks natural to anyone familiar with `data-*` attributes or Alpine.js, and Datastar silently ignores it. This is the single most common source of "it doesn't work" debugging sessions.

### Problem 2: CSRF Token Omission

Mutating operations (POST, PUT, PATCH, DELETE) must include a CSRF token. Forgetting it produces a 403 that's easy to misdiagnose. Developers must manually construct the header injection every time:

```html
<button data-on:click="@post('/api/create', {headers: {'X-CSRF-Token': document.querySelector('meta[name=csrf-token]').content}})">
```

This is verbose, error-prone, and easy to forget on one of twenty endpoints.

### Problem 3: Signal Namespace Mismatch

Datastar namespaces signals by component ID. But HTML allows hyphens in IDs (`my-form`) while JavaScript property access doesn't — Datastar silently converts hyphens to underscores internally. This creates a mismatch:

```html
<!-- Frontend: component ID has hyphens -->
<div data-signals-my-form="{name: ''}">
  <input data-bind:my-form.name />
</div>
```

```go
// Backend: must use underscores to read the signal
type Signals struct {
    Name string `json:"name"`
}
// What namespace? "my-form"? "my_form"?
```

Getting this wrong means the backend silently reads empty values.

### Problem 4: Repetitive SSE Response Patterns

Common UI operations — showing modals, drawers, toasts, confirmation dialogs — require constructing HTML fragments, escaping them into JavaScript expressions, and pushing them via SSE. Each pattern involves 10-20 lines of boilerplate that's identical across every handler.

### Problem 5: Retry and Cancellation Defaults

Datastar retries failed requests by default. This is useful for SSE streams but harmful for one-shot mutations (a failed POST shouldn't retry and create duplicate records). Developers must remember to set `retryMaxCount: 0` on every non-streaming request.

## Decision

The `ds` package provides two layers of abstraction over Datastar:

1. **Frontend helpers** — Go functions that return `templ.Attributes` or expression strings with correct syntax guaranteed by construction
2. **Backend helpers** (`ds.Send`) — SSE response builders for common UI patterns (modals, drawers, toasts, patches)

### Design Principles

- **Impossible to write wrong syntax**: Functions generate `data-on:click` (colon), never `data-on-click` (hyphen)
- **Security by default**: Mutating verbs automatically include CSRF; read-only verbs never do
- **Transparent sanitization**: Hyphens in component IDs are converted to underscores on both frontend and backend
- **Zero-overhead defaults**: `GetOnce` / `PostOnce` disable retries; `Get` with `WithRequestCancellation("disabled")` is available for streams
- **Composable**: All helpers return strings or `templ.Attributes` — they compose with each other and with raw Datastar attributes

### Frontend Attribute Helpers

These return `templ.Attributes` and are spread onto elements in templ:

```go
// Event handlers
<button { ds.OnClick(ds.PostOnce("/api/create"))... }>Save</button>
<input { ds.On("keydown.enter", ds.PostOnce("/api/search"))... } />

// Two-way binding — automatically handles hyphen→underscore
<input { ds.Bind("my-form", "email")... } />
// Produces: data-bind:my_form.email

// Visibility, text, classes
<div { ds.Show("$drawer.open")... }>...</div>
<span { ds.Text("$counter.value")... }></span>
<li { ds.ClassToggle("active", "$item.selected")... }>...</li>

// Signals and lifecycle
<div { ds.RawSignals(`{"count": 0}`)... }></div>
<div { ds.Init(ds.GetOnce("/api/load"))... }></div>
<div { ds.Effect("console.log($count)")... }></div>
```

#### Action Expressions

Functions that return Datastar action strings (`@get(...)`, `@post(...)`, etc.):

| Function | CSRF | Retries | Use Case |
|----------|------|---------|----------|
| `Get(url)` | No | Yes | SSE streams, polling |
| `GetOnce(url)` | No | No | One-shot data loads |
| `Post(url)` | Yes | Yes | Streaming mutations |
| `PostOnce(url)` | Yes | No | Form submissions |
| `Put(url)` / `PutOnce(url)` | Yes | Yes / No | Updates |
| `Patch(url)` / `PatchOnce(url)` | Yes | Yes / No | Partial updates |
| `Delete(url)` / `DeleteOnce(url)` | Yes | Yes / No | Deletions |

CSRF injection reads from `<meta name="csrf-token">` in the page head, which the layout template sets from the WebX context.

#### Action Options

Fine-grained control via option functions:

```go
// Persistent SSE connection (no retry, no cancellation)
ds.Get("/stream", ds.WithRequestCancellation("disabled"))

// File upload
ds.Post("/api/upload", ds.WithContentType("form"))

// Only send this component's signals
ds.Post("/api/submit", ds.WithFilterSignals("my-form"))

// Custom retry count
ds.Get("/api/poll", ds.WithRetries(5))
```

### Signal Management

The `SignalManager` provides type-safe signal references:

```go
signals := ds.NewSignals("tabs", struct {
    Active string `json:"active"`
}{Active: "overview"})

// Templ usage
<div { signals.Attrs()... }>
    <button { ds.OnClick(signals.SetString("active", "overview"))... }>Overview</button>
    <button { ds.OnClick(signals.SetString("active", "settings"))... }>Settings</button>
    <div { ds.Show(signals.Equals("active", "'overview'"))... }>Overview content</div>
    <div { ds.Show(signals.Equals("active", "'settings'"))... }>Settings content</div>
</div>
```

Methods:

| Method | Output | Example |
|--------|--------|---------|
| `Signal(prop)` | `$id.prop` | `$tabs.active` |
| `Toggle(prop)` | `$id.prop = !$id.prop` | Toggle boolean |
| `Set(prop, val)` | `$id.prop = val` | `$tabs.active = 'settings'` |
| `SetString(prop, val)` | `$id.prop = 'val'` | String with quotes |
| `Equals(prop, val)` | `$id.prop === val` | Conditional check |
| `NotEquals(prop, val)` | `$id.prop !== val` | Inverse check |
| `Conditional(prop, t, f)` | `$id.prop ? t : f` | Ternary |

### Expression Builder

For multi-statement actions:

```go
expr := ds.NewExpression().
    SetSignal("$drawer.open", "false").
    Statement("@post('/api/save')").
    Build()
// → "$drawer.open = false; @post('/api/save')"
```

### Backend SSE Helpers (`ds.Send`)

All methods take a `*datastar.ServerSentEventGenerator` as first argument:

#### Patch

Renders a templ component and morphs it into the DOM by matching element IDs:

```go
func (h *handler) list(w http.ResponseWriter, r *http.Request) {
    sse := datastar.NewSSE(w, r)
    ds.Send.Patch(sse, CustomerTable(rows))
}
```

#### Toast

Shows a notification with level, message, and options:

```go
ds.Send.Toast(sse, ds.ToastSuccess, "Customer saved")
ds.Send.Toast(sse, ds.ToastError, "Validation failed",
    ds.WithToastPersistent(),
    ds.WithToastAction("Retry", "/api/retry"),
)
```

Levels: `ToastInfo`, `ToastSuccess`, `ToastWarning`, `ToastError`.

Options: `WithToastDuration(ms)`, `WithToastPersistent()`, `WithToastAction(label, url)`, `WithToastLink(text, url)`.

#### Modal

Shows a centered dialog with a templ component as content:

```go
ds.Send.Modal(sse, ConfirmDeleteDialog(id), ds.WithModalMaxWidth("max-w-sm"))
```

`ds.Send.HideModal(sse)` closes it.

#### Drawer

Shows a slide-in panel from the right:

```go
ds.Send.Drawer(sse, EditCustomerForm(customer), ds.WithDrawerMaxWidth("max-w-md"))
```

`ds.Send.HideDrawer(sse)` closes it.

Important: Drawer and Modal content renders with `context.Background()`, not the request context. Components inside them cannot use `dsx.FromContext(ctx)` — pass data as function parameters instead.

#### Confirm

Shows a confirmation dialog with customizable labels and action:

```go
ds.Send.Confirm(sse, "Delete this customer?", "/api/customers/42",
    ds.WithConfirmTitle("Delete Customer"),
    ds.WithConfirmLabel("Delete"),
    ds.WithConfirmClass("btn btn-error"),
)
```

#### Redirect and Download

```go
ds.Send.Redirect(sse, "/dashboard")           // Navigate browser
ds.Send.Download(sse, "/api/export.csv", "report.csv")  // Trigger download
```

### Signal Reading (Backend)

Read signals from incoming requests with automatic namespace handling:

```go
type FormSignals struct {
    Name    string `json:"name"`
    Email   string `json:"email"`
}

func handler(w http.ResponseWriter, r *http.Request) {
    var signals FormSignals
    ds.ReadSignals("my-form", r, &signals)
    // Reads from the "my_form" namespace (hyphen→underscore automatic)
}
```

`ReadRaw(r, &dest)` reads all namespaces as raw JSON for advanced use cases.

### Attribute Composition

Multiple attribute helpers compose via `ds.Merge`:

```go
<button { ds.Merge(
    ds.OnClick(ds.PostOnce("/api/save")),
    ds.ClassToggle("loading", "$form.submitting"),
    ds.Indicator("form.submitting"),
)... }>
    Save
</button>
```

`Merge` combines multiple `templ.Attributes` maps into one, with later values overriding earlier ones for the same key.

## Consequences

### Positive

- **Zero silent failures**: Attribute syntax is correct by construction — `data-on:click` is the only possible output
- **CSRF handled automatically**: Developers cannot forget it on mutations, and it's never sent on reads
- **Namespace transparency**: Hyphens in component IDs work seamlessly across frontend and backend
- **Consistent retry behavior**: `Once` variants prevent accidental duplicate mutations; stream variants keep retrying
- **Reduced boilerplate**: Modal/drawer/toast patterns go from 15+ lines to a single function call
- **Discoverable API**: IDE autocomplete on `ds.` surfaces all available Datastar operations
- **Composable with raw Datastar**: Helpers return standard `templ.Attributes` — they can be mixed with hand-written attributes

### Negative

- **Abstraction layer**: Developers need to learn the ds API in addition to (or instead of) raw Datastar
- **Drawer/Modal context limitation**: Content renders with `context.Background()`, requiring data to be passed as parameters rather than read from request context
- **String-based expressions**: Signal references like `$tabs.active` are still strings — no compile-time validation of signal names

### Trade-offs

- **Convenience vs. control**: The package covers ~95% of use cases. For the remaining 5% (custom Datastar plugins, non-standard attribute patterns), developers can use raw `templ.Attributes` alongside ds helpers
- **Implicit CSRF vs. explicit**: Automatic CSRF injection is safer but less visible. Developers should understand that `PostOnce` includes CSRF while `GetOnce` does not
- **Naming convention**: `ds` is short and convenient but could conflict with other packages. The brevity was chosen because it's used on nearly every line of templ code
