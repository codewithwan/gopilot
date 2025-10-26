package repository

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// Helper functions for converting between Go types and pgtype

func toNullString(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func fromNullString(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

func toNullTime(t *time.Time) pgtype.Timestamp {
	if t == nil {
		return pgtype.Timestamp{Valid: false}
	}
	return pgtype.Timestamp{Time: *t, Valid: true}
}

func fromNullTime(t pgtype.Timestamp) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}
