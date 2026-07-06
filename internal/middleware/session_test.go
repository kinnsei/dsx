package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testSecret() []byte {
	return []byte("test-secret-key-must-be-32-bytes!")
}

func TestCSRFTokenRoundTrip(t *testing.T) {
	token, err := generateCSRFToken(testSecret())
	if err != nil {
		t.Fatalf("generateCSRFToken: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
	if !verifyCSRFToken(token, testSecret()) {
		t.Fatal("token should verify against same secret")
	}
}

func TestCSRFTokenRejectsInvalid(t *testing.T) {
	cases := []struct {
		name  string
		token string
	}{
		{"empty", ""},
		{"no dot", "foobar"},
		{"bad signature", "dGVzdA.badsig"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if verifyCSRFToken(tc.token, testSecret()) {
				t.Fatalf("expected token %q to be rejected", tc.token)
			}
		})
	}
}

func TestMiddlewareMutatingWithoutCSRF(t *testing.T) {
	handler := Middleware(MiddlewareConfig{
		Secret: testSecret(),
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// POST without CSRF token → should be rejected
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", nil)
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestMiddlewareMutatingWithCSRF(t *testing.T) {
	handler := Middleware(MiddlewareConfig{
		Secret: testSecret(),
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// GET to set the CSRF cookie
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	handler.ServeHTTP(w, r)

	// Collect cookies from the response
	cookies := w.Header().Values("Set-Cookie")
	var csrfCookie string
	for _, c := range cookies {
		if len(c) > 10 && c[:9] == "dsx_csrf=" {
			csrfCookie = c[9:]
			if idx := strings.Index(csrfCookie, ";"); idx >= 0 {
				csrfCookie = csrfCookie[:idx]
			}
			break
		}
	}
	if csrfCookie == "" {
		t.Fatal("no CSRF cookie set")
	}

	// Now make a POST with the CSRF token
	w2 := httptest.NewRecorder()
	r2 := httptest.NewRequest("POST", "/", nil)
	r2.Header.Set("X-CSRF-Token", csrfCookie)
	handler.ServeHTTP(w2, r2)

	if w2.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w2.Code)
	}
}
