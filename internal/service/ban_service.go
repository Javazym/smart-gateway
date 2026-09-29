package service

import (
	"context"

	"smart-gateway/internal/clients"
	"smart-gateway/internal/model"
	"smart-gateway/internal/repository"
)

type BanService struct {
	repo *repository.BanRepo
}

func NewBanService(db *clients.MysqlDB) *BanService {
	return &BanService{repo: repository.NewBanRepo(db)}
}

func (s *BanService) List(ctx context.Context, dto model.BanDto) ([]model.Ban, int64, error) {
	return s.repo.FindAll(ctx, dto)
}

func (s *BanService) Get(ctx context.Context, id uint64) (*model.Ban, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *BanService) Create(ctx context.Context, ban *model.Ban) error {
	return s.repo.Create(ctx, ban)
}

func (s *BanService) Unban(ctx context.Context, id uint64) error {
	return s.repo.Unban(ctx, id)
}

func (s *BanService) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}

// IsBanned 判断来源 IP 或设备指纹当前是否被封禁。
func (s *BanService) IsBanned(ctx context.Context, ip, fingerprintID string) (bool, error) {
	for _, target := range []string{ip, fingerprintID} {
		if target == "" {
			continue
		}
		bans, err := s.repo.FindActiveByTarget(ctx, target)
		if err != nil {
			return false, err
		}
		if len(bans) > 0 {
			return true, nil
		}
	}
	return false, nil
}
