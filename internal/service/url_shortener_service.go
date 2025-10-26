package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/codewithwan/gopilot/internal/domain"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// URLShortenerRepository defines the interface for URL shortener storage
type URLShortenerRepository interface {
	CreateShortURL(ctx context.Context, shortURL *domain.ShortURL) error
	GetShortURLByCode(ctx context.Context, code string) (*domain.ShortURL, error)
	IncrementClicks(ctx context.Context, id int64) error
	LogClick(ctx context.Context, click *domain.URLClickLog) error
	DeleteExpiredURLs(ctx context.Context) error
}

// URLShortenerService handles URL shortening operations
type URLShortenerService struct {
	repo URLShortenerRepository
}

// NewURLShortenerService creates a new URL shortener service
func NewURLShortenerService(repo URLShortenerRepository) *URLShortenerService {
	return &URLShortenerService{repo: repo}
}

// CreateShortURL creates a new short URL
func (s *URLShortenerService) CreateShortURL(ctx context.Context, req *domain.CreateShortURLRequest) (*domain.ShortURL, error) {
	var code string
	if req.Alias != nil && *req.Alias != "" {
		code = *req.Alias
	} else {
		code = s.generateBase62Code(8)
	}

	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}

	var expiresAt *time.Time
	if req.ExpireIn != nil {
		expiry := time.Now().Add(time.Duration(*req.ExpireIn) * time.Hour)
		expiresAt = &expiry
	}

	shortURL := &domain.ShortURL{
		Code:        code,
		OriginalURL: req.OriginalURL,
		Alias:       req.Alias,
		IsPublic:    isPublic,
		ExpiresAt:   expiresAt,
		Clicks:      0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.CreateShortURL(ctx, shortURL); err != nil {
		return nil, fmt.Errorf("failed to create short URL: %w", err)
	}

	return shortURL, nil
}

// GetShortURL retrieves a short URL by code
func (s *URLShortenerService) GetShortURL(ctx context.Context, code string) (*domain.ShortURL, error) {
	shortURL, err := s.repo.GetShortURLByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get short URL: %w", err)
	}

	// Check if expired
	if shortURL.ExpiresAt != nil && time.Now().After(*shortURL.ExpiresAt) {
		return nil, fmt.Errorf("short URL has expired")
	}

	return shortURL, nil
}

// RecordClick records a click on a short URL
func (s *URLShortenerService) RecordClick(ctx context.Context, shortURL *domain.ShortURL, referrer, userAgent, ipAddress string) error {
	if err := s.repo.IncrementClicks(ctx, shortURL.ID); err != nil {
		return fmt.Errorf("failed to increment clicks: %w", err)
	}

	click := &domain.URLClickLog{
		ShortURLID: shortURL.ID,
		ClickedAt:  time.Now(),
	}

	if referrer != "" {
		click.Referrer = &referrer
	}
	if userAgent != "" {
		click.UserAgent = &userAgent
	}
	if ipAddress != "" {
		click.IPAddress = &ipAddress
	}

	if err := s.repo.LogClick(ctx, click); err != nil {
		return fmt.Errorf("failed to log click: %w", err)
	}

	return nil
}

// generateBase62Code generates a random base62 code
func (s *URLShortenerService) generateBase62Code(length int) string {
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(base62Chars))))
		if err != nil {
			return ""
		}
		result[i] = base62Chars[num.Int64()]
	}
	return string(result)
}
