package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/codewithwan/gopilot/internal/domain"
)

// PastebinRepository defines the interface for pastebin storage
type PastebinRepository interface {
	CreatePaste(ctx context.Context, paste *domain.Paste) error
	GetPasteByID(ctx context.Context, id string) (*domain.Paste, error)
	DeletePaste(ctx context.Context, id string) error
	ListRecentPastes(ctx context.Context, limit int) ([]*domain.Paste, error)
	DeleteExpiredPastes(ctx context.Context) error
}

// PastebinService handles pastebin operations
type PastebinService struct {
	repo PastebinRepository
}

// NewPastebinService creates a new pastebin service
func NewPastebinService(repo PastebinRepository) *PastebinService {
	return &PastebinService{repo: repo}
}

// CreatePaste creates a new paste
func (s *PastebinService) CreatePaste(ctx context.Context, req *domain.CreatePasteRequest) (*domain.Paste, error) {
	id := s.generateID(10)

	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}

	isCompressed := false
	if req.Compressed != nil {
		isCompressed = *req.Compressed
	}

	var expiresAt *time.Time
	if req.ExpireIn != nil {
		expiry := time.Now().Add(time.Duration(*req.ExpireIn) * time.Hour)
		expiresAt = &expiry
	} else {
		// Default 24h expiry
		expiry := time.Now().Add(24 * time.Hour)
		expiresAt = &expiry
	}

	paste := &domain.Paste{
		ID:           id,
		Title:        req.Title,
		Content:      req.Content,
		Syntax:       req.Syntax,
		IsPublic:     isPublic,
		IsCompressed: isCompressed,
		ExpiresAt:    expiresAt,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.CreatePaste(ctx, paste); err != nil {
		return nil, fmt.Errorf("failed to create paste: %w", err)
	}

	return paste, nil
}

// GetPaste retrieves a paste by ID
func (s *PastebinService) GetPaste(ctx context.Context, id string) (*domain.Paste, error) {
	paste, err := s.repo.GetPasteByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get paste: %w", err)
	}

	// Check if expired
	if paste.ExpiresAt != nil && time.Now().After(*paste.ExpiresAt) {
		return nil, fmt.Errorf("paste has expired")
	}

	return paste, nil
}

// DeletePaste deletes a paste by ID
func (s *PastebinService) DeletePaste(ctx context.Context, id string) error {
	if err := s.repo.DeletePaste(ctx, id); err != nil {
		return fmt.Errorf("failed to delete paste: %w", err)
	}
	return nil
}

// ListRecentPastes lists recent public pastes
func (s *PastebinService) ListRecentPastes(ctx context.Context, limit int) ([]*domain.Paste, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	pastes, err := s.repo.ListRecentPastes(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list recent pastes: %w", err)
	}

	return pastes, nil
}

// generateID generates a random ID
func (s *PastebinService) generateID(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)[:length]
}
