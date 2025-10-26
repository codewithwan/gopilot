package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/codewithwan/gopilot/internal/domain"
	"github.com/skip2/go-qrcode"
)

// QRCodeRepository defines the interface for QR code storage
type QRCodeRepository interface {
	CreateQRCode(ctx context.Context, qr *domain.QRCode) error
	GetQRCodeByID(ctx context.Context, id string) (*domain.QRCode, error)
}

// QRCodeService handles QR code generation
type QRCodeService struct {
	repo QRCodeRepository
}

// NewQRCodeService creates a new QR code service
func NewQRCodeService(repo QRCodeRepository) *QRCodeService {
	return &QRCodeService{repo: repo}
}

// GenerateQR generates a QR code
func (s *QRCodeService) GenerateQR(ctx context.Context, req *domain.GenerateQRRequest) (*domain.QRCode, error) {
	size := 256
	if req.Size != nil {
		size = *req.Size
	}

	format := "png"
	if req.Format != nil {
		format = *req.Format
	}

	var imageData []byte
	var err error

	if format == "png" {
		imageData, err = qrcode.Encode(req.Text, qrcode.Medium, size)
		if err != nil {
			return nil, fmt.Errorf("failed to generate QR code: %w", err)
		}
	} else {
		// For SVG, we'll store a placeholder (full SVG support would require additional library)
		return nil, fmt.Errorf("SVG format not yet supported")
	}

	id := s.generateID(10)

	qr := &domain.QRCode{
		ID:        id,
		Text:      req.Text,
		Format:    format,
		Size:      size,
		ImageData: imageData,
	}

	if err := s.repo.CreateQRCode(ctx, qr); err != nil {
		return nil, fmt.Errorf("failed to save QR code: %w", err)
	}

	return qr, nil
}

// GetQRCode retrieves a QR code by ID
func (s *QRCodeService) GetQRCode(ctx context.Context, id string) (*domain.QRCode, error) {
	qr, err := s.repo.GetQRCodeByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get QR code: %w", err)
	}

	return qr, nil
}

// generateID generates a random ID
func (s *QRCodeService) generateID(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:length]
}
