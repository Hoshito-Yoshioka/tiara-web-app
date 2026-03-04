package usecase

import (
	"context"
	"tiara-web-app/backend/internal/domain"
)

// StaffRepository はスタッフデータの永続化を抽象化するインターフェース。
type StaffRepository interface {
	ListStaffs(ctx context.Context) ([]domain.Staff, error)
	GetStaffByID(ctx context.Context, id string) (domain.Staff, error)
	ListSchedulesByStaffID(ctx context.Context, staffID string) ([]domain.StaffSchedule, error)
	ListAllSchedules(ctx context.Context) ([]domain.StaffSchedule, error)
	CreateStaff(ctx context.Context, input domain.CreateStaffInput) (domain.Staff, error)
	UpdateStaff(ctx context.Context, id string, input domain.UpdateStaffInput) (domain.Staff, error)
	DeleteStaff(ctx context.Context, id string) error
	ReplaceSchedules(ctx context.Context, staffID string, schedules []domain.ScheduleInput) ([]domain.StaffSchedule, error)
}

// StaffUsecase はスタッフに関するビジネスロジックを定義するインターフェース。
type StaffUsecase interface {
	ListStaffs(ctx context.Context) ([]domain.Staff, error)
	GetStaffWithSchedules(ctx context.Context, id string) (domain.StaffWithSchedules, error)
	ListAllStaffsWithSchedules(ctx context.Context) ([]domain.StaffWithSchedules, error)
	CreateStaff(ctx context.Context, input domain.CreateStaffInput) (domain.StaffWithSchedules, error)
	UpdateStaff(ctx context.Context, id string, input domain.UpdateStaffInput) (domain.StaffWithSchedules, error)
	DeleteStaff(ctx context.Context, id string) error
}

type staffUsecase struct {
	staffRepo StaffRepository
}

// NewStaffUsecase は新しいStaffUsecaseのインスタンスを作成する。
func NewStaffUsecase(repo StaffRepository) StaffUsecase {
	return &staffUsecase{staffRepo: repo}
}

// ListStaffs はすべてのスタッフを取得する。
func (u *staffUsecase) ListStaffs(ctx context.Context) ([]domain.Staff, error) {
	return u.staffRepo.ListStaffs(ctx)
}

// GetStaffWithSchedules はスタッフ詳細と出勤スケジュールを合わせて取得する。
// ドメイン集約 StaffWithSchedules を構成するのが Usecase 層の責務。
func (u *staffUsecase) GetStaffWithSchedules(ctx context.Context, id string) (domain.StaffWithSchedules, error) {
	staff, err := u.staffRepo.GetStaffByID(ctx, id)
	if err != nil {
		return domain.StaffWithSchedules{}, err
	}

	schedules, err := u.staffRepo.ListSchedulesByStaffID(ctx, id)
	if err != nil {
		return domain.StaffWithSchedules{}, err
	}

	return domain.StaffWithSchedules{
		Staff:     staff,
		Schedules: schedules,
	}, nil
}

// ListAllStaffsWithSchedules は全スタッフとその出勤スケジュールをまとめて取得する。
// Schedule ページ用。スタッフ一覧 + 全スケジュールを取得し、スタッフごとに集約する。
func (u *staffUsecase) ListAllStaffsWithSchedules(ctx context.Context) ([]domain.StaffWithSchedules, error) {
	staffs, err := u.staffRepo.ListStaffs(ctx)
	if err != nil {
		return nil, err
	}

	allSchedules, err := u.staffRepo.ListAllSchedules(ctx)
	if err != nil {
		return nil, err
	}

	// スタッフIDをキーにスケジュールをグルーピング
	scheduleMap := make(map[string][]domain.StaffSchedule)
	for _, s := range allSchedules {
		key := s.StaffID.String()
		scheduleMap[key] = append(scheduleMap[key], s)
	}

	result := make([]domain.StaffWithSchedules, len(staffs))
	for i, staff := range staffs {
		result[i] = domain.StaffWithSchedules{
			Staff:     staff,
			Schedules: scheduleMap[staff.ID.String()],
		}
	}

	return result, nil
}

// CreateStaff は新しいスタッフを作成し、スケジュールも登録する。
func (u *staffUsecase) CreateStaff(ctx context.Context, input domain.CreateStaffInput) (domain.StaffWithSchedules, error) {
	staff, err := u.staffRepo.CreateStaff(ctx, input)
	if err != nil {
		return domain.StaffWithSchedules{}, err
	}

	var schedules []domain.StaffSchedule
	if len(input.Schedules) > 0 {
		schedules, err = u.staffRepo.ReplaceSchedules(ctx, staff.ID.String(), input.Schedules)
		if err != nil {
			return domain.StaffWithSchedules{}, err
		}
	}

	return domain.StaffWithSchedules{
		Staff:     staff,
		Schedules: schedules,
	}, nil
}

// UpdateStaff は指定されたIDのスタッフ情報とスケジュールを更新する。
func (u *staffUsecase) UpdateStaff(ctx context.Context, id string, input domain.UpdateStaffInput) (domain.StaffWithSchedules, error) {
	staff, err := u.staffRepo.UpdateStaff(ctx, id, input)
	if err != nil {
		return domain.StaffWithSchedules{}, err
	}

	schedules, err := u.staffRepo.ReplaceSchedules(ctx, id, input.Schedules)
	if err != nil {
		return domain.StaffWithSchedules{}, err
	}

	return domain.StaffWithSchedules{
		Staff:     staff,
		Schedules: schedules,
	}, nil
}

// DeleteStaff は指定されたIDのスタッフを削除する。
func (u *staffUsecase) DeleteStaff(ctx context.Context, id string) error {
	return u.staffRepo.DeleteStaff(ctx, id)
}
