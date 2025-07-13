package merchant

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresMerchantRepository struct {
	db *pgxpool.Pool
}

func NewPostgresMerchantRepository(db *pgxpool.Pool) MerchantRepository {
	return &PostgresMerchantRepository{db: db}
}

func (r *PostgresMerchantRepository) Save(ctx context.Context, merchant *Merchant) error {
	query := `
		INSERT INTO merchants (id, name, description, is_active)
		VALUES ($1, $2, $3, $4);
	`
	_, err := r.db.Exec(ctx, query, merchant.ID, merchant.Name, merchant.Description, merchant.IsActive)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresMerchantRepository) GetByID(ctx context.Context, id uuid.UUID) (*Merchant, error) {
	query := `
		SELECT id, name, description, is_active
		FROM menu_items
		WHERE id = $1;
	`
	merchant := &Merchant{}

	err := r.db.QueryRow(ctx, query, id).Scan(&merchant.ID, &merchant.Name, &merchant.Description, &merchant.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	return merchant, nil
}
