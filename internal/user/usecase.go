package user

import (
	"context"
	"minimart/internal/shared/eventbus"
	"time"

	"github.com/google/uuid"
)

type UserUsecase interface {
	RegisterUser(ctx context.Context, name, email string) (*User, error)
}

type userUsecase struct {
	repo     UserRepository
	eventBus eventbus.EventBus
}

func NewUserUsecase(repo UserRepository, eventBus eventbus.EventBus) UserUsecase {
	return &userUsecase{repo: repo, eventBus: eventBus}
}

func (u *userUsecase) RegisterUser(ctx context.Context, name, email string) (*User, error) {
	user := &User{ID: uuid.New(), Name: name, Email: email}
	if err := u.repo.Save(ctx, user); err != nil {
		return nil, err
	}

	event := UserCreatedEvent{
		UserID:    user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: time.Now(),
	}

	if err := u.eventBus.Publish(ctx, event); err != nil {
		// Depending on your design, you might want to handle the error differently
		return nil, err
	}
	return user, nil
}
