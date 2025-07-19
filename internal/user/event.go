package user

import "time"

const UserCreatedTopic = "user.created"

type UserCreatedEvent struct {
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (e UserCreatedEvent) Topic() string {
	return UserCreatedTopic
}
