package service

import (
	"context"

	"smart-gateway/internal/clients"
	"smart-gateway/internal/model"
	"smart-gateway/internal/repository"
)

type AuditLogService struct {
	repo *repository.AuditLogRepo
}

func NewAuditLogService(db *clients.MysqlDB) *AuditLogService {
	return &AuditLogService{repo: repository.NewAuditLogRepo(db)}
}

func (s *AuditLogService) List(ctx context.Context, dto model.AuditLogDto) ([]model.AuditLog, int64, error) {
	return s.repo.FindAll(ctx, dto)
}

func (s *AuditLogService) Create(ctx context.Context, log *model.AuditLog) error {
	return s.repo.Create(ctx, log)
}
