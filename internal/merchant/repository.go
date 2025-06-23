package merchant

import (
	"context"

	"github.com/google/uuid"
)

type MerchantRepository interface {
	GetMerchantByID(ctx context.Context, id uuid.UUID) (*Merchant, error)
	Save(ctx context.Context, merchant *Merchant) error
	Create(ctx context.Context, name, description string) (*Merchant, error)
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

func (r *InMemoryMerchantRepository) GetMerchantByID(ctx context.Context, id uuid.UUID) (*Merchant, error) {
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

func (r *InMemoryMerchantRepository) Create(ctx context.Context, name, description string) (*Merchant, error) {
	id := uuid.New()
	merchant := &Merchant{
		ID:          id,
		Name:        name,
		Description: description,
		IsActive:    true,
	}

	r.merchants[id] = merchant
	return merchant, nil
}
