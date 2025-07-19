package eventbus

import (
	"context"
	"log"
	"sync"
)

// InMemoryEventBus is a simple in-memory implementation of the EventBus interface.
// It uses a map to store event handlers and a mutex for concurrent access.
type InMemoryEventBus struct {
	mu          sync.RWMutex
	subscribers map[string][]Handler
}

// NewInMemoryEventBus create a new InMemoryEventBus.
func NewInMemoryEventBus() EventBus {
	return &InMemoryEventBus{
		subscribers: make(map[string][]Handler),
	}
}

// Publish sends an event to all registered handlers for its topic.
func (b *InMemoryEventBus) Publish(ctx context.Context, event Event) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	handlers, found := b.subscribers[event.Topic()]
	if !found {
		return nil // No subscribers for this topic
	}

	for _, handler := range handlers {
		// Run each handler in its own goroutine
		go func(h Handler) {
			// In a real system, you'd want more robust error handling,
			// perhaps logging the error to a central service.
			if err := h(ctx, event); err != nil {
				log.Printf("Error handling event %s: %v", event.Topic(), err)
			}
		}(handler)
	}
	return nil
}

// Subscribe adds a new handler for a given topic.
func (b *InMemoryEventBus) Subscribe(topic string, handler Handler) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.subscribers[topic] = append(b.subscribers[topic], handler)
	return nil
}
