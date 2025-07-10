package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Save(ctx context.Context, user *User) error {
	query := `INSERT INTO users (id, name, email) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, user.ID, user.Name, user.Email)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // Unique violation
				return errors.New("user with this email already exists")
			}
		}
		return err
	}
	return nil
}
