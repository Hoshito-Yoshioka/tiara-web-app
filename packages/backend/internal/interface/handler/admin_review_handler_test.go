package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"tiara-web-app/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// mockAdminReviewUsecase は AdminReviewUsecase のテスト用モック。
type mockAdminReviewUsecase struct {
	profileDrafts  []domain.StaffProfileDraft
	profileDraft   domain.StaffProfileDraft
	scheduleDrafts []domain.StaffScheduleDraft
	scheduleDraft  domain.StaffScheduleDraft
	staffName      string
	images         []domain.StaffImage
	err            error
	publishErr     error
}

func (m *mockAdminReviewUsecase) ListPendingProfileDrafts(_ context.Context) ([]domain.StaffProfileDraft, error) {
	return m.profileDrafts, m.err
}
func (m *mockAdminReviewUsecase) ListPendingScheduleDrafts(_ context.Context) ([]domain.StaffScheduleDraft, error) {
	return m.scheduleDrafts, m.err
}
func (m *mockAdminReviewUsecase) ListApprovedScheduleDrafts(_ context.Context) ([]domain.StaffScheduleDraft, error) {
	return m.scheduleDrafts, m.err
}
func (m *mockAdminReviewUsecase) ReviewProfileDraft(_ context.Context, _ uuid.UUID, _ domain.ReviewDraftInput) (domain.StaffProfileDraft, error) {
	return m.profileDraft, m.err
}
func (m *mockAdminReviewUsecase) ReviewScheduleDraft(_ context.Context, _ uuid.UUID, _ domain.ReviewDraftInput) (domain.StaffScheduleDraft, error) {
	return m.scheduleDraft, m.err
}
func (m *mockAdminReviewUsecase) PublishScheduleDraft(_ context.Context, _ uuid.UUID) error {
	return m.publishErr
}
func (m *mockAdminReviewUsecase) GetProfileDraft(_ context.Context, _ uuid.UUID) (domain.StaffProfileDraft, error) {
	return m.profileDraft, m.err
}
func (m *mockAdminReviewUsecase) GetScheduleDraft(_ context.Context, _ uuid.UUID) (domain.StaffScheduleDraft, error) {
	return m.scheduleDraft, m.err
}
func (m *mockAdminReviewUsecase) UpdateProfileDraftContent(_ context.Context, _ uuid.UUID, _ domain.SaveProfileDraftInput) (domain.StaffProfileDraft, error) {
	return m.profileDraft, m.err
}
func (m *mockAdminReviewUsecase) UpdateScheduleDraftContent(_ context.Context, _ uuid.UUID, _ []domain.ScheduleDraftItem) (domain.StaffScheduleDraft, error) {
	return m.scheduleDraft, m.err
}
func (m *mockAdminReviewUsecase) GetStaffName(_ context.Context, _ uuid.UUID) (string, error) {
	return m.staffName, nil
}
func (m *mockAdminReviewUsecase) ListImagesByStaffID(_ context.Context, _ uuid.UUID) ([]domain.StaffImage, error) {
	return m.images, nil
}

// --- Profile Draft Tests ---

func TestAdminReviewHandler_ListPendingProfileDrafts(t *testing.T) {
	staffID := uuid.New()

	t.Run("正常系: 承認待ちプロフィール一覧を返す", func(t *testing.T) {
		now := time.Now()
		mock := &mockAdminReviewUsecase{
			profileDrafts: []domain.StaffProfileDraft{
				{ID: uuid.New(), StaffID: staffID, Name: "Yuki", Status: domain.DraftStatusPending, SubmittedAt: &now, CreatedAt: now},
			},
			staffName: "Yuki",
			images:    []domain.StaffImage{{ID: uuid.New(), StaffID: staffID, ImageURL: "/img.jpg"}},
		}
		h := NewAdminReviewHandler(mock)

		req := httptest.NewRequest(http.MethodGet, "/admin/reviews/profiles", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := h.ListPendingProfileDrafts(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp []PendingProfileDraftResponse
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp, 1)
		assert.Equal(t, "Yuki", resp[0].StaffName)
		assert.Len(t, resp[0].Images, 1)
	})

	t.Run("異常系: usecase エラー → 500", func(t *testing.T) {
		mock := &mockAdminReviewUsecase{err: domain.ErrInternal}
		h := NewAdminReviewHandler(mock)

		req := httptest.NewRequest(http.MethodGet, "/admin/reviews/profiles", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		_ = h.ListPendingProfileDrafts(c)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestAdminReviewHandler_ReviewProfileDraft(t *testing.T) {
	draftID := uuid.New()

	t.Run("正常系: プロフィール承認 → 200", func(t *testing.T) {
		mock := &mockAdminReviewUsecase{
			profileDraft: domain.StaffProfileDraft{ID: draftID, Status: domain.DraftStatusApproved},
		}
		h := NewAdminReviewHandler(mock)

		body := `{"status":"approved","adminComment":"OK"}`
		req := httptest.NewRequest(http.MethodPut, "/admin/reviews/profiles/:id", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(draftID.String())

		err := h.ReviewProfileDraft(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("異常系: 不正なドラフトID → 400", func(t *testing.T) {
		h := NewAdminReviewHandler(&mockAdminReviewUsecase{})

		body := `{"status":"approved","adminComment":"OK"}`
		req := httptest.NewRequest(http.MethodPut, "/admin/reviews/profiles/:id", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("bad")

		_ = h.ReviewProfileDraft(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: 不正なボディ → 400", func(t *testing.T) {
		h := NewAdminReviewHandler(&mockAdminReviewUsecase{})

		req := httptest.NewRequest(http.MethodPut, "/admin/reviews/profiles/:id", strings.NewReader("bad"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(draftID.String())

		_ = h.ReviewProfileDraft(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: 存在しないドラフト → 404", func(t *testing.T) {
		mock := &mockAdminReviewUsecase{err: domain.ErrNotFound}
		h := NewAdminReviewHandler(mock)

		body := `{"status":"approved","adminComment":"OK"}`
		req := httptest.NewRequest(http.MethodPut, "/admin/reviews/profiles/:id", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(uuid.New().String())

		_ = h.ReviewProfileDraft(c)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestAdminReviewHandler_GetProfileDraft(t *testing.T) {
	draftID := uuid.New()

	t.Run("正常系: プロフィール下書き取得 → 200", func(t *testing.T) {
		mock := &mockAdminReviewUsecase{
			profileDraft: domain.StaffProfileDraft{ID: draftID, Name: "Yuki"},
		}
		h := NewAdminReviewHandler(mock)

		req := httptest.NewRequest(http.MethodGet, "/admin/reviews/profiles/:id", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(draftID.String())

		err := h.GetProfileDraft(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("異常系: 不正なID → 400", func(t *testing.T) {
		h := NewAdminReviewHandler(&mockAdminReviewUsecase{})

		req := httptest.NewRequest(http.MethodGet, "/admin/reviews/profiles/:id", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("bad")

		_ = h.GetProfileDraft(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestAdminReviewHandler_UpdateProfileDraftContent(t *testing.T) {
	draftID := uuid.New()

	t.Run("正常系: プロフィール内容修正 → 200", func(t *testing.T) {
		mock := &mockAdminReviewUsecase{
			profileDraft: domain.StaffProfileDraft{ID: draftID, Name: "Updated"},
		}
		h := NewAdminReviewHandler(mock)

		body := `{"name":"Updated","role":"Cast","bio":"bio","imageUrl":"","imageCropPosition":""}`
		req := httptest.NewRequest(http.MethodPut, "/admin/reviews/profiles/:id/content", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(draftID.String())

		err := h.UpdateProfileDraftContent(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("異常系: 不正なID → 400", func(t *testing.T) {
		h := NewAdminReviewHandler(&mockAdminReviewUsecase{})

		body := `{"name":"X","role":"Cast","bio":"","imageUrl":"","imageCropPosition":""}`
		req := httptest.NewRequest(http.MethodPut, "/admin/reviews/profiles/:id/content", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("bad")

		_ = h.UpdateProfileDraftContent(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: 不正なボディ → 400", func(t *testing.T) {
		h := NewAdminReviewHandler(&mockAdminReviewUsecase{})

		req := httptest.NewRequest(http.MethodPut, "/admin/reviews/profiles/:id/content", strings.NewReader("bad"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(draftID.String())

		_ = h.UpdateProfileDraftContent(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

// --- Schedule Draft Tests ---

func TestAdminReviewHandler_ListPendingScheduleDrafts(t *testing.T) {
	staffID := uuid.New()

	t.Run("正常系: 承認待ちスケジュール一覧を返す", func(t *testing.T) {
		now := time.Now()
		startTime, _ := time.Parse("15:04", "20:00")
		endTime, _ := time.Parse("15:04", "02:00")
		mock := &mockAdminReviewUsecase{
			scheduleDrafts: []domain.StaffScheduleDraft{
				{
					ID: uuid.New(), StaffID: staffID, Status: domain.DraftStatusPending,
					SubmittedAt: &now, CreatedAt: now,
					Items: []domain.ScheduleDraftItem{
						{ID: uuid.New(), DayOfWeek: 1, StartTime: startTime, EndTime: endTime},
					},
				},
			},
			staffName: "Yuki",
		}
		h := NewAdminReviewHandler(mock)

		req := httptest.NewRequest(http.MethodGet, "/admin/reviews/schedules", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := h.ListPendingScheduleDrafts(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp []PendingScheduleDraftResponse
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp, 1)
		assert.Equal(t, "Yuki", resp[0].StaffName)
		assert.Len(t, resp[0].Items, 1)
		assert.Equal(t, "20:00", resp[0].Items[0].StartTime)
	})

	t.Run("異常系: usecase エラー → 500", func(t *testing.T) {
		mock := &mockAdminReviewUsecase{err: domain.ErrInternal}
		h := NewAdminReviewHandler(mock)

		req := httptest.NewRequest(http.MethodGet, "/admin/reviews/schedules", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		_ = h.ListPendingScheduleDrafts(c)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestAdminReviewHandler_ListApprovedScheduleDrafts(t *testing.T) {
	t.Run("正常系: 承認済みスケジュール一覧を返す", func(t *testing.T) {
		mock := &mockAdminReviewUsecase{
			scheduleDrafts: []domain.StaffScheduleDraft{},
			staffName:      "Yuki",
		}
		h := NewAdminReviewHandler(mock)

		req := httptest.NewRequest(http.MethodGet, "/admin/reviews/schedules/approved", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := h.ListApprovedScheduleDrafts(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestAdminReviewHandler_PublishScheduleDraft(t *testing.T) {
	draftID := uuid.New()

	t.Run("正常系: スケジュール公開 → 200", func(t *testing.T) {
		mock := &mockAdminReviewUsecase{}
		h := NewAdminReviewHandler(mock)

		req := httptest.NewRequest(http.MethodPut, "/admin/reviews/schedules/:id/publish", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(draftID.String())

		err := h.PublishScheduleDraft(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("異常系: 不正なドラフトID → 400", func(t *testing.T) {
		h := NewAdminReviewHandler(&mockAdminReviewUsecase{})

		req := httptest.NewRequest(http.MethodPut, "/admin/reviews/schedules/:id/publish", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("bad")

		_ = h.PublishScheduleDraft(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: 存在しないドラフト → 404", func(t *testing.T) {
		mock := &mockAdminReviewUsecase{publishErr: domain.ErrNotFound}
		h := NewAdminReviewHandler(mock)

		req := httptest.NewRequest(http.MethodPut, "/admin/reviews/schedules/:id/publish", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(uuid.New().String())

		_ = h.PublishScheduleDraft(c)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestAdminReviewHandler_ReviewScheduleDraft(t *testing.T) {
	draftID := uuid.New()

	t.Run("正常系: スケジュール承認 → 200", func(t *testing.T) {
		mock := &mockAdminReviewUsecase{
			scheduleDraft: domain.StaffScheduleDraft{ID: draftID, Status: domain.DraftStatusApproved},
		}
		h := NewAdminReviewHandler(mock)

		body := `{"status":"approved","adminComment":"OK"}`
		req := httptest.NewRequest(http.MethodPut, "/admin/reviews/schedules/:id", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(draftID.String())

		err := h.ReviewScheduleDraft(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("異常系: 不正なID → 400", func(t *testing.T) {
		h := NewAdminReviewHandler(&mockAdminReviewUsecase{})

		body := `{"status":"approved","adminComment":"OK"}`
		req := httptest.NewRequest(http.MethodPut, "/admin/reviews/schedules/:id", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("bad")

		_ = h.ReviewScheduleDraft(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: 不正なボディ → 400", func(t *testing.T) {
		h := NewAdminReviewHandler(&mockAdminReviewUsecase{})

		req := httptest.NewRequest(http.MethodPut, "/admin/reviews/schedules/:id", strings.NewReader("bad"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(draftID.String())

		_ = h.ReviewScheduleDraft(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestAdminReviewHandler_GetScheduleDraft(t *testing.T) {
	draftID := uuid.New()

	t.Run("正常系: スケジュール下書き取得 → 200", func(t *testing.T) {
		mock := &mockAdminReviewUsecase{
			scheduleDraft: domain.StaffScheduleDraft{ID: draftID},
		}
		h := NewAdminReviewHandler(mock)

		req := httptest.NewRequest(http.MethodGet, "/admin/reviews/schedules/:id", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(draftID.String())

		err := h.GetScheduleDraft(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("異常系: 不正なID → 400", func(t *testing.T) {
		h := NewAdminReviewHandler(&mockAdminReviewUsecase{})

		req := httptest.NewRequest(http.MethodGet, "/admin/reviews/schedules/:id", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("bad")

		_ = h.GetScheduleDraft(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestAdminReviewHandler_UpdateScheduleDraftContent(t *testing.T) {
	draftID := uuid.New()

	t.Run("正常系: スケジュール内容修正 → 200", func(t *testing.T) {
		mock := &mockAdminReviewUsecase{
			scheduleDraft: domain.StaffScheduleDraft{ID: draftID},
		}
		h := NewAdminReviewHandler(mock)

		body := `{"items":[{"dayOfWeek":1,"startTime":"20:00","endTime":"02:00"}]}`
		req := httptest.NewRequest(http.MethodPut, "/admin/reviews/schedules/:id/content", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(draftID.String())

		err := h.UpdateScheduleDraftContent(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("異常系: 不正なID → 400", func(t *testing.T) {
		h := NewAdminReviewHandler(&mockAdminReviewUsecase{})

		body := `{"items":[]}`
		req := httptest.NewRequest(http.MethodPut, "/admin/reviews/schedules/:id/content", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("bad")

		_ = h.UpdateScheduleDraftContent(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: 不正なボディ → 400", func(t *testing.T) {
		h := NewAdminReviewHandler(&mockAdminReviewUsecase{})

		req := httptest.NewRequest(http.MethodPut, "/admin/reviews/schedules/:id/content", strings.NewReader("bad"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(draftID.String())

		_ = h.UpdateScheduleDraftContent(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: 不正な時刻フォーマット → 400", func(t *testing.T) {
		h := NewAdminReviewHandler(&mockAdminReviewUsecase{})

		body := `{"items":[{"dayOfWeek":1,"startTime":"invalid","endTime":"02:00"}]}`
		req := httptest.NewRequest(http.MethodPut, "/admin/reviews/schedules/:id/content", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(draftID.String())

		_ = h.UpdateScheduleDraftContent(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
