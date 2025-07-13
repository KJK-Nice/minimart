package menu

import "github.com/google/uuid"

// MenuItem represents a product or service that can be ordered.
type MenuItem struct {
	ID          uuid.UUID
	MerchantID  uuid.UUID
	Name        string
	Description string
	Price       int
	InStock     bool
}
