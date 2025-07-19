package eventbus

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// RedisEventBus is an implementation of the EventBus interface that uses Redis Pub/Sub.
type RedisEventBus struct {
	client *redis.Client
}

// NewRedisEventBus creates a new RedisEventBus.
func NewRedisEventBus(client *redis.Client) EventBus {
	return &RedisEventBus{client: client}
}

// Publish sends an event to a Redis channel.
func (b *RedisEventBus) Publish(ctx context.Context, event Event) error {
	// First, we convert our event struct into JSON.
	// THis is called "marshaling"
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Then, we publish the JSON payload to a Redis channel.
	// The channel name is the event's topic.
	return b.client.Publish(ctx, event.Topic(), payload).Err()
}

// Subscribe is more complex for an external broker like Redis.
// It typically runs in a separate, long-running process or goroutine.
func (b *RedisEventBus) Subscribe(topic string, handler Handler) error {
	// This method does't fit the model of a long-running subsciber worker.
	// We will handle subscriptions differently for Redis.
	return fmt.Errorf("Subscribe is not implemented for RedisEventBus; subscriptions should be handled by a dedicated worker process")
}
