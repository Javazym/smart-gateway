package service

import (
	"context"

	"smart-gateway/internal/clients"
	"smart-gateway/internal/model"
	"smart-gateway/internal/repository"
)

type UpstreamService struct {
	repo *repository.UpstreamRepo
}

func NewUpstreamService(db *clients.MysqlDB) *UpstreamService {
	return &UpstreamService{repo: repository.NewUpstreamRepo(db)}
}

// 分组

func (s *UpstreamService) ListGroup(ctx context.Context, dto model.GroupDto) ([]model.UpstreamGroup, int64, error) {
	return s.repo.FindAddGroup(ctx, dto)
}

func (s *UpstreamService) GetGroup(ctx context.Context, id uint64) (*model.UpstreamGroup, error) {
	return s.repo.FindUpstreamGroupById(ctx, id)
}

func (s *UpstreamService) CreateGroup(ctx context.Context, group *model.UpstreamGroup) error {
	return s.repo.CreateGroup(ctx, group)
}

func (s *UpstreamService) UpdateGroup(ctx context.Context, group *model.UpstreamGroup) error {
	existing, err := s.repo.FindUpstreamGroupById(ctx, group.ID)
	if err != nil {
		return err
	}
	group.CreatedAt = existing.CreatedAt
	return s.repo.UpdateGroup(ctx, group)
}

func (s *UpstreamService) DeleteGroup(ctx context.Context, id uint64) error {
	return s.repo.DeleteGroup(ctx, id)
}

// 节点

func (s *UpstreamService) ListUpstream(ctx context.Context, groupID uint64) ([]model.Upstream, error) {
	return s.repo.FindUpstreamsByGroupID(ctx, groupID)
}

func (s *UpstreamService) GetUpstream(ctx context.Context, id uint64) (*model.Upstream, error) {
	return s.repo.FindUpstreamById(ctx, int64(id))
}

func (s *UpstreamService) CreateUpstream(ctx context.Context, upstream *model.Upstream) error {
	return s.repo.CreateUpstream(ctx, upstream)
}

func (s *UpstreamService) UpdateUpstream(ctx context.Context, upstream *model.Upstream) error {
	existing, err := s.repo.FindUpstreamById(ctx, int64(upstream.ID))
	if err != nil {
		return err
	}
	upstream.CreatedAt = existing.CreatedAt
	return s.repo.UpdateUpstream(ctx, upstream)
}

func (s *UpstreamService) DeleteUpstream(ctx context.Context, id uint64) error {
	return s.repo.DeleteUpstream(ctx, id)
}
