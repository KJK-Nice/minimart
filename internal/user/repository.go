package user

import (
	"context"
	"errors"
)

type UserRepository interface {
	Save(ctx context.Context, user *User) error
}

type InMemoryUserRepository struct {
	users map[string]*User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*User),
	}
}

func (r *InMemoryUserRepository) Save(ctx context.Context, user *User) error {
	if _, exists := r.users[user.Email]; exists {
		return errors.New("User with this email already exists")
	}
	r.users[user.Email] = user
	return nil
}
