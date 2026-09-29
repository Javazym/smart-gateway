package repository

import (
	"context"

	"smart-gateway/internal/clients"
	"smart-gateway/internal/model"
)

type AuditLogRepo struct {
	db *clients.MysqlDB
}

func NewAuditLogRepo(db *clients.MysqlDB) *AuditLogRepo {
	return &AuditLogRepo{db: db}
}

func (auditLogRepo *AuditLogRepo) Create(ctx context.Context, log *model.AuditLog) error {
	return auditLogRepo.db.WithContext(ctx).Create(log).Error
}

func (auditLogRepo *AuditLogRepo) FindAll(ctx context.Context, dto model.AuditLogDto) ([]model.AuditLog, int64, error) {
	var list []model.AuditLog
	var total int64
	q := auditLogRepo.db.WithContext(ctx).Model(&model.AuditLog{})
	if dto.Username != "" {
		q = q.Where("username = ?", dto.Username)
	}
	if dto.Action != "" {
		q = q.Where("action = ?", dto.Action)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("created_at DESC").
		Offset((dto.Page - 1) * dto.Size).
		Limit(dto.Size).
		Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
