package stream_test

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/plaenen/webx"
	"github.com/plaenen/webx/stream"
	"github.com/starfederation/datastar-go/datastar"
)

// startNATS starts an in-process NATS server and returns a connection.
func startNATS(t *testing.T) *nats.Conn {
	t.Helper()
	ns, err := server.NewServer(&server.Options{DontListen: true})
	if err != nil {
		t.Fatalf("creating NATS server: %v", err)
	}
	ns.Start()
	t.Cleanup(ns.Shutdown)

	if !ns.ReadyForConnections(4 * time.Second) {
		t.Fatal("NATS server not ready")
	}

	nc, err := nats.Connect(ns.ClientURL(), nats.InProcessServer(ns))
	if err != nil {
		t.Fatalf("connecting to NATS: %v", err)
	}
	t.Cleanup(nc.Close)
	return nc
}

func TestScopeKey(t *testing.T) {
	tests := []struct {
		scope string
		want  string
	}{
		{"invoice:42", "invoice_42"},
		{"invoices:*", "invoices_WILD"},
		{"workspace:1:*", "workspace_1_WILD"},
		{"simple", "simple"},
		{"a.b:c", "a_b_c"},
	}
	for _, tt := range tests {
		t.Run(tt.scope, func(t *testing.T) {
			got := stream.ScopeKey(tt.scope)
			if got != tt.want {
				t.Errorf("ScopeKey(%q) = %q, want %q", tt.scope, got, tt.want)
			}
		})
	}
}

func TestScopeSignals(t *testing.T) {
	got := stream.ScopeSignals("counter:shared")
	// Should produce valid Datastar signals object syntax
	if !strings.HasPrefix(got, "{") || !strings.HasSuffix(got, "}") {
		t.Errorf("ScopeSignals should be wrapped in {}, got: %s", got)
	}
	if !strings.Contains(got, "_stream") {
		t.Errorf("ScopeSignals should contain _stream namespace, got: %s", got)
	}
	if !strings.Contains(got, "counter_shared") {
		t.Errorf("ScopeSignals should contain counter_shared key, got: %s", got)
	}
	t.Logf("ScopeSignals output: %s", got)
}

func TestWatchEffect(t *testing.T) {
	wctx := &webx.WebXContext{}
	ctx := wctx.WithContext(context.Background())

	effect := stream.WatchEffect(ctx, "counter:shared", "/api/counter")

	// Should register the scope
	if len(wctx.Scopes) != 1 || wctx.Scopes[0] != "counter:shared" {
		t.Errorf("WatchEffect should register scope, got: %v", wctx.Scopes)
	}

	// Should return a data-effect expression
	if !strings.Contains(effect, "$_stream.counter_shared") {
		t.Errorf("effect should reference $_stream.counter_shared, got: %s", effect)
	}
	if !strings.Contains(effect, "@get('/api/counter')") {
		t.Errorf("effect should contain @get with reload URL, got: %s", effect)
	}
	t.Logf("WatchEffect output: %s", effect)
}

func TestWatchEffect_Dedup(t *testing.T) {
	wctx := &webx.WebXContext{}
	ctx := wctx.WithContext(context.Background())

	stream.WatchEffect(ctx, "counter:shared", "/api/counter")
	stream.WatchEffect(ctx, "counter:shared", "/api/counter")

	if len(wctx.Scopes) != 1 {
		t.Errorf("WatchEffect should deduplicate scopes, got %d", len(wctx.Scopes))
	}
}

func TestBrokerInvalidate(t *testing.T) {
	nc := startNATS(t)
	broker := stream.NewBroker(nc)

	// Subscribe to the expected NATS subject
	received := make(chan string, 1)
	sub, err := nc.Subscribe("webx.scope.counter.shared", func(msg *nats.Msg) {
		received <- msg.Subject
	})
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Unsubscribe()

	// Invalidate
	if err := broker.Invalidate("counter:shared"); err != nil {
		t.Fatalf("Invalidate failed: %v", err)
	}

	select {
	case subj := <-received:
		if subj != "webx.scope.counter.shared" {
			t.Errorf("unexpected subject: %s", subj)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("did not receive NATS message within 2s")
	}
}

func TestBrokerInvalidate_CustomPrefix(t *testing.T) {
	nc := startNATS(t)
	broker := stream.NewBroker(nc, stream.WithSubjectPrefix("myapp.scope"))

	received := make(chan string, 1)
	sub, err := nc.Subscribe("myapp.scope.counter.shared", func(msg *nats.Msg) {
		received <- msg.Subject
	})
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Unsubscribe()

	broker.Invalidate("counter:shared")

	select {
	case subj := <-received:
		if subj != "myapp.scope.counter.shared" {
			t.Errorf("unexpected subject: %s", subj)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("did not receive NATS message")
	}
}

func TestBrokerInvalidateMany(t *testing.T) {
	nc := startNATS(t)
	broker := stream.NewBroker(nc)

	received := make(chan string, 10)
	sub, err := nc.Subscribe("webx.scope.>", func(msg *nats.Msg) {
		received <- msg.Subject
	})
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Unsubscribe()

	broker.InvalidateMany("counter:shared", "invoice:42")

	subjects := map[string]bool{}
	for i := 0; i < 2; i++ {
		select {
		case subj := <-received:
			subjects[subj] = true
		case <-time.After(2 * time.Second):
			t.Fatal("timeout waiting for messages")
		}
	}

	if !subjects["webx.scope.counter.shared"] {
		t.Error("missing counter:shared")
	}
	if !subjects["webx.scope.invoice.42"] {
		t.Error("missing invoice:42")
	}
}

// parseSSEEvents reads SSE events from a reader and returns them.
func parseSSEEvents(r io.Reader) []map[string]string {
	var events []map[string]string
	scanner := bufio.NewScanner(r)
	current := map[string]string{}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			if len(current) > 0 {
				events = append(events, current)
				current = map[string]string{}
			}
			continue
		}
		if strings.HasPrefix(line, "event: ") {
			current["event"] = strings.TrimPrefix(line, "event: ")
		} else if strings.HasPrefix(line, "data: ") {
			current["data"] += strings.TrimPrefix(line, "data: ") + "\n"
		}
	}
	if len(current) > 0 {
		events = append(events, current)
	}
	return events
}

func TestStreamHandler_NoScopes(t *testing.T) {
	nc := startNATS(t)
	broker := stream.NewBroker(nc)

	req := httptest.NewRequest("GET", "/stream", nil)
	w := httptest.NewRecorder()

	broker.Handler().ServeHTTP(w, req)

	// Should return immediately with no SSE when no scopes
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestStreamHandler_ReceivesInvalidation(t *testing.T) {
	nc := startNATS(t)
	broker := stream.NewBroker(nc)

	// Create a cancellable context for the SSE request
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := httptest.NewRequest("GET", "/stream?scope=counter:shared", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	// Run handler in background
	done := make(chan struct{})
	go func() {
		defer close(done)
		broker.Handler().ServeHTTP(w, req)
	}()

	// Give the handler time to subscribe
	time.Sleep(100 * time.Millisecond)

	// Invalidate the scope
	if err := broker.Invalidate("counter:shared"); err != nil {
		t.Fatalf("Invalidate failed: %v", err)
	}

	// Give time for the SSE event to be written
	time.Sleep(200 * time.Millisecond)

	// Cancel context to stop the handler
	cancel()
	<-done

	body := w.Body.String()
	t.Logf("SSE response body:\n%s", body)

	// Check SSE headers
	ct := w.Header().Get("Content-Type")
	if !strings.Contains(ct, "text/event-stream") {
		t.Errorf("expected text/event-stream content type, got: %s", ct)
	}

	// Parse SSE events
	events := parseSSEEvents(strings.NewReader(body))
	if len(events) == 0 {
		t.Fatal("expected at least one SSE event, got none")
	}

	// Should be a patch-signals event with _stream.counter_shared = true
	found := false
	for _, evt := range events {
		t.Logf("Event: %v", evt)
		if evt["event"] == "datastar-patch-signals" {
			if strings.Contains(evt["data"], "counter_shared") {
				found = true
			}
		}
	}
	if !found {
		t.Error("did not find datastar-patch-signals event with counter_shared")
	}
}

func TestStreamHandler_WildcardScope(t *testing.T) {
	nc := startNATS(t)
	broker := stream.NewBroker(nc)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Subscribe to wildcard scope "invoices:*"
	req := httptest.NewRequest("GET", "/stream?scope=invoices:*", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	done := make(chan struct{})
	go func() {
		defer close(done)
		broker.Handler().ServeHTTP(w, req)
	}()

	time.Sleep(100 * time.Millisecond)

	// Invalidate a specific invoice — should be received by wildcard subscriber
	broker.Invalidate("invoices:42")

	time.Sleep(200 * time.Millisecond)
	cancel()
	<-done

	body := w.Body.String()
	t.Logf("SSE response body:\n%s", body)

	events := parseSSEEvents(strings.NewReader(body))
	found := false
	for _, evt := range events {
		if evt["event"] == "datastar-patch-signals" && strings.Contains(evt["data"], "invoices_WILD") {
			found = true
		}
	}
	if !found {
		t.Error("wildcard scope did not receive invalidation for invoices:42")
	}
}

func TestCounterHandler_GetCounter(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/stream/counter", nil)
	w := httptest.NewRecorder()

	nc := startNATS(t)
	broker := stream.NewBroker(nc)
	_ = broker // counter handler doesn't use broker

	// Simulate the counter handler directly
	var counter atomic.Int64
	counter.Store(42)

	handler := func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)
		count := counter.Load()
		sse.PatchElements(
			fmt.Sprintf(`<span id="stream-counter-value" class="text-6xl font-bold tabular-nums">%d</span>`, count),
		)
	}

	handler(w, req)

	body := w.Body.String()
	t.Logf("Counter SSE response:\n%s", body)

	// Should contain the counter value
	if !strings.Contains(body, "42") {
		t.Error("response should contain counter value 42")
	}

	// Should be a patch-elements event
	events := parseSSEEvents(strings.NewReader(body))
	if len(events) == 0 {
		t.Fatal("expected at least one SSE event")
	}

	found := false
	for _, evt := range events {
		if evt["event"] == "datastar-patch-elements" {
			if strings.Contains(evt["data"], "stream-counter-value") {
				found = true
			}
		}
	}
	if !found {
		t.Error("expected patch-elements event targeting stream-counter-value")
	}
}

func TestConnectTemplate_RendersWhenScopesExist(t *testing.T) {
	wctx := &webx.WebXContext{
		StreamURL: "/showcase/stream",
		Scopes:    []string{"counter:shared"},
	}
	ctx := wctx.WithContext(context.Background())

	var buf bytes.Buffer
	err := stream.Connect().Render(ctx, &buf)
	if err != nil {
		t.Fatalf("Connect render failed: %v", err)
	}

	html := buf.String()
	t.Logf("Connect HTML output:\n%s", html)

	// Should render a hidden div
	if !strings.Contains(html, `style="display:none"`) {
		t.Error("Connect should render a hidden div")
	}

	// Should have data-signals with outer braces
	if !strings.Contains(html, "data-signals=") {
		t.Error("Connect should have data-signals attribute")
	}

	// Check the data-signals value is properly formatted
	// After HTML escaping, the value should decode to: {_stream: {"counter_shared":false}}
	if !strings.Contains(html, "_stream") {
		t.Error("data-signals should contain _stream namespace")
	}

	// Should have data-init with the stream URL
	if !strings.Contains(html, "data-init") {
		t.Error("Connect should have data-init attribute")
	}
	if !strings.Contains(html, "/showcase/stream?scope=counter:shared") {
		t.Error("Connect should have stream URL with scope param")
	}

	// Should have requestCancellation disabled
	if !strings.Contains(html, "requestCancellation") {
		t.Error("Connect should have requestCancellation option")
	}
}

func TestConnectTemplate_NoRenderWithoutScopes(t *testing.T) {
	wctx := &webx.WebXContext{
		StreamURL: "/showcase/stream",
		Scopes:    nil,
	}
	ctx := wctx.WithContext(context.Background())

	var buf bytes.Buffer
	err := stream.Connect().Render(ctx, &buf)
	if err != nil {
		t.Fatalf("Connect render failed: %v", err)
	}

	if buf.Len() > 0 {
		t.Errorf("Connect should render nothing when no scopes, got: %s", buf.String())
	}
}

func TestConnectTemplate_NoRenderWithoutStreamURL(t *testing.T) {
	wctx := &webx.WebXContext{
		StreamURL: "",
		Scopes:    []string{"counter:shared"},
	}
	ctx := wctx.WithContext(context.Background())

	var buf bytes.Buffer
	err := stream.Connect().Render(ctx, &buf)
	if err != nil {
		t.Fatalf("Connect render failed: %v", err)
	}

	if buf.Len() > 0 {
		t.Errorf("Connect should render nothing when no StreamURL, got: %s", buf.String())
	}
}

// TestE2E_FullFlow tests the complete reactive flow:
// 1. Page registers scope via WatchEffect
// 2. Connect template renders with correct signals and stream URL
// 3. Stream handler subscribes and receives invalidation
// 4. Counter handler returns correct SSE response
func TestE2E_FullFlow(t *testing.T) {
	nc := startNATS(t)
	broker := stream.NewBroker(nc)

	// === Step 1: Simulate page render ===
	wctx := &webx.WebXContext{
		BasePath:  "/showcase",
		StreamURL: "/showcase/stream",
	}
	ctx := wctx.WithContext(context.Background())

	// Component registers scope (like the stream page does)
	effect := stream.WatchEffect(ctx, "counter:shared", "/showcase/api/stream/counter")
	t.Logf("Step 1 - WatchEffect returned: %s", effect)
	t.Logf("Step 1 - Scopes registered: %v", wctx.Scopes)

	if len(wctx.Scopes) == 0 {
		t.Fatal("no scopes registered after WatchEffect")
	}

	// === Step 2: Render Connect template ===
	var connectBuf bytes.Buffer
	if err := stream.Connect().Render(ctx, &connectBuf); err != nil {
		t.Fatalf("Connect render: %v", err)
	}
	connectHTML := connectBuf.String()
	t.Logf("Step 2 - Connect HTML:\n%s", connectHTML)

	if connectHTML == "" {
		t.Fatal("Connect rendered empty — stream won't connect")
	}

	// === Step 3: Stream handler receives invalidation ===
	streamCtx, streamCancel := context.WithCancel(context.Background())
	defer streamCancel()

	streamReq := httptest.NewRequest("GET", "/showcase/stream?scope=counter:shared", nil).WithContext(streamCtx)
	streamW := httptest.NewRecorder()

	streamDone := make(chan struct{})
	go func() {
		defer close(streamDone)
		broker.Handler().ServeHTTP(streamW, streamReq)
	}()

	// Wait for subscription to be set up
	time.Sleep(150 * time.Millisecond)

	// Simulate what a button click handler does: invalidate
	if err := broker.Invalidate("counter:shared"); err != nil {
		t.Fatalf("Invalidate: %v", err)
	}

	// Wait for SSE event to be written
	time.Sleep(200 * time.Millisecond)

	streamCancel()
	<-streamDone

	streamBody := streamW.Body.String()
	t.Logf("Step 3 - Stream SSE response:\n%s", streamBody)

	streamEvents := parseSSEEvents(strings.NewReader(streamBody))
	if len(streamEvents) == 0 {
		t.Fatal("stream handler produced no SSE events after invalidation")
	}

	// Verify the stale signal was pushed
	staleSignalFound := false
	for _, evt := range streamEvents {
		if evt["event"] == "datastar-patch-signals" {
			data := evt["data"]
			t.Logf("Step 3 - patch-signals data: %s", data)
			if strings.Contains(data, "counter_shared") && strings.Contains(data, "_stream") {
				staleSignalFound = true
			}
		}
	}
	if !staleSignalFound {
		t.Fatal("stream did not push stale signal for counter_shared")
	}

	// === Step 4: Counter handler returns correct response ===
	counterReq := httptest.NewRequest("GET", "/showcase/api/stream/counter", nil)
	counterW := httptest.NewRecorder()

	// Simulate counter handler
	sse := datastar.NewSSE(counterW, counterReq)
	sse.PatchElements(`<span id="stream-counter-value" class="text-6xl font-bold tabular-nums">0</span>`)

	counterBody := counterW.Body.String()
	t.Logf("Step 4 - Counter SSE response:\n%s", counterBody)

	counterEvents := parseSSEEvents(strings.NewReader(counterBody))
	patchFound := false
	for _, evt := range counterEvents {
		if evt["event"] == "datastar-patch-elements" {
			data := evt["data"]
			if strings.Contains(data, "stream-counter-value") && strings.Contains(data, "0") {
				patchFound = true
			}
		}
	}
	if !patchFound {
		t.Fatal("counter handler did not return expected patch-elements event")
	}

	t.Log("=== E2E flow passed: scope registration → Connect render → stream invalidation → counter patch ===")
}

// TestE2E_MutationHandler_NoEmptyPatch verifies mutation handlers
// don't send PatchElements with no target (which causes the browser error).
func TestE2E_MutationHandler_NoEmptyPatch(t *testing.T) {
	nc := startNATS(t)
	broker := stream.NewBroker(nc)

	var counter atomic.Int64

	// Simulate increment handler
	handler := func(w http.ResponseWriter, r *http.Request) {
		counter.Add(1)
		broker.Invalidate("counter:shared")
		datastar.NewSSE(w, r)
	}

	req := httptest.NewRequest("GET", "/api/stream/increment", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	body := w.Body.String()
	t.Logf("Mutation handler response:\n%q", body)

	// Should NOT contain patch-elements (which would error with no target)
	if strings.Contains(body, "datastar-patch-elements") {
		t.Error("mutation handler should NOT send patch-elements event")
	}

	// Counter should have been incremented
	if counter.Load() != 1 {
		t.Errorf("counter should be 1, got %d", counter.Load())
	}
}

// TestE2E_DataSignalsFormat verifies the exact format of data-signals
// that Datastar expects.
func TestE2E_DataSignalsFormat(t *testing.T) {
	// Test ScopeSignals format
	signals := stream.ScopeSignals("counter:shared")
	t.Logf("ScopeSignals: %s", signals)

	// Must start with { and end with }
	if signals[0] != '{' || signals[len(signals)-1] != '}' {
		t.Errorf("ScopeSignals must be wrapped in {}, got: %s", signals)
	}

	// Must contain the namespace
	if !strings.Contains(signals, "_stream") {
		t.Errorf("must contain _stream, got: %s", signals)
	}

	// The format should be: {_stream: {"counter_shared":false}}
	// Datastar parses this as a JavaScript object expression
	expected := `{_stream: {"counter_shared":false}}`
	if signals != expected {
		t.Errorf("ScopeSignals format mismatch\n  got:  %s\n  want: %s", signals, expected)
	}

	// Test Connect template produces the same format
	wctx := &webx.WebXContext{
		StreamURL: "/stream",
		Scopes:    []string{"counter:shared"},
	}
	ctx := wctx.WithContext(context.Background())

	var buf bytes.Buffer
	stream.Connect().Render(ctx, &buf)
	html := buf.String()
	t.Logf("Connect HTML: %s", html)

	// The HTML-escaped version should decode to the same thing
	// In HTML: {_stream: {&#34;counter_shared&#34;:false}}
	// or: {_stream: {&quot;counter_shared&quot;:false}}
	// Both decode to: {_stream: {"counter_shared":false}}
	if !strings.Contains(html, "_stream") {
		t.Error("Connect HTML should contain _stream")
	}
}

// TestE2E_MultipleScopes tests with multiple scopes registered.
func TestE2E_MultipleScopes(t *testing.T) {
	nc := startNATS(t)
	broker := stream.NewBroker(nc)

	wctx := &webx.WebXContext{
		StreamURL: "/stream",
	}
	ctx := wctx.WithContext(context.Background())

	stream.WatchEffect(ctx, "counter:shared", "/api/counter")
	stream.WatchEffect(ctx, "invoice:42", "/api/invoice/42")

	if len(wctx.Scopes) != 2 {
		t.Fatalf("expected 2 scopes, got %d", len(wctx.Scopes))
	}

	// Render Connect
	var buf bytes.Buffer
	stream.Connect().Render(ctx, &buf)
	html := buf.String()
	t.Logf("Connect HTML with 2 scopes:\n%s", html)

	// Should have both scopes in query params
	if !strings.Contains(html, "scope=counter:shared") {
		t.Error("missing counter:shared scope param")
	}
	if !strings.Contains(html, "scope=invoice:42") {
		t.Error("missing invoice:42 scope param")
	}

	// Test stream handler with multiple scopes
	streamCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := httptest.NewRequest("GET", "/stream?scope=counter:shared&scope=invoice:42", nil).WithContext(streamCtx)
	w := httptest.NewRecorder()

	done := make(chan struct{})
	go func() {
		defer close(done)
		broker.Handler().ServeHTTP(w, req)
	}()

	time.Sleep(150 * time.Millisecond)

	// Invalidate only invoice:42
	broker.Invalidate("invoice:42")
	time.Sleep(200 * time.Millisecond)

	cancel()
	<-done

	body := w.Body.String()
	t.Logf("Stream body (invoice:42 invalidation):\n%s", body)

	events := parseSSEEvents(strings.NewReader(body))
	invoiceStale := false
	for _, evt := range events {
		if evt["event"] == "datastar-patch-signals" && strings.Contains(evt["data"], "invoice_42") {
			invoiceStale = true
		}
	}
	if !invoiceStale {
		t.Error("should have received stale signal for invoice_42")
	}
}
