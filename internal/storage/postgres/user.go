package postgres

import (
	"context"
	"database/sql"
	"errors"

	"pvz_service/internal/domain"

	"github.com/jmoiron/sqlx"
)

const createUserQuery = `
INSERT INTO users (email, password_hash, role) 
VALUES ($1, $2, $3)
ON CONFLICT (email) DO NOTHING
RETURNING id, created_at`

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
    err := r.db.QueryRowxContext(ctx, createUserQuery,
        user.Email,
        user.PasswordHash,
        user.Role,
    ).StructScan(user)

    if errors.Is(err, sql.ErrNoRows) {
        return nil, domain.ErrUserExists
    }
    
    return user, err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, email, password_hash, role, created_at FROM users WHERE email = $1`
	err := r.db.GetContext(ctx, &user, query, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	return &user, err
}