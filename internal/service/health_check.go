package service

import (
	v1 "bubble/api/healthcheck/v1"
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

type HealthCheckService struct {
	v1.UnimplementedHealthCheckServiceServer
}

func NewHealthCheckService() *HealthCheckService {
	return &HealthCheckService{}
}

func (s *HealthCheckService) HealthCheck(ctx context.Context, _ *emptypb.Empty) (*v1.HealthCheckResponse, error) {
	return &v1.HealthCheckResponse{
		Status: v1.HealthCheckResponse_SERVING,
	}, nil
}
