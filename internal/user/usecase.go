package user

import (
	"context"
	"errors"
	"minimart/internal/shared/eventbus"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// ErrInvalidCredentials is a specific error for login failures.
var ErrInvalidCredentials = errors.New("Invalid email or password")

type UserUsecase interface {
	RegisterUser(ctx context.Context, name, email string, password string) (*User, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type userUsecase struct {
	repo      UserRepository
	eventBus  eventbus.EventBus
	jwtSecret string
}

func NewUserUsecase(repo UserRepository, eventBus eventbus.EventBus, jwtSecret string) UserUsecase {
	return &userUsecase{
		repo:      repo,
		eventBus:  eventBus,
		jwtSecret: jwtSecret,
	}
}

func (u *userUsecase) RegisterUser(ctx context.Context, name, email, password string) (*User, error) {
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		Password:  string(hasedPassword),
		CreatedAt: time.Now(),
	}
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

// Login handles the user authentication and JWT generation.
func (u *userUsecase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		// Use a generic error to avoid revealing if the user exists.
		return "", ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	claims := jwt.MapClaims{
		"sub":   user.ID,
		"name":  user.Name,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create the token object
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret key
	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
