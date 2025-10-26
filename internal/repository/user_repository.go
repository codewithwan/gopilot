package repository

import (
	"context"

	"github.com/codewithwan/gopilot/internal/domain"
	"github.com/codewithwan/gopilot/internal/repository/db"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByID(ctx context.Context, id int64) (*domain.User, error)
}

type userRepository struct {
	queries *db.Queries
}

func NewUserRepository(queries *db.Queries) UserRepository {
	return &userRepository{queries: queries}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	result, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Username: user.Username,
		Password: user.Password,
	})
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:        result.ID,
		Username:  result.Username,
		CreatedAt: result.CreatedAt.Time,
		UpdatedAt: result.UpdatedAt.Time,
	}, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	result, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:        result.ID,
		Username:  result.Username,
		Password:  result.Password,
		CreatedAt: result.CreatedAt.Time,
		UpdatedAt: result.UpdatedAt.Time,
	}, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	result, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:        result.ID,
		Username:  result.Username,
		CreatedAt: result.CreatedAt.Time,
		UpdatedAt: result.UpdatedAt.Time,
	}, nil
}
