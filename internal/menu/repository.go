package menu

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

// MenuRepository defines the interface for intreacting with menu item storage.
type MenuRepository interface {
	Save(ctx context.Context, item *MenuItem) error
	GetByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]*MenuItem, error)
}

// InMemoryMenuRepository is a simple in-memory implementation of MenuRepository.
type InMemoryMenuRepository struct {
	mu    sync.RWMutex
	items map[uuid.UUID][]*MenuItem
}

func NewInMemoryMenuRepository() MenuRepository {
	return &InMemoryMenuRepository{
		items: make(map[uuid.UUID][]*MenuItem),
	}
}

func (r *InMemoryMenuRepository) Save(ctx context.Context, item *MenuItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.items[item.MerchantID] = append(r.items[item.MerchantID], item)
	return nil
}

func (r *InMemoryMenuRepository) GetByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]*MenuItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.items[merchantID], nil
}
