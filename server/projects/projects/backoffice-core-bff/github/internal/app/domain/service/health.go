package service

import "github.com/pismo/backoffice-core-bff/internal/app/domain/entity"

type healthService struct {
	Message string
}

func (s *healthService) GetMessage() (*entity.Health, error) {
	return &entity.Health{Message: s.Message}, nil
}

func NewHealthService() HealthService {
	return &healthService{
		Message: "it works",
	}
}

