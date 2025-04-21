package grpc

import (
	"context"

	"pvz_service/internal/domain"
	"pvz_service/pvz/pvz_v1"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PVZController struct {
	pvzRepo domain.PVZRepository
	pvz_v1.UnimplementedPVZServiceServer
}

func NewPVZController(pvzRepo domain.PVZRepository) *PVZController {
	return &PVZController{pvzRepo: pvzRepo}
}

func (c *PVZController) Register(s *grpc.Server) {
	pvz_v1.RegisterPVZServiceServer(s, c)
}

func (c *PVZController) GetPVZList(ctx context.Context, req *pvz_v1.GetPVZListRequest) (*pvz_v1.GetPVZListResponse, error) {
	pvzs, err := c.pvzRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	response := &pvz_v1.GetPVZListResponse{
		Pvzs: make([]*pvz_v1.PVZ, 0, len(pvzs)),
	}

	for _, p := range pvzs {
		response.Pvzs = append(response.Pvzs, &pvz_v1.PVZ{
			Id:               p.ID.String(),
			RegistrationDate: timestamppb.New(p.RegistrationDate),
			City:             p.City,
		})
	}

	return response, nil
}