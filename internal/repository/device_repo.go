package repository

import (
	"context"

	"gorm.io/gorm/clause"
	"smart-gateway/internal/clients"
	"smart-gateway/internal/model"
)

type DeviceRepo struct {
	db *clients.MysqlDB
}

func NewDeviceRepo(db *clients.MysqlDB) *DeviceRepo {
	return &DeviceRepo{db: db}
}

func (deviceRepo *DeviceRepo) Create(ctx context.Context, device *model.Device) error {
	return deviceRepo.db.WithContext(ctx).Create(device).Error
}

func (deviceRepo *DeviceRepo) Update(ctx context.Context, device *model.Device) error {
	return deviceRepo.db.WithContext(ctx).Save(device).Error
}

func (deviceRepo *DeviceRepo) Delete(ctx context.Context, id uint64) error {
	return deviceRepo.db.WithContext(ctx).Delete(&model.Device{}, id).Error
}

func (deviceRepo *DeviceRepo) FindByID(ctx context.Context, id uint64) (*model.Device, error) {
	var device model.Device
	if err := deviceRepo.db.WithContext(ctx).Where("id = ?", id).First(&device).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

func (deviceRepo *DeviceRepo) FindByFingerprintID(ctx context.Context, fingerprintID string) (*model.Device, error) {
	var device model.Device
	if err := deviceRepo.db.WithContext(ctx).Where("fingerprint_id = ?", fingerprintID).First(&device).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

// Upsert 按 fingerprint_id 唯一键写入或更新设备：首次出现即插入，已存在则刷新最近出现信息。
// 用 OnConflict 避免并发首请求撞唯一索引；device_status / first_seen_at 不在此更新（分别由管理端与首见决定）。
func (deviceRepo *DeviceRepo) Upsert(ctx context.Context, device *model.Device) error {
	return deviceRepo.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "fingerprint_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"http_header_hash", "user_agent", "last_seen_ip", "last_seen_at", "updated_at"}),
	}).Create(device).Error
}

func (deviceRepo *DeviceRepo) FindAll(ctx context.Context, dto model.DeviceDto) ([]model.Device, int64, error) {
	var devices []model.Device
	var total int64
	q := deviceRepo.db.WithContext(ctx).Model(&model.Device{})
	if dto.Status != "" {
		q = q.Where("device_status = ?", dto.Status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("last_seen_at DESC").
		Offset((dto.Page - 1) * dto.Size).
		Limit(dto.Size).
		Find(&devices).Error
	if err != nil {
		return nil, 0, err
	}
	return devices, total, nil
}

func (deviceRepo *DeviceRepo) UpdateStatus(ctx context.Context, id uint64, status string) error {
	return deviceRepo.db.WithContext(ctx).Model(&model.Device{}).
		Where("id = ?", id).
		Update("device_status", status).Error
}
