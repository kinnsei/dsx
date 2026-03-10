// Package chanpubsub provides an in-process PubSub implementation backed by
// Go channels. It supports dot-separated topics with wildcards (* = one
// segment, > = rest). Ideal for development and testing — zero external deps.
package chanpubsub

import (
	"strings"
	"sync"
	"sync/atomic"

	"github.com/laenen-partners/dsx/pubsub"
)

// ChanPubSub is an in-process fan-out pub/sub with wildcard topic matching.
type ChanPubSub struct {
	mu     sync.RWMutex
	subs   map[uint64]*sub
	nextID atomic.Uint64
	closed atomic.Bool
}

type sub struct {
	pattern string
	handler func([]byte)
}

// New creates a new in-process PubSub.
func New() *ChanPubSub {
	return &ChanPubSub{subs: make(map[uint64]*sub)}
}

// Publish sends data to all subscribers whose pattern matches the topic.
// Handlers are called in separate goroutines.
func (c *ChanPubSub) Publish(topic string, data []byte) error {
	if c.closed.Load() {
		return nil
	}

	c.mu.RLock()
	matches := make([]func([]byte), 0, len(c.subs))
	for _, s := range c.subs {
		if matchTopic(s.pattern, topic) {
			matches = append(matches, s.handler)
		}
	}
	c.mu.RUnlock()

	for _, h := range matches {
		// Copy data to avoid races if the publisher reuses the slice.
		cp := make([]byte, len(data))
		copy(cp, data)
		go h(cp)
	}
	return nil
}

// Subscribe registers a handler for topics matching the given pattern.
func (c *ChanPubSub) Subscribe(topic string, handler func([]byte)) (pubsub.Subscription, error) {
	id := c.nextID.Add(1)
	c.mu.Lock()
	c.subs[id] = &sub{pattern: topic, handler: handler}
	c.mu.Unlock()
	return &chanSub{ps: c, id: id}, nil
}

// Close removes all subscriptions.
func (c *ChanPubSub) Close() error {
	c.closed.Store(true)
	c.mu.Lock()
	c.subs = make(map[uint64]*sub)
	c.mu.Unlock()
	return nil
}

type chanSub struct {
	ps   *ChanPubSub
	id   uint64
	once sync.Once
}

func (s *chanSub) Unsubscribe() error {
	s.once.Do(func() {
		s.ps.mu.Lock()
		delete(s.ps.subs, s.id)
		s.ps.mu.Unlock()
	})
	return nil
}

// matchTopic checks if a concrete topic matches a subscription pattern.
// Segments are dot-separated. * matches exactly one segment, > matches
// one or more remaining segments (must be last token).
func matchTopic(pattern, topic string) bool {
	patParts := strings.Split(pattern, ".")
	topParts := strings.Split(topic, ".")

	pi, ti := 0, 0
	for pi < len(patParts) && ti < len(topParts) {
		pp := patParts[pi]
		if pp == ">" {
			return true // matches rest
		}
		if pp != "*" && pp != topParts[ti] {
			return false
		}
		pi++
		ti++
	}
	return pi == len(patParts) && ti == len(topParts)
}
