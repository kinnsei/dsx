# ds — Datastar Helpers for WebX

> **Import**: `"github.com/plaenen/webx/ds"`

The `ds` package is split into two namespaces:

- **`ds.XXX`** — Frontend attribute helpers. Return `templ.Attributes` or expression strings for use in `.templ` files.
- **`ds.Send.XXX`** — Backend SSE operations. Called from Go HTTP handlers to send events to the browser.

---

## Frontend Attributes (`ds.XXX`)

These helpers generate Datastar `data-*` HTML attributes. They return `templ.Attributes` and are spread onto elements in templ files using `{ ds.XXX(...)... }`.

### Critical: Colon syntax

Datastar uses **colons** for parameterized attributes (`data-on:click`), NOT hyphens (`data-on-click`). A hyphen is silently ignored. The `ds` helpers make this mistake impossible.

### Parameterized attributes

Return `templ.Attributes` — spread with `{ ... }`:

```go
ds.On("click", expr)       // data-on:click="expr"
ds.OnClick(expr)            // shorthand for ds.On("click", expr)
ds.Bind("name")             // data-bind:name=""
ds.ClassToggle("bold", "$x")// data-class:bold="$x"
ds.Attr("disabled", "$x")   // data-attr:disabled="$x"
ds.Style("color", "$x")     // data-style:color="$x"
ds.Computed("total", "$a*$b")// data-computed:total="$a*$b"
ds.Indicator("loading")      // data-indicator:loading=""
ds.Ref("myEl")              // data-ref:myEl=""
```

### Standalone attributes

```go
ds.Signals("{count: 0}")    // data-signals="{count: 0}"
ds.Show("$visible")         // data-show="$visible"
ds.Text("$name")            // data-text="$name"
ds.Class("{'bold': $x}")    // data-class="{'bold': $x}"
ds.Init(expr)               // data-init="expr"
ds.Effect(expr)             // data-effect="expr"
```

### Backend action expressions

Return **strings** (not `templ.Attributes`). Use inside `ds.On(...)`, `ds.OnClick(...)`, or `ds.Init(...)`:

```go
ds.Get("/api/data")                     // @get('/api/data')
ds.GetOnce("/api/data")                 // @get('/api/data', {retryMaxCount: 0})
ds.Post("/api/submit")                  // @post('/api/submit', {headers: {'X-CSRF-Token': ...}})
ds.PostOnce("/api/submit")              // same with retryMaxCount: 0
ds.Put("/api/item")                     // @put with CSRF
ds.Patch("/api/item")                   // @patch with CSRF
ds.Delete("/api/item")                  // @delete with CSRF
```

- `Get` does NOT include CSRF. `Post`, `Put`, `Patch`, `Delete` auto-include the CSRF token from `<meta name="csrf-token">`.
- `*Once` variants disable retries (single-shot requests for button clicks, form submits).
- Use `Get`/`Post` (with retries) for SSE streams and long-polling.

**Options:**

```go
ds.Get(url, ds.WithRetries(3))                // custom retry count
ds.Get(url, ds.WithRequestCancellation("disabled")) // persistent SSE streams
ds.Post(url, ds.WithContentType("form"))      // multipart/form-data (file uploads)
```

### Combining attributes

```go
ds.Merge(ds.OnClick(expr), ds.Show("$x"))  // combine multiple templ.Attributes
```

### Templ usage patterns

**Button triggering an SSE action:**
```templ
<button
    class="btn btn-primary"
    { ds.OnClick(ds.GetOnce(wctx.APIPath("/api/items/load")))... }
>
    Load Items
</button>
```

**Element with initial data load:**
```templ
<div
    id="user-profile"
    { ds.Init(ds.GetOnce(wctx.APIPath("/api/profile")))... }
>
    Loading...
</div>
```

**Reactive show/hide:**
```templ
<div { ds.Show("$isOpen")... } style="display: none;">
    Content shown when $isOpen is true
</div>
```

**Two-way binding:**
```templ
<input type="text" { ds.Bind("search")... } placeholder="Search..." />
```

**Combining multiple attributes on one element:**
```templ
<button
    class="btn"
    { ds.Merge(
        ds.OnClick(ds.PostOnce(wctx.APIPath("/api/save"))),
        ds.Indicator("saving"),
        ds.Attr("disabled", "$saving"),
    )... }
>
    Save
</button>
```

---

## Backend SSE Operations (`ds.Send.XXX`)

Called from Go HTTP handlers. These create SSE events that update the browser DOM.

### Setup requirement

The base layout (`layouts/base.templ`) must include container elements:

```html
<body>
    { children... }
    <div id="drawer-panel"></div>        <!-- drawer container -->
    <div id="toast-container" class="toast toast-end toast-bottom z-50"></div>  <!-- toast container -->
</body>
```

### Drawer

Opens a slide-in panel from the right side. Content is a templ component rendered server-side.

```go
func (s *Sender) Drawer(sse *datastar.ServerSentEventGenerator, content templ.Component, opts ...DrawerOption) error
func (s *Sender) HideDrawer(sse *datastar.ServerSentEventGenerator) error
```

**Handler example:**

```go
func (h *handler) showDetail() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id := chi.URLParam(r, "id")
        item := loadItem(id)

        sse := datastar.NewSSE(w, r)
        ds.Send.Drawer(sse, itemDetail(item))
    }
}
```

**With options:**

```go
// Wider panel
ds.Send.Drawer(sse, content, ds.WithDrawerMaxWidth("max-w-2xl"))
```

**Server-initiated close** (e.g. after form submit inside drawer):

```go
func (h *handler) saveAndClose() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // ... save data ...
        sse := datastar.NewSSE(w, r)
        ds.Send.HideDrawer(sse)
        toast.Show(sse, toast.LevelSuccess, "Saved!", 3000)
    }
}
```

**How close works:** The close button and overlay use an inline Datastar expression (`document.getElementById('drawer-panel').innerHTML=''`). No close endpoint is needed for user-initiated close. `HideDrawer` is only needed when the server wants to close the drawer (e.g. after a successful form submission).

**Templ side — triggering the drawer:**

```templ
<button
    class="btn btn-primary"
    { ds.OnClick(ds.GetOnce(wctx.APIPath("/api/items/" + id)))... }
>
    View Details
</button>
```

### Toast (current API — will move to `ds.Send.Toast` in future)

Toasts currently live in `ui/toast/handler.go` as `toast.Show(...)`. The API will move to `ds.Send.Toast(...)` in a future release. Current usage:

```go
import "github.com/plaenen/webx/ui/toast"

// Auto-dismiss (3 seconds)
toast.Show(sse, toast.LevelSuccess, "Saved!", 3000)

// Persistent (stays until user closes)
toast.ShowPersistent(sse, toast.LevelWarning, "Check your input.")

// With action button
toast.ShowAction(sse, toast.LevelError, "Deleted.", "Undo", wctx.APIPath("/api/undo"))

// With link
toast.ShowLink(sse, toast.LevelInfo, "New version available.", "Update", "/changelog", 5000)

// Custom templ component
toast.ShowComponent(sse, myCustomToast(data))
```

**Levels:** `toast.LevelInfo`, `toast.LevelSuccess`, `toast.LevelWarning`, `toast.LevelError`

---

## Best Practices

### 1. Always use `wctx.APIPath()` for URLs

The app may be mounted at a base path (e.g. `/showcase`). Always construct URLs via `wctx.APIPath()`:

```go
// Correct
ds.GetOnce(wctx.APIPath("/api/items/1"))
// Produces: @get('/showcase/api/items/1', {retryMaxCount: 0})

// Wrong — will break if basePath is set
ds.GetOnce("/api/items/1")
```

### 2. Use `*Once` for button clicks, `Get`/`Post` for streams

- `ds.GetOnce()` / `ds.PostOnce()` — single-shot requests (buttons, form submits). Retries disabled.
- `ds.Get()` / `ds.Post()` — SSE streams, long-polling. Retries enabled by default.

### 3. Mutating actions need CSRF

`Post`, `Put`, `Patch`, `Delete` auto-include CSRF. `Get` does not. Never use `Get` for mutations — use the appropriate HTTP method.

### 4. Every PatchElements target needs an ID

When calling `sse.PatchElements(html)`, the root element in the HTML **must have an `id`**. Datastar uses the ID to find the target element in the DOM. Without an ID, you get `PatchElementsNoTargetsFound`.

```go
// Correct
sse.PatchElements(`<div id="my-element">Updated content</div>`)

// Wrong — no target
sse.PatchElements(`<div>Updated content</div>`)
```

### 5. Drawer content is just a templ component

The `ds.Send.Drawer` helper wraps your component in the drawer shell (overlay, panel, close button). Your templ component should only contain the **content** — no wrapper divs, no close buttons:

```templ
templ itemDetail(item Item) {
    <h2 class="text-2xl font-bold">{ item.Name }</h2>
    <p class="text-base-content/70">{ item.Description }</p>
    // ... more content ...
}
```

### 6. Combine Send operations in one handler

A single SSE response can contain multiple events. Show a drawer and a toast in the same handler:

```go
sse := datastar.NewSSE(w, r)
ds.Send.Drawer(sse, editForm(item))
toast.Show(sse, toast.LevelInfo, "Editing item", 2000)
```

### 7. Drawer close + follow-up action

After a form submission inside a drawer, close the drawer and show feedback:

```go
sse := datastar.NewSSE(w, r)
ds.Send.HideDrawer(sse)
toast.Show(sse, toast.LevelSuccess, "Changes saved!", 3000)
// Optionally patch another element to refresh the list
sse.PatchElements(`<div id="item-list">...updated list...</div>`)
```

### 8. Don't patch elements that have Datastar attributes you need to keep

Datastar's morph preserves attributes on morphed elements, but if you replace an element entirely (e.g. via PatchElements), the replacement HTML must include any `data-*` attributes the element needs. If a `data-effect` or `data-on:click` disappears after a patch, you've likely replaced the element that had it. Solution: patch an inner child element instead.

---

## Package Structure

```
ds/
  ds.go             — Frontend attribute helpers (ds.OnClick, ds.Show, ds.Get, etc.)
  send.go           — Sender type + Send var
  send_drawer.go    — ds.Send.Drawer, ds.Send.HideDrawer
  send_toast.go     — (future) ds.Send.Toast, ds.Send.ToastPersistent, etc.
```

## Container IDs

| Container | ID | Purpose |
|---|---|---|
| Drawer | `drawer-panel` | Slide-in detail panel |
| Toast | `toast-container` | Stacked notifications |

Both containers are in the base layout and must exist in the DOM before any `Send` operations target them.
