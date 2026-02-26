// Package pubsubtest provides a conformance test suite for PubSub
// implementations. Call RunSuite from each adapter's test file.
package pubsubtest

import (
	"sync"
	"testing"
	"time"

	"github.com/plaenen/webx/pubsub"
)

// RunSuite runs the full PubSub conformance test suite against the given
// implementation.
func RunSuite(t *testing.T, ps pubsub.PubSub) {
	t.Helper()

	t.Run("ExactPublishSubscribe", func(t *testing.T) {
		got := make(chan string, 1)
		sub, err := ps.Subscribe("suite.exact.topic", func(data []byte) {
			got <- string(data)
		})
		if err != nil {
			t.Fatal(err)
		}
		defer sub.Unsubscribe()

		time.Sleep(100 * time.Millisecond)

		if err := ps.Publish("suite.exact.topic", []byte("hello")); err != nil {
			t.Fatal(err)
		}

		select {
		case v := <-got:
			if v != "hello" {
				t.Errorf("got %q, want %q", v, "hello")
			}
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for message")
		}
	})

	t.Run("WildcardStar", func(t *testing.T) {
		got := make(chan string, 2)
		sub, err := ps.Subscribe("suite.wild.*", func(data []byte) {
			got <- string(data)
		})
		if err != nil {
			t.Fatal(err)
		}
		defer sub.Unsubscribe()

		time.Sleep(100 * time.Millisecond)

		ps.Publish("suite.wild.one", []byte("a"))
		ps.Publish("suite.wild.two", []byte("b"))

		for range 2 {
			select {
			case <-got:
			case <-time.After(5 * time.Second):
				t.Fatal("timeout")
			}
		}
	})

	t.Run("WildcardGt", func(t *testing.T) {
		got := make(chan string, 3)
		sub, err := ps.Subscribe("suite.gt.>", func(data []byte) {
			got <- string(data)
		})
		if err != nil {
			t.Fatal(err)
		}
		defer sub.Unsubscribe()

		time.Sleep(100 * time.Millisecond)

		ps.Publish("suite.gt.a", []byte("1"))
		ps.Publish("suite.gt.a.b", []byte("2"))
		ps.Publish("suite.gt.a.b.c", []byte("3"))

		for range 3 {
			select {
			case <-got:
			case <-time.After(5 * time.Second):
				t.Fatal("timeout")
			}
		}
	})

	t.Run("NoMatchDifferentTopic", func(t *testing.T) {
		got := make(chan struct{}, 1)
		sub, err := ps.Subscribe("suite.nomatch.a", func([]byte) {
			got <- struct{}{}
		})
		if err != nil {
			t.Fatal(err)
		}
		defer sub.Unsubscribe()

		time.Sleep(100 * time.Millisecond)

		ps.Publish("suite.nomatch.b", nil)

		select {
		case <-got:
			t.Fatal("should not have received message on different topic")
		case <-time.After(500 * time.Millisecond):
		}
	})

	t.Run("FanOut", func(t *testing.T) {
		var mu sync.Mutex
		counts := map[string]int{}

		sub1, _ := ps.Subscribe("suite.fanout.topic", func([]byte) {
			mu.Lock()
			counts["a"]++
			mu.Unlock()
		})
		defer sub1.Unsubscribe()

		sub2, _ := ps.Subscribe("suite.fanout.topic", func([]byte) {
			mu.Lock()
			counts["b"]++
			mu.Unlock()
		})
		defer sub2.Unsubscribe()

		time.Sleep(100 * time.Millisecond)

		ps.Publish("suite.fanout.topic", nil)

		time.Sleep(500 * time.Millisecond)

		mu.Lock()
		defer mu.Unlock()
		if counts["a"] != 1 || counts["b"] != 1 {
			t.Errorf("expected both handlers called once, got %v", counts)
		}
	})

	t.Run("Unsubscribe", func(t *testing.T) {
		got := make(chan struct{}, 1)
		sub, _ := ps.Subscribe("suite.unsub.topic", func([]byte) {
			got <- struct{}{}
		})

		time.Sleep(100 * time.Millisecond)
		sub.Unsubscribe()
		time.Sleep(100 * time.Millisecond)

		ps.Publish("suite.unsub.topic", nil)

		select {
		case <-got:
			t.Fatal("should not receive after unsubscribe")
		case <-time.After(500 * time.Millisecond):
		}
	})
}
