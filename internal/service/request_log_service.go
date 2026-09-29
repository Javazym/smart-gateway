package service

import (
	"context"

	"smart-gateway/internal/clients"
	"smart-gateway/internal/model"
	"smart-gateway/internal/repository"
)

type RequestLogService struct {
	repo *repository.RequestLogRepo
}

func NewRequestLogService(db *clients.MysqlDB) *RequestLogService {
	return &RequestLogService{repo: repository.NewRequestLogRepo(db)}
}

func (s *RequestLogService) List(ctx context.Context, dto model.RequestLogDto) ([]model.RequestLog, int64, error) {
	return s.repo.FindAll(ctx, dto)
}

func (s *RequestLogService) Create(ctx context.Context, log *model.RequestLog) error {
	return s.repo.Create(ctx, log)
}
