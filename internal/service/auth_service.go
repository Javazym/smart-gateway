package service

import (
	"context"
	"errors"
	"time"

	"smart-gateway/internal/clients"
	"smart-gateway/internal/model"
	"smart-gateway/internal/repository"
	"smart-gateway/internal/util"
)

var (
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrUserDisabled       = errors.New("用户已被禁用")
)

type AuthService struct {
	userRepo  *repository.UserRepo
	jwtSecret string
	jwtExpire time.Duration
}

func NewAuthService(mysqlDB *clients.MysqlDB, jwtSecret string, jwtExpire time.Duration) *AuthService {
	return &AuthService{
		userRepo:  repository.NewUserRepo(mysqlDB),
		jwtSecret: jwtSecret,
		jwtExpire: jwtExpire,
	}
}

func (s *AuthService) Login(ctx context.Context, req *model.UserRequest) (*model.LoginResponse, error) {
	user, err := s.userRepo.FindByUserName(ctx, req.Username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if !util.CheckPassword(user.Password, req.Password) {
		return nil, ErrInvalidCredentials
	}
	if user.Status != 1 {
		return nil, ErrUserDisabled
	}
	token, err := util.GenerateToken(s.jwtSecret, s.jwtExpire, user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}
	_ = s.userRepo.UpdateLastLogin(ctx, user.ID, time.Now())
	return &model.LoginResponse{
		Token:    token,
		Username: user.Username,
		Role:     user.Role,
	}, nil
}

func (s *AuthService) CreateUser(ctx context.Context, req *model.UserRequest) error {
	hash, err := util.HashPassword(req.Password)
	if err != nil {
		return err
	}
	user := &model.User{
		Username: req.Username,
		Password: hash,
		Role:     "viewer",
		Status:   1,
	}
	return s.userRepo.Create(ctx, user)
}

// EnsureAdmin 确保存在初始 admin/admin 账号（密码以 bcrypt 哈希存储）。
func EnsureAdmin(db *clients.MysqlDB) error {
	userRepo := repository.NewUserRepo(db)
	ctx := context.Background()
	if _, err := userRepo.FindByUserName(ctx, "admin"); err == nil {
		return nil
	}
	hash, err := util.HashPassword("admin")
	if err != nil {
		return err
	}
	return userRepo.Create(ctx, &model.User{
		Username: "admin",
		Password: hash,
		Role:     "admin",
		Status:   1,
	})
}
