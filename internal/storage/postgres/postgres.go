package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"pvz_service/internal/domain"
)

type PVZRepository struct {
	db *sqlx.DB
}

func NewPVZRepository(db *sqlx.DB) *PVZRepository {
	return &PVZRepository{db: db}
}

func (r *PVZRepository) Create(ctx context.Context, pvz domain.PVZ) (*domain.PVZ, error) {
	query := `
		INSERT INTO pvz (id, registration_date, city)
		VALUES (:id, :registration_date, :city)
		RETURNING id, registration_date, city
	`

	// Генерация UUID правильного типа
	pvz.ID = uuid.New()
	pvz.RegistrationDate = time.Now().UTC()

	_, err := r.db.NamedExecContext(ctx, query, pvz)
	if err != nil {
		return nil, err
	}

	return &pvz, nil
}

func (r *PVZRepository) GetAll(ctx context.Context) ([]domain.PVZ, error) {
	query := `SELECT * FROM pvz`
	var pvzs []domain.PVZ
	err := r.db.SelectContext(ctx, &pvzs, query)
	return pvzs, err
}

func (r *PVZRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.PVZ, error) {
	query := `SELECT * FROM pvz WHERE id = $1`
	var pvz domain.PVZ
	err := r.db.GetContext(ctx, &pvz, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return &pvz, err
}

func NewPostgresDB(connStr string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}
	
	// Настройка маппинга для работы с UUID
	db.MapperFunc(func(s string) string { return s })
	return db, nil
}

