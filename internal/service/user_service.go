package service

import (
	"context"

	"smart-gateway/internal/clients"
	"smart-gateway/internal/model"
	"smart-gateway/internal/repository"
	"smart-gateway/internal/util"
)

type UserService struct {
	repo *repository.UserRepo
}

func NewUserService(db *clients.MysqlDB) *UserService {
	return &UserService{repo: repository.NewUserRepo(db)}
}

func (s *UserService) List(ctx context.Context, dto model.UserDto) ([]model.User, int64, error) {
	return s.repo.FindAll(ctx, dto)
}

func (s *UserService) Get(ctx context.Context, id uint64) (*model.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *UserService) Create(ctx context.Context, req *model.UserCreateRequest) error {
	hash, err := util.HashPassword(req.Password)
	if err != nil {
		return err
	}
	role := req.Role
	if role == "" {
		role = "viewer"
	}
	status := req.Status
	if status == 0 {
		status = 1
	}
	return s.repo.Create(ctx, &model.User{
		Username: req.Username,
		Password: hash,
		Role:     role,
		Status:   status,
	})
}

func (s *UserService) Update(ctx context.Context, id uint64, req *model.UserUpdateRequest) error {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Status != nil {
		user.Status = *req.Status
	}
	if req.Password != "" {
		hash, err := util.HashPassword(req.Password)
		if err != nil {
			return err
		}
		user.Password = hash
	}
	return s.repo.Update(ctx, user)
}

func (s *UserService) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}
