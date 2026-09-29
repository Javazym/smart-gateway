package repository

import (
	"context"

	"smart-gateway/internal/clients"
	"smart-gateway/internal/model"
)

type RequestLogRepo struct {
	db *clients.MysqlDB
}

func NewRequestLogRepo(db *clients.MysqlDB) *RequestLogRepo {
	return &RequestLogRepo{db: db}
}

func (requestLogRepo *RequestLogRepo) Create(ctx context.Context, log *model.RequestLog) error {
	return requestLogRepo.db.WithContext(ctx).Create(log).Error
}

// CreateBatch 批量写入请求日志，用于网关异步落库。
func (requestLogRepo *RequestLogRepo) CreateBatch(ctx context.Context, logs []model.RequestLog) error {
	if len(logs) == 0 {
		return nil
	}
	return requestLogRepo.db.WithContext(ctx).Create(&logs).Error
}

func (requestLogRepo *RequestLogRepo) FindAll(ctx context.Context, dto model.RequestLogDto) ([]model.RequestLog, int64, error) {
	var list []model.RequestLog
	var total int64
	q := requestLogRepo.db.WithContext(ctx).Model(&model.RequestLog{})
	if dto.DeviceFingerprintID != "" {
		q = q.Where("device_fingerprint_id = ?", dto.DeviceFingerprintID)
	}
	if dto.UpstreamID > 0 {
		q = q.Where("upstream_id = ?", dto.UpstreamID)
	}
	if dto.RouteID > 0 {
		q = q.Where("route_id = ?", dto.RouteID)
	}
	if dto.SourceIP != "" {
		q = q.Where("source_ip = ?", dto.SourceIP)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("received_at DESC").
		Offset((dto.Page - 1) * dto.Size).
		Limit(dto.Size).
		Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
