package handlers

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/laenen-partners/dsx/ds"
	"github.com/laenen-partners/dsx/ui/combobox"
	"github.com/starfederation/datastar-go/datastar"
)

type comboboxHandlers struct{}

func newComboboxHandlers() *comboboxHandlers {
	return &comboboxHandlers{}
}

func (h *comboboxHandlers) register(r chi.Router) {
	r.Get("/combobox/countries", h.searchCountries())
	r.Get("/combobox/fruits", h.searchFruits())
}

// Mock country data for the showcase.
var mockCountries = []combobox.Item{
	{Value: "BE", Label: "🇧🇪 Belgium", Description: "Europe"},
	{Value: "BR", Label: "🇧🇷 Brazil", Description: "South America"},
	{Value: "CA", Label: "🇨🇦 Canada", Description: "North America"},
	{Value: "CN", Label: "🇨🇳 China", Description: "Asia"},
	{Value: "DE", Label: "🇩🇪 Germany", Description: "Europe"},
	{Value: "ES", Label: "🇪🇸 Spain", Description: "Europe"},
	{Value: "FR", Label: "🇫🇷 France", Description: "Europe"},
	{Value: "GB", Label: "🇬🇧 United Kingdom", Description: "Europe"},
	{Value: "IN", Label: "🇮🇳 India", Description: "Asia"},
	{Value: "IT", Label: "🇮🇹 Italy", Description: "Europe"},
	{Value: "JP", Label: "🇯🇵 Japan", Description: "Asia"},
	{Value: "KR", Label: "🇰🇷 South Korea", Description: "Asia"},
	{Value: "MX", Label: "🇲🇽 Mexico", Description: "North America"},
	{Value: "NL", Label: "🇳🇱 Netherlands", Description: "Europe"},
	{Value: "PT", Label: "🇵🇹 Portugal", Description: "Europe"},
	{Value: "SE", Label: "🇸🇪 Sweden", Description: "Europe"},
	{Value: "US", Label: "🇺🇸 United States", Description: "North America"},
	{Value: "ZA", Label: "🇿🇦 South Africa", Description: "Africa"},
}

var mockFruits = []combobox.Item{
	{Value: "apple", Label: "Apple"},
	{Value: "banana", Label: "Banana"},
	{Value: "blueberry", Label: "Blueberry"},
	{Value: "cherry", Label: "Cherry"},
	{Value: "grape", Label: "Grape"},
	{Value: "kiwi", Label: "Kiwi"},
	{Value: "lemon", Label: "Lemon"},
	{Value: "mango", Label: "Mango"},
	{Value: "orange", Label: "Orange"},
	{Value: "peach", Label: "Peach"},
	{Value: "pear", Label: "Pear"},
	{Value: "pineapple", Label: "Pineapple"},
	{Value: "strawberry", Label: "Strawberry"},
	{Value: "watermelon", Label: "Watermelon"},
}

func (h *comboboxHandlers) searchCountries() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
		id := r.URL.Query().Get("id")

		var results []combobox.Item
		for _, c := range mockCountries {
			if q == "" || strings.Contains(strings.ToLower(c.Label), q) || strings.Contains(strings.ToLower(c.Value), q) {
				results = append(results, c)
			}
		}

		sse := datastar.NewSSE(w, r)
		_ = ds.Send.Patch(sse, combobox.Options(combobox.OptionsProps{
			ComboboxID: id,
			Items:      results,
		}))
	}
}

func (h *comboboxHandlers) searchFruits() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
		id := r.URL.Query().Get("id")

		var results []combobox.Item
		for _, f := range mockFruits {
			if q == "" || strings.Contains(strings.ToLower(f.Label), q) {
				results = append(results, f)
			}
		}

		sse := datastar.NewSSE(w, r)
		_ = ds.Send.Patch(sse, combobox.Options(combobox.OptionsProps{
			ComboboxID: id,
			Items:      results,
		}))
	}
}
