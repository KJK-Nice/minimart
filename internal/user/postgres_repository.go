package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	query := `INSERT INTO users (id, name, email, password, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query, user.ID, user.Name, user.Email, user.Password, user.CreatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // Unique violation
				return errors.New("User with this email already exists")
			}
		}
		return err
	}
	return nil
}

// FindByID retrives a user from the database by their ID.
func (r *PostgresUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	query := `SELECT id, name, email, password, created_at FROM users WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("User not found")
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail retrives a user from the database by their email.
func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, name, email, password, created_at FROM users WHERE email = $1`
	row := r.db.QueryRow(ctx, query, email)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("User not found")
		}
		return nil, err
	}
	return &user, nil
}
