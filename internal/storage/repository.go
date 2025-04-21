package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"pvz_service/internal/domain"
)

type PVZRepository struct {
	db *sqlx.DB
}

func NewPVZRepository(db *sqlx.DB) *PVZRepository {
	return &PVZRepository{db: db}
}

func (r *PVZRepository) Create(ctx context.Context, pvz domain.PVZ) (*domain.PVZ, error) {
	if !isValidCity(pvz.City) {
		return nil, errors.New("invalid city")
	}

	query := `
		INSERT INTO pvz (id, registration_date, city)
		VALUES ($1, $2, $3)
		RETURNING id, registration_date, city
	`

	pvz.ID = uuid.New()
	pvz.RegistrationDate = time.Now().UTC()

	_, err := r.db.ExecContext(ctx, query, pvz.ID, pvz.RegistrationDate, pvz.City)
	if err != nil {
		return nil, err
	}

	return &pvz, nil
}

func isValidCity(city string) bool {
	switch city {
	case "Москва", "Санкт-Петербург", "Казань":
		return true
	}
	return false
}

func (r *PVZRepository) GetActiveReception(ctx context.Context, pvzID uuid.UUID) (*domain.Reception, error) {
	query := `
		SELECT id, pvz_id, start_time, end_time, status
		FROM receptions
		WHERE pvz_id = $1 AND status = 'in_progress'
		ORDER BY start_time DESC
		LIMIT 1
	`

	var reception domain.Reception
	err := r.db.GetContext(ctx, &reception, query, pvzID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &reception, nil
}

func (r *PVZRepository) AddProduct(ctx context.Context, product domain.Product) error {
	query := `
		INSERT INTO products (id, reception_id, added_at, type)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(ctx, query, 
		uuid.New(), 
		product.ReceptionID, 
		product.AddedAt.UTC(), 
		product.Type,
	)
	return err
}