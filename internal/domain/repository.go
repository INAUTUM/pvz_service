package domain

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type PVZRepository interface {
	Create(ctx context.Context, pvz PVZ) (*PVZ, error)
	GetAll(ctx context.Context) ([]PVZ, error)
	GetByID(ctx context.Context, id uuid.UUID) (*PVZ, error) // Исправлен тип ID
}

type ReceptionRepository interface {
	Create(ctx context.Context, reception Reception) (*Reception, error)
	GetActive(ctx context.Context, pvzID uuid.UUID) (*Reception, error)
	Close(ctx context.Context, id uuid.UUID) error
}

type ProductRepository interface {
	Create(ctx context.Context, product Product) error
	DeleteLast(ctx context.Context, receptionID uuid.UUID) error // Добавлен метод
}

