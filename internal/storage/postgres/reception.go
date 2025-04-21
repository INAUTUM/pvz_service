package postgres

import (
	"context"

	"pvz_service/internal/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ReceptionRepository struct {
	db *sqlx.DB
}

func NewReceptionRepository(db *sqlx.DB) *ReceptionRepository {
	return &ReceptionRepository{db: db}
}

func (r *ReceptionRepository) Create(ctx context.Context, reception domain.Reception) (*domain.Reception, error) {
	query := `
		INSERT INTO receptions (id, pvz_id, start_time, status)
		VALUES (:id, :pvz_id, :start_time, :status)
		RETURNING id, pvz_id, start_time, status
	`
	
	reception.ID = uuid.New()
	_, err := r.db.NamedExecContext(ctx, query, reception)
	return &reception, err
}

func (r *ReceptionRepository) GetActive(ctx context.Context, pvzID uuid.UUID) (*domain.Reception, error) {
    query := `
        SELECT * FROM receptions 
        WHERE pvz_id = $1 AND status = 'in_progress'
        ORDER BY start_time DESC 
        LIMIT 1
    `
    var reception domain.Reception
    err := r.db.GetContext(ctx, &reception, query, pvzID)
    return &reception, err
}

func (r *ReceptionRepository) Close(ctx context.Context, id uuid.UUID) error {
    query := `UPDATE receptions SET status = 'closed', end_time = NOW() WHERE id = $1`
    _, err := r.db.ExecContext(ctx, query, id)
    return err
}