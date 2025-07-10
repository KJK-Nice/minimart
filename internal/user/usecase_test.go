package user

import (
	"context"
	"errors"
	"testing"
)

type mockUserRepository struct {
	errToReturn error
}

func (m *mockUserRepository) Save(ctx context.Context, user *User) error {
	if m.errToReturn != nil {
		return m.errToReturn
	}
	return nil
}

func TestUserUseCase_RegisterUser(t *testing.T) {
	t.Run("should register a user succsessfully", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		mockRepo := &mockUserRepository{}
		userUsecase := NewUserUsecase(mockRepo)

		// Act
		userName := "John Wick"
		userEmail := "john.wick@example.com"
		createdUser, err := userUsecase.RegisterUser(ctx, userName, userEmail)

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

	t.Run("should return an error when repository fails", func(t *testing.T) {
		// Arrange
		expectedError := errors.New("database is down")
		mockRepo := &mockUserRepository{errToReturn: expectedError}
		userUsecase := NewUserUsecase(mockRepo)
		ctx := context.Background()

		// Act
		createdUser, err := userUsecase.RegisterUser(ctx, "Jane Doe", "jane.doe@example.com")

		// Assert
		if err == nil {
			t.Fatal("expected an error, but got nil")
		}
		if !errors.Is(err, expectedError) {
			t.Errorf("expected error to be %q, but got %q", expectedError.Error(), err.Error())
		}
		if createdUser != nil {
			t.Error("expected nil user, but got a user")
		}
	})
}
