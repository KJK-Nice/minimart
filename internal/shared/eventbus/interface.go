package eventbus

import "context"

// Event represents a generic event that occurs in the system.
type Event interface {
	// Topic returns the name of ther event topic.
	Topic() string
}

// Handler defines the function signature for an event handler.
type Handler func(ctx context.Context, event Event) error

// EventBus defines the interface for a system-wide event bus.
type EventBus interface {
	// Publish sends an event to all subscribers of its topic.
	Publish(ctx context.Context, event Event) error

	// Subscribe registers a handler for a specific event topic.
	Subscribe(topic string, handler Handler) error
}
