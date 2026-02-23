package ui

import "github.com/go-chi/chi/v5"

// RouteOption configures handler registration on a chi.Router.
// Each component package exports a Route() function that returns one.
type RouteOption func(chi.Router)

// RegisterRoutes applies all given route options to the router.
//
//	r.Route(basePath, func(r chi.Router) {
//	    ui.RegisterRoutes(r,
//	        calendar.Route(),
//	        themecontroller.Route(false),
//	        markdown.Route(),
//	        moneyinput.DecimalRoute(),
//	        moneyinput.MoneyRoute("USD", "EUR"),
//	        fileupload.Route(store),
//	    )
//	})
func RegisterRoutes(r chi.Router, opts ...RouteOption) {
	for _, opt := range opts {
		opt(r)
	}
}
