package repository

import (
	"context"
	"time"

	"smart-gateway/internal/clients"
	"smart-gateway/internal/model"
)

type UserRepo struct {
	db *clients.MysqlDB
}

func NewUserRepo(db *clients.MysqlDB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (userRepo *UserRepo) FindByUserName(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := userRepo.db.WithContext(ctx).
		Where("username = ?", username).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (userRepo *UserRepo) FindAll(ctx context.Context, dto model.UserDto) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	q := userRepo.db.WithContext(ctx).Model(&model.User{}).
		Where("username LIKE ?", "%"+dto.Username+"%")
	if dto.Status != nil {
		q = q.Where("status = ?", *dto.Status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Offset((dto.Page - 1) * dto.Size).
		Limit(dto.Size).
		Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (userRepo *UserRepo) Create(ctx context.Context, user *model.User) error {
	return userRepo.db.WithContext(ctx).Create(user).Error
}

func (userRepo *UserRepo) FindByID(ctx context.Context, id uint64) (*model.User, error) {
	var user model.User
	if err := userRepo.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (userRepo *UserRepo) Update(ctx context.Context, user *model.User) error {
	return userRepo.db.WithContext(ctx).Save(user).Error
}

func (userRepo *UserRepo) Delete(ctx context.Context, id uint64) error {
	return userRepo.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

func (userRepo *UserRepo) UpdateLastLogin(ctx context.Context, id uint64, t time.Time) error {
	return userRepo.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", id).
		Update("last_login_at", t).Error
}
