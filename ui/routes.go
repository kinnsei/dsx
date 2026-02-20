package ui

import (
	"github.com/go-chi/chi/v5"
	"github.com/plaenen/webx/ui/calendar"
	"github.com/plaenen/webx/ui/fileupload"
	"github.com/plaenen/webx/ui/markdown"
	"github.com/plaenen/webx/ui/moneyinput"
	"github.com/plaenen/webx/ui/themecontroller"
)

// RouteOption configures optional handler registration.
type RouteOption func(chi.Router)

// WithMarkdownPreview registers the markdown preview handler at
// markdown.PreviewPath (POST /api/preview/markdown).
func WithMarkdownPreview() RouteOption {
	return func(r chi.Router) {
		r.Post(markdown.PreviewPath, markdown.PreviewHandler())
	}
}

// WithDecimalParser registers the decimal parser handler at
// moneyinput.DecimalPath (GET /api/parse/decimal).
func WithDecimalParser() RouteOption {
	return func(r chi.Router) {
		r.Get(moneyinput.DecimalPath, moneyinput.DecimalHandler())
	}
}

// WithMoneyParser registers the money parser handler at
// moneyinput.MoneyPath (GET /api/parse/money). If allowedCurrencies
// is provided, only those currencies are accepted.
func WithMoneyParser(allowedCurrencies ...string) RouteOption {
	return func(r chi.Router) {
		r.Get(moneyinput.MoneyPath, moneyinput.MoneyHandler(allowedCurrencies...))
	}
}

// WithFileUpload registers upload and remove handlers at
// fileupload.UploadPath (POST /api/upload/files) and
// fileupload.RemovePath (POST /api/upload/remove).
func WithFileUpload(store *fileupload.Store, opts ...fileupload.HandlerOption) RouteOption {
	return func(r chi.Router) {
		r.Post(fileupload.UploadPath, fileupload.UploadHandler(store, opts...))
		r.Post(fileupload.RemovePath, fileupload.RemoveHandler(store))
	}
}

// RegisterRoutes registers all SSE/API handlers from UI component packages.
// Calendar navigation and theme persistence are always registered.
// Use options to enable additional handlers:
//
//	r.Route(basePath, func(r chi.Router) {
//	    ui.RegisterRoutes(r,
//	        ui.WithMarkdownPreview(),
//	        ui.WithDecimalParser(),
//	        ui.WithMoneyParser(),
//	        ui.WithFileUpload(store),
//	    )
//	})
func RegisterRoutes(r chi.Router, opts ...RouteOption) {
	r.Get(calendar.NavigatePath, calendar.NavigateHandlerFromQuery())
	r.Post(themecontroller.SetThemePath, themecontroller.SetThemeHandler())

	for _, opt := range opts {
		opt(r)
	}
}
