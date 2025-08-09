package service

import (
	"context"
	"errors"
	"time"

	"nebula-live/internal/domain/entity"
	"nebula-live/internal/domain/repository"
	"nebula-live/pkg/logger"
	"nebula-live/pkg/security"
	
	"go.uber.org/zap"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserBanned         = errors.New("user is banned")
	ErrUserInactive       = errors.New("user is inactive")
)

// UserService 用户领域服务接口
type UserService interface {
	// CreateUser 创建用户
	CreateUser(ctx context.Context, username, email, password, nickname string) (*entity.User, error)

	// GetUserByID 根据ID获取用户
	GetUserByID(ctx context.Context, id uint) (*entity.User, error)

	// GetUserByUsername 根据用户名获取用户
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)

	// GetUserByEmail 根据邮箱获取用户
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)

	// UpdateUser 更新用户信息
	UpdateUser(ctx context.Context, user *entity.User) error

	// DeleteUser 删除用户
	DeleteUser(ctx context.Context, id uint) error

	// ListUsers 获取用户列表
	ListUsers(ctx context.Context, offset, limit int) ([]*entity.User, error)

	// CountUsers 获取用户总数
	CountUsers(ctx context.Context) (int64, error)

	// ValidateUser 验证用户凭证
	ValidateUser(ctx context.Context, username, password string) (*entity.User, error)

	// ActivateUser 激活用户
	ActivateUser(ctx context.Context, id uint) error

	// DeactivateUser 停用用户
	DeactivateUser(ctx context.Context, id uint) error

	// BanUser 禁用用户
	BanUser(ctx context.Context, id uint) error
}

// userService 用户领域服务实现
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// CreateUser 创建用户
func (s *userService) CreateUser(ctx context.Context, username, email, password, nickname string) (*entity.User, error) {
	logger.Info("Creating new user", 
		zap.String("username", username), 
		zap.String("email", email))

	// 检查用户名是否已存在
	exists, err := s.userRepo.ExistsByUsername(ctx, username)
	if err != nil {
		logger.Error("Failed to check username existence", 
			zap.String("username", username), 
			zap.Error(err))
		return nil, err
	}
	if exists {
		logger.Warn("User creation failed - username already exists", 
			zap.String("username", username))
		return nil, ErrUserAlreadyExists
	}

	// 检查邮箱是否已存在
	exists, err = s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		logger.Error("Failed to check email existence", 
			zap.String("email", email), 
			zap.Error(err))
		return nil, err
	}
	if exists {
		logger.Warn("User creation failed - email already exists", 
			zap.String("email", email))
		return nil, ErrUserAlreadyExists
	}

	// 哈希密码
	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		logger.Error("Failed to hash password", 
			zap.String("username", username), 
			zap.Error(err))
		return nil, err
	}

	// 创建用户实体
	user := &entity.User{
		Username:  username,
		Email:     email,
		Password:  hashedPassword,
		Nickname:  nickname,
		Status:    entity.UserStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 保存用户
	err = s.userRepo.Create(ctx, user)
	if err != nil {
		logger.Error("Failed to create user", 
			zap.String("username", username), 
			zap.String("email", email), 
			zap.Error(err))
		return nil, err
	}

	logger.Info("User created successfully", 
		zap.Uint("user_id", user.ID), 
		zap.String("username", username))

	return user, nil
}

// GetUserByID 根据ID获取用户
func (s *userService) GetUserByID(ctx context.Context, id uint) (*entity.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// GetUserByUsername 根据用户名获取用户
func (s *userService) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	return s.userRepo.GetByUsername(ctx, username)
}

// GetUserByEmail 根据邮箱获取用户
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

// UpdateUser 更新用户信息
func (s *userService) UpdateUser(ctx context.Context, user *entity.User) error {
	user.UpdatedAt = time.Now()
	return s.userRepo.Update(ctx, user)
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	return s.userRepo.Delete(ctx, id)
}

// ListUsers 获取用户列表
func (s *userService) ListUsers(ctx context.Context, offset, limit int) ([]*entity.User, error) {
	return s.userRepo.List(ctx, offset, limit)
}

// CountUsers 获取用户总数
func (s *userService) CountUsers(ctx context.Context) (int64, error) {
	return s.userRepo.Count(ctx)
}

// ValidateUser 验证用户凭证
func (s *userService) ValidateUser(ctx context.Context, username, password string) (*entity.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// 验证密码
	valid, err := security.VerifyPassword(password, user.Password)
	if err != nil {
		logger.Error("Failed to verify password", 
			zap.String("username", username), 
			zap.Error(err))
		return nil, ErrInvalidCredentials
	}
	if !valid {
		return nil, ErrInvalidCredentials
	}

	// 检查用户状态
	if user.IsBanned() {
		return nil, ErrUserBanned
	}

	if !user.IsActive() {
		return nil, ErrUserInactive
	}

	return user, nil
}

// ActivateUser 激活用户
func (s *userService) ActivateUser(ctx context.Context, id uint) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	user.Activate()
	return s.userRepo.Update(ctx, user)
}

// DeactivateUser 停用用户
func (s *userService) DeactivateUser(ctx context.Context, id uint) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	user.Deactivate()
	return s.userRepo.Update(ctx, user)
}

// BanUser 禁用用户
func (s *userService) BanUser(ctx context.Context, id uint) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	user.Ban()
	return s.userRepo.Update(ctx, user)
}
