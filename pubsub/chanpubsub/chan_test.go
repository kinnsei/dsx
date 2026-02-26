package chanpubsub

import (
	"sync"
	"testing"
	"time"
)

func TestExactMatch(t *testing.T) {
	ps := New()
	defer ps.Close()

	got := make(chan string, 1)
	ps.Subscribe("foo.bar", func(data []byte) { got <- string(data) })

	ps.Publish("foo.bar", []byte("hello"))

	select {
	case v := <-got:
		if v != "hello" {
			t.Errorf("got %q, want %q", v, "hello")
		}
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
}

func TestNoMatch(t *testing.T) {
	ps := New()
	defer ps.Close()

	got := make(chan struct{}, 1)
	ps.Subscribe("foo.bar", func([]byte) { got <- struct{}{} })

	ps.Publish("foo.baz", nil)

	select {
	case <-got:
		t.Fatal("should not have matched")
	case <-time.After(100 * time.Millisecond):
		// ok
	}
}

func TestWildcardStar(t *testing.T) {
	ps := New()
	defer ps.Close()

	got := make(chan string, 2)
	ps.Subscribe("foo.*", func(data []byte) { got <- string(data) })

	ps.Publish("foo.bar", []byte("a"))
	ps.Publish("foo.baz", []byte("b"))
	ps.Publish("foo.bar.deep", []byte("nope")) // should not match

	var results []string
	for i := 0; i < 2; i++ {
		select {
		case v := <-got:
			results = append(results, v)
		case <-time.After(time.Second):
			t.Fatal("timeout")
		}
	}

	select {
	case <-got:
		t.Fatal("deep topic should not match single wildcard")
	case <-time.After(100 * time.Millisecond):
	}

	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestWildcardGt(t *testing.T) {
	ps := New()
	defer ps.Close()

	got := make(chan string, 3)
	ps.Subscribe("foo.>", func(data []byte) { got <- string(data) })

	ps.Publish("foo.bar", []byte("a"))
	ps.Publish("foo.bar.baz", []byte("b"))
	ps.Publish("foo.x.y.z", []byte("c"))

	for i := 0; i < 3; i++ {
		select {
		case <-got:
		case <-time.After(time.Second):
			t.Fatalf("timeout on message %d", i)
		}
	}
}

func TestFanOut(t *testing.T) {
	ps := New()
	defer ps.Close()

	var mu sync.Mutex
	counts := map[string]int{}

	ps.Subscribe("topic", func([]byte) {
		mu.Lock()
		counts["a"]++
		mu.Unlock()
	})
	ps.Subscribe("topic", func([]byte) {
		mu.Lock()
		counts["b"]++
		mu.Unlock()
	})

	ps.Publish("topic", nil)
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if counts["a"] != 1 || counts["b"] != 1 {
		t.Errorf("expected both handlers called once, got %v", counts)
	}
}

func TestUnsubscribe(t *testing.T) {
	ps := New()
	defer ps.Close()

	got := make(chan struct{}, 1)
	sub, _ := ps.Subscribe("topic", func([]byte) { got <- struct{}{} })

	sub.Unsubscribe()

	ps.Publish("topic", nil)

	select {
	case <-got:
		t.Fatal("should not receive after unsubscribe")
	case <-time.After(100 * time.Millisecond):
	}
}

func TestUnsubscribeIdempotent(t *testing.T) {
	ps := New()
	defer ps.Close()

	sub, _ := ps.Subscribe("topic", func([]byte) {})
	if err := sub.Unsubscribe(); err != nil {
		t.Fatal(err)
	}
	if err := sub.Unsubscribe(); err != nil {
		t.Fatal(err)
	}
}

func TestCloseStopsDelivery(t *testing.T) {
	ps := New()

	got := make(chan struct{}, 1)
	ps.Subscribe("topic", func([]byte) { got <- struct{}{} })
	ps.Close()

	ps.Publish("topic", nil)

	select {
	case <-got:
		t.Fatal("should not deliver after close")
	case <-time.After(100 * time.Millisecond):
	}
}

func TestMatchTopic(t *testing.T) {
	tests := []struct {
		pattern string
		topic   string
		want    bool
	}{
		{"a.b", "a.b", true},
		{"a.b", "a.c", false},
		{"a.*", "a.b", true},
		{"a.*", "a.b.c", false},
		{"a.>", "a.b", true},
		{"a.>", "a.b.c", true},
		{"a.>", "a.b.c.d", true},
		{">", "a", true},
		{">", "a.b", true},
		{"a.*.c", "a.b.c", true},
		{"a.*.c", "a.b.d", false},
		{"a", "a.b", false},
		{"a.b", "a", false},
	}
	for _, tt := range tests {
		t.Run(tt.pattern+"_"+tt.topic, func(t *testing.T) {
			if got := matchTopic(tt.pattern, tt.topic); got != tt.want {
				t.Errorf("matchTopic(%q, %q) = %v, want %v", tt.pattern, tt.topic, got, tt.want)
			}
		})
	}
}
