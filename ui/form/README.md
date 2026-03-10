# Form

Datastar-powered form component with SSE submission, signal-based validation, and automatic error handling.

## Import

```go
import "github.com/laenen-partners/dsx/ui/form"
```

## Quick start

```go
type LoginSignals struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

@form.Form(form.Props{
    ID:      "login",
    Action:  "/api/auth/login",
    Signals: LoginSignals{},
}) {
    @form.Field() {
        @form.Label() { Email }
        <input type="email" class="input" { ds.Bind("login", "email")... } />
        @form.Error("$login.email_error")
    }
    @form.Field() {
        @form.Label() { Password }
        <input type="password" class="input" { ds.Bind("login", "password")... } />
        @form.Error("$login.password_error")
    }
    @form.FormError("login")
    @form.Submit(form.SubmitProps{FormID: "login", Variant: button.VariantPrimary}) {
        Sign In
    }
}
```

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `ID` | `string` | auto-generated | Unique form identifier. Used for signal namespacing. |
| `Class` | `string` | `""` | Additional CSS classes on the `<form>`. |
| `Attributes` | `templ.Attributes` | `nil` | Arbitrary HTML attributes. |
| `Action` | `string` | `""` | Backend endpoint for form submission. |
| `Method` | `string` | `"post"` | HTTP method: `"post"`, `"put"`, `"patch"`, `"delete"`. |
| `Multipart` | `bool` | `false` | Send as `multipart/form-data`. Required for file uploads. |
| `Signals` | `any` | `nil` | Initial form signal state (your custom signals struct). |

## Sub-components

### Field

Wraps a label + input + error into a `<fieldset>` group.

```go
@form.Field() { ... }
@form.Field(form.FieldProps{Class: "gap-1"}) { ... }
```

### Label

Renders a `<legend>` label inside a field.

```go
@form.Label() { Email }
@form.Label(form.LabelProps{Class: "text-lg"}) { Email }
```

### Description

Helper text below a form field.

```go
@form.Description() { We'll never share your email. }
```

### Error

Dynamic error message bound to a signal. Shows when the signal is non-empty.

```go
@form.Error("$login.email_error")
```

### ErrorStatic

Static error message that shows when a signal is non-empty. Use when the error text is known at render time.

```go
@form.ErrorStatic("$login.email_error") { Please enter a valid email }
```

### Success

Success message bound to a signal. Shows when the signal is non-empty.

```go
@form.Success("$login.success")
```

### Submit

Submit button with automatic loading spinner while the form is submitting.

```go
@form.Submit(form.SubmitProps{
    FormID:  "login",
    Variant: button.VariantPrimary,
    Class:   "w-full",
}) { Sign In }
```

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `FormID` | `string` | `""` | The form's ID for accessing its signals. |
| `Variant` | `button.Variant` | none | Button color variant (e.g. `button.VariantPrimary`). |
| `Class` | `string` | `""` | Additional CSS classes. |
| `Attributes` | `templ.Attributes` | `nil` | Arbitrary HTML attributes. |

### FormError

Form-level error banner. Shows when the form's `"error"` signal is non-empty.

```go
@form.FormError("login")
```

## Backend handler

`form.Handler` returns an `http.HandlerFunc` that processes submissions via SSE. Error signal names are derived automatically by appending `_error` to each JSON tag in your signals struct.

```go
r.Post("/api/auth/login", form.Handler(
    LoginSignals{},
    func(formID string, r *http.Request) []form.FieldError {
        email := ds.ReadSignal[string](formID, r, "email")
        if email == "" {
            return []form.FieldError{{Field: "email_error", Message: "Email is required"}}
        }
        return nil
    },
    func(formID string, sse *datastar.ServerSentEventGenerator) {
        ds.Send.Redirect(sse, "/dashboard")
    },
))
```

Return a `FieldError` with `Field: "error"` to show a toast instead of an inline error:

```go
return []form.FieldError{{Field: "error", Message: "Something went wrong"}}
```

## Multipart / file uploads

Set `Multipart: true` to submit the form as `multipart/form-data`. This adds `enctype="multipart/form-data"` to the HTML and tells Datastar to use `contentType: 'form'`.

```go
@form.Form(form.Props{
    ID:        "upload",
    Action:    "/api/documents/upload",
    Multipart: true,
    Signals:   UploadSignals{},
}) {
    @form.Field() {
        @form.Label() { Document }
        <input type="file" name="document" class="file-input" />
    }
    @form.Submit(form.SubmitProps{FormID: "upload", Variant: button.VariantPrimary}) {
        Upload
    }
}
```

In the handler, access files via the standard `*http.Request` methods:

```go
func(formID string, r *http.Request) []form.FieldError {
    file, header, err := r.FormFile("document")
    if err != nil {
        return []form.FieldError{{Field: "error", Message: "Please select a file"}}
    }
    defer file.Close()
    // process file...
    return nil
}
```

## How it works

1. The form renders with `data-signals` for your signals struct plus `FormSignals` (submitting, error).
2. On submit, Datastar fires an SSE request to the `Action` endpoint with the form's signals.
3. The backend handler validates and responds with signal patches (errors or success state).
4. Error signals are automatically cleared on each submission before applying new errors.
5. The submit button disables itself and shows a spinner while `submitting` is true.
