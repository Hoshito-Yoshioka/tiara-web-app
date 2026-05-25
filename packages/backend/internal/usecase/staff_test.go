package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"tiara-web-app/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// mockStaffRepository は StaffRepository のテスト用モック。
type mockStaffRepository struct {
	staffs    []domain.Staff
	staff     domain.Staff
	schedules []domain.StaffSchedule
	images    []domain.StaffImage
	image     domain.StaffImage
	err       error
}

func (m *mockStaffRepository) ListStaffs(_ context.Context) ([]domain.Staff, error) {
	return m.staffs, m.err
}

func (m *mockStaffRepository) GetStaffByID(_ context.Context, _ string) (domain.Staff, error) {
	return m.staff, m.err
}

func (m *mockStaffRepository) ListSchedulesByStaffID(_ context.Context, _ string) ([]domain.StaffSchedule, error) {
	return m.schedules, m.err
}

func (m *mockStaffRepository) ListAllSchedules(_ context.Context) ([]domain.StaffSchedule, error) {
	return m.schedules, m.err
}

func (m *mockStaffRepository) CreateStaff(_ context.Context, _ domain.CreateStaffInput) (domain.Staff, error) {
	return m.staff, m.err
}

func (m *mockStaffRepository) UpdateStaff(_ context.Context, _ string, _ domain.UpdateStaffInput) (domain.Staff, error) {
	return m.staff, m.err
}

func (m *mockStaffRepository) DeleteStaff(_ context.Context, _ string) error {
	return m.err
}

func (m *mockStaffRepository) ReplaceSchedules(_ context.Context, _ string, _ []domain.ScheduleInput) ([]domain.StaffSchedule, error) {
	return m.schedules, m.err
}

func (m *mockStaffRepository) ListImagesByStaffID(_ context.Context, _ string) ([]domain.StaffImage, error) {
	return m.images, m.err
}

func (m *mockStaffRepository) ListAllStaffImages(_ context.Context) ([]domain.StaffImage, error) {
	return m.images, m.err
}

func (m *mockStaffRepository) CreateStaffImage(_ context.Context, _ string, _ string, _ bool, _ int) (domain.StaffImage, error) {
	return m.image, m.err
}

func (m *mockStaffRepository) DeleteStaffImage(_ context.Context, _ string) error {
	return m.err
}

func (m *mockStaffRepository) SetMainImage(_ context.Context, _ string, _ string) (domain.StaffImage, error) {
	return m.image, m.err
}

func (m *mockStaffRepository) UpdateImageCropPosition(_ context.Context, _ string, _ string) (domain.StaffImage, error) {
	return m.image, m.err
}

func (m *mockStaffRepository) ListStaffsPaginated(_ context.Context, limit, offset int) ([]domain.Staff, error) {
	// ページネーションをシミュレート
	start := offset
	end := offset + limit
	if start >= len(m.staffs) {
		return []domain.Staff{}, m.err
	}
	if end > len(m.staffs) {
		end = len(m.staffs)
	}
	return m.staffs[start:end], m.err
}

func (m *mockStaffRepository) CountStaffs(_ context.Context) (int, error) {
	return len(m.staffs), m.err
}

func TestStaffUsecase_ListStaffs_Success(t *testing.T) {
	staffs := []domain.Staff{
		{ID: uuid.New(), Name: "Staff A"},
		{ID: uuid.New(), Name: "Staff B"},
	}

	repo := &mockStaffRepository{staffs: staffs}
	uc := NewStaffUsecase(repo)

	result, err := uc.ListStaffs(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestStaffUsecase_GetStaffWithSchedules_Success(t *testing.T) {
	staffID := uuid.New()
	staff := domain.Staff{ID: staffID, Name: "Staff A"}
	schedules := []domain.StaffSchedule{
		{ID: uuid.New(), StaffID: staffID, DayOfWeek: 1, StartTime: time.Now(), EndTime: time.Now()},
	}
	images := []domain.StaffImage{
		{ID: uuid.New(), StaffID: staffID, ImageURL: "/uploads/staff/test.jpg", IsMain: true},
	}

	repo := &mockStaffRepository{staff: staff, schedules: schedules, images: images}
	uc := NewStaffUsecase(repo)

	result, err := uc.GetStaffWithSchedules(context.Background(), staffID.String())

	assert.NoError(t, err)
	assert.Equal(t, "Staff A", result.Staff.Name)
	assert.Len(t, result.Schedules, 1)
	assert.Len(t, result.Images, 1)
}

func TestStaffUsecase_GetStaffWithSchedules_StaffNotFound(t *testing.T) {
	repo := &mockStaffRepository{err: errors.New("not found")}
	uc := NewStaffUsecase(repo)

	_, err := uc.GetStaffWithSchedules(context.Background(), uuid.New().String())

	assert.Error(t, err)
}

func TestStaffUsecase_ListAllStaffsWithSchedules_Success(t *testing.T) {
	staffA := uuid.New()
	staffB := uuid.New()

	staffs := []domain.Staff{
		{ID: staffA, Name: "Staff A"},
		{ID: staffB, Name: "Staff B"},
	}
	schedules := []domain.StaffSchedule{
		{ID: uuid.New(), StaffID: staffA, DayOfWeek: 1},
		{ID: uuid.New(), StaffID: staffB, DayOfWeek: 3},
	}
	images := []domain.StaffImage{
		{ID: uuid.New(), StaffID: staffA, ImageURL: "/test.jpg"},
	}

	repo := &mockStaffRepository{staffs: staffs, schedules: schedules, images: images}
	uc := NewStaffUsecase(repo)

	result, err := uc.ListAllStaffsWithSchedules(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	// スタッフAにはスケジュールと画像が1件ずつ
	assert.Len(t, result[0].Schedules, 1)
	assert.Len(t, result[0].Images, 1)
	// スタッフBにはスケジュール1件、画像なし
	assert.Len(t, result[1].Schedules, 1)
	assert.Len(t, result[1].Images, 0)
}

func TestStaffUsecase_DeleteStaff_Success(t *testing.T) {
	repo := &mockStaffRepository{}
	uc := NewStaffUsecase(repo)

	err := uc.DeleteStaff(context.Background(), uuid.New().String())

	assert.NoError(t, err)
}

func TestStaffUsecase_UploadStaffImage_Success(t *testing.T) {
	img := domain.StaffImage{
		ID:       uuid.New(),
		ImageURL: "/uploads/staff/new.jpg",
		IsMain:   false,
	}

	repo := &mockStaffRepository{image: img}
	uc := NewStaffUsecase(repo)

	result, err := uc.UploadStaffImage(context.Background(), uuid.New().String(), "/uploads/staff/new.jpg", false, 1)

	assert.NoError(t, err)
	assert.Equal(t, img.ImageURL, result.ImageURL)
}

func TestStaffUsecase_CreateStaff_Success(t *testing.T) {
	staffID := uuid.New()
	staff := domain.Staff{ID: staffID, Name: "New Staff", Role: "キャスト"}
	schedules := []domain.StaffSchedule{
		{ID: uuid.New(), StaffID: staffID, DayOfWeek: 1},
	}

	repo := &mockStaffRepository{staff: staff, schedules: schedules, images: nil}
	uc := NewStaffUsecase(repo)

	result, err := uc.CreateStaff(context.Background(), domain.CreateStaffInput{
		ShopID: uuid.New().String(),
		Name:   "New Staff",
		Role:   "キャスト",
		Schedules: []domain.ScheduleInput{
			{DayOfWeek: 1, StartTime: "20:00", EndTime: "05:00"},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, "New Staff", result.Staff.Name)
}

func TestStaffUsecase_CreateStaff_NoSchedules(t *testing.T) {
	staff := domain.Staff{ID: uuid.New(), Name: "New Staff"}

	repo := &mockStaffRepository{staff: staff}
	uc := NewStaffUsecase(repo)

	result, err := uc.CreateStaff(context.Background(), domain.CreateStaffInput{
		ShopID: uuid.New().String(),
		Name:   "New Staff",
	})

	assert.NoError(t, err)
	assert.Equal(t, "New Staff", result.Staff.Name)
	assert.Nil(t, result.Schedules)
}

func TestStaffUsecase_CreateStaff_Error(t *testing.T) {
	repo := &mockStaffRepository{err: errors.New("db error")}
	uc := NewStaffUsecase(repo)

	_, err := uc.CreateStaff(context.Background(), domain.CreateStaffInput{Name: "Test"})

	assert.Error(t, err)
}

func TestStaffUsecase_UpdateStaff_Success(t *testing.T) {
	staffID := uuid.New()
	staff := domain.Staff{ID: staffID, Name: "Updated"}
	schedules := []domain.StaffSchedule{
		{ID: uuid.New(), StaffID: staffID, DayOfWeek: 2},
	}

	repo := &mockStaffRepository{staff: staff, schedules: schedules}
	uc := NewStaffUsecase(repo)

	result, err := uc.UpdateStaff(context.Background(), staffID.String(), domain.UpdateStaffInput{
		Name: "Updated",
		Schedules: []domain.ScheduleInput{
			{DayOfWeek: 2, StartTime: "20:00", EndTime: "05:00"},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, "Updated", result.Staff.Name)
}

func TestStaffUsecase_UpdateStaff_Error(t *testing.T) {
	repo := &mockStaffRepository{err: errors.New("not found")}
	uc := NewStaffUsecase(repo)

	_, err := uc.UpdateStaff(context.Background(), uuid.New().String(), domain.UpdateStaffInput{})

	assert.Error(t, err)
}

func TestStaffUsecase_DeleteStaffImage_Success(t *testing.T) {
	repo := &mockStaffRepository{}
	uc := NewStaffUsecase(repo)

	err := uc.DeleteStaffImage(context.Background(), uuid.New().String())

	assert.NoError(t, err)
}

func TestStaffUsecase_DeleteStaffImage_Error(t *testing.T) {
	repo := &mockStaffRepository{err: errors.New("not found")}
	uc := NewStaffUsecase(repo)

	err := uc.DeleteStaffImage(context.Background(), uuid.New().String())

	assert.Error(t, err)
}

func TestStaffUsecase_SetMainImage_Success(t *testing.T) {
	img := domain.StaffImage{ID: uuid.New(), IsMain: true}

	repo := &mockStaffRepository{image: img}
	uc := NewStaffUsecase(repo)

	result, err := uc.SetMainImage(context.Background(), uuid.New().String(), img.ID.String())

	assert.NoError(t, err)
	assert.True(t, result.IsMain)
}

func TestStaffUsecase_SetMainImage_Error(t *testing.T) {
	repo := &mockStaffRepository{err: errors.New("not found")}
	uc := NewStaffUsecase(repo)

	_, err := uc.SetMainImage(context.Background(), uuid.New().String(), uuid.New().String())

	assert.Error(t, err)
}

func TestStaffUsecase_UpdateImageCropPosition_Success(t *testing.T) {
	img := domain.StaffImage{ID: uuid.New(), CropPosition: "30 70"}

	repo := &mockStaffRepository{image: img}
	uc := NewStaffUsecase(repo)

	result, err := uc.UpdateImageCropPosition(context.Background(), img.ID.String(), "30 70")

	assert.NoError(t, err)
	assert.Equal(t, "30 70", result.CropPosition)
}

func TestStaffUsecase_UpdateImageCropPosition_Error(t *testing.T) {
	repo := &mockStaffRepository{err: errors.New("not found")}
	uc := NewStaffUsecase(repo)

	_, err := uc.UpdateImageCropPosition(context.Background(), uuid.New().String(), "50 50")

	assert.Error(t, err)
}
