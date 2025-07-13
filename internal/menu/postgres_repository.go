package menu

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresRepository is the PostgreSQL implmentation of the MenuRepository.
type PostgresRepository struct {
	db *pgxpool.Pool
}

// NewPostgresMenuRepository creates a new PostgresRepository.
func NewPostgresMenuRepository(db *pgxpool.Pool) MenuRepository {
	return &PostgresRepository{
		db: db,
	}
}

// Save inserts a new menu item into the database.
func (r *PostgresRepository) Save(ctx context.Context, item *MenuItem) error {
	query := `
		INSERT INTO menu_items (id, merchant_id, name, description, price, in_stock)
		VALUES ($1, $2, $3, $4, $5, $6);
	`
	_, err := r.db.Exec(ctx, query, item.ID, item.MerchantID, item.Name, item.Description, item.Price, item.InStock)
	return err
}

// GetByMerchantID retrieves all menu items for a specific merchant.
func (r *PostgresRepository) GetByMerchantID(ctx context.Context, merchantID uuid.UUID) ([]*MenuItem, error) {
	query := `
		SELECT id, merchant_id, name, description, price, in_stock
		FROM menu_items
		WHERE merchant_id = $1
		ORDER BY name;
	`
	rows, err := r.db.Query(ctx, query, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*MenuItem
	for rows.Next() {
		item := &MenuItem{}
		err := rows.Scan(
			&item.ID,
			&item.MerchantID,
			&item.Name,
			&item.Description,
			&item.Price,
			&item.InStock,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	// rows.Err() checks for any errors that may have occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
