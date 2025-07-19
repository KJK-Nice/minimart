package notifications

import (
	"context"
	"fmt"

	"minimart/internal/shared/eventbus"
	"minimart/internal/user"
)

// UserSubscriber is a dedicated subscriber for user-related events.
type UserSubscriber struct {
	// dependencies here. an email service, a logger, etc.
}

// NewUserSubscriber creates a new instance of UserSubscriber.
func NewUserSubscriber() *UserSubscriber {
	return &UserSubscriber{}
}

// HandleUserCreatedEvent is the handler for the UserCreatedEvent.
func (s *UserSubscriber) HandleUserCreatedEvent(ctx context.Context, event eventbus.Event) error {
	// Type assert the event to the specifiic UserCreatedEvent
	userEvent, ok := event.(user.UserCreatedEvent)
	if !ok {
		// This should not happen if the event bus is working correctly,
		// but it's good practice to handle it.
		return fmt.Errorf("unexpected event type: %T", event)
	}

	fmt.Printf("--> [Notifications] New user created: Name: %s, Email: %s\n", userEvent.Name, userEvent.Email)

	return nil
}
