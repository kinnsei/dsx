package form_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/laenen-partners/dsx/ui/form"
)

// loginSignals defines the expected form fields for testing.
type loginSignals struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func TestHandler_ReturnsErrorOnMissingID(t *testing.T) {
	h := form.Handler(loginSignals{}, func(_ string, r *http.Request) []form.FieldError {
		return nil
	}, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/submit", nil)
	h(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandler_ValidationErrors(t *testing.T) {
	h := form.Handler(loginSignals{}, func(_ string, r *http.Request) []form.FieldError {
		return []form.FieldError{
			{Field: "email_error", Message: "Email is required"},
			{Field: "password_error", Message: "Too short"},
		}
	}, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/submit?id=login-form", nil)
	r.Header.Set("Content-Type", "application/json")

	h(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Email is required") {
		t.Errorf("response should contain 'Email is required', got: %s", body[:min(len(body), 300)])
	}
	if !strings.Contains(body, "Too short") {
		t.Errorf("response should contain 'Too short', got: %s", body[:min(len(body), 300)])
	}
}

func TestHandler_ErrorFieldDerivation(t *testing.T) {
	// The handler should derive error fields from the JSON tags:
	// email → email_error, password → password_error
	h := form.Handler(loginSignals{}, func(_ string, r *http.Request) []form.FieldError {
		return []form.FieldError{
			{Field: "email_error", Message: "required"},
		}
	}, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/submit?id=login-form", nil)
	r.Header.Set("Content-Type", "application/json")

	h(w, r)

	body := w.Body.String()
	if !strings.Contains(body, "email_error") {
		t.Errorf("response should include email_error signal patch, got: %s", body[:min(len(body), 300)])
	}
}
