package menu

import (
	"context"

	"github.com/google/uuid"
)

// MenuUsecase defines the interface for menu-related business logic.
type MenuUsecase interface {
	CreateMenuItem(ctx context.Context, merchantID uuid.UUID, name, description string, price int) (*MenuItem, error)
	GetMenuForMerchant(ctx context.Context, merchantID uuid.UUID) ([]*MenuItem, error)
}

type menuUsecase struct {
	repo MenuRepository
}

// NewMenuUsecase creates a new instance of MenuUsecase.
func NewMenuUsecase(repo MenuRepository) MenuUsecase {
	return &menuUsecase{
		repo: repo,
	}
}

func (u *menuUsecase) CreateMenuItem(ctx context.Context, merchantID uuid.UUID, name, description string, price int) (*MenuItem, error) {
	item := &MenuItem{
		ID:          uuid.New(),
		MerchantID:  merchantID,
		Name:        name,
		Description: description,
		Price:       price,
		InStock:     true, // New items are in stock by default
	}

	if err := u.repo.Save(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (u *menuUsecase) GetMenuForMerchant(ctx context.Context, merchantID uuid.UUID) ([]*MenuItem, error) {
	return u.repo.GetByMerchantID(ctx, merchantID)
}
