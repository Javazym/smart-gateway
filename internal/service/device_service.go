package service

import (
	"context"

	"smart-gateway/internal/clients"
	"smart-gateway/internal/model"
	"smart-gateway/internal/repository"
)

type DeviceService struct {
	repo *repository.DeviceRepo
}

func NewDeviceService(db *clients.MysqlDB) *DeviceService {
	return &DeviceService{repo: repository.NewDeviceRepo(db)}
}

// Register 登记一次设备出现：新设备入库，已有设备刷新最近出现信息。
func (s *DeviceService) Register(ctx context.Context, device *model.Device) error {
	return s.repo.Upsert(ctx, device)
}

func (s *DeviceService) List(ctx context.Context, dto model.DeviceDto) ([]model.Device, int64, error) {
	return s.repo.FindAll(ctx, dto)
}

func (s *DeviceService) Get(ctx context.Context, id uint64) (*model.Device, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *DeviceService) UpdateStatus(ctx context.Context, id uint64, status string) error {
	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *DeviceService) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}
