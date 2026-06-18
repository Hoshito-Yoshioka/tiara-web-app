package db

import (
	"context"
	"tiara-web-app/backend/internal/domain"
	"tiara-web-app/backend/internal/usecase"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// staffRefreshTokenRepository は usecase.StaffRefreshTokenRepository インターフェースの実装。
type staffRefreshTokenRepository struct {
	q *Queries
}

// NewStaffRefreshTokenRepository は新しいStaffRefreshTokenRepositoryのインスタンスを作成する。
func NewStaffRefreshTokenRepository(q *Queries) usecase.StaffRefreshTokenRepository {
	return &staffRefreshTokenRepository{q: q}
}

func convertToRefreshTokenDomain(row StaffRefreshToken) domain.StaffRefreshToken {
	return domain.StaffRefreshToken{
		ID:        uuid.UUID(row.ID.Bytes),
		StaffID:   uuid.UUID(row.StaffID.Bytes),
		Token:     row.Token,
		ExpiresAt: row.ExpiresAt.Time,
		CreatedAt: row.CreatedAt.Time,
	}
}

// CreateRefreshToken はリフレッシュトークンをDBに保存する。
func (r *staffRefreshTokenRepository) CreateRefreshToken(ctx context.Context, staffID uuid.UUID, token string, expiresAt time.Time) (domain.StaffRefreshToken, error) {
	row, err := r.q.CreateRefreshToken(ctx, CreateRefreshTokenParams{
		StaffID:   pgtype.UUID{Bytes: staffID, Valid: true},
		Token:     token,
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	})
	if err != nil {
		return domain.StaffRefreshToken{}, err
	}
	return convertToRefreshTokenDomain(row), nil
}

// GetRefreshToken はトークン文字列でリフレッシュトークンを取得する。
func (r *staffRefreshTokenRepository) GetRefreshToken(ctx context.Context, token string) (domain.StaffRefreshToken, error) {
	row, err := r.q.GetRefreshToken(ctx, token)
	if err != nil {
		return domain.StaffRefreshToken{}, err
	}
	return convertToRefreshTokenDomain(row), nil
}

// DeleteRefreshToken はトークンをDBから削除する。
func (r *staffRefreshTokenRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	return r.q.DeleteRefreshToken(ctx, token)
}

// DeleteRefreshTokensByStaffID はスタッフIDに紐づく全トークンを削除する。
func (r *staffRefreshTokenRepository) DeleteRefreshTokensByStaffID(ctx context.Context, staffID uuid.UUID) error {
	return r.q.DeleteRefreshTokensByStaffID(ctx, pgtype.UUID{Bytes: staffID, Valid: true})
}
