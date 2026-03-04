package db

import (
	"context"
	"tiara-web-app/backend/internal/domain"
	"tiara-web-app/backend/internal/usecase"

	"github.com/google/uuid"
)

// adminRepository は usecase.AdminRepository インターフェースの実装。
type adminRepository struct {
	q *Queries
}

// NewAdminRepository は新しいAdminRepositoryのインスタンスを作成する。
func NewAdminRepository(q *Queries) usecase.AdminRepository {
	return &adminRepository{q: q}
}

// GetAdminByUsername はユーザー名で管理者を取得する。
func (r *adminRepository) GetAdminByUsername(ctx context.Context, username string) (domain.Admin, error) {
	row, err := r.q.GetAdminByUsername(ctx, username)
	if err != nil {
		return domain.Admin{}, err
	}

	return domain.Admin{
		ID:           uuid.UUID(row.ID.Bytes),
		Username:     row.Username,
		PasswordHash: row.PasswordHash,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}, nil
}
