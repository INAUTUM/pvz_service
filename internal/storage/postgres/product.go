// internal/storage/postgres/product.go
package postgres

import (
	"context"
	"time"

	"pvz_service/internal/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, product domain.Product) error {
	query := `
		INSERT INTO products (id, reception_id, added_at, type)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, query,
		uuid.New().String(),
		product.ReceptionID,
		time.Now(),
		product.Type,
	)
	return err
}

func (r *ProductRepository) DeleteLast(ctx context.Context, receptionID uuid.UUID) error {
    query := `
        DELETE FROM products
        WHERE id = (
            SELECT id FROM products 
            WHERE reception_id = $1 
            ORDER BY added_at DESC 
            LIMIT 1
        )`
    _, err := r.db.ExecContext(ctx, query, receptionID)
    return err
}

