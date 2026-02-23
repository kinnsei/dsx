package webx

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
)

const (
	sessionCookieName = "webx_session"
	csrfCookieName    = "webx_csrf"
	themeCookieName   = "webx_theme"
)

// MiddlewareConfig configures the webx middleware.
type MiddlewareConfig struct {
	Secret []byte // HMAC-SHA256 key, minimum 32 bytes — panics if shorter
	Secure bool   // Secure flag on cookies (true for HTTPS / production)
}

// Middleware reads or creates session/CSRF/theme cookies and populates
// Context on each request. On mutating methods it validates the signed
// double-submit CSRF token.
func Middleware(cfg MiddlewareConfig) func(http.Handler) http.Handler {
	if len(cfg.Secret) < 32 {
		panic("webx: MiddlewareConfig.Secret must be at least 32 bytes")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// --- Session cookie ---
			sessionID, isNewSession, err := sessionIDFromRequest(r)
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			if isNewSession {
				http.SetCookie(w, &http.Cookie{
					Name:     sessionCookieName,
					Value:    sessionID,
					Path:     "/",
					HttpOnly: true,
					SameSite: http.SameSiteLaxMode,
					Secure:   cfg.Secure,
				})
			}

			// --- CSRF cookie (signed double-submit) ---
			csrfToken := csrfTokenFromCookie(r, cfg.Secret)
			if csrfToken == "" {
				var err error
				csrfToken, err = generateCSRFToken(cfg.Secret)
				if err != nil {
					http.Error(w, fmt.Sprintf("failed to generate CSRF token: %v", err), http.StatusInternalServerError)
					return
				}
				http.SetCookie(w, &http.Cookie{
					Name:     csrfCookieName,
					Value:    csrfToken,
					Path:     "/",
					HttpOnly: true,
					SameSite: http.SameSiteLaxMode,
					Secure:   cfg.Secure,
				})
			}

			// Validate CSRF on mutating requests.
			if r.Method != http.MethodGet && r.Method != http.MethodHead && r.Method != http.MethodOptions {
				headerToken := r.Header.Get("X-CSRF-Token")
				if subtle.ConstantTimeCompare([]byte(headerToken), []byte(csrfToken)) != 1 {
					http.Error(w, "invalid or missing CSRF token", http.StatusForbidden)
					return
				}
			}

			// --- Theme cookie ---
			theme := ""
			if c, err := r.Cookie(themeCookieName); err == nil {
				theme = c.Value
			}

			ctx := FromContext(r.Context())
			ctx.SessionID = sessionID
			ctx.CSRFToken = csrfToken
			ctx.Theme = theme

			next.ServeHTTP(w, r.WithContext(ctx.WithContext(r.Context())))
		})
	}
}

// csrfTokenFromCookie reads the CSRF cookie and verifies the HMAC signature.
// Returns the token value if valid, or empty string if missing/invalid.
func csrfTokenFromCookie(r *http.Request, secret []byte) string {
	c, err := r.Cookie(csrfCookieName)
	if err != nil || c.Value == "" {
		return ""
	}
	if verifyCSRFToken(c.Value, secret) {
		return c.Value
	}
	return ""
}

// generateCSRFToken creates a signed CSRF token: base64url(nonce).base64url(hmac).
func generateCSRFToken(secret []byte) (string, error) {
	nonce := make([]byte, 32)
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("reading random bytes: %w", err)
	}
	nonceEnc := base64.RawURLEncoding.EncodeToString(nonce)
	mac := hmac.New(sha256.New, secret)
	mac.Write(nonce)
	sigEnc := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return nonceEnc + "." + sigEnc, nil
}

// verifyCSRFToken checks that the token has a valid HMAC signature.
func verifyCSRFToken(token string, secret []byte) bool {
	parts := strings.SplitN(token, ".", 2)
	if len(parts) != 2 {
		return false
	}
	nonce, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}
	sig, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}
	mac := hmac.New(sha256.New, secret)
	mac.Write(nonce)
	return hmac.Equal(mac.Sum(nil), sig)
}

// sessionIDFromRequest returns the session ID from the cookie, or generates a
// new one. The bool indicates whether the ID is new.
func sessionIDFromRequest(r *http.Request) (string, bool, error) {
	if c, err := r.Cookie(sessionCookieName); err == nil && c.Value != "" {
		return c.Value, false, nil
	}
	id, err := randomHex(16)
	if err != nil {
		return "", true, fmt.Errorf("generating session ID: %w", err)
	}
	return id, true, nil
}

// SecurityHeadersMiddleware sets common security response headers.
// When secure is true (i.e. HTTPS / production), HSTS is also added.
func SecurityHeadersMiddleware(secure ...bool) func(http.Handler) http.Handler {
	addHSTS := len(secure) > 0 && secure[0]
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data:")
			if addHSTS {
				w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
			}
			next.ServeHTTP(w, r)
		})
	}
}

func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("reading random bytes: %w", err)
	}
	return hex.EncodeToString(b), nil
}
