package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type PVZ struct {
	ID               uuid.UUID
	RegistrationDate time.Time
	City             string
}

type Reception struct {
	ID        uuid.UUID      `db:"id"`
	PVZID     uuid.UUID      `db:"pvz_id"`
	StartTime time.Time      `db:"start_time"`
	EndTime   *time.Time     `db:"end_time"`
	Status    ReceptionStatus `db:"status"`
}

type Product struct {
	ID          uuid.UUID
	ReceptionID uuid.UUID
	AddedAt     time.Time
	Type        ProductType
}

type ReceptionStatus string

const (
	ReceptionStatusInProgress ReceptionStatus = "in_progress"
	ReceptionStatusClosed     ReceptionStatus = "closed"
)

type ProductType string

const (
	ProductTypeElectronics ProductType = "электроника"
	ProductTypeClothing    ProductType = "одежда"
	ProductTypeShoes       ProductType = "обувь"
)

func (pt ProductType) IsValid() bool {
	switch pt {
	case ProductTypeElectronics, ProductTypeClothing, ProductTypeShoes:
		return true
	}
	return false
}


func ParseProductType(s string) (ProductType, error) {
	switch s {
	case "электроника":
		return ProductTypeElectronics, nil
	case "одежда":
		return ProductTypeClothing, nil
	case "обувь":
		return ProductTypeShoes, nil
	default:
		return "", errors.New("invalid product type")
	}
}

type PVZFilter struct {
	StartDate time.Time
	EndDate   time.Time
	Page      int
	Limit     int
}

type UserRole string

const (
	UserRoleEmployee  UserRole = "employee"
	UserRoleModerator UserRole = "moderator"
	UserRoleClient    UserRole = "client"
)

func (r UserRole) IsValid() bool {
	switch r {
	case UserRoleEmployee, UserRoleModerator, UserRoleClient:
		return true
	}
	return false
}


type User struct {
    ID           uuid.UUID  `db:"id"`
    Email        string     `db:"email"`
    PasswordHash string     `db:"password_hash"`
    Role         UserRole   `db:"role"`
    CreatedAt    time.Time  `db:"created_at"`
}

func NewUUID() uuid.UUID {
	return uuid.New()
}