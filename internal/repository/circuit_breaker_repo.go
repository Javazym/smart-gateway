package repository

import (
	"context"

	"smart-gateway/internal/clients"
	"smart-gateway/internal/model"
)

type CircuitBreakerRepo struct {
	db *clients.MysqlDB
}

func NewCircuitBreakerRepo(db *clients.MysqlDB) *CircuitBreakerRepo {
	return &CircuitBreakerRepo{db: db}
}

func (circuitBreakerRepo *CircuitBreakerRepo) Create(ctx context.Context, cb *model.CircuitBreaker) error {
	return circuitBreakerRepo.db.WithContext(ctx).Create(cb).Error
}

func (circuitBreakerRepo *CircuitBreakerRepo) Update(ctx context.Context, cb *model.CircuitBreaker) error {
	return circuitBreakerRepo.db.WithContext(ctx).Save(cb).Error
}

func (circuitBreakerRepo *CircuitBreakerRepo) Delete(ctx context.Context, id uint64) error {
	return circuitBreakerRepo.db.WithContext(ctx).Delete(&model.CircuitBreaker{}, id).Error
}

func (circuitBreakerRepo *CircuitBreakerRepo) FindByID(ctx context.Context, id uint64) (*model.CircuitBreaker, error) {
	var cb model.CircuitBreaker
	if err := circuitBreakerRepo.db.WithContext(ctx).Where("id = ?", id).First(&cb).Error; err != nil {
		return nil, err
	}
	return &cb, nil
}

func (circuitBreakerRepo *CircuitBreakerRepo) FindAll(ctx context.Context, dto model.CircuitBreakerDto) ([]model.CircuitBreaker, int64, error) {
	var list []model.CircuitBreaker
	var total int64
	q := circuitBreakerRepo.db.WithContext(ctx).Model(&model.CircuitBreaker{})
	if dto.Status != nil {
		q = q.Where("status = ?", *dto.Status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Offset((dto.Page - 1) * dto.Size).Limit(dto.Size).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
