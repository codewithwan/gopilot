package service

import (
	"context"
	"errors"

	"github.com/codewithwan/gopilot/internal/domain"
	"github.com/codewithwan/gopilot/internal/repository"
	"go.uber.org/zap"
)

type TodoService interface {
	Create(ctx context.Context, req *domain.CreateTodoRequest, userID int64) (*domain.Todo, error)
	GetByID(ctx context.Context, id, userID int64) (*domain.Todo, error)
	List(ctx context.Context, userID int64, limit, offset int32) ([]*domain.Todo, int64, error)
	Update(ctx context.Context, id int64, req *domain.UpdateTodoRequest, userID int64) (*domain.Todo, error)
	Delete(ctx context.Context, id, userID int64) error
}

type todoService struct {
	repo   repository.TodoRepository
	logger *zap.Logger
}

func NewTodoService(repo repository.TodoRepository, logger *zap.Logger) TodoService {
	return &todoService{
		repo:   repo,
		logger: logger,
	}
}

func (s *todoService) Create(ctx context.Context, req *domain.CreateTodoRequest, userID int64) (*domain.Todo, error) {
	todo := &domain.Todo{
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
		UserID:      userID,
	}

	created, err := s.repo.Create(ctx, todo)
	if err != nil {
		s.logger.Error("failed to create todo", zap.Error(err), zap.Int64("user_id", userID))
		return nil, err
	}

	s.logger.Info("todo created", zap.Int64("todo_id", created.ID), zap.Int64("user_id", userID))
	return created, nil
}

func (s *todoService) GetByID(ctx context.Context, id, userID int64) (*domain.Todo, error) {
	todo, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		s.logger.Error("failed to get todo", zap.Error(err), zap.Int64("todo_id", id), zap.Int64("user_id", userID))
		return nil, err
	}

	return todo, nil
}

func (s *todoService) List(ctx context.Context, userID int64, limit, offset int32) ([]*domain.Todo, int64, error) {
	todos, err := s.repo.List(ctx, userID, limit, offset)
	if err != nil {
		s.logger.Error("failed to list todos", zap.Error(err), zap.Int64("user_id", userID))
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx, userID)
	if err != nil {
		s.logger.Error("failed to count todos", zap.Error(err), zap.Int64("user_id", userID))
		return nil, 0, err
	}

	return todos, count, nil
}

func (s *todoService) Update(ctx context.Context, id int64, req *domain.UpdateTodoRequest, userID int64) (*domain.Todo, error) {
	existing, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		s.logger.Error("failed to get todo for update", zap.Error(err), zap.Int64("todo_id", id), zap.Int64("user_id", userID))
		return nil, err
	}

	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.Completed != nil {
		existing.Completed = *req.Completed
	}

	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		s.logger.Error("failed to update todo", zap.Error(err), zap.Int64("todo_id", id), zap.Int64("user_id", userID))
		return nil, err
	}

	s.logger.Info("todo updated", zap.Int64("todo_id", id), zap.Int64("user_id", userID))
	return updated, nil
}

func (s *todoService) Delete(ctx context.Context, id, userID int64) error {
	if err := s.repo.Delete(ctx, id, userID); err != nil {
		s.logger.Error("failed to delete todo", zap.Error(err), zap.Int64("todo_id", id), zap.Int64("user_id", userID))
		return err
	}

	s.logger.Info("todo deleted", zap.Int64("todo_id", id), zap.Int64("user_id", userID))
	return nil
}

var ErrTodoNotFound = errors.New("todo not found")
