package repository

import (
	"context"
	"smart-gateway/internal/clients"
	"smart-gateway/internal/model"
)

type UpstreamRepo struct {
	db *clients.MysqlDB
}

func NewUpstreamRepo(db *clients.MysqlDB) *UpstreamRepo {
	return &UpstreamRepo{db: db}
}

func (s *UpstreamRepo) FindAddGroup(ctx context.Context, dto model.GroupDto) ([]model.UpstreamGroup, int64, error) {
	var upstreamGroups []model.UpstreamGroup
	var total int64
	q := s.db.WithContext(ctx).Model(&model.UpstreamGroup{})
	if dto.Status != nil {
		q = q.Where("status = ?", *dto.Status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Offset((dto.Page - 1) * dto.Size).
		Limit(dto.Size).
		Find(&upstreamGroups).Error
	if err != nil {
		return nil, 0, err
	}
	return upstreamGroups, total, nil
}

func (s *UpstreamRepo) FindUpstreamGroupById(ctx context.Context, id uint64) (*model.UpstreamGroup, error) {
	var upstreamGroup model.UpstreamGroup
	if err := s.db.WithContext(ctx).
		Preload("Upstreams").
		Where("id = ?", id).
		First(&upstreamGroup).Error; err != nil {
		return nil, err
	}
	return &upstreamGroup, nil
}

func (s *UpstreamRepo) FindUpstreamById(ctx context.Context, id int64) (*model.Upstream, error) {
	var upstream model.Upstream
	if err := s.db.WithContext(ctx).
		Where("id = ?", id).
		First(&upstream, id).Error; err != nil {
		return nil, err
	}
	return &upstream, nil
}

func (s *UpstreamRepo) CreateGroup(ctx context.Context, group *model.UpstreamGroup) error {
	return s.db.WithContext(ctx).Create(group).Error
}

func (s *UpstreamRepo) UpdateGroup(ctx context.Context, group *model.UpstreamGroup) error {
	return s.db.WithContext(ctx).Save(group).Error
}

func (s *UpstreamRepo) DeleteGroup(ctx context.Context, id uint64) error {
	return s.db.WithContext(ctx).Delete(&model.UpstreamGroup{}, id).Error
}

func (s *UpstreamRepo) CreateUpstream(ctx context.Context, upstream *model.Upstream) error {
	return s.db.WithContext(ctx).Create(upstream).Error
}

func (s *UpstreamRepo) UpdateUpstream(ctx context.Context, upstream *model.Upstream) error {
	return s.db.WithContext(ctx).Save(upstream).Error
}

func (s *UpstreamRepo) DeleteUpstream(ctx context.Context, id uint64) error {
	return s.db.WithContext(ctx).Delete(&model.Upstream{}, id).Error
}

func (s *UpstreamRepo) FindUpstreamsByGroupID(ctx context.Context, groupID uint64) ([]model.Upstream, error) {
	var upstreams []model.Upstream
	if err := s.db.WithContext(ctx).Where("group_id = ?", groupID).Find(&upstreams).Error; err != nil {
		return nil, err
	}
	return upstreams, nil
}
