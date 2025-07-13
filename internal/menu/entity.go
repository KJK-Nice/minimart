package menu

import "github.com/google/uuid"

type MenuItem struct {
	ID         uuid.UUID
	MerchantID uuid.UUID
	Name       string
	Price      float64
	InStock    bool
}
