package merchant

import (
	"context"

	"github.com/google/uuid"
)

type MerchantRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Merchant, error)
	Save(ctx context.Context, merchant *Merchant) error
}

type InMemoryMerchantRepository struct {
	merchants map[uuid.UUID]*Merchant
}

func NewInMemoryMerchantRepository() *InMemoryMerchantRepository {
	merchant1 := &Merchant{
		ID:          uuid.New(),
		Name:        "Merchant 1",
		Description: "Merchant 1 description",
		IsActive:    false,
	}

	merchant2 := &Merchant{
		ID:          uuid.New(),
		Name:        "Merchant 2",
		Description: "Merchant 2 description",
		IsActive:    true,
	}

	return &InMemoryMerchantRepository{
		merchants: map[uuid.UUID]*Merchant{
			merchant1.ID: merchant1,
			merchant2.ID: merchant2,
		},
	}
}

func (r *InMemoryMerchantRepository) GetByID(ctx context.Context, id uuid.UUID) (*Merchant, error) {
	merchant, exists := r.merchants[id]
	if !exists {
		return nil, nil
	}
	return merchant, nil
}

func (r *InMemoryMerchantRepository) Save(ctx context.Context, merchant *Merchant) error {
	r.merchants[merchant.ID] = merchant
	return nil
}
