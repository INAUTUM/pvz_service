package grpc

import (
	"net"
	"pvz_service/pvz/pvz_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func StartGRPCServer(port string, pvzCtrl *PVZController) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	srv := grpc.NewServer()
	pvz_v1.RegisterPVZServiceServer(srv, pvzCtrl)
	
	// Health check
	grpc_health_v1.RegisterHealthServer(srv, health.NewServer())
	
	return srv.Serve(lis)
}