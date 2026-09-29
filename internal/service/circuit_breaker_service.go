package service

import (
	"context"

	"smart-gateway/internal/clients"
	"smart-gateway/internal/model"
	"smart-gateway/internal/repository"
)

type CircuitBreakerService struct {
	repo *repository.CircuitBreakerRepo
}

func NewCircuitBreakerService(db *clients.MysqlDB) *CircuitBreakerService {
	return &CircuitBreakerService{repo: repository.NewCircuitBreakerRepo(db)}
}

func (s *CircuitBreakerService) List(ctx context.Context, dto model.CircuitBreakerDto) ([]model.CircuitBreaker, int64, error) {
	return s.repo.FindAll(ctx, dto)
}

func (s *CircuitBreakerService) Get(ctx context.Context, id uint64) (*model.CircuitBreaker, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *CircuitBreakerService) Create(ctx context.Context, cb *model.CircuitBreaker) error {
	return s.repo.Create(ctx, cb)
}

func (s *CircuitBreakerService) Update(ctx context.Context, cb *model.CircuitBreaker) error {
	existing, err := s.repo.FindByID(ctx, cb.ID)
	if err != nil {
		return err
	}
	cb.CreatedAt = existing.CreatedAt
	return s.repo.Update(ctx, cb)
}

func (s *CircuitBreakerService) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}
