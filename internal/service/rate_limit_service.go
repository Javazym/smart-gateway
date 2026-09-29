package service

import (
	"context"

	"smart-gateway/internal/clients"
	"smart-gateway/internal/model"
	"smart-gateway/internal/repository"
)

type RateLimitService struct {
	repo *repository.RateLimitPolicyRepo
}

func NewRateLimitService(db *clients.MysqlDB) *RateLimitService {
	return &RateLimitService{repo: repository.NewRateLimitPolicyRepo(db)}
}

func (s *RateLimitService) List(ctx context.Context, dto model.RateLimitPolicyDto) ([]model.RateLimitPolicy, int64, error) {
	return s.repo.FindAll(ctx, dto)
}

func (s *RateLimitService) Get(ctx context.Context, id uint64) (*model.RateLimitPolicy, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *RateLimitService) Create(ctx context.Context, policy *model.RateLimitPolicy) error {
	return s.repo.Create(ctx, policy)
}

func (s *RateLimitService) Update(ctx context.Context, policy *model.RateLimitPolicy) error {
	existing, err := s.repo.FindByID(ctx, policy.ID)
	if err != nil {
		return err
	}
	policy.CreatedAt = existing.CreatedAt
	return s.repo.Update(ctx, policy)
}

func (s *RateLimitService) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}
