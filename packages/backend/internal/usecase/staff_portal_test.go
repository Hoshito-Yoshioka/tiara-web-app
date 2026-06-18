package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"tiara-web-app/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// --- Mock Repositories ---

type mockStaffAccountRepository struct {
	account  domain.StaffAccount
	accounts []domain.StaffAccount
	err      error
}

func (m *mockStaffAccountRepository) GetStaffAccountByUsername(_ context.Context, _ string) (domain.StaffAccount, error) {
	return m.account, m.err
}

func (m *mockStaffAccountRepository) GetStaffAccountByStaffID(_ context.Context, _ uuid.UUID) (domain.StaffAccount, error) {
	return m.account, m.err
}

func (m *mockStaffAccountRepository) ListStaffAccounts(_ context.Context) ([]domain.StaffAccount, error) {
	return m.accounts, m.err
}

func (m *mockStaffAccountRepository) CreateStaffAccount(_ context.Context, _ uuid.UUID, _ string, _ string) (domain.StaffAccount, error) {
	return m.account, m.err
}

func (m *mockStaffAccountRepository) UpdateStaffAccount(_ context.Context, _ uuid.UUID, _ string, _ string) (domain.StaffAccount, error) {
	return m.account, m.err
}

func (m *mockStaffAccountRepository) DeleteStaffAccount(_ context.Context, _ uuid.UUID) error {
	return m.err
}

type mockStaffRefreshTokenRepository struct {
	token domain.StaffRefreshToken
	err   error
}

func (m *mockStaffRefreshTokenRepository) CreateRefreshToken(_ context.Context, _ uuid.UUID, _ string, _ time.Time) (domain.StaffRefreshToken, error) {
	return m.token, m.err
}
func (m *mockStaffRefreshTokenRepository) GetRefreshToken(_ context.Context, _ string) (domain.StaffRefreshToken, error) {
	return m.token, m.err
}
func (m *mockStaffRefreshTokenRepository) DeleteRefreshToken(_ context.Context, _ string) error {
	return m.err
}
func (m *mockStaffRefreshTokenRepository) DeleteRefreshTokensByStaffID(_ context.Context, _ uuid.UUID) error {
	return m.err
}

type mockStaffDraftRepository struct {
	profileDraft   domain.StaffProfileDraft
	profileDrafts  []domain.StaffProfileDraft
	scheduleDraft  domain.StaffScheduleDraft
	scheduleDrafts []domain.StaffScheduleDraft
	err            error
	// 特定メソッド用のエラーを分離
	getByStaffIDErr error
}

func (m *mockStaffDraftRepository) GetProfileDraftByStaffID(_ context.Context, _ uuid.UUID) (domain.StaffProfileDraft, error) {
	if m.getByStaffIDErr != nil {
		return domain.StaffProfileDraft{}, m.getByStaffIDErr
	}
	return m.profileDraft, m.err
}

func (m *mockStaffDraftRepository) GetProfileDraftByID(_ context.Context, _ uuid.UUID) (domain.StaffProfileDraft, error) {
	return m.profileDraft, m.err
}

func (m *mockStaffDraftRepository) ListPendingProfileDrafts(_ context.Context) ([]domain.StaffProfileDraft, error) {
	return m.profileDrafts, m.err
}

func (m *mockStaffDraftRepository) CreateProfileDraft(_ context.Context, _ uuid.UUID, _ domain.SaveProfileDraftInput) (domain.StaffProfileDraft, error) {
	return m.profileDraft, m.err
}

func (m *mockStaffDraftRepository) UpdateProfileDraft(_ context.Context, _ uuid.UUID, _ domain.SaveProfileDraftInput, _ time.Time) (domain.StaffProfileDraft, error) {
	return m.profileDraft, m.err
}

func (m *mockStaffDraftRepository) SubmitProfileDraft(_ context.Context, _ uuid.UUID, _ time.Time) (domain.StaffProfileDraft, error) {
	return m.profileDraft, m.err
}

func (m *mockStaffDraftRepository) ReviewProfileDraft(_ context.Context, _ uuid.UUID, _ domain.ReviewDraftInput, _ time.Time) (domain.StaffProfileDraft, error) {
	return m.profileDraft, m.err
}

func (m *mockStaffDraftRepository) DeleteProfileDraft(_ context.Context, _ uuid.UUID) error {
	return m.err
}

func (m *mockStaffDraftRepository) GetScheduleDraftByStaffID(_ context.Context, _ uuid.UUID) (domain.StaffScheduleDraft, error) {
	if m.getByStaffIDErr != nil {
		return domain.StaffScheduleDraft{}, m.getByStaffIDErr
	}
	return m.scheduleDraft, m.err
}

func (m *mockStaffDraftRepository) GetScheduleDraftByID(_ context.Context, _ uuid.UUID) (domain.StaffScheduleDraft, error) {
	return m.scheduleDraft, m.err
}

func (m *mockStaffDraftRepository) ListPendingScheduleDrafts(_ context.Context) ([]domain.StaffScheduleDraft, error) {
	return m.scheduleDrafts, m.err
}

func (m *mockStaffDraftRepository) ListApprovedScheduleDrafts(_ context.Context) ([]domain.StaffScheduleDraft, error) {
	return m.scheduleDrafts, m.err
}

func (m *mockStaffDraftRepository) CreateScheduleDraft(_ context.Context, _ uuid.UUID, _ []domain.ScheduleDraftItem) (domain.StaffScheduleDraft, error) {
	return m.scheduleDraft, m.err
}

func (m *mockStaffDraftRepository) UpdateScheduleDraftItems(_ context.Context, _ uuid.UUID, _ []domain.ScheduleDraftItem, _ time.Time) (domain.StaffScheduleDraft, error) {
	return m.scheduleDraft, m.err
}

func (m *mockStaffDraftRepository) SubmitScheduleDraft(_ context.Context, _ uuid.UUID, _ time.Time) (domain.StaffScheduleDraft, error) {
	return m.scheduleDraft, m.err
}

func (m *mockStaffDraftRepository) ReviewScheduleDraft(_ context.Context, _ uuid.UUID, _ domain.ReviewDraftInput, _ time.Time) (domain.StaffScheduleDraft, error) {
	return m.scheduleDraft, m.err
}

func (m *mockStaffDraftRepository) DeleteScheduleDraft(_ context.Context, _ uuid.UUID) error {
	return m.err
}

func (m *mockStaffDraftRepository) UpdateProfileDraftContent(_ context.Context, _ uuid.UUID, _ domain.SaveProfileDraftInput, _ time.Time) (domain.StaffProfileDraft, error) {
	return m.profileDraft, m.err
}

func (m *mockStaffDraftRepository) ReplaceScheduleDraftItems(_ context.Context, _ uuid.UUID, _ []domain.ScheduleDraftItem, _ time.Time) (domain.StaffScheduleDraft, error) {
	return m.scheduleDraft, m.err
}

// ============================================================
// StaffAuthUsecase Tests
// ============================================================

func TestStaffAuthUsecase_Login_Success(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	account := domain.StaffAccount{
		ID:           uuid.New(),
		StaffID:      uuid.New(),
		Username:     "staff1",
		PasswordHash: string(hash),
	}

	repo := &mockStaffAccountRepository{account: account}
	uc := NewStaffAuthUsecase(repo, &mockStaffRefreshTokenRepository{})

	result, err := uc.Login(context.Background(), "staff1", "password123")

	assert.NoError(t, err)
	assert.Equal(t, "staff1", result.Username)
}

func TestStaffAuthUsecase_Login_UserNotFound(t *testing.T) {
	repo := &mockStaffAccountRepository{err: errors.New("not found")}
	uc := NewStaffAuthUsecase(repo, &mockStaffRefreshTokenRepository{})

	_, err := uc.Login(context.Background(), "unknown", "pass")

	assert.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrUnauthorized))
}

func TestStaffAuthUsecase_Login_WrongPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	account := domain.StaffAccount{
		ID:           uuid.New(),
		Username:     "staff1",
		PasswordHash: string(hash),
	}

	repo := &mockStaffAccountRepository{account: account}
	uc := NewStaffAuthUsecase(repo, &mockStaffRefreshTokenRepository{})

	_, err := uc.Login(context.Background(), "staff1", "wrong")

	assert.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrUnauthorized))
}

// ============================================================
// StaffPortalUsecase Tests
// ============================================================

func TestStaffPortalUsecase_GetMyProfileDraft_Exists(t *testing.T) {
	staffID := uuid.New()
	draft := domain.StaffProfileDraft{
		ID:      uuid.New(),
		StaffID: staffID,
		Name:    "Test",
		Status:  domain.DraftStatusDraft,
	}

	draftRepo := &mockStaffDraftRepository{profileDraft: draft}
	staffRepo := &mockStaffRepository{}
	uc := NewStaffPortalUsecase(draftRepo, staffRepo)

	result, err := uc.GetMyProfileDraft(context.Background(), staffID)

	assert.NoError(t, err)
	assert.Equal(t, "Test", result.Name)
	assert.Equal(t, domain.DraftStatusDraft, result.Status)
}

func TestStaffPortalUsecase_GetMyProfileDraft_NotExists_FallbackToStaff(t *testing.T) {
	staffID := uuid.New()
	staff := domain.Staff{ID: staffID, Name: "Current Name", Role: "キャスト"}

	draftRepo := &mockStaffDraftRepository{getByStaffIDErr: errors.New("not found")}
	staffRepo := &mockStaffRepository{staff: staff}
	uc := NewStaffPortalUsecase(draftRepo, staffRepo)

	result, err := uc.GetMyProfileDraft(context.Background(), staffID)

	assert.NoError(t, err)
	assert.Equal(t, "Current Name", result.Name)
	assert.Equal(t, domain.DraftStatus(""), result.Status) // 未作成を示す
}

func TestStaffPortalUsecase_GetMyProfileDraft_StaffNotFound(t *testing.T) {
	staffID := uuid.New()

	draftRepo := &mockStaffDraftRepository{getByStaffIDErr: errors.New("not found")}
	staffRepo := &mockStaffRepository{err: errors.New("staff not found")}
	uc := NewStaffPortalUsecase(draftRepo, staffRepo)

	_, err := uc.GetMyProfileDraft(context.Background(), staffID)

	assert.Error(t, err)
}

func TestStaffPortalUsecase_SaveProfileDraft_Create(t *testing.T) {
	staffID := uuid.New()
	draft := domain.StaffProfileDraft{
		ID:      uuid.New(),
		StaffID: staffID,
		Name:    "New Name",
		Status:  domain.DraftStatusDraft,
	}

	draftRepo := &mockStaffDraftRepository{
		getByStaffIDErr: errors.New("not found"),
		profileDraft:    draft,
	}
	staffRepo := &mockStaffRepository{}
	uc := NewStaffPortalUsecase(draftRepo, staffRepo)

	result, err := uc.SaveProfileDraft(context.Background(), staffID, domain.SaveProfileDraftInput{
		Name: "New Name",
		Role: "キャスト",
	}, time.Time{})

	assert.NoError(t, err)
	assert.Equal(t, "New Name", result.Name)
}

func TestStaffPortalUsecase_SaveProfileDraft_Update(t *testing.T) {
	staffID := uuid.New()
	existingDraft := domain.StaffProfileDraft{
		ID:      uuid.New(),
		StaffID: staffID,
		Name:    "Old Name",
		Status:  domain.DraftStatusDraft,
	}
	updatedDraft := domain.StaffProfileDraft{
		ID:      existingDraft.ID,
		StaffID: staffID,
		Name:    "Updated Name",
		Status:  domain.DraftStatusDraft,
	}

	draftRepo := &mockStaffDraftRepository{profileDraft: updatedDraft}
	staffRepo := &mockStaffRepository{}
	uc := NewStaffPortalUsecase(draftRepo, staffRepo)

	result, err := uc.SaveProfileDraft(context.Background(), staffID, domain.SaveProfileDraftInput{
		Name: "Updated Name",
	}, time.Time{})

	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", result.Name)
}

func TestStaffPortalUsecase_SubmitProfileDraft_Success(t *testing.T) {
	staffID := uuid.New()
	draftID := uuid.New()
	draft := domain.StaffProfileDraft{
		ID:      draftID,
		StaffID: staffID,
		Status:  domain.DraftStatusDraft,
	}

	draftRepo := &mockStaffDraftRepository{profileDraft: draft}
	staffRepo := &mockStaffRepository{}
	uc := NewStaffPortalUsecase(draftRepo, staffRepo)

	result, err := uc.SubmitProfileDraft(context.Background(), staffID, draftID, time.Time{})

	assert.NoError(t, err)
	assert.Equal(t, draftID, result.ID)
}

func TestStaffPortalUsecase_SubmitProfileDraft_Forbidden(t *testing.T) {
	staffID := uuid.New()
	otherStaffID := uuid.New()
	draftID := uuid.New()
	draft := domain.StaffProfileDraft{
		ID:      draftID,
		StaffID: otherStaffID, // 他人の下書き
		Status:  domain.DraftStatusDraft,
	}

	draftRepo := &mockStaffDraftRepository{profileDraft: draft}
	staffRepo := &mockStaffRepository{}
	uc := NewStaffPortalUsecase(draftRepo, staffRepo)

	_, err := uc.SubmitProfileDraft(context.Background(), staffID, draftID, time.Time{})

	assert.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrForbidden))
}

func TestStaffPortalUsecase_SubmitProfileDraft_InvalidStatus(t *testing.T) {
	staffID := uuid.New()
	draftID := uuid.New()
	draft := domain.StaffProfileDraft{
		ID:      draftID,
		StaffID: staffID,
		Status:  domain.DraftStatusPending, // already pending
	}

	draftRepo := &mockStaffDraftRepository{profileDraft: draft}
	staffRepo := &mockStaffRepository{}
	uc := NewStaffPortalUsecase(draftRepo, staffRepo)

	_, err := uc.SubmitProfileDraft(context.Background(), staffID, draftID, time.Time{})

	assert.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
}

func TestStaffPortalUsecase_GetMyScheduleDraft_Exists(t *testing.T) {
	staffID := uuid.New()
	draft := domain.StaffScheduleDraft{
		ID:      uuid.New(),
		StaffID: staffID,
		Status:  domain.DraftStatusDraft,
		Items:   []domain.ScheduleDraftItem{{DayOfWeek: 1}},
	}

	draftRepo := &mockStaffDraftRepository{scheduleDraft: draft}
	staffRepo := &mockStaffRepository{}
	uc := NewStaffPortalUsecase(draftRepo, staffRepo)

	result, err := uc.GetMyScheduleDraft(context.Background(), staffID)

	assert.NoError(t, err)
	assert.Len(t, result.Items, 1)
}

func TestStaffPortalUsecase_GetMyScheduleDraft_NotExists_FallbackToSchedules(t *testing.T) {
	staffID := uuid.New()
	now := time.Now()
	schedules := []domain.StaffSchedule{
		{ID: uuid.New(), StaffID: staffID, DayOfWeek: 2, StartTime: now, EndTime: now},
	}

	draftRepo := &mockStaffDraftRepository{getByStaffIDErr: errors.New("not found")}
	staffRepo := &mockStaffRepository{schedules: schedules}
	uc := NewStaffPortalUsecase(draftRepo, staffRepo)

	result, err := uc.GetMyScheduleDraft(context.Background(), staffID)

	assert.NoError(t, err)
	assert.Equal(t, domain.DraftStatus(""), result.Status)
	assert.Len(t, result.Items, 1)
	assert.Equal(t, 2, result.Items[0].DayOfWeek)
}

func TestStaffPortalUsecase_SaveScheduleDraft_Create(t *testing.T) {
	staffID := uuid.New()
	draft := domain.StaffScheduleDraft{
		ID:      uuid.New(),
		StaffID: staffID,
		Status:  domain.DraftStatusDraft,
	}

	draftRepo := &mockStaffDraftRepository{
		getByStaffIDErr: errors.New("not found"),
		scheduleDraft:   draft,
	}
	staffRepo := &mockStaffRepository{}
	uc := NewStaffPortalUsecase(draftRepo, staffRepo)

	result, err := uc.SaveScheduleDraft(context.Background(), staffID, []domain.ScheduleDraftItem{}, time.Time{})

	assert.NoError(t, err)
	assert.Equal(t, domain.DraftStatusDraft, result.Status)
}

func TestStaffPortalUsecase_SaveScheduleDraft_Update(t *testing.T) {
	staffID := uuid.New()
	existing := domain.StaffScheduleDraft{
		ID:      uuid.New(),
		StaffID: staffID,
		Status:  domain.DraftStatusDraft,
	}

	draftRepo := &mockStaffDraftRepository{scheduleDraft: existing}
	staffRepo := &mockStaffRepository{}
	uc := NewStaffPortalUsecase(draftRepo, staffRepo)

	result, err := uc.SaveScheduleDraft(context.Background(), staffID, []domain.ScheduleDraftItem{}, time.Time{})

	assert.NoError(t, err)
	assert.Equal(t, staffID, result.StaffID)
}

func TestStaffPortalUsecase_SubmitScheduleDraft_Success(t *testing.T) {
	staffID := uuid.New()
	draftID := uuid.New()
	draft := domain.StaffScheduleDraft{
		ID:      draftID,
		StaffID: staffID,
		Status:  domain.DraftStatusDraft,
	}

	draftRepo := &mockStaffDraftRepository{scheduleDraft: draft}
	staffRepo := &mockStaffRepository{}
	uc := NewStaffPortalUsecase(draftRepo, staffRepo)

	result, err := uc.SubmitScheduleDraft(context.Background(), staffID, draftID, time.Time{})

	assert.NoError(t, err)
	assert.Equal(t, draftID, result.ID)
}

func TestStaffPortalUsecase_SubmitScheduleDraft_Forbidden(t *testing.T) {
	staffID := uuid.New()
	draftID := uuid.New()
	draft := domain.StaffScheduleDraft{
		ID:      draftID,
		StaffID: uuid.New(), // 他のスタッフの下書き
		Status:  domain.DraftStatusDraft,
	}

	draftRepo := &mockStaffDraftRepository{scheduleDraft: draft}
	staffRepo := &mockStaffRepository{}
	uc := NewStaffPortalUsecase(draftRepo, staffRepo)

	_, err := uc.SubmitScheduleDraft(context.Background(), staffID, draftID, time.Time{})

	assert.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrForbidden))
}

func TestStaffPortalUsecase_SubmitScheduleDraft_InvalidStatus(t *testing.T) {
	staffID := uuid.New()
	draftID := uuid.New()
	draft := domain.StaffScheduleDraft{
		ID:      draftID,
		StaffID: staffID,
		Status:  domain.DraftStatusApproved,
	}

	draftRepo := &mockStaffDraftRepository{scheduleDraft: draft}
	staffRepo := &mockStaffRepository{}
	uc := NewStaffPortalUsecase(draftRepo, staffRepo)

	_, err := uc.SubmitScheduleDraft(context.Background(), staffID, draftID, time.Time{})

	assert.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
}

// ============================================================
// AdminReviewUsecase Tests
// ============================================================

func TestAdminReviewUsecase_ListPendingProfileDrafts_Success(t *testing.T) {
	drafts := []domain.StaffProfileDraft{
		{ID: uuid.New(), Name: "Draft 1"},
		{ID: uuid.New(), Name: "Draft 2"},
	}

	draftRepo := &mockStaffDraftRepository{profileDrafts: drafts}
	staffRepo := &mockStaffRepository{}
	uc := NewAdminReviewUsecase(draftRepo, staffRepo)

	result, err := uc.ListPendingProfileDrafts(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestAdminReviewUsecase_ListPendingScheduleDrafts_Success(t *testing.T) {
	drafts := []domain.StaffScheduleDraft{
		{ID: uuid.New()},
	}

	draftRepo := &mockStaffDraftRepository{scheduleDrafts: drafts}
	staffRepo := &mockStaffRepository{}
	uc := NewAdminReviewUsecase(draftRepo, staffRepo)

	result, err := uc.ListPendingScheduleDrafts(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestAdminReviewUsecase_ListApprovedScheduleDrafts_Success(t *testing.T) {
	drafts := []domain.StaffScheduleDraft{
		{ID: uuid.New(), Status: domain.DraftStatusApproved},
	}

	draftRepo := &mockStaffDraftRepository{scheduleDrafts: drafts}
	staffRepo := &mockStaffRepository{}
	uc := NewAdminReviewUsecase(draftRepo, staffRepo)

	result, err := uc.ListApprovedScheduleDrafts(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestAdminReviewUsecase_ReviewProfileDraft_Approve(t *testing.T) {
	staffID := uuid.New()
	draftID := uuid.New()
	draft := domain.StaffProfileDraft{
		ID:                draftID,
		StaffID:           staffID,
		Name:              "New Name",
		Role:              "キャスト",
		Status:            domain.DraftStatusPending,
		ImageCropPosition: "50 50",
	}

	draftRepo := &mockStaffDraftRepository{profileDraft: draft}
	staffRepo := &mockStaffRepository{staff: domain.Staff{ID: staffID}}
	uc := NewAdminReviewUsecase(draftRepo, staffRepo)

	result, err := uc.ReviewProfileDraft(context.Background(), draftID, domain.ReviewDraftInput{
		Status: domain.DraftStatusApproved,
	}, time.Time{})

	assert.NoError(t, err)
	assert.Equal(t, draftID, result.ID)
}

func TestAdminReviewUsecase_ReviewProfileDraft_InvalidStatus(t *testing.T) {
	draftRepo := &mockStaffDraftRepository{}
	staffRepo := &mockStaffRepository{}
	uc := NewAdminReviewUsecase(draftRepo, staffRepo)

	_, err := uc.ReviewProfileDraft(context.Background(), uuid.New(), domain.ReviewDraftInput{
		Status: domain.DraftStatus("invalid"),
	}, time.Time{})

	assert.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
}

func TestAdminReviewUsecase_ReviewProfileDraft_NotPending(t *testing.T) {
	draft := domain.StaffProfileDraft{
		ID:     uuid.New(),
		Status: domain.DraftStatusDraft, // not pending
	}

	draftRepo := &mockStaffDraftRepository{profileDraft: draft}
	staffRepo := &mockStaffRepository{}
	uc := NewAdminReviewUsecase(draftRepo, staffRepo)

	_, err := uc.ReviewProfileDraft(context.Background(), draft.ID, domain.ReviewDraftInput{
		Status: domain.DraftStatusApproved,
	}, time.Time{})

	assert.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
}

func TestAdminReviewUsecase_ReviewProfileDraft_NotFound(t *testing.T) {
	draftRepo := &mockStaffDraftRepository{err: errors.New("not found")}
	staffRepo := &mockStaffRepository{}
	uc := NewAdminReviewUsecase(draftRepo, staffRepo)

	_, err := uc.ReviewProfileDraft(context.Background(), uuid.New(), domain.ReviewDraftInput{
		Status: domain.DraftStatusApproved,
	}, time.Time{})

	assert.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
}

func TestAdminReviewUsecase_ReviewScheduleDraft_Approve(t *testing.T) {
	draftID := uuid.New()
	draft := domain.StaffScheduleDraft{
		ID:     draftID,
		Status: domain.DraftStatusPending,
	}

	draftRepo := &mockStaffDraftRepository{scheduleDraft: draft}
	staffRepo := &mockStaffRepository{}
	uc := NewAdminReviewUsecase(draftRepo, staffRepo)

	result, err := uc.ReviewScheduleDraft(context.Background(), draftID, domain.ReviewDraftInput{
		Status: domain.DraftStatusApproved,
	}, time.Time{})

	assert.NoError(t, err)
	assert.Equal(t, draftID, result.ID)
}

func TestAdminReviewUsecase_ReviewScheduleDraft_InvalidStatus(t *testing.T) {
	draftRepo := &mockStaffDraftRepository{}
	staffRepo := &mockStaffRepository{}
	uc := NewAdminReviewUsecase(draftRepo, staffRepo)

	_, err := uc.ReviewScheduleDraft(context.Background(), uuid.New(), domain.ReviewDraftInput{
		Status: domain.DraftStatusPending, // pending is not a valid review action
	}, time.Time{})

	assert.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
}

func TestAdminReviewUsecase_ReviewScheduleDraft_NotPending(t *testing.T) {
	draft := domain.StaffScheduleDraft{
		ID:     uuid.New(),
		Status: domain.DraftStatusDraft,
	}

	draftRepo := &mockStaffDraftRepository{scheduleDraft: draft}
	staffRepo := &mockStaffRepository{}
	uc := NewAdminReviewUsecase(draftRepo, staffRepo)

	_, err := uc.ReviewScheduleDraft(context.Background(), draft.ID, domain.ReviewDraftInput{
		Status: domain.DraftStatusRejected,
	}, time.Time{})

	assert.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
}

func TestAdminReviewUsecase_PublishScheduleDraft_Success(t *testing.T) {
	draftID := uuid.New()
	staffID := uuid.New()
	now := time.Now()
	draft := domain.StaffScheduleDraft{
		ID:      draftID,
		StaffID: staffID,
		Status:  domain.DraftStatusApproved,
		Items: []domain.ScheduleDraftItem{
			{DayOfWeek: 1, StartTime: now, EndTime: now},
		},
	}

	draftRepo := &mockStaffDraftRepository{scheduleDraft: draft}
	staffRepo := &mockStaffRepository{}
	uc := NewAdminReviewUsecase(draftRepo, staffRepo)

	err := uc.PublishScheduleDraft(context.Background(), draftID)

	assert.NoError(t, err)
}

func TestAdminReviewUsecase_PublishScheduleDraft_NotApproved(t *testing.T) {
	draft := domain.StaffScheduleDraft{
		ID:     uuid.New(),
		Status: domain.DraftStatusPending,
	}

	draftRepo := &mockStaffDraftRepository{scheduleDraft: draft}
	staffRepo := &mockStaffRepository{}
	uc := NewAdminReviewUsecase(draftRepo, staffRepo)

	err := uc.PublishScheduleDraft(context.Background(), draft.ID)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
}

func TestAdminReviewUsecase_PublishScheduleDraft_NotFound(t *testing.T) {
	draftRepo := &mockStaffDraftRepository{err: errors.New("not found")}
	staffRepo := &mockStaffRepository{}
	uc := NewAdminReviewUsecase(draftRepo, staffRepo)

	err := uc.PublishScheduleDraft(context.Background(), uuid.New())

	assert.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
}

func TestAdminReviewUsecase_GetProfileDraft_Success(t *testing.T) {
	draft := domain.StaffProfileDraft{ID: uuid.New(), Name: "Test"}

	draftRepo := &mockStaffDraftRepository{profileDraft: draft}
	staffRepo := &mockStaffRepository{}
	uc := NewAdminReviewUsecase(draftRepo, staffRepo)

	result, err := uc.GetProfileDraft(context.Background(), draft.ID)

	assert.NoError(t, err)
	assert.Equal(t, "Test", result.Name)
}

func TestAdminReviewUsecase_GetScheduleDraft_Success(t *testing.T) {
	draft := domain.StaffScheduleDraft{ID: uuid.New(), Status: domain.DraftStatusPending}

	draftRepo := &mockStaffDraftRepository{scheduleDraft: draft}
	staffRepo := &mockStaffRepository{}
	uc := NewAdminReviewUsecase(draftRepo, staffRepo)

	result, err := uc.GetScheduleDraft(context.Background(), draft.ID)

	assert.NoError(t, err)
	assert.Equal(t, domain.DraftStatusPending, result.Status)
}

func TestAdminReviewUsecase_UpdateProfileDraftContent_Success(t *testing.T) {
	draft := domain.StaffProfileDraft{ID: uuid.New(), Name: "Admin Updated"}

	draftRepo := &mockStaffDraftRepository{profileDraft: draft}
	staffRepo := &mockStaffRepository{}
	uc := NewAdminReviewUsecase(draftRepo, staffRepo)

	result, err := uc.UpdateProfileDraftContent(context.Background(), draft.ID, domain.SaveProfileDraftInput{Name: "Admin Updated"}, time.Time{})

	assert.NoError(t, err)
	assert.Equal(t, "Admin Updated", result.Name)
}

func TestAdminReviewUsecase_UpdateScheduleDraftContent_Success(t *testing.T) {
	draft := domain.StaffScheduleDraft{ID: uuid.New()}

	draftRepo := &mockStaffDraftRepository{scheduleDraft: draft}
	staffRepo := &mockStaffRepository{}
	uc := NewAdminReviewUsecase(draftRepo, staffRepo)

	result, err := uc.UpdateScheduleDraftContent(context.Background(), draft.ID, []domain.ScheduleDraftItem{}, time.Time{})

	assert.NoError(t, err)
	assert.Equal(t, draft.ID, result.ID)
}

func TestAdminReviewUsecase_GetStaffName_Success(t *testing.T) {
	staffID := uuid.New()
	staff := domain.Staff{ID: staffID, Name: "テスト太郎"}

	draftRepo := &mockStaffDraftRepository{}
	staffRepo := &mockStaffRepository{staff: staff}
	uc := NewAdminReviewUsecase(draftRepo, staffRepo)

	name, err := uc.GetStaffName(context.Background(), staffID)

	assert.NoError(t, err)
	assert.Equal(t, "テスト太郎", name)
}

func TestAdminReviewUsecase_GetStaffName_NotFound(t *testing.T) {
	draftRepo := &mockStaffDraftRepository{}
	staffRepo := &mockStaffRepository{err: errors.New("not found")}
	uc := NewAdminReviewUsecase(draftRepo, staffRepo)

	_, err := uc.GetStaffName(context.Background(), uuid.New())

	assert.Error(t, err)
}

func TestAdminReviewUsecase_ListImagesByStaffID_Success(t *testing.T) {
	images := []domain.StaffImage{
		{ID: uuid.New(), ImageURL: "/img1.jpg"},
		{ID: uuid.New(), ImageURL: "/img2.jpg"},
	}

	draftRepo := &mockStaffDraftRepository{}
	staffRepo := &mockStaffRepository{images: images}
	uc := NewAdminReviewUsecase(draftRepo, staffRepo)

	result, err := uc.ListImagesByStaffID(context.Background(), uuid.New())

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

// ============================================================
// AdminAccountUsecase Tests
// ============================================================

func TestAdminAccountUsecase_ListStaffAccounts_Success(t *testing.T) {
	accounts := []domain.StaffAccount{
		{ID: uuid.New(), Username: "staff1"},
		{ID: uuid.New(), Username: "staff2"},
	}

	repo := &mockStaffAccountRepository{accounts: accounts}
	uc := NewAdminAccountUsecase(repo)

	result, err := uc.ListStaffAccounts(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestAdminAccountUsecase_GetStaffAccountByStaffID_Success(t *testing.T) {
	account := domain.StaffAccount{ID: uuid.New(), Username: "staff1"}

	repo := &mockStaffAccountRepository{account: account}
	uc := NewAdminAccountUsecase(repo)

	result, err := uc.GetStaffAccountByStaffID(context.Background(), uuid.New())

	assert.NoError(t, err)
	assert.Equal(t, "staff1", result.Username)
}

func TestAdminAccountUsecase_CreateStaffAccount_Success(t *testing.T) {
	account := domain.StaffAccount{
		ID:       uuid.New(),
		StaffID:  uuid.New(),
		Username: "newstaff",
	}

	repo := &mockStaffAccountRepository{account: account}
	uc := NewAdminAccountUsecase(repo)

	result, err := uc.CreateStaffAccount(context.Background(), account.StaffID, "newstaff", "password123")

	assert.NoError(t, err)
	assert.Equal(t, "newstaff", result.Username)
}

func TestAdminAccountUsecase_CreateStaffAccount_EmptyUsername(t *testing.T) {
	repo := &mockStaffAccountRepository{}
	uc := NewAdminAccountUsecase(repo)

	_, err := uc.CreateStaffAccount(context.Background(), uuid.New(), "", "password")

	assert.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
}

func TestAdminAccountUsecase_CreateStaffAccount_EmptyPassword(t *testing.T) {
	repo := &mockStaffAccountRepository{}
	uc := NewAdminAccountUsecase(repo)

	_, err := uc.CreateStaffAccount(context.Background(), uuid.New(), "user", "")

	assert.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
}

func TestAdminAccountUsecase_UpdateStaffAccount_WithNewPassword(t *testing.T) {
	account := domain.StaffAccount{ID: uuid.New(), Username: "updated"}

	repo := &mockStaffAccountRepository{account: account}
	uc := NewAdminAccountUsecase(repo)

	result, err := uc.UpdateStaffAccount(context.Background(), account.ID, "updated", "newpass")

	assert.NoError(t, err)
	assert.Equal(t, "updated", result.Username)
}

func TestAdminAccountUsecase_UpdateStaffAccount_KeepPassword(t *testing.T) {
	accountID := uuid.New()
	accounts := []domain.StaffAccount{
		{ID: accountID, Username: "staff1", PasswordHash: "$2a$10$existing"},
	}

	repo := &mockStaffAccountRepository{accounts: accounts, account: domain.StaffAccount{ID: accountID, Username: "renamed"}}
	uc := NewAdminAccountUsecase(repo)

	result, err := uc.UpdateStaffAccount(context.Background(), accountID, "renamed", "")

	assert.NoError(t, err)
	assert.Equal(t, "renamed", result.Username)
}

func TestAdminAccountUsecase_UpdateStaffAccount_EmptyUsername(t *testing.T) {
	repo := &mockStaffAccountRepository{}
	uc := NewAdminAccountUsecase(repo)

	_, err := uc.UpdateStaffAccount(context.Background(), uuid.New(), "", "pass")

	assert.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrInvalidInput))
}

func TestAdminAccountUsecase_UpdateStaffAccount_AccountNotFound(t *testing.T) {
	repo := &mockStaffAccountRepository{accounts: []domain.StaffAccount{}}
	uc := NewAdminAccountUsecase(repo)

	_, err := uc.UpdateStaffAccount(context.Background(), uuid.New(), "user", "")

	assert.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
}

func TestAdminAccountUsecase_DeleteStaffAccount_Success(t *testing.T) {
	repo := &mockStaffAccountRepository{}
	uc := NewAdminAccountUsecase(repo)

	err := uc.DeleteStaffAccount(context.Background(), uuid.New())

	assert.NoError(t, err)
}

func TestAdminAccountUsecase_DeleteStaffAccount_Error(t *testing.T) {
	repo := &mockStaffAccountRepository{err: errors.New("db error")}
	uc := NewAdminAccountUsecase(repo)

	err := uc.DeleteStaffAccount(context.Background(), uuid.New())

	assert.Error(t, err)
}
