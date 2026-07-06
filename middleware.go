package dsx

import (
	"github.com/laenen-partners/dsx/internal/middleware"
)

// MiddlewareConfig configures the dsx middleware.
type MiddlewareConfig = middleware.MiddlewareConfig

// Middleware reads or creates session/CSRF/theme cookies and populates
// Context on each request. On mutating methods it validates the signed
// double-submit CSRF token.
var Middleware = middleware.Middleware

// SecurityHeadersMiddleware sets common security response headers.
var SecurityHeadersMiddleware = middleware.SecurityHeadersMiddleware
