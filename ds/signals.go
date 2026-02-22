package ds

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/valyala/bytebufferpool"
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

	// Read the raw request body (same as datastar.ReadSignals does).
	var raw []byte
	if r.Method == "GET" {
		dsJSON := r.URL.Query().Get("datastar")
		if dsJSON == "" {
			return nil
		}
		raw = []byte(dsJSON)
	} else {
		buf := bytebufferpool.Get()
		defer bytebufferpool.Put(buf)
		if _, err := buf.ReadFrom(r.Body); err != nil {
			return fmt.Errorf("read signals for %q: read body: %w", componentID, err)
		}
		raw = buf.Bytes()
	}

	// Decode the top-level JSON object into raw messages per key.
	var top map[string]json.RawMessage
	if err := json.Unmarshal(raw, &top); err != nil {
		return fmt.Errorf("read signals for %q: unmarshal top-level: %w", componentID, err)
	}

	// Extract and unmarshal the namespaced portion into dest.
	nsRaw, ok := top[sanitizedID]
	if !ok {
		return nil // namespace not present — leave dest at zero values
	}
	if err := json.Unmarshal(nsRaw, dest); err != nil {
		return fmt.Errorf("read signals for %q: unmarshal namespace: %w", componentID, err)
	}
	return nil
}
