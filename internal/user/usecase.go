package user

import (
	"context"

	"github.com/google/uuid"
)

type UserUsecase interface {
	RegisterUser(ctx context.Context, name, email string) (*User, error)
}

type userUsecase struct {
	repo UserRepository
}

func NewUserUsecase(repo UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}

func (u *userUsecase) RegisterUser(ctx context.Context, name, email string) (*User, error) {
	user := &User{ID: uuid.New(), Name: name, Email: email}
	if err := u.repo.Save(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
