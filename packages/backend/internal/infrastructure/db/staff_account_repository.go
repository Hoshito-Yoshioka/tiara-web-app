package db

import (
	"context"
	"tiara-web-app/backend/internal/domain"
	"tiara-web-app/backend/internal/usecase"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// staffAccountRepository は usecase.StaffAccountRepository インターフェースの実装。
type staffAccountRepository struct {
	q *Queries
}

// NewStaffAccountRepository は新しいStaffAccountRepositoryのインスタンスを作成する。
func NewStaffAccountRepository(q *Queries) usecase.StaffAccountRepository {
	return &staffAccountRepository{q: q}
}

// convertToStaffAccountDomain は sqlc 生成の StaffAccount を domain.StaffAccount に変換する。
func convertToStaffAccountDomain(row StaffAccount) domain.StaffAccount {
	return domain.StaffAccount{
		ID:           uuid.UUID(row.ID.Bytes),
		StaffID:      uuid.UUID(row.StaffID.Bytes),
		Username:     row.Username,
		PasswordHash: row.PasswordHash,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}
}

// GetStaffAccountByUsername はユーザー名でスタッフアカウントを取得する。
func (r *staffAccountRepository) GetStaffAccountByUsername(ctx context.Context, username string) (domain.StaffAccount, error) {
	row, err := r.q.GetStaffAccountByUsername(ctx, username)
	if err != nil {
		return domain.StaffAccount{}, err
	}
	return convertToStaffAccountDomain(row), nil
}

// ListStaffAccounts は全スタッフアカウント一覧を取得する。
func (r *staffAccountRepository) ListStaffAccounts(ctx context.Context) ([]domain.StaffAccount, error) {
	rows, err := r.q.ListStaffAccounts(ctx)
	if err != nil {
		return nil, err
	}
	accounts := make([]domain.StaffAccount, len(rows))
	for i, row := range rows {
		accounts[i] = convertToStaffAccountDomain(row)
	}
	return accounts, nil
}

// GetStaffAccountByStaffID はスタッフIDでアカウントを取得する。
func (r *staffAccountRepository) GetStaffAccountByStaffID(ctx context.Context, staffID uuid.UUID) (domain.StaffAccount, error) {
	row, err := r.q.GetStaffAccountByStaffID(ctx, pgtype.UUID{Bytes: staffID, Valid: true})
	if err != nil {
		return domain.StaffAccount{}, err
	}
	return convertToStaffAccountDomain(row), nil
}

// CreateStaffAccount は新しいスタッフアカウントを作成する。
func (r *staffAccountRepository) CreateStaffAccount(ctx context.Context, staffID uuid.UUID, username, passwordHash string) (domain.StaffAccount, error) {
	row, err := r.q.CreateStaffAccount(ctx, CreateStaffAccountParams{
		StaffID:      pgtype.UUID{Bytes: staffID, Valid: true},
		Username:     username,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return domain.StaffAccount{}, err
	}
	return convertToStaffAccountDomain(row), nil
}

// UpdateStaffAccount はスタッフアカウントのユーザー名とパスワードハッシュを更新する。
func (r *staffAccountRepository) UpdateStaffAccount(ctx context.Context, id uuid.UUID, username, passwordHash string) (domain.StaffAccount, error) {
	row, err := r.q.UpdateStaffAccount(ctx, UpdateStaffAccountParams{
		ID:           pgtype.UUID{Bytes: id, Valid: true},
		Username:     username,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return domain.StaffAccount{}, err
	}
	return convertToStaffAccountDomain(row), nil
}

// DeleteStaffAccount はスタッフアカウントを削除する。
func (r *staffAccountRepository) DeleteStaffAccount(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteStaffAccount(ctx, pgtype.UUID{Bytes: id, Valid: true})
}
