package db

import (
	"context"
	"tiara-web-app/backend/internal/domain"
	"tiara-web-app/backend/internal/usecase"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// staffRepository は usecase.StaffRepository インターフェースの実装。
type staffRepository struct {
	q *Queries
}

// NewStaffRepository は新しいStaffRepositoryのインスタンスを作成する。
func NewStaffRepository(q *Queries) usecase.StaffRepository {
	return &staffRepository{q: q}
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
