// Package pubsub defines the interface for publish/subscribe messaging
// used by the stream package. Adapters for NATS, Redis, and in-process
// channels are provided in sub-packages.
package pubsub

// PubSub is the minimal publish/subscribe interface required by the
// stream Broker. Topics use dot-separated segments with optional
// wildcards: * matches one segment, > matches the rest.
type PubSub interface {
	Publish(topic string, data []byte) error
	Subscribe(topic string, handler func(data []byte)) (Subscription, error)
	Close() error
}

// Subscription represents an active topic subscription.
type Subscription interface {
	Unsubscribe() error
}
