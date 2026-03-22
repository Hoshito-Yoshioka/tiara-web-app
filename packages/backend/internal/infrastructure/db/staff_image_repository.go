package db

import (
	"context"
	"tiara-web-app/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// convertToStaffImageDomain は sqlc 生成の StaffImage モデルを domain.StaffImage に変換する。
func convertToStaffImageDomain(row StaffImage) domain.StaffImage {
	return domain.StaffImage{
		ID:        uuid.UUID(row.ID.Bytes),
		StaffID:   uuid.UUID(row.StaffID.Bytes),
		ImageURL:  row.ImageUrl,
		IsMain:    row.IsMain,
		SortOrder: int(row.SortOrder),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}

// ListImagesByStaffID は指定されたスタッフIDの画像を取得する。
func (r *staffRepository) ListImagesByStaffID(ctx context.Context, staffID string) ([]domain.StaffImage, error) {
	uid, err := uuid.Parse(staffID)
	if err != nil {
		return nil, err
	}

	rows, err := r.q.ListImagesByStaffID(ctx, pgtype.UUID{Bytes: uid, Valid: true})
	if err != nil {
		return nil, err
	}

	images := make([]domain.StaffImage, len(rows))
	for i, row := range rows {
		images[i] = convertToStaffImageDomain(row)
	}
	return images, nil
}

// ListAllStaffImages は全スタッフの画像を取得する。
func (r *staffRepository) ListAllStaffImages(ctx context.Context) ([]domain.StaffImage, error) {
	rows, err := r.q.ListAllStaffImages(ctx)
	if err != nil {
		return nil, err
	}

	images := make([]domain.StaffImage, len(rows))
	for i, row := range rows {
		images[i] = convertToStaffImageDomain(row)
	}
	return images, nil
}

// CreateStaffImage はスタッフ画像を作成する。
func (r *staffRepository) CreateStaffImage(ctx context.Context, staffID string, imageURL string, isMain bool, sortOrder int) (domain.StaffImage, error) {
	uid, err := uuid.Parse(staffID)
	if err != nil {
		return domain.StaffImage{}, err
	}

	row, err := r.q.CreateStaffImage(ctx, CreateStaffImageParams{
		StaffID:   pgtype.UUID{Bytes: uid, Valid: true},
		ImageUrl:  imageURL,
		IsMain:    isMain,
		SortOrder: int32(sortOrder),
	})
	if err != nil {
		return domain.StaffImage{}, err
	}

	return convertToStaffImageDomain(row), nil
}

// DeleteStaffImage は指定されたIDの画像を削除する。
func (r *staffRepository) DeleteStaffImage(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.q.DeleteStaffImage(ctx, pgtype.UUID{Bytes: uid, Valid: true})
}

// SetMainImage は指定された画像をメイン画像に設定する（他のメインフラグをクリアしてから）。
func (r *staffRepository) SetMainImage(ctx context.Context, staffID string, imageID string) (domain.StaffImage, error) {
	staffUID, err := uuid.Parse(staffID)
	if err != nil {
		return domain.StaffImage{}, err
	}
	imgUID, err := uuid.Parse(imageID)
	if err != nil {
		return domain.StaffImage{}, err
	}

	// トランザクション開始
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.StaffImage{}, err
	}
	defer tx.Rollback(ctx)

	qtx := r.q.WithTx(tx)

	// 既存のメインフラグをクリア
	err = qtx.ClearMainFlagByStaffID(ctx, pgtype.UUID{Bytes: staffUID, Valid: true})
	if err != nil {
		return domain.StaffImage{}, err
	}

	// 新しいメイン画像を設定
	row, err := qtx.SetMainImage(ctx, pgtype.UUID{Bytes: imgUID, Valid: true})
	if err != nil {
		return domain.StaffImage{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.StaffImage{}, err
	}

	return convertToStaffImageDomain(row), nil
}
