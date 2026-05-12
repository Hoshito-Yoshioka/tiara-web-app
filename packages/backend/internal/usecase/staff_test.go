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
