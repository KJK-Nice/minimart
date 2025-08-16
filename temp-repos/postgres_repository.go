package merchant

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

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

// MerchantRecord represents the database structure for merchants
type MerchantRecord struct {
	ID              uuid.UUID `db:"id"`
	Name            string    `db:"name"`
	Description     string    `db:"description"`
	IsActive        bool      `db:"is_active"`
	OperatingHours  *string   `db:"operating_hours"` // JSON string
	PreparationTime int       `db:"preparation_time"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

type OperatingHoursRecord struct {
	OpenTime  int             `json:"open_time"`   // Hour of day (0-23)
	CloseTime int             `json:"close_time"`  // Hour of day (0-23)  
	DaysOpen  []time.Weekday  `json:"days_open"`
}

// Save stores or updates a merchant in the database
func (r *PostgresMerchantRepository) Save(ctx context.Context, merchant *Merchant) error {
	// Serialize operating hours to JSON
	var operatingHoursJSON *string
	if merchant.OperatingHours() != nil {
		hours := merchant.OperatingHours()
		record := OperatingHoursRecord{
			OpenTime:  int(hours.OpenTime / time.Hour),
			CloseTime: int(hours.CloseTime / time.Hour),
			DaysOpen:  hours.DaysOpen,
		}
		
		hoursBytes, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("failed to marshal operating hours: %w", err)
		}
		hoursStr := string(hoursBytes)
		operatingHoursJSON = &hoursStr
	}

	query := `
		INSERT INTO merchants (
			id, name, description, is_active, operating_hours,
			preparation_time, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			is_active = EXCLUDED.is_active,
			operating_hours = EXCLUDED.operating_hours,
			preparation_time = EXCLUDED.preparation_time,
			updated_at = EXCLUDED.updated_at`

	_, err := r.db.Exec(ctx, query,
		merchant.ID(),
		merchant.Name(),
		merchant.Description(),
		merchant.IsActive(),
		operatingHoursJSON,
		merchant.PreparationTime(),
		merchant.CreatedAt(),
		merchant.UpdatedAt())

	if err != nil {
		return fmt.Errorf("failed to save merchant: %w", err)
	}

	return nil
}

// FindByID retrieves a merchant by ID
func (r *PostgresMerchantRepository) FindByID(ctx context.Context, id uuid.UUID) (*Merchant, error) {
	query := `
		SELECT id, name, description, is_active, operating_hours,
			   preparation_time, created_at, updated_at
		FROM merchants WHERE id = $1`

	var record MerchantRecord
	var operatingHoursStr *string

	err := r.db.QueryRow(ctx, query, id).Scan(
		&record.ID, &record.Name, &record.Description, &record.IsActive,
		&operatingHoursStr, &record.PreparationTime,
		&record.CreatedAt, &record.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Merchant not found
		}
		return nil, fmt.Errorf("failed to query merchant: %w", err)
	}

	// Parse operating hours
	var operatingHours *OperatingHours
	if operatingHoursStr != nil {
		var hoursRecord OperatingHoursRecord
		err = json.Unmarshal([]byte(*operatingHoursStr), &hoursRecord)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal operating hours: %w", err)
		}
		
		hours, err := NewOperatingHours(hoursRecord.OpenTime, hoursRecord.CloseTime, hoursRecord.DaysOpen)
		if err != nil {
			return nil, fmt.Errorf("failed to create operating hours: %w", err)
		}
		operatingHours = &hours
	}

	// Reconstruct Merchant entity
	return ReconstructMerchant(
		record.ID,
		record.Name,
		record.Description,
		record.IsActive,
		operatingHours,
		record.PreparationTime,
		record.CreatedAt,
		record.UpdatedAt,
	)
}

// FindActive retrieves all active merchants
func (r *PostgresMerchantRepository) FindActive(ctx context.Context) ([]*Merchant, error) {
	query := `
		SELECT id, name, description, is_active, operating_hours,
			   preparation_time, created_at, updated_at
		FROM merchants 
		WHERE is_active = true 
		ORDER BY name`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query active merchants: %w", err)
	}
	defer rows.Close()

	var merchants []*Merchant
	for rows.Next() {
		var record MerchantRecord
		var operatingHoursStr *string

		err := rows.Scan(
			&record.ID, &record.Name, &record.Description, &record.IsActive,
			&operatingHoursStr, &record.PreparationTime,
			&record.CreatedAt, &record.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan merchant: %w", err)
		}

		// Parse operating hours
		var operatingHours *OperatingHours
		if operatingHoursStr != nil {
			var hoursRecord OperatingHoursRecord
			err = json.Unmarshal([]byte(*operatingHoursStr), &hoursRecord)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal operating hours: %w", err)
			}
			
			hours, err := NewOperatingHours(hoursRecord.OpenTime, hoursRecord.CloseTime, hoursRecord.DaysOpen)
			if err != nil {
				return nil, fmt.Errorf("failed to create operating hours: %w", err)
			}
			operatingHours = &hours
		}

		merchant, err := ReconstructMerchant(
			record.ID,
			record.Name,
			record.Description,
			record.IsActive,
			operatingHours,
			record.PreparationTime,
			record.CreatedAt,
			record.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to reconstruct merchant: %w", err)
		}

		merchants = append(merchants, merchant)
	}

	return merchants, nil
}

// Delete soft deletes a merchant by setting it as inactive
func (r *PostgresMerchantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE merchants 
		SET is_active = false, updated_at = $2 
		WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to delete merchant: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("merchant with id %s not found", id)
	}

	return nil
}
