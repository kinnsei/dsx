# Stream API Migration Guide — v0.9.1

This guide covers migrating from the old render-time stream API (`stream.Attrs`, `stream.WatchEffect`, `stream.Connect`, `stream.Reload`) to the new DOM-driven, action-first API introduced in v0.9.1.

## What changed

The old API accumulated watchers during server-side render and required `stream.Connect()` in the base layout. The new API is fully DOM-driven — components declare `data-watch` attributes, and a JS watch worker manages SSE connections automatically.

Key improvements:
- **No more `stream.Connect()`** — the watch worker handles SSE lifecycle
- **No more render-order dependency** — DOM scanning replaces context accumulation
- **Per-domain signals** — `_ds_customers` instead of global `_dsEvent`, avoiding O(N) re-evaluations
- **Action-first builder** — type-safe, fluent API replaces fragile string-based actions
- **Reconnect protection** — automatic catch-up on SSE reconnect
- **Debounce support** — opt-in for bulk operations

## API mapping

### Removed

| Old | Replacement |
|-----|-------------|
| `stream.Connect()` | Removed. The watch worker JS (`static/js/watch-worker.js`) handles SSE automatically. Add `<meta name="stream-url" content="/stream">` to your base layout instead. |
| `stream.Attrs(ctx, scope, url)` | `stream.Watch(ctx, domain, reactions...)` |
| `stream.WatchEffect(ctx, scope, url)` | `stream.Watch(ctx, domain, reactions...)` |
| `stream.EventSignals()` | Removed. Each `Watch()` initializes its own per-domain signal. |
| `stream.Reload("actions", url)` | `stream.Structural.Get(url)`, `stream.Any.Get(url)`, etc. |
| `stream.WithID(id)` | `.ID(id)` on the builder chain |
| `stream.Debounce(d)` | `.Debounce(d)` on the builder chain |
| `dsx.Context.Watchers` | Removed. No server-side watcher accumulation. |

### Actions

| Old (string) | New (type-safe) | Use case |
|--------------|-----------------|----------|
| `"*"` | `stream.Any` | Counters, dashboards — any change |
| `"created,deleted"` | `stream.Structural` | Lists, tables — items added/removed |
| `"created"` | `stream.Created` | Only new entities |
| `"updated"` | `stream.Updated` | In-place changes (rows, details) |
| `"deleted"` | `stream.Deleted` | Only removals |
| `"created,deleted"` | `stream.Created.Or(stream.Deleted)` | Explicit combine (same as Structural) |
| custom string | `stream.Action("archived")` | App-specific actions |

### Reaction builder

| Old | New |
|-----|-----|
| `stream.Reload("*", url)` | `stream.Any.Get(url)` |
| `stream.Reload("created,deleted", url)` | `stream.Structural.Get(url)` |
| `stream.Reload("updated", url, stream.WithID(42))` | `stream.Updated.ID(42).Get(url)` |
| `stream.Reload("created,deleted", url, stream.Debounce(300*time.Millisecond))` | `stream.Structural.Debounce(300*time.Millisecond).Get(url)` |
| `stream.Reload("updated", url)` | `stream.Updated.Get(url)` |
| `stream.On(stream.Created, stream.Deleted).Get(url)` | `stream.Structural.Get(url)` |

## Migration steps

### 1. Add stream-url meta tag to base layout

Replace `@stream.Connect()` with a meta tag. The watch worker reads this to know where to connect.

```diff
- @stream.Connect()
+ <meta name="stream-url" content={ wxctx.StreamURL() }/>
```

Include `static/js/watch-worker.js` as a script in your base layout:

```html
<script src="/static/js/watch-worker.js"></script>
```

### 2. Replace stream.Attrs / stream.WatchEffect

**Old — scope-based, string pattern:**
```go
<div { stream.Attrs(ctx, "customers:*", wxctx.APIPath("/customers/list"))... }>
```

**New — domain + typed actions:**
```go
<div { stream.Watch(ctx, "customers",
    stream.Structural.Get(wxctx.APIPath("/customers/list")))... }>
```

The scope format changed: `"customers:*"` becomes domain `"customers"` with action filtering via the builder.

### 3. Replace stream.WatchEffect with Watch

**Old:**
```go
@stream.WatchEffect(ctx, "invoice:42", "/invoice/42")
```

**New:**
```go
<div { stream.Watch(ctx, "invoice",
    stream.Updated.ID(42).Get("/invoice/42"))... }>
```

### 4. Replace Reload string actions with builder

**Old:**
```go
stream.Watch(ctx, "customers",
    stream.Reload("created,deleted", url))
stream.Watch(ctx, "customers",
    stream.Reload("*", url))
stream.Watch(ctx, "customers",
    stream.Reload("updated", url, stream.WithID(42)))
```

**New:**
```go
stream.Watch(ctx, "customers",
    stream.Structural.Get(url))
stream.Watch(ctx, "customers",
    stream.Any.Get(url))
stream.Watch(ctx, "customers",
    stream.Updated.ID(42).Get(url))
```

### 5. Remove EventSignals / data-signals setup

If you had manual `data-signals` with `_dsEvent`, remove it. Each `Watch()` now initializes its own per-domain signal automatically.

**Old (in JS or templ):**
```go
data-signals={stream.EventSignals()}
```

**New:** Nothing — `Watch()` handles it.

### 6. Remove context-based watcher accumulation

If you had middleware calling `stream.WithPreRegister(ctx)` or similar context setup for watcher accumulation, remove it. The DOM-driven approach doesn't need server-side registration.

### 7. Update SSE signal references

If you have any code that references the `_dsEvent` signal directly:

| Old | New |
|-----|-----|
| `_dsEvent` | `_ds_{domain}` (e.g. `_ds_customers`) |
| `_dsEvent.domain` | Field removed — the signal key IS the domain |
| `_dsEvent.action` | `_ds_customers.action` |
| `_dsEvent.id` | `_ds_customers.id` |
| `_dsEvent.ts` | `_ds_customers.ts` |

## Example: full component migration

### Old

```go
templ CustomerList() {
    {{ wxctx := dsx.FromContext(ctx) }}
    <div id="customer-list"
        data-init={ds.GetOnce(wxctx.APIPath("/customers/list"))}
        { stream.Attrs(ctx, "customers:*", wxctx.APIPath("/customers/list"))... }>
    </div>
}
```

With `@stream.Connect()` in the base layout.

### New

```go
templ CustomerList() {
    {{ wxctx := dsx.FromContext(ctx) }}
    <div id="customer-list"
        data-init={ds.GetOnce(wxctx.APIPath("/customers/list"))}
        { stream.Watch(ctx, "customers",
            stream.Structural.Get(wxctx.APIPath("/customers/list")))... }>
    </div>
}
```

With `<meta name="stream-url">` and `watch-worker.js` in the base layout. No `stream.Connect()` needed.

## Multiple reactions

```go
// List reloads on structural changes, count reloads on any change
stream.Watch(ctx, "customers",
    stream.Structural.Get(listURL),
    stream.Any.Get(countURL))

// Row reloads on update for specific ID
stream.Watch(ctx, "customers",
    stream.Updated.ID(c.ID).Get(rowURL))

// With debounce for bulk operations
stream.Watch(ctx, "customers",
    stream.Structural.Debounce(300*time.Millisecond).Get(listURL))

// Custom actions
stream.Watch(ctx, "orders",
    stream.Action("shipped").Or(stream.Action("delivered")).Get(trackingURL))
```
