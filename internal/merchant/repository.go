package merchant

import (
	"context"

	"github.com/google/uuid"
)

// MerchantRepository defines the persistence interface for merchants
// This interface is pure persistence with no business logic
type MerchantRepository interface {
	// Save stores or updates a merchant
	Save(ctx context.Context, merchant *Merchant) error
	
	// FindByID retrieves a merchant by ID
	FindByID(ctx context.Context, id uuid.UUID) (*Merchant, error)
	
	// FindActive retrieves all active merchants
	FindActive(ctx context.Context) ([]*Merchant, error)
	
	// Delete removes a merchant (soft delete by deactivating)
	Delete(ctx context.Context, id uuid.UUID) error
}
