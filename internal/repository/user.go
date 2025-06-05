package repository

import (
	"context"
	"nebulaLive/internal/entity/ent"
)

// UserRepository provides methods for User entity operations.
type UserRepository struct {
	client *Client
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(client *Client) *UserRepository {
	return &UserRepository{client: client}
}

// CreateUser creates a new user.
func (r *UserRepository) CreateUser(ctx context.Context, name, email string) (*ent.User, error) {
	return r.client.User.Create().
		SetUsername(name).
		SetEmail(email).
		Save(ctx)
}

// GetUserByID gets a user by ID.
func (r *UserRepository) GetUserByID(ctx context.Context, id uint32) (*ent.User, error) {
	return r.client.User.Get(ctx, id)
}
