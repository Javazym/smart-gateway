package repository

import (
	"context"

	"smart-gateway/internal/clients"
	"smart-gateway/internal/model"
)

type RouteRepo struct {
	db *clients.MysqlDB
}

func NewRouteRepo(db *clients.MysqlDB) *RouteRepo {
	return &RouteRepo{db: db}
}

func (routeRepo *RouteRepo) Create(ctx context.Context, route *model.Route) error {
	return routeRepo.db.WithContext(ctx).Create(route).Error
}

func (routeRepo *RouteRepo) Update(ctx context.Context, route *model.Route) error {
	return routeRepo.db.WithContext(ctx).Save(route).Error
}

func (routeRepo *RouteRepo) Delete(ctx context.Context, id uint64) error {
	return routeRepo.db.WithContext(ctx).Delete(&model.Route{}, id).Error
}

func (routeRepo *RouteRepo) FindByID(ctx context.Context, id uint64) (*model.Route, error) {
	var route model.Route
	if err := routeRepo.db.WithContext(ctx).Where("id = ?", id).First(&route).Error; err != nil {
		return nil, err
	}
	return &route, nil
}

func (routeRepo *RouteRepo) FindAll(ctx context.Context, dto model.RouteDto) ([]model.Route, int64, error) {
	var routes []model.Route
	var total int64
	q := routeRepo.db.WithContext(ctx).Model(&model.Route{})
	if dto.Status != nil {
		q = q.Where("status = ?", *dto.Status)
	}
	if dto.Name != "" {
		q = q.Where("name LIKE ?", "%"+dto.Name+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("priority ASC").
		Offset((dto.Page - 1) * dto.Size).
		Limit(dto.Size).
		Find(&routes).Error
	if err != nil {
		return nil, 0, err
	}
	return routes, total, nil
}

// FindAllEnabled 返回所有启用的路由，供网关运行时热加载路由表。
func (routeRepo *RouteRepo) FindAllEnabled(ctx context.Context) ([]model.Route, error) {
	var routes []model.Route
	if err := routeRepo.db.WithContext(ctx).
		Where("status = ?", 1).
		Order("priority ASC").
		Find(&routes).Error; err != nil {
		return nil, err
	}
	return routes, nil
}
