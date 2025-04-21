package service

import (
	"context"
	"errors"

	"pvz_service/internal/domain"
)

type PVZService struct {
	repo domain.PVZRepository
}

func NewPVZService(repo domain.PVZRepository) *PVZService {
	return &PVZService{repo: repo}
}

func (s *PVZService) CreatePVZ(ctx context.Context, pvz domain.PVZ) (*domain.PVZ, error) {
	if !isValidCity(pvz.City) {
		return nil, errors.New("invalid city")
	}
	return s.repo.Create(ctx, pvz)
}

func isValidCity(city string) bool {
	switch city {
	case "Москва", "Санкт-Петербург", "Казань":
		return true
	}
	return false
}

func (s *PVZService) GetPVZs(ctx context.Context, filter domain.PVZFilter) ([]domain.PVZ, error) {
	return s.repo.GetAll(ctx)
}
