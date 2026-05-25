package db

import (
	"context"
	"errors"
	"fmt"
	"tiara-web-app/backend/internal/domain"
	"tiara-web-app/backend/internal/usecase"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// staffDraftRepository は usecase.StaffDraftRepository インターフェースの実装。
// pool はトランザクション（スケジュール下書きのアイテム一括操作）に使用する。
type staffDraftRepository struct {
	q    *Queries
	pool *pgxpool.Pool
}

// NewStaffDraftRepository は新しいStaffDraftRepositoryのインスタンスを作成する。
func NewStaffDraftRepository(q *Queries, pool *pgxpool.Pool) usecase.StaffDraftRepository {
	return &staffDraftRepository{q: q, pool: pool}
}

// --- Helper ---

// detectConflict は pgx.ErrNoRows を楽観的ロックの競合エラー（ErrConflict）に変換する。
// :one クエリで WHERE updated_at = $N が一致しない場合、pgx は ErrNoRows を返す。
func detectConflict(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("他のユーザーによって更新されています: %w", domain.ErrConflict)
	}
	return err
}

// convertNullableTimestamptz は pgtype.Timestamptz（nullable）を *time.Time に変換する。
func convertNullableTimestamptz(ts pgtype.Timestamptz) *time.Time {
	if !ts.Valid {
		return nil
	}
	t := ts.Time
	return &t
}

// uuidToPgtype は uuid.UUID を pgtype.UUID に変換する。
func uuidToPgtype(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

// timeToPgTimestamptz は time.Time を pgtype.Timestamptz に変換する。
func timeToPgTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

// --- Profile Draft ---

// convertToProfileDraftDomain は sqlc 生成の StaffProfileDraft を domain に変換する。
func convertToProfileDraftDomain(row StaffProfileDraft) domain.StaffProfileDraft {
	return domain.StaffProfileDraft{
		ID:                uuid.UUID(row.ID.Bytes),
		StaffID:           uuid.UUID(row.StaffID.Bytes),
		Name:              row.Name,
		Role:              row.Role,
		Bio:               row.Bio,
		ImageURL:          row.ImageUrl,
		ImageCropPosition: row.ImageCropPosition,
		Status:            domain.DraftStatus(row.Status),
		AdminComment:      row.AdminComment,
		SubmittedAt:       convertNullableTimestamptz(row.SubmittedAt),
		ReviewedAt:        convertNullableTimestamptz(row.ReviewedAt),
		CreatedAt:         row.CreatedAt.Time,
		UpdatedAt:         row.UpdatedAt.Time,
	}
}

// GetProfileDraftByStaffID はスタッフIDで最新のアクティブなプロフィール下書きを取得する。
func (r *staffDraftRepository) GetProfileDraftByStaffID(ctx context.Context, staffID uuid.UUID) (domain.StaffProfileDraft, error) {
	row, err := r.q.GetProfileDraftByStaffID(ctx, uuidToPgtype(staffID))
	if err != nil {
		return domain.StaffProfileDraft{}, err
	}
	return convertToProfileDraftDomain(row), nil
}

// GetProfileDraftByID はIDでプロフィール下書きを取得する。
func (r *staffDraftRepository) GetProfileDraftByID(ctx context.Context, id uuid.UUID) (domain.StaffProfileDraft, error) {
	row, err := r.q.GetProfileDraftByID(ctx, uuidToPgtype(id))
	if err != nil {
		return domain.StaffProfileDraft{}, err
	}
	return convertToProfileDraftDomain(row), nil
}

// ListPendingProfileDrafts は承認待ちのプロフィール下書き一覧を取得する。
func (r *staffDraftRepository) ListPendingProfileDrafts(ctx context.Context) ([]domain.StaffProfileDraft, error) {
	rows, err := r.q.ListPendingProfileDrafts(ctx)
	if err != nil {
		return nil, err
	}
	drafts := make([]domain.StaffProfileDraft, len(rows))
	for i, row := range rows {
		drafts[i] = convertToProfileDraftDomain(row)
	}
	return drafts, nil
}

// CreateProfileDraft は新しいプロフィール下書きを作成する。
func (r *staffDraftRepository) CreateProfileDraft(ctx context.Context, staffID uuid.UUID, input domain.SaveProfileDraftInput) (domain.StaffProfileDraft, error) {
	row, err := r.q.CreateProfileDraft(ctx, CreateProfileDraftParams{
		StaffID:           uuidToPgtype(staffID),
		Name:              input.Name,
		Role:              input.Role,
		Bio:               input.Bio,
		ImageUrl:          input.ImageURL,
		ImageCropPosition: input.ImageCropPosition,
		Status:            string(domain.DraftStatusDraft),
	})
	if err != nil {
		return domain.StaffProfileDraft{}, err
	}
	return convertToProfileDraftDomain(row), nil
}

// UpdateProfileDraft はプロフィール下書きを更新する。
// 楽観的ロック: updated_at が一致しない場合は ErrConflict を返す。
func (r *staffDraftRepository) UpdateProfileDraft(ctx context.Context, id uuid.UUID, input domain.SaveProfileDraftInput, updatedAt time.Time) (domain.StaffProfileDraft, error) {
	row, err := r.q.UpdateProfileDraft(ctx, UpdateProfileDraftParams{
		ID:                uuidToPgtype(id),
		Name:              input.Name,
		Role:              input.Role,
		Bio:               input.Bio,
		ImageUrl:          input.ImageURL,
		ImageCropPosition: input.ImageCropPosition,
		Status:            string(domain.DraftStatusDraft),
		UpdatedAt:         timeToPgTimestamptz(updatedAt),
	})
	if err != nil {
		return domain.StaffProfileDraft{}, detectConflict(err)
	}
	return convertToProfileDraftDomain(row), nil
}

// SubmitProfileDraft はプロフィール下書きを承認申請（pending）に変更する。
// 楽観的ロック: updated_at が一致しない場合は ErrConflict を返す。
func (r *staffDraftRepository) SubmitProfileDraft(ctx context.Context, id uuid.UUID, updatedAt time.Time) (domain.StaffProfileDraft, error) {
	row, err := r.q.SubmitProfileDraft(ctx, SubmitProfileDraftParams{
		ID:        uuidToPgtype(id),
		UpdatedAt: timeToPgTimestamptz(updatedAt),
	})
	if err != nil {
		return domain.StaffProfileDraft{}, detectConflict(err)
	}
	return convertToProfileDraftDomain(row), nil
}

// ReviewProfileDraft は管理者がプロフィール下書きをレビューする。
// 楽観的ロック: updated_at が一致しない場合は ErrConflict を返す。
func (r *staffDraftRepository) ReviewProfileDraft(ctx context.Context, id uuid.UUID, input domain.ReviewDraftInput, updatedAt time.Time) (domain.StaffProfileDraft, error) {
	row, err := r.q.ReviewProfileDraft(ctx, ReviewProfileDraftParams{
		ID:           uuidToPgtype(id),
		Status:       string(input.Status),
		AdminComment: input.AdminComment,
		UpdatedAt:    timeToPgTimestamptz(updatedAt),
	})
	if err != nil {
		return domain.StaffProfileDraft{}, detectConflict(err)
	}
	return convertToProfileDraftDomain(row), nil
}

// DeleteProfileDraft はプロフィール下書きを削除する。
func (r *staffDraftRepository) DeleteProfileDraft(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteProfileDraft(ctx, uuidToPgtype(id))
}

// --- Schedule Draft ---

// convertToScheduleDraftDomain は sqlc 生成の StaffScheduleDraft を domain に変換する。
func convertToScheduleDraftDomain(row StaffScheduleDraft) domain.StaffScheduleDraft {
	return domain.StaffScheduleDraft{
		ID:           uuid.UUID(row.ID.Bytes),
		StaffID:      uuid.UUID(row.StaffID.Bytes),
		Status:       domain.DraftStatus(row.Status),
		AdminComment: row.AdminComment,
		SubmittedAt:  convertNullableTimestamptz(row.SubmittedAt),
		ReviewedAt:   convertNullableTimestamptz(row.ReviewedAt),
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}
}

// convertToScheduleDraftItemDomain は sqlc 生成の StaffScheduleDraftItem を domain に変換する。
func convertToScheduleDraftItemDomain(row StaffScheduleDraftItem) domain.ScheduleDraftItem {
	return domain.ScheduleDraftItem{
		ID:        uuid.UUID(row.ID.Bytes),
		DraftID:   uuid.UUID(row.DraftID.Bytes),
		DayOfWeek: int(row.DayOfWeek),
		StartTime: convertPgtypeTimeToTime(row.StartTime),
		EndTime:   convertPgtypeTimeToTime(row.EndTime),
	}
}

// GetScheduleDraftByStaffID はスタッフIDで最新のアクティブなスケジュール下書きを取得する。
func (r *staffDraftRepository) GetScheduleDraftByStaffID(ctx context.Context, staffID uuid.UUID) (domain.StaffScheduleDraft, error) {
	row, err := r.q.GetScheduleDraftByStaffID(ctx, uuidToPgtype(staffID))
	if err != nil {
		return domain.StaffScheduleDraft{}, err
	}
	draft := convertToScheduleDraftDomain(row)

	// アイテムを取得して結合
	items, err := r.q.ListScheduleDraftItems(ctx, row.ID)
	if err != nil {
		return domain.StaffScheduleDraft{}, err
	}
	draft.Items = make([]domain.ScheduleDraftItem, len(items))
	for i, item := range items {
		draft.Items[i] = convertToScheduleDraftItemDomain(item)
	}
	return draft, nil
}

// GetScheduleDraftByID はIDでスケジュール下書きを取得する（アイテム含む）。
func (r *staffDraftRepository) GetScheduleDraftByID(ctx context.Context, id uuid.UUID) (domain.StaffScheduleDraft, error) {
	row, err := r.q.GetScheduleDraftByID(ctx, uuidToPgtype(id))
	if err != nil {
		return domain.StaffScheduleDraft{}, err
	}
	draft := convertToScheduleDraftDomain(row)

	items, err := r.q.ListScheduleDraftItems(ctx, row.ID)
	if err != nil {
		return domain.StaffScheduleDraft{}, err
	}
	draft.Items = make([]domain.ScheduleDraftItem, len(items))
	for i, item := range items {
		draft.Items[i] = convertToScheduleDraftItemDomain(item)
	}
	return draft, nil
}

// ListPendingScheduleDrafts は承認待ちのスケジュール下書き一覧を取得する（アイテム含む）。
func (r *staffDraftRepository) ListPendingScheduleDrafts(ctx context.Context) ([]domain.StaffScheduleDraft, error) {
	rows, err := r.q.ListPendingScheduleDrafts(ctx)
	if err != nil {
		return nil, err
	}
	drafts := make([]domain.StaffScheduleDraft, len(rows))
	for i, row := range rows {
		drafts[i] = convertToScheduleDraftDomain(row)
		// 各ドラフトのアイテムも取得
		items, err := r.q.ListScheduleDraftItems(ctx, row.ID)
		if err != nil {
			return nil, err
		}
		drafts[i].Items = make([]domain.ScheduleDraftItem, len(items))
		for j, item := range items {
			drafts[i].Items[j] = convertToScheduleDraftItemDomain(item)
		}
	}
	return drafts, nil
}

// ListApprovedScheduleDrafts は承認済みのスケジュール下書き一覧を取得する（アイテム含む）。
func (r *staffDraftRepository) ListApprovedScheduleDrafts(ctx context.Context) ([]domain.StaffScheduleDraft, error) {
	rows, err := r.q.ListApprovedScheduleDrafts(ctx)
	if err != nil {
		return nil, err
	}
	drafts := make([]domain.StaffScheduleDraft, len(rows))
	for i, row := range rows {
		drafts[i] = convertToScheduleDraftDomain(row)
		items, err := r.q.ListScheduleDraftItems(ctx, row.ID)
		if err != nil {
			return nil, err
		}
		drafts[i].Items = make([]domain.ScheduleDraftItem, len(items))
		for j, item := range items {
			drafts[i].Items[j] = convertToScheduleDraftItemDomain(item)
		}
	}
	return drafts, nil
}

// CreateScheduleDraft は新しいスケジュール下書きを作成する（アイテム含む、トランザクション使用）。
func (r *staffDraftRepository) CreateScheduleDraft(ctx context.Context, staffID uuid.UUID, items []domain.ScheduleDraftItem) (domain.StaffScheduleDraft, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.StaffScheduleDraft{}, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck // rollback after commit is a no-op

	qtx := r.q.WithTx(tx)

	row, err := qtx.CreateScheduleDraft(ctx, CreateScheduleDraftParams{
		StaffID: uuidToPgtype(staffID),
		Status:  string(domain.DraftStatusDraft),
	})
	if err != nil {
		return domain.StaffScheduleDraft{}, err
	}

	draft := convertToScheduleDraftDomain(row)
	draft.Items = make([]domain.ScheduleDraftItem, 0, len(items))

	for _, item := range items {
		createdItem, err := qtx.CreateScheduleDraftItem(ctx, CreateScheduleDraftItemParams{
			DraftID:   row.ID,
			DayOfWeek: int32(item.DayOfWeek),
			StartTime: pgtype.Time{Microseconds: item.StartTime.Sub(time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)).Microseconds(), Valid: true},
			EndTime:   pgtype.Time{Microseconds: item.EndTime.Sub(time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)).Microseconds(), Valid: true},
		})
		if err != nil {
			return domain.StaffScheduleDraft{}, err
		}
		draft.Items = append(draft.Items, convertToScheduleDraftItemDomain(createdItem))
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.StaffScheduleDraft{}, err
	}
	return draft, nil
}

// UpdateScheduleDraftItems はスケジュール下書きのアイテムを全置換する（トランザクション使用）。
// 楽観的ロック: updated_at が一致しない場合は ErrConflict を返す。
func (r *staffDraftRepository) UpdateScheduleDraftItems(ctx context.Context, draftID uuid.UUID, items []domain.ScheduleDraftItem, updatedAt time.Time) (domain.StaffScheduleDraft, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.StaffScheduleDraft{}, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck // rollback after commit is a no-op

	qtx := r.q.WithTx(tx)

	// ステータスをdraftに戻す
	row, err := qtx.UpdateScheduleDraftStatus(ctx, UpdateScheduleDraftStatusParams{
		ID:        uuidToPgtype(draftID),
		Status:    string(domain.DraftStatusDraft),
		UpdatedAt: timeToPgTimestamptz(updatedAt),
	})
	if err != nil {
		return domain.StaffScheduleDraft{}, detectConflict(err)
	}

	// 既存アイテム削除
	if err := qtx.DeleteScheduleDraftItems(ctx, uuidToPgtype(draftID)); err != nil {
		return domain.StaffScheduleDraft{}, err
	}

	draft := convertToScheduleDraftDomain(row)
	draft.Items = make([]domain.ScheduleDraftItem, 0, len(items))

	// 新しいアイテム挿入
	for _, item := range items {
		createdItem, err := qtx.CreateScheduleDraftItem(ctx, CreateScheduleDraftItemParams{
			DraftID:   uuidToPgtype(draftID),
			DayOfWeek: int32(item.DayOfWeek),
			StartTime: pgtype.Time{Microseconds: item.StartTime.Sub(time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)).Microseconds(), Valid: true},
			EndTime:   pgtype.Time{Microseconds: item.EndTime.Sub(time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)).Microseconds(), Valid: true},
		})
		if err != nil {
			return domain.StaffScheduleDraft{}, err
		}
		draft.Items = append(draft.Items, convertToScheduleDraftItemDomain(createdItem))
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.StaffScheduleDraft{}, err
	}
	return draft, nil
}

// SubmitScheduleDraft はスケジュール下書きを承認申請（pending）に変更する。
// 楽観的ロック: updated_at が一致しない場合は ErrConflict を返す。
func (r *staffDraftRepository) SubmitScheduleDraft(ctx context.Context, id uuid.UUID, updatedAt time.Time) (domain.StaffScheduleDraft, error) {
	row, err := r.q.SubmitScheduleDraft(ctx, SubmitScheduleDraftParams{
		ID:        uuidToPgtype(id),
		UpdatedAt: timeToPgTimestamptz(updatedAt),
	})
	if err != nil {
		return domain.StaffScheduleDraft{}, detectConflict(err)
	}
	draft := convertToScheduleDraftDomain(row)
	// アイテムも返す
	items, err := r.q.ListScheduleDraftItems(ctx, row.ID)
	if err != nil {
		return domain.StaffScheduleDraft{}, err
	}
	draft.Items = make([]domain.ScheduleDraftItem, len(items))
	for i, item := range items {
		draft.Items[i] = convertToScheduleDraftItemDomain(item)
	}
	return draft, nil
}

// ReviewScheduleDraft は管理者がスケジュール下書きをレビューする。
// 楽観的ロック: updated_at が一致しない場合は ErrConflict を返す。
func (r *staffDraftRepository) ReviewScheduleDraft(ctx context.Context, id uuid.UUID, input domain.ReviewDraftInput, updatedAt time.Time) (domain.StaffScheduleDraft, error) {
	row, err := r.q.ReviewScheduleDraft(ctx, ReviewScheduleDraftParams{
		ID:           uuidToPgtype(id),
		Status:       string(input.Status),
		AdminComment: input.AdminComment,
		UpdatedAt:    timeToPgTimestamptz(updatedAt),
	})
	if err != nil {
		return domain.StaffScheduleDraft{}, detectConflict(err)
	}
	draft := convertToScheduleDraftDomain(row)
	items, err := r.q.ListScheduleDraftItems(ctx, row.ID)
	if err != nil {
		return domain.StaffScheduleDraft{}, err
	}
	draft.Items = make([]domain.ScheduleDraftItem, len(items))
	for i, item := range items {
		draft.Items[i] = convertToScheduleDraftItemDomain(item)
	}
	return draft, nil
}

// DeleteScheduleDraft はスケジュール下書きを削除する（カスケードでアイテムも削除）。
func (r *staffDraftRepository) DeleteScheduleDraft(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteScheduleDraft(ctx, uuidToPgtype(id))
}

// --- Admin用メソッド（ステータスを変更せずに内容のみ更新） ---

// UpdateProfileDraftContent はプロフィール下書きの内容のみ更新する（ステータス変更なし）。
// 管理者がレビュー時に内容を修正する際に使用する。
// 楽観的ロック: updated_at が一致しない場合は ErrConflict を返す。
func (r *staffDraftRepository) UpdateProfileDraftContent(ctx context.Context, id uuid.UUID, input domain.SaveProfileDraftInput, updatedAt time.Time) (domain.StaffProfileDraft, error) {
	row, err := r.q.UpdateProfileDraftContent(ctx, UpdateProfileDraftContentParams{
		ID:                uuidToPgtype(id),
		Name:              input.Name,
		Role:              input.Role,
		Bio:               input.Bio,
		ImageUrl:          input.ImageURL,
		ImageCropPosition: input.ImageCropPosition,
		UpdatedAt:         timeToPgTimestamptz(updatedAt),
	})
	if err != nil {
		return domain.StaffProfileDraft{}, detectConflict(err)
	}
	return convertToProfileDraftDomain(row), nil
}

// ReplaceScheduleDraftItems はスケジュール下書きのアイテムのみ全置換する（ステータス変更なし）。
// 管理者がレビュー時にスケジュールを修正する際に使用する。
// 楽観的ロック: updated_at が一致しない場合は ErrConflict を返す。
func (r *staffDraftRepository) ReplaceScheduleDraftItems(ctx context.Context, draftID uuid.UUID, items []domain.ScheduleDraftItem, updatedAt time.Time) (domain.StaffScheduleDraft, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.StaffScheduleDraft{}, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck // rollback after commit is a no-op

	qtx := r.q.WithTx(tx)

	// ドラフト本体を取得し、楽観的ロックを検証
	row, err := qtx.GetScheduleDraftByID(ctx, uuidToPgtype(draftID))
	if err != nil {
		return domain.StaffScheduleDraft{}, err
	}
	if !row.UpdatedAt.Time.Equal(updatedAt) {
		return domain.StaffScheduleDraft{}, fmt.Errorf("他のユーザーによって更新されています: %w", domain.ErrConflict)
	}

	// 既存アイテム削除
	if err := qtx.DeleteScheduleDraftItems(ctx, uuidToPgtype(draftID)); err != nil {
		return domain.StaffScheduleDraft{}, err
	}

	draft := convertToScheduleDraftDomain(row)
	draft.Items = make([]domain.ScheduleDraftItem, 0, len(items))

	// 新しいアイテム挿入
	for _, item := range items {
		createdItem, err := qtx.CreateScheduleDraftItem(ctx, CreateScheduleDraftItemParams{
			DraftID:   uuidToPgtype(draftID),
			DayOfWeek: int32(item.DayOfWeek),
			StartTime: pgtype.Time{Microseconds: item.StartTime.Sub(time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)).Microseconds(), Valid: true},
			EndTime:   pgtype.Time{Microseconds: item.EndTime.Sub(time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)).Microseconds(), Valid: true},
		})
		if err != nil {
			return domain.StaffScheduleDraft{}, err
		}
		draft.Items = append(draft.Items, convertToScheduleDraftItemDomain(createdItem))
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.StaffScheduleDraft{}, err
	}
	return draft, nil
}
