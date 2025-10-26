package repository

import (
	"context"

	"github.com/codewithwan/gopilot/internal/domain"
	"github.com/codewithwan/gopilot/internal/repository/db"
)

type URLShortenerRepository struct {
	queries *db.Queries
}

func NewURLShortenerRepository(queries *db.Queries) *URLShortenerRepository {
	return &URLShortenerRepository{queries: queries}
}

func (r *URLShortenerRepository) CreateShortURL(ctx context.Context, shortURL *domain.ShortURL) error {
	params := db.CreateShortURLParams{
		Code:        shortURL.Code,
		OriginalUrl: shortURL.OriginalURL,
		Alias:       toNullString(shortURL.Alias),
		Clicks:      shortURL.Clicks,
		IsPublic:    shortURL.IsPublic,
		ExpiresAt:   toNullTime(shortURL.ExpiresAt),
	}

	result, err := r.queries.CreateShortURL(ctx, params)
	if err != nil {
		return err
	}

	shortURL.ID = result.ID
	shortURL.CreatedAt = result.CreatedAt.Time
	shortURL.UpdatedAt = result.UpdatedAt.Time

	return nil
}

func (r *URLShortenerRepository) GetShortURLByCode(ctx context.Context, code string) (*domain.ShortURL, error) {
	result, err := r.queries.GetShortURLByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	return &domain.ShortURL{
		ID:          result.ID,
		Code:        result.Code,
		OriginalURL: result.OriginalUrl,
		Alias:       fromNullString(result.Alias),
		Clicks:      result.Clicks,
		IsPublic:    result.IsPublic,
		ExpiresAt:   fromNullTime(result.ExpiresAt),
		CreatedAt:   result.CreatedAt.Time,
		UpdatedAt:   result.UpdatedAt.Time,
	}, nil
}

func (r *URLShortenerRepository) IncrementClicks(ctx context.Context, id int64) error {
	return r.queries.IncrementShortURLClicks(ctx, id)
}

func (r *URLShortenerRepository) LogClick(ctx context.Context, click *domain.URLClickLog) error {
	params := db.CreateURLClickParams{
		ShortUrlID: click.ShortURLID,
		Referrer:   toNullString(click.Referrer),
		UserAgent:  toNullString(click.UserAgent),
		IpAddress:  toNullString(click.IPAddress),
	}

	return r.queries.CreateURLClick(ctx, params)
}

func (r *URLShortenerRepository) DeleteExpiredURLs(ctx context.Context) error {
	return r.queries.DeleteExpiredShortURLs(ctx)
}
