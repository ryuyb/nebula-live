package service

import (
	"context"
	"nebulaLive/internal/entity/ent"
	"nebulaLive/internal/repository"
)

// UserService provides business logic for User entity.
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService creates a new UserService.
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// CreateUser creates a new user.
func (s *UserService) CreateUser(ctx context.Context, name, email string) (*ent.User, error) {
	return s.userRepo.CreateUser(ctx, name, email)
}

// GetUserByID gets a user by ID.
func (s *UserService) GetUserByID(ctx context.Context, id uint32) (*ent.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}
