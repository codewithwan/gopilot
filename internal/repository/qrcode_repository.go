package repository

import (
	"context"

	"github.com/codewithwan/gopilot/internal/domain"
	"github.com/codewithwan/gopilot/internal/repository/db"
)

type QRCodeRepository struct {
	queries *db.Queries
}

func NewQRCodeRepository(queries *db.Queries) *QRCodeRepository {
	return &QRCodeRepository{queries: queries}
}

func (r *QRCodeRepository) CreateQRCode(ctx context.Context, qr *domain.QRCode) error {
	// Security: Validate size is within safe range for int32 conversion
	if qr.Size < 0 || qr.Size > 2048 {
		qr.Size = 256 // default safe size
	}

	params := db.CreateQRCodeParams{
		ID:        qr.ID,
		Text:      qr.Text,
		Format:    qr.Format,
		Size:      int32(qr.Size), // #nosec G115 - size is validated to be within safe range
		ImageData: qr.ImageData,
	}

	result, err := r.queries.CreateQRCode(ctx, params)
	if err != nil {
		return err
	}

	qr.CreatedAt = result.CreatedAt.Time

	return nil
}

func (r *QRCodeRepository) GetQRCodeByID(ctx context.Context, id string) (*domain.QRCode, error) {
	result, err := r.queries.GetQRCodeByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.QRCode{
		ID:        result.ID,
		Text:      result.Text,
		Format:    result.Format,
		Size:      int(result.Size),
		ImageData: result.ImageData,
		CreatedAt: result.CreatedAt.Time,
	}, nil
}
