package menu

import (
	"context"

	"github.com/google/uuid"
)

// MenuRepository defines the persistence interface for menu items
// This interface is pure persistence with no business logic
type MenuRepository interface {
	// Save stores or updates a menu item
	Save(ctx context.Context, item *MenuItem) error
	
	// FindByID retrieves a menu item by ID
	FindByID(ctx context.Context, id uuid.UUID) (*MenuItem, error)
	
	// FindByMerchantID retrieves all menu items for a merchant
	FindByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]*MenuItem, error)
	
	// FindAvailableByMerchantID retrieves only available menu items for a merchant
	FindAvailableByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]*MenuItem, error)
	
	// FindByIDs retrieves multiple menu items by their IDs
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*MenuItem, error)
	
	// Delete removes a menu item (soft delete by setting unavailable)
	Delete(ctx context.Context, id uuid.UUID) error
}
