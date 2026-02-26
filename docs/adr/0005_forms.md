# ADR-0005: Forms — Signal-Namespaced Server-Side Validation

## Status

Accepted

## Context

Forms in server-rendered applications traditionally follow the POST-redirect-GET cycle: submit a form, server validates, either redirect on success or re-render the page with errors. This requires full page reloads, loses client state, and makes inline field-level error display awkward.

WebX uses Datastar for frontend interactivity and SSE for server-to-client communication. This enables a different model: forms submit via SSE, the server validates and patches signals back, and the client updates inline — no page reload, no lost state, no flash-of-error.

The form system must solve three problems:

1. **Multiple forms per page**: A dashboard page might have a login form, a contact form, and a search bar. Their signals must not collide.
2. **Inline error display**: Each field needs its own error signal that shows/hides a message without affecting other fields.
3. **Consistent submit lifecycle**: Loading spinners, disabled buttons, error clearing, and success feedback must follow the same pattern everywhere.

## Decision

### Architecture Overview

The form system has three parts:

```
form.Form (templ component)          form.Handler (Go handler wrapper)
─────────────────────────            ──────────────────────────────────
Renders <form> with:                 Wraps validate + onSuccess:
 • namespaced data-signals           1. Read form ID from ?id=
 • data-on:submit → @post            2. Call validate(formID, r)
 • novalidate (no HTML5)             3. If errors → patch error signals
 • submit button with spinner        4. If success → patch submitting=false,
                                        call onSuccess(formID, sse)
         │                                          │
         │     POST /action?id=login                │
         │     Body: {"login": {email, password}}    │
         └──────────────────────────────────────────┘
                           │
                           ▼
                    ds.ReadSignals("login", r, &dest)
                    Reads only the "login" namespace
```

### Signal Namespacing

Every form gets a unique ID. All its signals live under that ID as a namespace. This is the key mechanism that enables multiple forms per page.

A form with `ID: "login"` and signals `{email: "", password: ""}` renders:

```html
<form id="login" data-signals="{login: {email: '', password: ''}}">
```

A second form with `ID: "contact"` on the same page renders:

```html
<form id="contact" data-signals="{contact: {name: '', email: '', message: ''}}">
```

Datastar's signal store becomes:

```json
{
  "login": {"email": "", "password": "", "submitting": false, "error": ""},
  "contact": {"name": "", "email": "", "message": "", "submitting": false, "error": ""}
}
```

Each form's signals are isolated. Setting `login.email` does not affect `contact.email`. The submit button references its form's signals specifically: `$login.submitting`, not a global `submitting`.

#### ID sanitization

HTML allows hyphens in IDs (`new-customer`), but JavaScript property access does not (`$new-customer.email` is a syntax error). Datastar internally converts hyphens to underscores. The form system mirrors this on both sides:

- **Frontend**: `ds.NewSignals("new-customer", ...)` sanitizes to `new_customer` in the `data-signals` output
- **Backend**: `ds.ReadSignals("new-customer", r, &dest)` sanitizes to `new_customer` before looking up the namespace in the JSON body
- **Handler**: `form.Handler` sanitizes the form ID before patching response signals

This means developers can use natural hyphenated IDs everywhere and the system handles the conversion transparently.

### The Form Component

`form.Form` renders a `<form>` element with two signal layers:

```go
templ Form(props Props) {
    signals := ds.NewSignals(props.ID, props.Signals)       // user's field signals
    formSignals := ds.NewSignals(props.ID, FormSignals{})   // built-in: submitting, error

    <form id={props.ID}
          data-signals={signals.DataSignals}
          data-on:submit__prevent={submitAction}
          novalidate>
        <div data-signals={formSignals.DataSignals} class="contents">
            {children...}
        </div>
    </form>
}
```

**Two `data-signals` blocks**: The user's field signals (`email`, `password`) and the form's built-in signals (`submitting`, `error`) are rendered separately. This is because they come from different structs — the user defines their fields, while `FormSignals` is always the same. Datastar merges them into a single namespace at runtime.

**`novalidate`**: HTML5 validation is disabled. All validation happens server-side. This ensures the server is always the single source of truth and validation logic isn't duplicated between client and server.

**`data-on:submit__prevent`**: The `__prevent` modifier calls `preventDefault()`, stopping the browser's default form submission. Instead, Datastar executes the action expression — an SSE POST to the form's action URL.

**Action URL includes form ID**: The submit action builds `{Action}?id={FormID}`. The handler reads this to know which signal namespace to patch in the response.

#### FormSignals

Every form automatically gets two built-in signals:

```go
type FormSignals struct {
    Submitting bool   `json:"submitting"`
    Error      string `json:"error"`
}
```

- `submitting` — `true` while the request is in flight. The submit button reads this to show a spinner and disable itself.
- `error` — form-level error message (displayed by `form.FormError` or as a toast). Distinct from field-level errors.

### The Submit Button

`form.Submit` renders a button that references the form's `submitting` signal:

```go
templ Submit(props SubmitProps) {
    signals := ds.NewSignals(props.FormID, FormSignals{})

    <button type="submit"
            {ds.Attr("disabled", signals.Signal("submitting"))...}>
        <span class="loading loading-spinner loading-sm"
              {ds.Show(signals.Signal("submitting"))...}></span>
        {children...}
    </button>
}
```

This produces:

```html
<button type="submit" data-attr:disabled="$login.submitting">
    <span class="loading loading-spinner loading-sm"
          data-show="$login.submitting"></span>
    Sign In
</button>
```

When `$login.submitting` becomes `true`, the button disables and shows a spinner. When the server patches `submitting: false`, the button re-enables. This prevents double-submission without any custom JavaScript.

### Field Error Display

The `form.Error` component shows/hides based on a signal value:

```go
templ Error(signal string) {
    <p class="label text-xs text-error"
       {ds.Show(signal + " !== ''")...}>
        <span {ds.Text(signal)...}></span>
    </p>
}
```

Usage in a form:

```go
@form.Field() {
    @form.Label() { Email }
    <input class="input" { ds.Bind("login", "email")... } />
    @form.Error("$login.email_error")
}
```

The error message is:
- **Hidden** when `$login.email_error === ""` (empty string)
- **Visible** when the server patches a non-empty value
- **Dynamic**: the text comes from the signal, so the server controls the exact message

`form.ErrorStatic` is a variant where the error text is hardcoded in the template and the signal only controls visibility.

`form.Success` follows the same pattern but with `text-success` styling.

### The Form Handler

`form.Handler` is a Go function that wraps the validate/success pattern into an `http.HandlerFunc`:

```go
func Handler(
    signals any,
    validate SubmitFunc,
    onSuccess func(formID string, sse *datastar.ServerSentEventGenerator),
) http.HandlerFunc
```

The handler's lifecycle:

```
POST /form/login?id=login
Body: {"login": {"email": "test", "password": ""}}
│
├─ 1. Read form ID from ?id= query param
│
├─ 2. Call validate(formID, r)
│     └─ Validate reads signals: ds.ReadSignals(formID, r, &signals)
│     └─ Returns []FieldError or nil
│
├─ 3. Create SSE writer: datastar.NewSSE(w, r)
│     ⚠️ MUST be after ReadSignals — SSE creation consumes the body
│
├─ 4a. If errors:
│      ├─ Clear all declared errorFields to ""
│      ├─ Field errors → patch as signals: {login: {email_error: "...", submitting: false}}
│      └─ Form errors (Field=="error") → send as toast via ds.Send.Toast()
│
└─ 4b. If success:
       ├─ Patch {login: {submitting: false}} + clear all errorFields
       └─ Call onSuccess(formID, sse) — hide drawer, show toast, redirect, etc.
```

#### Reading signals

`ds.ReadSignals` extracts the form's namespace from the request body:

```go
func ReadSignals(componentID string, r *http.Request, dest any) error {
    sanitizedID := strings.ReplaceAll(componentID, "-", "_")

    // Read and parse full body: {"login": {...}, "contact": {...}, ...}
    var top map[string]json.RawMessage
    json.Unmarshal(raw, &top)

    // Extract only our namespace
    nsRaw := top[sanitizedID]
    json.Unmarshal(nsRaw, dest)
}
```

This is critical for multiple forms: when the request body contains signals from all forms on the page, `ReadSignals` extracts only the namespace matching the form ID. A login handler never accidentally reads contact form signals.

**Body consumption order**: `ReadSignals` must be called **before** `datastar.NewSSE()`. The SSE constructor flushes HTTP headers, which can consume the request body. The form handler enforces this order internally — `validate` runs first, then SSE is created.

#### Error routing

The handler routes errors based on the `Field` value:

| `FieldError.Field` | Behavior |
|---------------------|----------|
| Any field name (e.g. `"email_error"`) | Patched as a signal under the form's namespace |
| `"error"` (special) | Sent as a toast notification via `ds.Send.Toast()` |

This dual routing allows handlers to show both inline field errors and form-level error banners from a single return value:

```go
return []form.FieldError{
    {Field: "email_error", Message: "Email is required"},
    {Field: "error", Message: "Service temporarily unavailable"},
}
// → inline error under email field
// → toast notification at bottom-right
```

### Complete Submission Flow

Here is the full round-trip for a login form with a validation error:

```
Browser                              Server
───────                              ──────

1. User types "bad" in email,
   leaves password empty,
   clicks "Sign In"
        │
2. Datastar fires:                   3. Handler receives:
   data-on:submit__prevent              POST /form/login?id=login
   @post('/form/login?id=login')        Body: {"login": {"email":"bad","password":"",
                                               "submitting":false,"error":"",
                                               "email_error":"","password_error":""}}
        │                                      │
        │                               4. validate(formID, r):
        │                                  ds.ReadSignals("login", r, &signals)
        │                                  → signals.Email = "bad"
        │                                  → signals.Password = ""
        │                                  return [
        │                                    {Field:"email_error", Message:"Invalid email"},
        │                                    {Field:"password_error", Message:"Password is required"},
        │                                  ]
        │                                      │
        │                               5. SSE response:
        │                                  event: datastar-patch-signals
        │◄─────────────────────────────    data: {login: {
                                               email_error: "Invalid email",
                                               password_error: "Password is required",
                                               submitting: false
                                           }}
        │
6. Datastar patches signals:
   $login.email_error = "Invalid email"
   $login.password_error = "Password is required"
   $login.submitting = false
        │
7. data-show fires:
   email error message appears
   password error message appears
   submit button re-enables
```

On success, the flow is similar but step 5 patches `submitting: false` and the `onSuccess` callback fires — typically closing a drawer, showing a toast, or redirecting.

### Forms in Drawers and Modals

Drawers and modals render their content with `context.Background()`, not the request context (see ADR-0004). This means forms inside drawers cannot call `webx.FromContext(ctx)` to read the base path.

The solution: pass the action URL as a parameter:

```go
// Handler opens the drawer with the action URL pre-built
func (h *handlers) newDrawer() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        wxctx := webx.FromContext(r.Context())
        sse := datastar.NewSSE(w, r)
        ds.Send.Drawer(sse, CustomerDrawer(wxctx.APIPath("/customers/create")))
    }
}

// Drawer template receives the URL, doesn't need context
templ CustomerDrawer(actionURL string) {
    @form.Form(form.Props{
        ID:      "new-customer",
        Action:  actionURL,
        Signals: newCustomerSignals{},
    }) {
        // fields...
    }
}
```

The form works identically in a drawer as on a page. Signal namespacing, submit lifecycle, and error patching are all the same. The only difference is how the form gets its action URL.

Success callbacks in drawer forms typically hide the drawer:

```go
func(formID string, sse *datastar.ServerSentEventGenerator) {
    ds.Send.HideDrawer(sse)
    ds.Send.Toast(sse, ds.ToastSuccess, "Customer created")
}
```

### Multiple Forms on One Page

With the namespacing model, multiple forms coexist without configuration:

```go
// Form 1
@form.Form(form.Props{ID: "login", Action: "/auth/login", Signals: loginSignals{}}) {
    <input { ds.Bind("login", "email")... } />
    @form.Error("$login.email_error")
    @form.Submit(form.SubmitProps{FormID: "login"}) { Sign In }
}

// Form 2
@form.Form(form.Props{ID: "contact", Action: "/contact", Signals: contactSignals{}}) {
    <input { ds.Bind("contact", "email")... } />
    @form.Error("$contact.email_error")
    @form.Submit(form.SubmitProps{FormID: "contact"}) { Send }
}
```

What happens when the login form submits:

1. Datastar serializes **all signals** on the page into the request body:
   ```json
   {"login": {"email": "...", "password": "..."}, "contact": {"name": "", "email": "", "message": ""}}
   ```

2. The login handler calls `ds.ReadSignals("login", r, &signals)`, which extracts only `{"email": "...", "password": "..."}`. The contact form's signals are ignored.

3. The response patches only `login.*` signals. The contact form's state is untouched.

4. Each submit button is scoped to its own form: `$login.submitting` vs `$contact.submitting`. Clicking "Sign In" doesn't disable "Send".

There is no cross-contamination because:
- Signal references are always fully qualified (`$login.email`, not `$email`)
- `ds.ReadSignals` reads only one namespace
- The handler patches back to the same namespace
- `ds.Bind`, `ds.Show`, and `ds.Text` all use the full `$namespace.field` path

### Error Clearing

The `signals` parameter in `form.Handler` declares the form's field shape. Error signal names are derived automatically by appending `_error` to each json tag in the struct. On every submission, all derived error fields are cleared to `""` before applying actual errors. This prevents stale errors from a previous submission from persisting when the user fixes some fields and resubmits.

```go
type loginSignals struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

// Derives error fields: ["email_error", "password_error"]
// Both are cleared to "" on every submit before applying actual errors.
form.Handler(loginSignals{}, validateLogin, onSuccess)
```

On the **error path**: all derived error fields are set to `""`, then actual errors overwrite their fields. Fields that are now valid stay cleared.

On the **success path**: all derived error fields are set to `""` alongside `submitting: false`, ensuring no error messages linger after a successful submission.

### Inline Validators vs. Form Submission

WebX has two separate validation mechanisms that work alongside each other:

| Mechanism | When it runs | What it returns | Scope |
|-----------|-------------|-----------------|-------|
| `form.Handler` | On form submit | `[]FieldError` → signal patches | Full form |
| `validator.Input` | On field input (debounced) | Signal patch for one field | Single field |

Inline validators (`ui/validator`) provide immediate feedback as the user types. They call a backend endpoint per field (e.g. `GET /validate/email`) and patch a field-specific signal (`valid`, `error`). They do not submit the form and do not interact with the form handler.

Both can coexist on the same field. The inline validator gives live feedback, and the form handler re-validates on submit for the definitive check.

## Consequences

### Positive

- **Multiple forms per page work by default**: Signal namespacing prevents collisions without any extra configuration
- **Consistent UX**: Every form gets the same submit lifecycle — loading spinner, disabled button, inline errors, toast fallback
- **Server as single source of truth**: `novalidate` + server-side `SubmitFunc` means validation logic lives in one place
- **No client-side JavaScript**: Error display, loading states, and submit disabling are all signal-driven via Datastar attributes
- **Composable**: Forms work identically in pages, drawers, and modals — the only difference is how they receive their action URL

### Negative

- **~~All signals sent on every submit~~** (fixed): The form component now applies `ds.WithFilterSignals(props.ID)` to every submit action. This scopes the request payload to only the submitting form's signal namespace, avoiding sending the full signal tree on pages with multiple forms.
- **~~Error clearing is implicit~~** (fixed): `form.Handler` now takes the signals struct as its first parameter and derives error field names automatically (appending `_error` to each json tag). On every submission, all derived error fields are cleared to `""` before applying actual errors (error path) or as part of the success patch (success path). No manual error field declaration needed.
- **~~Body consumption ordering~~** (fixed): The `ds.ReadAndSSE(componentID, w, r, &dest)` helper reads signals and creates the SSE writer in the correct order, making it impossible to call them in the wrong sequence. `form.Handler` still enforces this internally, but custom handlers can now use `ReadAndSSE` instead of remembering the constraint.

### Trade-offs

- **Server-only validation vs. client-side**: No client-side validation means every validation requires a round-trip. For simple checks (required fields, email format), this adds latency. The inline `validator.Input` component mitigates this for fields where instant feedback matters, but it's opt-in per field.
- **Flat error signals vs. nested**: Error fields are siblings of data fields in the signal namespace (`email` and `email_error` live side by side). An alternative would be nested error objects (`errors.email`), but the flat model is simpler to reference in `data-show` expressions and matches how `FieldError.Field` maps to signal keys.
- **Form ID in query string vs. body**: The form ID is passed as `?id=login` in the URL rather than inside the signal body. This keeps the handler routing simple (one handler per action URL, form ID as a parameter) and avoids parsing the body twice. The trade-off is a slightly longer action URL.
