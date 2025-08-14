package order

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresOrderRepository struct {
	db *pgxpool.Pool
}

func NewPostgresOrderRepository(db *pgxpool.Pool) OrderRepository {
	return &PostgresOrderRepository{db: db}
}

// Save inserts a new order and its items into the database within a transaction.
func (r *PostgresOrderRepository) Save(ctx context.Context, order *Order) error {
	// Begin a transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	// Defer a rollback in case of panic or error.
	// If tx.Commit() is called, the rollback will be a no-op.
	defer tx.Rollback(ctx)

	// Insert into the 'orders' table
	orderQuery := "INSERT INTO orders (id, customer_id, status, created_at) VALUES ($1, $2, $3, $4)"

	_, err = tx.Exec(ctx, orderQuery, order.ID, order.CustomerID, order.Status, order.CreatedAt)
	if err != nil {
		return err
	}

	// Insert each item into the 'order_items' table
	for _, item := range order.Items {
		itemQuery := "INSERT INTO order_items (order_id, menu_item_id, quantity) VALUES ($1, $2, $3)"
		_, err = tx.Exec(ctx, itemQuery, order.ID, item.MenuItemID, item.Quantity)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *PostgresOrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*Order, error) {
	order := &Order{}
	orderQuery := "SELECT id, customer_id, status, created_at FROM orders WHERE id = $1"

	err := r.db.QueryRow(ctx, orderQuery, id).Scan(&order.ID, &order.CustomerID, &order.Status, &order.Status, &order.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	itemsQuery := "SELECT menu_item_id, quantity FROM order_items WHERE order_id = $1"
	rows, err := r.db.Query(ctx, itemsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []OrderItem
	for rows.Next() {
		var item OrderItem
		if err := rows.Scan(&item.MenuItemID, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	order.Items = items
	return order, nil
}
