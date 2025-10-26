package repository

import (
	"context"

	"github.com/codewithwan/gopilot/internal/domain"
	"github.com/codewithwan/gopilot/internal/repository/db"
)

type PastebinRepository struct {
	queries *db.Queries
}

func NewPastebinRepository(queries *db.Queries) *PastebinRepository {
	return &PastebinRepository{queries: queries}
}

func (r *PastebinRepository) CreatePaste(ctx context.Context, paste *domain.Paste) error {
	params := db.CreatePasteParams{
		ID:           paste.ID,
		Title:        toNullString(paste.Title),
		Content:      paste.Content,
		Syntax:       toNullString(paste.Syntax),
		IsPublic:     paste.IsPublic,
		IsCompressed: paste.IsCompressed,
		ExpiresAt:    toNullTime(paste.ExpiresAt),
	}

	result, err := r.queries.CreatePaste(ctx, params)
	if err != nil {
		return err
	}

	paste.CreatedAt = result.CreatedAt.Time
	paste.UpdatedAt = result.UpdatedAt.Time

	return nil
}

func (r *PastebinRepository) GetPasteByID(ctx context.Context, id string) (*domain.Paste, error) {
	result, err := r.queries.GetPasteByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.Paste{
		ID:           result.ID,
		Title:        fromNullString(result.Title),
		Content:      result.Content,
		Syntax:       fromNullString(result.Syntax),
		IsPublic:     result.IsPublic,
		IsCompressed: result.IsCompressed,
		ExpiresAt:    fromNullTime(result.ExpiresAt),
		CreatedAt:    result.CreatedAt.Time,
		UpdatedAt:    result.UpdatedAt.Time,
	}, nil
}

func (r *PastebinRepository) DeletePaste(ctx context.Context, id string) error {
	return r.queries.DeletePaste(ctx, id)
}

func (r *PastebinRepository) ListRecentPastes(ctx context.Context, limit int) ([]*domain.Paste, error) {
	// Security: Validate limit before conversion
	if limit < 0 || limit > 100 {
		limit = 20
	}

	results, err := r.queries.ListRecentPastes(ctx, int32(limit)) // #nosec G115 - limit is validated to be within safe range
	if err != nil {
		return nil, err
	}

	pastes := make([]*domain.Paste, len(results))
	for i, result := range results {
		pastes[i] = &domain.Paste{
			ID:           result.ID,
			Title:        fromNullString(result.Title),
			Content:      result.Content,
			Syntax:       fromNullString(result.Syntax),
			IsPublic:     result.IsPublic,
			IsCompressed: result.IsCompressed,
			ExpiresAt:    fromNullTime(result.ExpiresAt),
			CreatedAt:    result.CreatedAt.Time,
			UpdatedAt:    result.UpdatedAt.Time,
		}
	}

	return pastes, nil
}

func (r *PastebinRepository) DeleteExpiredPastes(ctx context.Context) error {
	return r.queries.DeleteExpiredPastes(ctx)
}
