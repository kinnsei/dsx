# Spec: dsx — Reference Data Components

> **Library**: `github.com/kinnsei/dsx`
> **Package**: `ui/countrypicker`, `ui/phoneinput`, `ui/currencypicker`, `ui/languagepicker`
> **Depends on**: `github.com/laenen-partners/domains/ref` for data + search + formatting

---

## Context

The domains SDK provides a comprehensive reference data registry (`ref.Registry`) with 103 countries, 35+ currencies, 41 languages, and 100+ timezones — all cross-referenced with search, phone prefixes, and locale-aware formatting. DSX has form inputs (`SelectInput`, `MoneyInput`) but no reference-data-aware components. Every app that collects addresses, phone numbers, or financial data needs these pickers.

These components use the Datastar SSE pattern (like `MoneyInput`): the initial render is server-side with full data, search/filtering happens via Datastar signals client-side or via SSE endpoints for larger datasets.

---

## Components

### 1. `countrypicker` — Country selection with search

**Package**: `ui/countrypicker`

**Props**:
```go
type Props struct {
    ID          string
    Class       string
    Attributes  templ.Attributes
    Name        string           // form field name
    Value       string           // selected ISO code (e.g. "BE")
    Placeholder string           // e.g. "Select country"
    Size        Size
    Variant     Variant
    Disabled    bool
    Countries   []CountryOption  // pre-rendered from ref.Registry
}

type CountryOption struct {
    Code string // "BE"
    Name string // "Belgium"
    Flag string // "🇧🇪" (derived from ISO code)
}
```

**Rendering**:
- Uses `SelectInput` internally with `<option>` children
- Each option: `🇧🇪 Belgium (BE)`
- Value is the ISO 3166-1 alpha-2 code

**Signals**:
```go
type CountryPickerSignals struct {
    Selected string `json:"selected"` // ISO code
}
```

**Server-side helper** (for populating options from the registry):
```go
func OptionsFromRegistry(reg *ref.Registry) []CountryOption {
    var opts []CountryOption
    for _, c := range reg.Countries() {
        opts = append(opts, CountryOption{
            Code: c.Code,
            Name: c.Name,
            Flag: countryFlag(c.Code), // converts "BE" → "🇧🇪"
        })
    }
    return opts
}
```

**Flag derivation**: ISO alpha-2 → regional indicator symbols. `B` = `🇧` (U+1F1E7), `E` = `🇪` (U+1F1EA) → `🇧🇪`. Pure function, no external dependency.

---

### 2. `phoneinput` — Phone number with country prefix

**Package**: `ui/phoneinput`

**Props**:
```go
type Props struct {
    ID          string
    Class       string
    Attributes  templ.Attributes
    Name        string           // form field name (receives full E.164 number)
    Value       string           // e.g. "+32 470 123 456"
    Placeholder string           // e.g. "Phone number"
    Size        Size
    Variant     Variant
    Disabled    bool
    DefaultCountry string        // ISO code for initial prefix (e.g. "BE" → "+32")
    Countries   []PhoneCountry   // available countries with prefixes
}

type PhoneCountry struct {
    Code   string // "BE"
    Name   string // "Belgium"
    Flag   string // "🇧🇪"
    Prefix string // "+32"
}
```

**Rendering**:
- Two-part input: country prefix dropdown (left) + number field (right)
- Uses DaisyUI `join` to group the dropdown and input
- Prefix dropdown shows: `🇧🇪 +32`
- Number field: plain text input
- Combined value emitted as E.164: `+32470123456`

**Signals**:
```go
type PhoneInputSignals struct {
    Country string `json:"country"` // selected ISO code
    Prefix  string `json:"prefix"`  // "+32"
    Number  string `json:"number"`  // "470123456"
    Full    string `json:"full"`    // "+32470123456" (computed)
}
```

**Behaviour**:
- When country changes → prefix updates automatically
- When user pastes a full E.164 number → auto-detect country from prefix, split into prefix + number
- `ds.Computed("full", "$prefix + $number.replace(/\\s/g, '')")` for the combined value

---

### 3. `currencypicker` — Currency selection

**Package**: `ui/currencypicker`

**Props**:
```go
type Props struct {
    ID          string
    Class       string
    Attributes  templ.Attributes
    Name        string
    Value       string            // ISO 4217 code (e.g. "EUR")
    Placeholder string
    Size        Size
    Variant     Variant
    Disabled    bool
    Currencies  []CurrencyOption  // pre-rendered from ref.Registry
}

type CurrencyOption struct {
    Code   string // "EUR"
    Name   string // "Euro"
    Symbol string // "€"
}
```

**Rendering**:
- SelectInput with options: `€ EUR — Euro`
- Value is ISO 4217 code

**Signals**:
```go
type CurrencyPickerSignals struct {
    Selected string `json:"selected"` // ISO code
}
```

**Linked country support**: Optional `CountryCode` prop — when set, auto-selects the country's primary currency and highlights it at the top of the list.

---

### 4. `languagepicker` — Language selection

**Package**: `ui/languagepicker`

**Props**:
```go
type Props struct {
    ID          string
    Class       string
    Attributes  templ.Attributes
    Name        string
    Value       string             // ISO 639-1 code (e.g. "fr")
    Placeholder string
    Size        Size
    Variant     Variant
    Disabled    bool
    Languages   []LanguageOption   // pre-rendered from ref.Registry
}

type LanguageOption struct {
    Code       string // "fr"
    Name       string // "French"
    NativeName string // "Français"
}
```

**Rendering**:
- SelectInput with options: `Français (French)`
- Native name first, English name in parentheses
- Value is ISO 639-1 code

**Linked country support**: Optional `CountryCode` prop — when set, the country's languages appear at the top of the list, separated by a divider.

---

## Shared patterns

### Server-side option builders

Each component has a `OptionsFromRegistry(*ref.Registry)` helper that converts the registry data into the component's option type. This keeps the component decoupled from `ref.Registry` — it just takes `[]XxxOption` slices.

### No client-side search for v1

For 103 countries / 35 currencies / 41 languages, a native `<select>` with all options is fast enough. Browser-native search (type to filter) works on `<select>`. Client-side search with Datastar signals can be added in v2 if needed.

### Flag emoji helper

Shared utility in `ui/countrypicker/flag.go`:
```go
func Flag(isoCode string) string {
    if len(isoCode) != 2 {
        return ""
    }
    r1 := rune(isoCode[0]) - 'A' + 0x1F1E6
    r2 := rune(isoCode[1]) - 'A' + 0x1F1E6
    return string([]rune{r1, r2})
}
```

Used by `countrypicker` and `phoneinput`.

### DaisyUI class conventions

All components use:
- `select` / `input` DaisyUI base classes
- `select-bordered` / `input-bordered` for border styling
- Variant and Size follow the same constants as `SelectInput`
- `join` for grouped inputs (phone: prefix + number)

---

## Files to create

| File | Purpose |
|---|---|
| `ui/countrypicker/countrypicker.templ` | Country select component |
| `ui/countrypicker/options.go` | `OptionsFromRegistry` + `Flag` helper |
| `ui/phoneinput/phoneinput.templ` | Phone input with prefix dropdown |
| `ui/phoneinput/options.go` | `PhoneOptionsFromRegistry` helper |
| `ui/currencypicker/currencypicker.templ` | Currency select component |
| `ui/currencypicker/options.go` | `OptionsFromRegistry` helper |
| `ui/languagepicker/languagepicker.templ` | Language select component |
| `ui/languagepicker/options.go` | `OptionsFromRegistry` helper |

## Acceptance criteria

- [ ] All components render as DaisyUI-styled form inputs
- [ ] Country picker shows flag emoji + name + code
- [ ] Phone input auto-fills prefix from selected country, combines into E.164
- [ ] Currency picker shows symbol + code + name
- [ ] Language picker shows native name + English name
- [ ] All components work with `ds.Bind` for form integration
- [ ] Options pre-rendered server-side from `ref.Registry` — no client-side data fetching
- [ ] Phone input handles paste of full E.164 numbers
