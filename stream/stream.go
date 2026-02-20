// Package stream provides a reactive SSE stream backed by NATS pub/sub.
//
// Components register scopes during render via [WatchEffect]. The stream
// handler subscribes to NATS subjects for those scopes and pushes stale
// signals when invalidations occur. Components watch the stale signal
// via data-effect and reload themselves.
//
// Scopes use colon-separated naming: "invoice:42", "invoices:*", "workspace:1:*".
// Wildcards map to NATS wildcards (* = single level, > via multi-level scope).
package stream

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/plaenen/webx"
	"github.com/starfederation/datastar-go/datastar"
)

const (
	// SignalNamespace is the Datastar signal namespace for stream stale flags.
	// The underscore prefix makes it local-only (never sent to backend).
	SignalNamespace = "_stream"

	defaultSubjectPrefix = "webx.scope"
)

// Broker wraps a NATS connection and provides publish/subscribe for scope
// invalidation. One Broker per application.
type Broker struct {
	conn   *nats.Conn
	prefix string
}

// Option configures the Broker.
type Option func(*Broker)

// WithSubjectPrefix overrides the default NATS subject prefix ("webx.scope").
func WithSubjectPrefix(prefix string) Option {
	return func(b *Broker) { b.prefix = prefix }
}

// NewBroker creates a Broker from an existing NATS connection.
// Use embedded NATS or connect to an external NATS server.
func NewBroker(conn *nats.Conn, opts ...Option) *Broker {
	b := &Broker{conn: conn, prefix: defaultSubjectPrefix}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

// Invalidate publishes an invalidation message for the given scope.
// Call this after mutations to notify all connected browsers.
//
//	broker.Invalidate("invoice:42")
func (b *Broker) Invalidate(scope string) error {
	subject := scopeToSubject(b.prefix, scope)
	if err := b.conn.Publish(subject, nil); err != nil {
		return fmt.Errorf("publishing invalidation for %q: %w", scope, err)
	}
	return nil
}

// InvalidateMany publishes invalidation for multiple scopes at once.
func (b *Broker) InvalidateMany(scopes ...string) error {
	for _, scope := range scopes {
		if err := b.Invalidate(scope); err != nil {
			return err
		}
	}
	return nil
}

// Handler returns an http.HandlerFunc that serves the persistent SSE stream.
// It reads scopes from the "scope" query parameter and subscribes to NATS
// subjects for each scope (supporting exact and wildcard patterns).
func (b *Broker) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		scopes := r.URL.Query()["scope"]
		if len(scopes) == 0 {
			return
		}

		sse := datastar.NewSSE(w, r)
		ctx := r.Context()

		// Subscribe to each scope's NATS subject.
		staleC := make(chan string, 64)
		var subs []*nats.Subscription

		for _, scope := range scopes {
			subject := scopeToSubject(b.prefix, scope)
			scopeKey := ScopeKey(scope)

			sub, err := b.conn.Subscribe(subject, func(msg *nats.Msg) {
				select {
				case staleC <- scopeKey:
				default:
					// Drop if channel full (backpressure).
				}
			})
			if err != nil {
				slog.Error("stream: subscribe failed", "scope", scope, "subject", subject, "error", err)
				continue
			}
			subs = append(subs, sub)
		}

		defer func() {
			for _, sub := range subs {
				_ = sub.Unsubscribe()
			}
		}()

		// Event loop: wait for NATS messages or client disconnect.
		for {
			select {
			case <-ctx.Done():
				return
			case key := <-staleC:
				err := sse.MarshalAndPatchSignals(map[string]any{
					SignalNamespace: map[string]any{key: true},
				})
				if err != nil {
					return
				}
			}
		}
	}
}

// ScopeKey converts a scope string to a safe signal property name.
// This is the key within the _stream signal namespace.
//
//	ScopeKey("invoice:42")  → "invoice_42"
//	ScopeKey("invoices:*")  → "invoices_WILD"
//	ScopeKey("workspace:1:*") → "workspace_1_WILD"
func ScopeKey(scope string) string {
	s := strings.ReplaceAll(scope, ":", "_")
	s = strings.ReplaceAll(s, ".", "_")
	s = strings.ReplaceAll(s, "*", "WILD")
	return s
}

// WatchEffect registers a scope on the context and returns a data-effect
// expression string that auto-reloads when the scope goes stale.
//
//	stream.WatchEffect(ctx, "invoice:42", "/showcase/api/invoice/42")
//	// registers scope, returns: "if($_stream.invoice_42) { $_stream.invoice_42 = false; @get('/showcase/api/invoice/42') }"
func WatchEffect(ctx context.Context, scope string, reloadURL string) string {
	wctx := webx.FromContext(ctx)
	wctx.WatchScope(scope)

	key := ScopeKey(scope)
	signal := fmt.Sprintf("$%s.%s", SignalNamespace, key)
	return fmt.Sprintf("if(%s) { %s = false; @get('%s') }", signal, signal, reloadURL)
}

// InitScope pushes a PatchSignals to initialize a stale signal for a scope
// that wasn't known at initial render (e.g. new rows from infinite scroll).
func InitScope(sse *datastar.ServerSentEventGenerator, scope string) error {
	key := ScopeKey(scope)
	return sse.MarshalAndPatchSignals(map[string]any{
		SignalNamespace: map[string]any{key: false},
	})
}

// ScopeSignals returns a data-signals value that initializes the stale flags
// for the given scopes. Place this on the component element so the signals
// exist before data-effect runs.
//
//	stream.ScopeSignals("counter:shared") → "_stream: {\"counter_shared\":false}"
func ScopeSignals(scopes ...string) string {
	m := make(map[string]any, len(scopes))
	for _, s := range scopes {
		m[ScopeKey(s)] = false
	}
	j, _ := json.Marshal(m)
	return "{" + SignalNamespace + ": " + string(j) + "}"
}

// scopeToSubject converts a scope string to a NATS subject.
// Colons become dots, * stays as NATS wildcard.
//
//	"invoice:42"      → "webx.scope.invoice.42"
//	"invoices:*"      → "webx.scope.invoices.*"
//	"workspace:1:*"   → "webx.scope.workspace.1.*"
func scopeToSubject(prefix, scope string) string {
	safe := strings.ReplaceAll(scope, ":", ".")
	return prefix + "." + safe
}
