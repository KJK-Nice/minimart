package merchant

import "github.com/google/uuid"

type Merchant struct {
	ID          uuid.UUID
	Name        string
	Description string
	IsActive    bool
}

func NewMerchant(name, description string) *Merchant {
	return &Merchant{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		IsActive:    true,
	}
}
