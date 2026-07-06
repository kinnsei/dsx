# CLAUDE.md

## Tech stack

Go + Chi (routing) + Templ (templating) + Tailwind CSS + DaisyUI (styling) + Datastar (frontend interactivity)

Module: `github.com/kinnsei/dsx`

## Commands

- `go tool templ generate` — generate Go code from .templ files
- `go build ./cmd/...` — build the application
- `go run ./cmd/...` — run the application
- `go tool task install:all` — install all dependencies (mod tidy + download DaisyUI)
- `go tool task install:daisyui` — download DaisyUI plugin files
- `go tool gotailwind` — build Tailwind CSS

## Project structure

```
cmd/            — main entry point (create when needed)
internal/       — internal packages (create when needed)
ui/             — DaisyUI components (one dir per component, e.g. ui/button/button.templ)
utils/          — shared templ utilities (TwMerge, If, RandomID, etc.)
static/css/     — Tailwind CSS + DaisyUI plugin files
docs/           — reference documentation
Taskfile.yaml   — task runner config
```

## Component pattern

Components use DaisyUI CSS classes and follow this pattern:

```go
package mycomponent

type Props struct {
    ID         string
    Class      string
    Attributes templ.Attributes
    // component-specific fields (Variant, Size, etc.)
}

templ MyComponent(props ...Props) {
    {{ var p Props }}
    if len(props) > 0 {
        {{ p = props[0] }}
    }
    // use utils.TwMerge() to combine DaisyUI + Tailwind classes
}
```

DaisyUI class conventions: `btn`, `btn-primary`, `card`, `card-body`, `card-title`, `card-actions`, `input`, `badge`, `alert`, `modal`, `drawer`, `tabs`, `dropdown`, etc. See [DaisyUI docs](https://daisyui.com/components/).

Import components as: `"github.com/kinnsei/dsx/ui/button"`
Import utils as: `"github.com/kinnsei/dsx/utils"`

## Rules

- **No Co-Author lines** in commits
- **No custom CSS/JS** — use DaisyUI classes + Datastar only
- **Wrap errors**: `fmt.Errorf("context: %w", err)`
- **Fix root cause** — when fixing a bug, fix the root cause, not the symptom
- **Use DaisyUI** for all UI styling — components live in `./ui`
- **Use Datastar for all frontend interactivity** — no raw JS
  - Read [Datastar Go reference](./docs/datastar-go-reference.md) before writing Go-side Datastar code
  - Read [Datastar HTML elements reference](./docs/html-datastar-elements-reference.md) before writing templ-side Datastar attributes
- **Run `go tool templ fmt` then `go tool templ generate`** after editing any `.templ` file, and always re-run both **before committing** to ensure formatting and `_templ.go` files are up to date and staged
- **Ask if backward compatibility is required** before making changes
- **Showcase should be production grade** — follow best practices for security, performance, and maintainability
- **Datastar On** - use all Datastar features correctly, the latest version uses for example `data-on:click` instead of `data-on-click`
- **App agnostic library** - the library should not be tied to any specific application, it should be usable by any application, cmd/dashboard is just an example of how to use the library and can show how to use it
- **Use the Theme** - use the theme tokens and colours do NOT use hardcoded colours / tokens
