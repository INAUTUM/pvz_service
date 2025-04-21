package service

import (
	"context"

	"pvz_service/internal/domain"

	"github.com/google/uuid"
)

type CompositeService struct {
	authService      AuthService
	pvzService       *PVZService
	receptionService *ReceptionService
}

func NewCompositeService(
	auth AuthService,
	pvz *PVZService,
	reception *ReceptionService,
) *CompositeService {
	return &CompositeService{
		authService:      auth,
		pvzService:       pvz,
		receptionService: reception,
	}
}

func (s *CompositeService) AddProduct(
	ctx context.Context, 
	receptionID uuid.UUID,  // Добавляем недостающий параметр
	product domain.Product,
) error {
	return s.receptionService.AddProduct(
		ctx, 
		receptionID,  // Передаем receptionID
		product,
	)
}