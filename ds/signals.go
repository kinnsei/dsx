package ds

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/starfederation/datastar-go/datastar"
)

// ReadSignals reads namespaced signals from a Datastar request.
// The componentID is sanitized (hyphens → underscores) to match the JS namespace.
// dest must be a pointer to a struct with json tags matching the signal shape.
//
// Call this BEFORE datastar.NewSSE() — SSE creation consumes the request body.
//
//	var signals commandbar.CommandBarSignals
//	if err := ds.ReadSignals("my-bar", r, &signals); err != nil { ... }
//	input := signals.Text
func ReadSignals(componentID string, r *http.Request, dest any) error {
	sanitizedID := strings.ReplaceAll(componentID, "-", "_")
	wrapper := map[string]any{sanitizedID: dest}
	if err := datastar.ReadSignals(r, &wrapper); err != nil {
		return fmt.Errorf("read signals for %q: %w", componentID, err)
	}
	return nil
}
