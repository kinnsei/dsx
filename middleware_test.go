package webx

import (
	"net/http"
	"net/http/httptest"
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

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rec.Code)
	}
}

func TestMiddlewareMutatingWithValidCSRF(t *testing.T) {
	var reached bool
	handler := Middleware(MiddlewareConfig{
		Secret: testSecret(),
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reached = true
		w.WriteHeader(http.StatusOK)
	}))

	token, err := generateCSRFToken(testSecret())
	if err != nil {
		t.Fatalf("generateCSRFToken: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	req.AddCookie(&http.Cookie{Name: csrfCookieName, Value: token})
	req.Header.Set("X-CSRF-Token", token)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !reached {
		t.Fatal("handler was not called")
	}
}

func TestMiddlewareThemeCookie(t *testing.T) {
	var gotTheme string
	handler := Middleware(MiddlewareConfig{
		Secret: testSecret(),
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := FromContext(r.Context())
		gotTheme = ctx.Theme
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(&http.Cookie{Name: themeCookieName, Value: "dark"})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if gotTheme != "dark" {
		t.Fatalf("expected theme 'dark', got %q", gotTheme)
	}
}

func TestMiddlewareSessionPersistence(t *testing.T) {
	handler := Middleware(MiddlewareConfig{
		Secret: testSecret(),
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// First request: should get a session cookie set.
	req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, req1)

	var sessionCookie *http.Cookie
	for _, c := range rec1.Result().Cookies() {
		if c.Name == sessionCookieName {
			sessionCookie = c
			break
		}
	}
	if sessionCookie == nil {
		t.Fatal("expected session cookie to be set")
	}

	// Second request with same cookie: should not get a new session cookie.
	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req2.AddCookie(sessionCookie)
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)

	for _, c := range rec2.Result().Cookies() {
		if c.Name == sessionCookieName {
			t.Fatal("session cookie should not be re-set on existing session")
		}
	}
}

func TestSecurityHeadersMiddleware(t *testing.T) {
	handler := SecurityHeadersMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	expected := map[string]string{
		"X-Content-Type-Options":  "nosniff",
		"X-Frame-Options":         "DENY",
		"Referrer-Policy":         "strict-origin-when-cross-origin",
		"Content-Security-Policy": "default-src 'self'; script-src 'self' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data:",
	}

	for header, want := range expected {
		got := rec.Header().Get(header)
		if got != want {
			t.Errorf("%s: got %q, want %q", header, got, want)
		}
	}

	// HSTS should NOT be present without secure flag
	if rec.Header().Get("Strict-Transport-Security") != "" {
		t.Error("HSTS should not be set without secure flag")
	}
}

func TestSecurityHeadersMiddlewareWithHSTS(t *testing.T) {
	handler := SecurityHeadersMiddleware(true)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	hsts := rec.Header().Get("Strict-Transport-Security")
	if hsts != "max-age=63072000; includeSubDomains" {
		t.Errorf("HSTS: got %q, want %q", hsts, "max-age=63072000; includeSubDomains")
	}
}
