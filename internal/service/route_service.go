package service

import (
	"context"

	"smart-gateway/internal/clients"
	"smart-gateway/internal/model"
	"smart-gateway/internal/repository"
)

type RouteService struct {
	repo *repository.RouteRepo
}

func NewRouteService(db *clients.MysqlDB) *RouteService {
	return &RouteService{repo: repository.NewRouteRepo(db)}
}

func (s *RouteService) List(ctx context.Context, dto model.RouteDto) ([]model.Route, int64, error) {
	return s.repo.FindAll(ctx, dto)
}

func (s *RouteService) Get(ctx context.Context, id uint64) (*model.Route, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *RouteService) Create(ctx context.Context, route *model.Route) error {
	return s.repo.Create(ctx, route)
}

func (s *RouteService) Update(ctx context.Context, route *model.Route) error {
	existing, err := s.repo.FindByID(ctx, route.ID)
	if err != nil {
		return err
	}
	route.CreatedAt = existing.CreatedAt
	return s.repo.Update(ctx, route)
}

func (s *RouteService) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}
