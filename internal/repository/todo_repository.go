package repository

import (
	"context"

	"github.com/codewithwan/gopilot/internal/domain"
	"github.com/codewithwan/gopilot/internal/repository/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type TodoRepository interface {
	Create(ctx context.Context, todo *domain.Todo) (*domain.Todo, error)
	GetByID(ctx context.Context, id, userID int64) (*domain.Todo, error)
	List(ctx context.Context, userID int64, limit, offset int32) ([]*domain.Todo, error)
	Update(ctx context.Context, todo *domain.Todo) (*domain.Todo, error)
	Delete(ctx context.Context, id, userID int64) error
	Count(ctx context.Context, userID int64) (int64, error)
}

type todoRepository struct {
	queries *db.Queries
}

func NewTodoRepository(queries *db.Queries) TodoRepository {
	return &todoRepository{queries: queries}
}

func (r *todoRepository) Create(ctx context.Context, todo *domain.Todo) (*domain.Todo, error) {
	result, err := r.queries.CreateTodo(ctx, db.CreateTodoParams{
		Title: todo.Title,
		Description: pgtype.Text{
			String: todo.Description,
			Valid:  true,
		},
		Completed: todo.Completed,
		UserID:    todo.UserID,
	})
	if err != nil {
		return nil, err
	}

	return &domain.Todo{
		ID:          result.ID,
		Title:       result.Title,
		Description: result.Description.String,
		Completed:   result.Completed,
		UserID:      result.UserID,
		CreatedAt:   result.CreatedAt.Time,
		UpdatedAt:   result.UpdatedAt.Time,
	}, nil
}

func (r *todoRepository) GetByID(ctx context.Context, id, userID int64) (*domain.Todo, error) {
	result, err := r.queries.GetTodoByID(ctx, db.GetTodoByIDParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	return &domain.Todo{
		ID:          result.ID,
		Title:       result.Title,
		Description: result.Description.String,
		Completed:   result.Completed,
		UserID:      result.UserID,
		CreatedAt:   result.CreatedAt.Time,
		UpdatedAt:   result.UpdatedAt.Time,
	}, nil
}

func (r *todoRepository) List(ctx context.Context, userID int64, limit, offset int32) ([]*domain.Todo, error) {
	results, err := r.queries.ListTodos(ctx, db.ListTodosParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	todos := make([]*domain.Todo, len(results))
	for i, result := range results {
		todos[i] = &domain.Todo{
			ID:          result.ID,
			Title:       result.Title,
			Description: result.Description.String,
			Completed:   result.Completed,
			UserID:      result.UserID,
			CreatedAt:   result.CreatedAt.Time,
			UpdatedAt:   result.UpdatedAt.Time,
		}
	}

	return todos, nil
}

func (r *todoRepository) Update(ctx context.Context, todo *domain.Todo) (*domain.Todo, error) {
	result, err := r.queries.UpdateTodo(ctx, db.UpdateTodoParams{
		Title: todo.Title,
		Description: pgtype.Text{
			String: todo.Description,
			Valid:  true,
		},
		Completed: todo.Completed,
		ID:        todo.ID,
		UserID:    todo.UserID,
	})
	if err != nil {
		return nil, err
	}

	return &domain.Todo{
		ID:          result.ID,
		Title:       result.Title,
		Description: result.Description.String,
		Completed:   result.Completed,
		UserID:      result.UserID,
		CreatedAt:   result.CreatedAt.Time,
		UpdatedAt:   result.UpdatedAt.Time,
	}, nil
}

func (r *todoRepository) Delete(ctx context.Context, id, userID int64) error {
	return r.queries.DeleteTodo(ctx, db.DeleteTodoParams{
		ID:     id,
		UserID: userID,
	})
}

func (r *todoRepository) Count(ctx context.Context, userID int64) (int64, error) {
	return r.queries.CountTodos(ctx, userID)
}
