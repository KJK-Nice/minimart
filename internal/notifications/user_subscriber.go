package notifications

import (
	"context"
	"fmt"
	"log/slog"
	"minimart/internal/shared/eventbus"
	"minimart/internal/user"
)

// UserSubscriber is a dedicated subscriber for user-related events.
type UserSubscriber struct {
	logger *slog.Logger
}

// NewUserSubscriber creates a new instance of UserSubscriber.
func NewUserSubscriber(logger *slog.Logger) *UserSubscriber {
	return &UserSubscriber{logger: logger}
}

// HandleUserCreatedEvent is the handler for the UserCreatedEvent.
func (s *UserSubscriber) HandleUserCreatedEvent(ctx context.Context, event eventbus.Event) error {
	// Type assert the event to the specifiic UserCreatedEvent
	userEvent, ok := event.(user.UserCreatedEvent)
	if !ok {
		s.logger.Error(
			"Unexpected event type received",
			"module", "notifications",
			"topic", event.Topic(),
			"event_type", fmt.Sprintf("%T", event),
		)
		return nil
	}

	s.logger.Info(
		"New user created",
		"module", "notifications",
		"user_id", userEvent.UserID,
		"email", userEvent.Email,
	)

	return nil
}
