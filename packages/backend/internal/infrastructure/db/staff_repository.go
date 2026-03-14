package db

import (
	"context"
	"tiara-web-app/backend/internal/domain"
	"tiara-web-app/backend/internal/usecase"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// staffRepository は usecase.StaffRepository インターフェースの実装。
// pool フィールドはトランザクション（スケジュール一括更新）に使用する。
type staffRepository struct {
	q    *Queries
	pool *pgxpool.Pool
}

// NewStaffRepository は新しいStaffRepositoryのインスタンスを作成する。
func NewStaffRepository(q *Queries, pool *pgxpool.Pool) usecase.StaffRepository {
	return &staffRepository{q: q, pool: pool}
}

// convertToStaffDomain は sqlc 生成の Staff モデルを domain.Staff に変換する。
func convertToStaffDomain(row Staff) domain.Staff {
	return domain.Staff{
		ID:        uuid.UUID(row.ID.Bytes),
		ShopID:    uuid.UUID(row.ShopID.Bytes),
		Name:      row.Name,
		Role:      row.Role,
		Bio:       row.Bio,
		ImageURL:  row.ImageUrl,
		SortOrder: int(row.SortOrder),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}

// convertToScheduleDomain は sqlc 生成の StaffSchedule モデルを domain.StaffSchedule に変換する。
func convertToScheduleDomain(row StaffSchedule) domain.StaffSchedule {
	return domain.StaffSchedule{
		ID:        uuid.UUID(row.ID.Bytes),
		StaffID:   uuid.UUID(row.StaffID.Bytes),
		DayOfWeek: int(row.DayOfWeek),
		StartTime: convertPgtypeTimeToTime(row.StartTime),
		EndTime:   convertPgtypeTimeToTime(row.EndTime),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}

// ListStaffs はすべてのスタッフをデータベースから取得する。
func (r *staffRepository) ListStaffs(ctx context.Context) ([]domain.Staff, error) {
	rows, err := r.q.ListStaffs(ctx)
	if err != nil {
		return nil, err
	}

	staffs := make([]domain.Staff, len(rows))
	for i, row := range rows {
		staffs[i] = convertToStaffDomain(row)
	}
	return staffs, nil
}

// GetStaffByID は指定されたIDのスタッフをデータベースから取得する。
func (r *staffRepository) GetStaffByID(ctx context.Context, id string) (domain.Staff, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return domain.Staff{}, err
	}

	pgtypeUUID := pgtype.UUID{Bytes: uid, Valid: true}

	row, err := r.q.GetStaffByID(ctx, pgtypeUUID)
	if err != nil {
		return domain.Staff{}, err
	}

	return convertToStaffDomain(row), nil
}

// ListSchedulesByStaffID は指定されたスタッフIDの出勤スケジュールを取得する。
func (r *staffRepository) ListSchedulesByStaffID(ctx context.Context, staffID string) ([]domain.StaffSchedule, error) {
	uid, err := uuid.Parse(staffID)
	if err != nil {
		return nil, err
	}

	pgtypeUUID := pgtype.UUID{Bytes: uid, Valid: true}

	rows, err := r.q.ListSchedulesByStaffID(ctx, pgtypeUUID)
	if err != nil {
		return nil, err
	}

	schedules := make([]domain.StaffSchedule, len(rows))
	for i, row := range rows {
		schedules[i] = convertToScheduleDomain(row)
	}
	return schedules, nil
}

// ListAllSchedules は全スタッフの出勤スケジュールを取得する。
func (r *staffRepository) ListAllSchedules(ctx context.Context) ([]domain.StaffSchedule, error) {
	rows, err := r.q.ListAllSchedules(ctx)
	if err != nil {
		return nil, err
	}

	schedules := make([]domain.StaffSchedule, len(rows))
	for i, row := range rows {
		schedules[i] = convertToScheduleDomain(row)
	}
	return schedules, nil
}

// CreateStaff は新しいスタッフをデータベースに作成する。
func (r *staffRepository) CreateStaff(ctx context.Context, input domain.CreateStaffInput) (domain.Staff, error) {
	shopUID, err := uuid.Parse(input.ShopID)
	if err != nil {
		return domain.Staff{}, err
	}

	row, err := r.q.CreateStaff(ctx, CreateStaffParams{
		ShopID:    pgtype.UUID{Bytes: shopUID, Valid: true},
		Name:      input.Name,
		Role:      input.Role,
		Bio:       input.Bio,
		ImageUrl:  input.ImageURL,
		SortOrder: int32(input.SortOrder),
	})
	if err != nil {
		return domain.Staff{}, err
	}

	return convertToStaffDomain(row), nil
}

// UpdateStaff は指定されたIDのスタッフ情報を更新する。
func (r *staffRepository) UpdateStaff(ctx context.Context, id string, input domain.UpdateStaffInput) (domain.Staff, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return domain.Staff{}, err
	}

	row, err := r.q.UpdateStaff(ctx, UpdateStaffParams{
		ID:        pgtype.UUID{Bytes: uid, Valid: true},
		Name:      input.Name,
		Role:      input.Role,
		Bio:       input.Bio,
		ImageUrl:  input.ImageURL,
		SortOrder: int32(input.SortOrder),
	})
	if err != nil {
		return domain.Staff{}, err
	}

	return convertToStaffDomain(row), nil
}

// DeleteStaff は指定されたIDのスタッフを削除する。
// CASCADE により関連するスケジュールも自動削除される。
func (r *staffRepository) DeleteStaff(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return r.q.DeleteStaff(ctx, pgtype.UUID{Bytes: uid, Valid: true})
}

// ReplaceSchedules はスタッフの出勤スケジュールをトランザクション内で一括置換する。
// 既存のスケジュールをすべて削除し、新しいスケジュールを挿入する。
func (r *staffRepository) ReplaceSchedules(ctx context.Context, staffID string, schedules []domain.ScheduleInput) ([]domain.StaffSchedule, error) {
	uid, err := uuid.Parse(staffID)
	if err != nil {
		return nil, err
	}

	pgtypeUID := pgtype.UUID{Bytes: uid, Valid: true}

	// トランザクション開始
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := r.q.WithTx(tx)

	// 既存のスケジュールをすべて削除
	err = qtx.DeleteSchedulesByStaffID(ctx, pgtypeUID)
	if err != nil {
		return nil, err
	}

	// 新しいスケジュールを挿入
	result := make([]domain.StaffSchedule, len(schedules))
	for i, s := range schedules {
		startTime, err := parseTimeToPgtype(s.StartTime)
		if err != nil {
			return nil, err
		}
		endTime, err := parseTimeToPgtype(s.EndTime)
		if err != nil {
			return nil, err
		}

		row, err := qtx.CreateSchedule(ctx, CreateScheduleParams{
			StaffID:   pgtypeUID,
			DayOfWeek: int32(s.DayOfWeek),
			StartTime: startTime,
			EndTime:   endTime,
		})
		if err != nil {
			return nil, err
		}

		result[i] = convertToScheduleDomain(row)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}
