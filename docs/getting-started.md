# Getting Started with WebX

> **Module**: `github.com/plaenen/webx`
> **Go Requirement**: Go 1.24+
> **Stack**: Go + Chi (routing) + Templ (templating) + Tailwind CSS + DaisyUI (styling) + Datastar (frontend interactivity)

---

## Installation

```bash
go get github.com/plaenen/webx
```

Install dependencies with the task runner:

```bash
go tool task install:all    # mod tidy + download DaisyUI
go tool task install:daisyui # DaisyUI plugin files only
```

Generate templ code after editing any `.templ` file:

```bash
go tool templ generate
```

Build Tailwind CSS:

```bash
go tool gotailwind
```

---

## Project Structure

```
cmd/            -- application entry points
internal/       -- internal packages
ui/             -- DaisyUI components (one dir per component, e.g. ui/button/)
ds/             -- Datastar helpers (frontend attributes + backend SSE helpers)
layouts/        -- base HTML layout with toast/modal/drawer containers
stream/         -- reactive streaming via NATS
utils/          -- shared templ utilities (TwMerge, If, RandomID, etc.)
static/css/     -- Tailwind CSS + DaisyUI plugin files
docs/           -- reference documentation
```

---

## Minimal Server Setup

```go
package main

import (
    "fmt"
    "net/http"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/nats-io/nats-server/v2/server"
    "github.com/nats-io/nats.go"
    "github.com/plaenen/webx"
    "github.com/plaenen/webx/stream"
    "github.com/plaenen/webx/ui"
)

func main() {
    // 1. Start embedded NATS (required for reactive streaming)
    ns, _ := server.NewServer(&server.Options{DontListen: true})
    ns.Start()
    defer ns.Shutdown()
    ns.ReadyForConnections(4 * time.Second)

    nc, _ := nats.Connect(ns.ClientURL(), nats.InProcessServer(ns))
    defer nc.Close()

    broker := stream.NewBroker(nc)

    // 2. Create router with middleware
    r := chi.NewRouter()

    store := NewMySessionStore() // implements webx.SessionStore
    r.Use(webx.SessionMiddleware(store))
    r.Use(webx.SecurityHeadersMiddleware())

    // 3. Configure WebX context
    const basePath = "/app"
    r.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            wctx := webx.FromContext(r.Context())
            wctx.BasePath = basePath
            wctx.StreamURL = basePath + "/stream"
            wctx.Stylesheets = []webx.Stylesheet{{Href: "/assets/css/output.css"}}
            wctx.Scripts = []webx.Script{{Src: "/assets/js/datastar.js"}}
            next.ServeHTTP(w, r.WithContext(wctx.WithContext(r.Context())))
        })
    })

    // 4. Register pages
    r.Get("/", templ.Handler(pages.Home()).ServeHTTP)

    // 5. Register SSE handlers
    r.Route(basePath, func(r chi.Router) {
        r.Get("/stream", broker.Handler())
        ui.RegisterRoutes(r,
            ui.WithMarkdownPreview(),
            ui.WithDecimalParser(),
            ui.WithMoneyParser(),
        )
    })

    http.ListenAndServe(":3000", r)
}
```

---

## WebX Context

Every request carries a `WebXContext` set up by `SessionMiddleware`. Access it in handlers and templ components:

```go
wctx := webx.FromContext(ctx)
```

| Field | Type | Description |
|---|---|---|
| `CSRFToken` | `string` | Auto-generated CSRF token (injected in base layout) |
| `SessionID` | `string` | Session cookie value |
| `BasePath` | `string` | Prefix for SSE handler routes (e.g. `"/app"`) |
| `StreamURL` | `string` | URL for the reactive SSE stream endpoint |
| `Theme` | `string` | DaisyUI theme name (persisted in session) |
| `Store` | `SessionStore` | Session store for reading/writing session data |
| `DevMode` | `bool` | Development mode flag |
| `Stylesheets` | `[]Stylesheet` | `<link>` tags injected in `<head>` |
| `Scripts` | `[]Script` | `<script>` tags injected in `<head>` |
| `BodyTags` | `[]BodyTag` | Elements injected at end of `<body>` |

Helper: `wctx.APIPath("/my-endpoint")` prepends `BasePath` to a path.

---

## Session Store

Implement the `webx.SessionStore` interface:

```go
type SessionStore interface {
    Get(sessionID string, key string) (string, error)
    Set(sessionID string, key string, value string) error
    Delete(sessionID string) error
}
```

WebX uses this for CSRF tokens, theme persistence, and any app-specific session data. See `cmd/showcase/internal/session/memory.go` for an in-memory reference implementation.

---

## Middleware

| Middleware | What it does |
|---|---|
| `webx.SessionMiddleware(store)` | Creates/reads session cookie, sets `WebXContext`, validates CSRF on non-GET requests |
| `webx.SecurityHeadersMiddleware()` | Sets `X-Content-Type-Options`, `X-Frame-Options`, `Referrer-Policy` headers |

---

## Handler Registration

`ui.RegisterRoutes` registers SSE handlers for UI components. Calendar navigation and theme persistence are always registered. Use options for additional handlers:

```go
r.Route(basePath, func(r chi.Router) {
    r.Get("/stream", broker.Handler())

    ui.RegisterRoutes(r,
        ui.WithMarkdownPreview(),                    // POST /api/preview/markdown
        ui.WithDecimalParser(),                      // GET  /api/parse/decimal
        ui.WithMoneyParser(),                        // GET  /api/parse/money
        ui.WithMoneyParser("USD", "EUR"),             // restricted currencies
        ui.WithFileUpload(store),                    // POST /api/upload/files + /remove
        ui.WithFileUpload(store,                     // with validation options
            fileupload.WithMaxFiles(3),
            fileupload.WithAllowedTypes("image/"),
        ),
    )
})
```

**Always registered (zero-config):**

| Handler | Path | Method |
|---|---|---|
| Calendar navigation | `/api/calendar/navigate` | GET |
| Theme persistence | `/api/theme` | POST |

**Handlers that require app-specific logic** (register manually):

```go
// Validator — each field needs its own validation function
r.Get("/api/validate/email", validator.Handler(emailValidator))

// Form — needs validation function + success callback
r.Post("/api/auth/login", form.Handler(loginValidator, onSuccess))
```

Each handler package exports a path constant (e.g. `markdown.PreviewPath`, `moneyinput.DecimalPath`, `fileupload.UploadPath`) for use when registering handlers manually.

---

## Components

Components live in `ui/` and follow a consistent pattern:

```go
import "github.com/plaenen/webx/ui/button"

// Use with optional props
@button.Button(button.Props{Variant: button.VariantPrimary}) {
    Click me
}

// Or with zero-value defaults
@button.Button() {
    Click me
}
```

All components accept `Class` (extra Tailwind classes) and `Attributes` (extra HTML attributes) props. Classes are merged with `utils.TwMerge()`.

See [Component Migration Guide](./component-migration-guide.md) for the full component authoring pattern.

---

## Frontend Interactivity (`ds` Package)

The `ds` package provides type-safe Go helpers for Datastar attributes. Use these in templ files instead of raw `data-*` strings.

**Attribute helpers** (return `templ.Attributes`):

```go
ds.On("click", expr)       // data-on:click="expr"
ds.OnClick(expr)            // shorthand
ds.Bind("fieldName")        // data-bind:fieldName
ds.Signals(`{"count": 0}`)  // data-signals
ds.Show("$count > 0")       // data-show
ds.Text("$count")           // data-text
```

**Action expressions** (return `string` for use in `ds.On`/`ds.OnClick`):

```go
ds.Get("/api/data")         // @get('/api/data') with retry
ds.Post("/api/submit")      // @post('/api/submit') with CSRF token
ds.Put("/api/update")       // @put with CSRF
ds.Delete("/api/remove")    // @delete with CSRF
```

See [Datastar Go Reference](./datastar-go-reference.md) for the full backend SDK API and [HTML Datastar Elements Reference](./html-datastar-elements-reference.md) for all frontend attributes.

---

## Backend SSE Helpers (`ds.Send`)

Send UI actions from SSE handlers:

```go
func myHandler(w http.ResponseWriter, r *http.Request) {
    sse := datastar.NewSSE(w, r)

    // Toast notification
    ds.Send.Toast(sse, ds.ToastSuccess, "Saved!")

    // Show modal with templ content
    ds.Send.Modal(sse, myModalContent())

    // Show drawer
    ds.Send.Drawer(sse, myDrawerContent())

    // Confirmation dialog
    ds.Send.Confirm(sse, "Delete this item?", "/api/items/delete")

    // Patch a DOM element with a templ component
    ds.Send.Patch(sse, myComponent())

    // Navigation
    ds.Send.Redirect(sse, "/dashboard")
    ds.Send.Download(sse, "/files/report.pdf", "report.pdf")
}
```

The base layout includes container elements (`#modal-panel`, `#drawer-panel`, `#toast-container`) that these helpers target.

---

## Reactive Streaming

The `stream` package enables real-time updates via NATS pub/sub:

```go
// In a templ component — watch a scope and re-fetch on invalidation
@stream.WatchEffect(ctx, "invoice:42", "/api/invoice/42")

// In a handler — invalidate after mutation
broker.Invalidate("invoice:42")
```

See [Datastar Go Reference](./datastar-go-reference.md) for streaming patterns and scope conventions.

---

## Further Reading

- [Datastar Go Reference](./datastar-go-reference.md) — full backend SDK API, SSE events, patterns
- [HTML Datastar Elements Reference](./html-datastar-elements-reference.md) — all HTML elements and Datastar attribute wiring
- [Component Migration Guide](./component-migration-guide.md) — converting DaisyUI patterns into WebX templ components
