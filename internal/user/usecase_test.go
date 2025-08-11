package user

import (
	"context"
	"minimart/internal/shared/eventbus"
	"testing"
)

func TestUserUseCase_RegisterUser(t *testing.T) {
	eventBus := eventbus.NewInMemoryEventBus()
	userRepo := NewInMemoryUserRepository()

	t.Run("should register a user succsessfully", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userUsecase := NewUserUsecase(userRepo, eventBus, "test-secret")

		// Act
		userName := "John Wick"
		userEmail := "john.wick@example.com"
		userPassword := "password"
		createdUser, err := userUsecase.RegisterUser(ctx, userName, userEmail, userPassword)

		// Assert
		// 1. Check taht no error occurred
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// 2. Check that a user object was created
		if createdUser == nil {
			t.Fatal("expected a user to be created, but got nil")
		}

		// 3. Check that theruser's details are correct
		if createdUser.Name != userName {
			t.Errorf("expected user name to be %s, but got %s", userName, createdUser.Name)
		}

		if createdUser.Email != userEmail {
			t.Errorf("expected user email to be %s, but got %s", userEmail, createdUser.Email)
		}

		// 4. Check that the user has an ID
		if createdUser.ID.String() == "" {
			t.Error("expected user to have an ID, but it was empty")
		}
	})
}
