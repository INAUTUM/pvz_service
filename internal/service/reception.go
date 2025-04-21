package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"pvz_service/internal/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ReceptionService struct {
	receptionRepo domain.ReceptionRepository
	productRepo   domain.ProductRepository
}

func NewReceptionService(
	receptionRepo domain.ReceptionRepository,
	productRepo domain.ProductRepository,
) *ReceptionService {
	return &ReceptionService{
		receptionRepo: receptionRepo,
		productRepo:   productRepo,
	}
}

func (s *ReceptionService) StartReception(ctx context.Context, pvzID uuid.UUID) (*domain.Reception, error) {
	active, err := s.receptionRepo.GetActive(ctx, pvzID)
	if err != nil {
		return nil, err
	}
	if active != nil {
		return nil, errors.New("active reception already exists")
	}

	return s.receptionRepo.Create(ctx, domain.Reception{
		PVZID:     pvzID,
		StartTime: time.Now(),
		Status:    domain.ReceptionStatusInProgress,
	})
}

func (s *ReceptionService) AddProduct(ctx context.Context, receptionID uuid.UUID, product domain.Product) error {
	product.ReceptionID = receptionID
	return s.productRepo.Create(ctx, product)
}

func (s *ReceptionService) CloseReception(ctx context.Context, pvzID uuid.UUID) error {
	reception, err := s.receptionRepo.GetActive(ctx, pvzID)
	if err != nil {
		return err
	}
	if reception == nil {
		return errors.New("no active reception")
	}
	return s.receptionRepo.Close(ctx, reception.ID)
}

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
	
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &reception, err
}

func (r *ReceptionRepository) Close(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE receptions SET status = 'closed', end_time = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (s *ReceptionService) GetActiveReception(
	ctx context.Context, 
	pvzID uuid.UUID,
) (*domain.Reception, error) {
	return s.receptionRepo.GetActive(ctx, pvzID)
}

func (s *ReceptionService) DeleteLastProduct(
	ctx context.Context, 
	receptionID uuid.UUID,
) error {
	return s.productRepo.DeleteLast(ctx, receptionID)
}

// func (s *ReceptionService) CloseReception(ctx context.Context, pvzID uuid.UUID) error {
// 	// Реализация закрытия приемки
// 	return s.receptionRepo.Close(ctx, pvzID)
// }