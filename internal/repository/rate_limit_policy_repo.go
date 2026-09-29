package repository

import (
	"context"

	"smart-gateway/internal/clients"
	"smart-gateway/internal/model"
)

type RateLimitPolicyRepo struct {
	db *clients.MysqlDB
}

func NewRateLimitPolicyRepo(db *clients.MysqlDB) *RateLimitPolicyRepo {
	return &RateLimitPolicyRepo{db: db}
}

func (rateLimitRepo *RateLimitPolicyRepo) Create(ctx context.Context, policy *model.RateLimitPolicy) error {
	return rateLimitRepo.db.WithContext(ctx).Create(policy).Error
}

func (rateLimitRepo *RateLimitPolicyRepo) Update(ctx context.Context, policy *model.RateLimitPolicy) error {
	return rateLimitRepo.db.WithContext(ctx).Save(policy).Error
}

func (rateLimitRepo *RateLimitPolicyRepo) Delete(ctx context.Context, id uint64) error {
	return rateLimitRepo.db.WithContext(ctx).Delete(&model.RateLimitPolicy{}, id).Error
}

func (rateLimitRepo *RateLimitPolicyRepo) FindByID(ctx context.Context, id uint64) (*model.RateLimitPolicy, error) {
	var policy model.RateLimitPolicy
	if err := rateLimitRepo.db.WithContext(ctx).Where("id = ?", id).First(&policy).Error; err != nil {
		return nil, err
	}
	return &policy, nil
}

func (rateLimitRepo *RateLimitPolicyRepo) FindAll(ctx context.Context, dto model.RateLimitPolicyDto) ([]model.RateLimitPolicy, int64, error) {
	var list []model.RateLimitPolicy
	var total int64
	q := rateLimitRepo.db.WithContext(ctx).Model(&model.RateLimitPolicy{})
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
