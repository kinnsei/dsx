package handlers

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/laenen-partners/dsx/ds"
	"github.com/laenen-partners/dsx/ui/multiselect"
	"github.com/starfederation/datastar-go/datastar"
)

type multiselectHandlers struct{}

func newMultiselectHandlers() *multiselectHandlers {
	return &multiselectHandlers{}
}

func (h *multiselectHandlers) register(r chi.Router) {
	r.Get("/multiselect/languages", h.searchLanguages())
}

var mockLanguages = []multiselect.Item{
	{Value: "ar", Label: "Arabic", Description: "العربية"},
	{Value: "zh", Label: "Chinese", Description: "中文"},
	{Value: "nl", Label: "Dutch", Description: "Nederlands"},
	{Value: "en", Label: "English"},
	{Value: "fr", Label: "French", Description: "Français"},
	{Value: "de", Label: "German", Description: "Deutsch"},
	{Value: "hi", Label: "Hindi", Description: "हिन्दी"},
	{Value: "it", Label: "Italian", Description: "Italiano"},
	{Value: "ja", Label: "Japanese", Description: "日本語"},
	{Value: "ko", Label: "Korean", Description: "한국어"},
	{Value: "pt", Label: "Portuguese", Description: "Português"},
	{Value: "ru", Label: "Russian", Description: "Русский"},
	{Value: "es", Label: "Spanish", Description: "Español"},
	{Value: "sv", Label: "Swedish", Description: "Svenska"},
	{Value: "tr", Label: "Turkish", Description: "Türkçe"},
}

func langLabel(value string) string {
	for _, l := range mockLanguages {
		if l.Value == value {
			return l.Label
		}
	}
	return value
}

func (h *multiselectHandlers) searchLanguages() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
		id := r.URL.Query().Get("id")
		selectedCSV := r.URL.Query().Get("selected")
		addValue := r.URL.Query().Get("add")
		removeValue := r.URL.Query().Get("remove")
		searchURL := "/showcase/multiselect/languages"

		// Parse selected values.
		selectedSet := map[string]bool{}
		var selectedList []string
		if selectedCSV != "" {
			for _, v := range strings.Split(selectedCSV, ",") {
				v = strings.TrimSpace(v)
				if v != "" {
					selectedSet[v] = true
					selectedList = append(selectedList, v)
				}
			}
		}

		// Apply add/remove.
		if addValue != "" && !selectedSet[addValue] {
			selectedSet[addValue] = true
			selectedList = append(selectedList, addValue)
		}
		if removeValue != "" {
			delete(selectedSet, removeValue)
			filtered := selectedList[:0]
			for _, v := range selectedList {
				if v != removeValue {
					filtered = append(filtered, v)
				}
			}
			selectedList = filtered
		}

		// Build selected items for tags.
		var selectedItems []multiselect.SelectedItem
		for _, v := range selectedList {
			selectedItems = append(selectedItems, multiselect.SelectedItem{
				Value: v,
				Label: langLabel(v),
			})
		}

		// Filter results: match query, exclude already selected.
		var results []multiselect.Item
		for _, l := range mockLanguages {
			if selectedSet[l.Value] {
				continue
			}
			if q == "" || strings.Contains(strings.ToLower(l.Label), q) || strings.Contains(strings.ToLower(l.Description), q) {
				results = append(results, l)
			}
		}

		newCSV := strings.Join(selectedList, ",")
		sanitizedID := strings.ReplaceAll(id, "-", "_")

		sse := datastar.NewSSE(w, r)

		// Patch the selected signal so the client stays in sync.
		_ = sse.MarshalAndPatchSignals(map[string]any{
			sanitizedID: map[string]any{
				"selected": newCSV,
				"search":   "",
			},
		})

		// Patch tags and results.
		_ = ds.Send.Patch(sse, multiselect.Tags(multiselect.TagsProps{
			MultiSelectID: id,
			SearchURL:     searchURL,
			Items:         selectedItems,
		}))
		_ = ds.Send.Patch(sse, multiselect.Options(multiselect.OptionsProps{
			MultiSelectID: id,
			SearchURL:     searchURL,
			Items:         results,
		}))
	}
}
