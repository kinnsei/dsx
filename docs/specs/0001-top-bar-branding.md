# 0023 — Separate topbar and sidebar branding in dsx dashboard layout

**Status:** Open
**Package:** `github.com/kinnsei/dsx` (`layouts/dashboard.templ`)

## Problem

The dashboard layout uses `AppBranding.Name` for both the topbar header and the sidebar title. There is no way to show different text in each location.

Current rendering:

```
┌─────────────┬──────────────────────────────────┐
│ Acme Workb. │  Acme Workbench                 │  ← both use App.Name
│─────────────│                                  │
│ Overview    │                                  │
│  Dashboard  │                                  │
│  Search     │                                  │
```

Desired:

```
┌─────────────┬──────────────────────────────────┐
│ Acme        │  Customers                       │  ← sidebar: App.Name, topbar: page Title
│ Workbench   │                                  │
│─────────────│                                  │
│ Overview    │                                  │
│  Dashboard  │                                  │
```

## Proposed change

Add a `Title` field to `DashboardProps` that overrides the topbar text independently of `App.Name`. When set, the topbar renders the title instead of the app name.

```go
type DashboardProps struct {
    // ...existing fields...

    // Title overrides the topbar text. When empty, falls back to App.Name.
    // Use for page-level context (e.g. "Customers", "Subscription — PLT-001").
    Title string
}
```

In `dashboard.templ`, the topbar link changes from:

```templ
<a href={ templ.SafeURL(props.App.DefaultHref()) } class="btn btn-ghost text-xl">
    { props.App.DefaultName() }
</a>
```

To:

```templ
<a href={ templ.SafeURL(props.App.DefaultHref()) } class="btn btn-ghost text-xl">
    if props.Title != "" {
        { props.Title }
    } else {
        { props.App.DefaultName() }
    }
</a>
```

The sidebar header (`dashboardSidebarHeader`) continues to use `App.Name` unchanged.

## Impact

- Non-breaking: `Title` defaults to empty, preserving current behavior
- App would set `Title` to the page title (already available in `DashboardProps.Title` on the app side) and `App.Name` to `"Acme Workbench"` for the sidebar
