package repository

import (
	"context"
	"time"

	"smart-gateway/internal/clients"
	"smart-gateway/internal/model"
)

type BanRepo struct {
	db *clients.MysqlDB
}

func NewBanRepo(db *clients.MysqlDB) *BanRepo {
	return &BanRepo{db: db}
}

func (banRepo *BanRepo) Create(ctx context.Context, ban *model.Ban) error {
	return banRepo.db.WithContext(ctx).Create(ban).Error
}

func (banRepo *BanRepo) Update(ctx context.Context, ban *model.Ban) error {
	return banRepo.db.WithContext(ctx).Save(ban).Error
}

func (banRepo *BanRepo) Delete(ctx context.Context, id uint64) error {
	return banRepo.db.WithContext(ctx).Delete(&model.Ban{}, id).Error
}

func (banRepo *BanRepo) FindByID(ctx context.Context, id uint64) (*model.Ban, error) {
	var ban model.Ban
	if err := banRepo.db.WithContext(ctx).Where("id = ?", id).First(&ban).Error; err != nil {
		return nil, err
	}
	return &ban, nil
}

func (banRepo *BanRepo) FindAll(ctx context.Context, dto model.BanDto) ([]model.Ban, int64, error) {
	var bans []model.Ban
	var total int64
	q := banRepo.db.WithContext(ctx).Model(&model.Ban{})
	if dto.Status != nil {
		q = q.Where("status = ?", *dto.Status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("banned_at DESC").
		Offset((dto.Page - 1) * dto.Size).
		Limit(dto.Size).
		Find(&bans).Error
	if err != nil {
		return nil, 0, err
	}
	return bans, total, nil
}

// Unban 解除封禁：置状态为 0 并记录解封时间。
func (banRepo *BanRepo) Unban(ctx context.Context, id uint64) error {
	now := time.Now()
	return banRepo.db.WithContext(ctx).Model(&model.Ban{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      0,
			"unbanned_at": now,
		}).Error
}

// FindActiveByTarget 查询指定目标当前生效的封禁记录，供网关请求拦截判断。
func (banRepo *BanRepo) FindActiveByTarget(ctx context.Context, targetValue string) ([]model.Ban, error) {
	var bans []model.Ban
	if err := banRepo.db.WithContext(ctx).
		Where("target_value = ? AND status = ?", targetValue, 1).
		Find(&bans).Error; err != nil {
		return nil, err
	}
	return bans, nil
}
