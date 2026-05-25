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
	"tiara-web-app/backend/internal/testutil"
	"tiara-web-app/backend/internal/usecase"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// --- Mock: StaffAuthUsecase ---

type mockStaffAuthUsecase struct {
	account domain.StaffAccount
	err     error
}

func (m *mockStaffAuthUsecase) Login(_ context.Context, _, _ string) (domain.StaffAccount, error) {
	return m.account, m.err
}

// --- Mock: StaffPortalUsecase ---

type mockStaffPortalUsecase struct {
	profileDraft  domain.StaffProfileDraft
	scheduleDraft domain.StaffScheduleDraft
	err           error
}

func (m *mockStaffPortalUsecase) GetMyProfileDraft(_ context.Context, _ uuid.UUID) (domain.StaffProfileDraft, error) {
	return m.profileDraft, m.err
}
func (m *mockStaffPortalUsecase) SaveProfileDraft(_ context.Context, _ uuid.UUID, _ domain.SaveProfileDraftInput) (domain.StaffProfileDraft, error) {
	return m.profileDraft, m.err
}
func (m *mockStaffPortalUsecase) SubmitProfileDraft(_ context.Context, _ uuid.UUID, _ uuid.UUID) (domain.StaffProfileDraft, error) {
	return m.profileDraft, m.err
}
func (m *mockStaffPortalUsecase) GetMyScheduleDraft(_ context.Context, _ uuid.UUID) (domain.StaffScheduleDraft, error) {
	return m.scheduleDraft, m.err
}
func (m *mockStaffPortalUsecase) SaveScheduleDraft(_ context.Context, _ uuid.UUID, _ []domain.ScheduleDraftItem) (domain.StaffScheduleDraft, error) {
	return m.scheduleDraft, m.err
}
func (m *mockStaffPortalUsecase) SubmitScheduleDraft(_ context.Context, _ uuid.UUID, _ uuid.UUID) (domain.StaffScheduleDraft, error) {
	return m.scheduleDraft, m.err
}

// newPortalHandler はテスト用の StaffPortalHandler を生成するヘルパー。
func newPortalHandler(auth usecase.StaffAuthUsecase, portal usecase.StaffPortalUsecase, staff usecase.StaffUsecase) *StaffPortalHandler {
	return NewStaffPortalHandler(auth, portal, staff, "test-secret", 2, "/tmp/test-uploads")
}

// --- Auth Tests ---

func TestStaffPortalHandler_Login(t *testing.T) {
	staffID := uuid.New()

	t.Run("正常系: ログイン成功 → トークン返却", func(t *testing.T) {
		mock := &mockStaffAuthUsecase{
			account: domain.StaffAccount{ID: uuid.New(), StaffID: staffID, Username: "staff1"},
		}
		h := newPortalHandler(mock, &mockStaffPortalUsecase{}, &mockStaffUsecase{})

		body := `{"username":"staff1","password":"pass"}`
		req := httptest.NewRequest(http.MethodPost, "/staff-auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := h.Login(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp StaffLoginResponse
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
		assert.Equal(t, staffID.String(), resp.StaffID)
	})

	t.Run("異常系: 不正なボディ → 400", func(t *testing.T) {
		h := newPortalHandler(&mockStaffAuthUsecase{}, &mockStaffPortalUsecase{}, &mockStaffUsecase{})

		req := httptest.NewRequest(http.MethodPost, "/staff-auth/login", strings.NewReader("bad"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		_ = h.Login(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: ユーザー名空 → 400", func(t *testing.T) {
		h := newPortalHandler(&mockStaffAuthUsecase{}, &mockStaffPortalUsecase{}, &mockStaffUsecase{})

		body := `{"username":"","password":"pass"}`
		req := httptest.NewRequest(http.MethodPost, "/staff-auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		_ = h.Login(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: パスワード空 → 400", func(t *testing.T) {
		h := newPortalHandler(&mockStaffAuthUsecase{}, &mockStaffPortalUsecase{}, &mockStaffUsecase{})

		body := `{"username":"staff1","password":""}`
		req := httptest.NewRequest(http.MethodPost, "/staff-auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		_ = h.Login(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: 認証失敗 → 401", func(t *testing.T) {
		mock := &mockStaffAuthUsecase{err: domain.ErrUnauthorized}
		h := newPortalHandler(mock, &mockStaffPortalUsecase{}, &mockStaffUsecase{})

		body := `{"username":"staff1","password":"wrong"}`
		req := httptest.NewRequest(http.MethodPost, "/staff-auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		_ = h.Login(c)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestStaffPortalHandler_Verify(t *testing.T) {
	staffID := uuid.New()

	t.Run("正常系: トークン有効 → 200", func(t *testing.T) {
		h := newPortalHandler(&mockStaffAuthUsecase{}, &mockStaffPortalUsecase{}, &mockStaffUsecase{})

		req := httptest.NewRequest(http.MethodGet, "/portal/auth/verify", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		testutil.SetStaffID(c, staffID)

		err := h.Verify(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "ok", resp["status"])
		assert.Equal(t, staffID.String(), resp["staffId"])
	})

	t.Run("異常系: staff_id なし → 401", func(t *testing.T) {
		h := newPortalHandler(&mockStaffAuthUsecase{}, &mockStaffPortalUsecase{}, &mockStaffUsecase{})

		req := httptest.NewRequest(http.MethodGet, "/portal/auth/verify", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		// staff_id をセットしない

		_ = h.Verify(c)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

// --- Profile Draft Tests ---

func TestStaffPortalHandler_GetMyProfileDraft(t *testing.T) {
	staffID := uuid.New()

	t.Run("正常系: プロフィール下書き取得 → 200", func(t *testing.T) {
		mock := &mockStaffPortalUsecase{
			profileDraft: domain.StaffProfileDraft{
				ID: uuid.New(), StaffID: staffID, Name: "Yuki", Status: domain.DraftStatusDraft,
			},
		}
		h := newPortalHandler(&mockStaffAuthUsecase{}, mock, &mockStaffUsecase{})

		req := httptest.NewRequest(http.MethodGet, "/portal/profile", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		testutil.SetStaffID(c, staffID)

		err := h.GetMyProfileDraft(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("異常系: staff_id なし → エラー", func(t *testing.T) {
		h := newPortalHandler(&mockStaffAuthUsecase{}, &mockStaffPortalUsecase{}, &mockStaffUsecase{})

		req := httptest.NewRequest(http.MethodGet, "/portal/profile", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := h.GetMyProfileDraft(c)
		assert.Error(t, err)
	})

	t.Run("異常系: usecase エラー → 500", func(t *testing.T) {
		mock := &mockStaffPortalUsecase{err: domain.ErrInternal}
		h := newPortalHandler(&mockStaffAuthUsecase{}, mock, &mockStaffUsecase{})

		req := httptest.NewRequest(http.MethodGet, "/portal/profile", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		testutil.SetStaffID(c, staffID)

		_ = h.GetMyProfileDraft(c)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestStaffPortalHandler_SaveMyProfileDraft(t *testing.T) {
	staffID := uuid.New()

	t.Run("正常系: プロフィール下書き保存 → 200", func(t *testing.T) {
		mock := &mockStaffPortalUsecase{
			profileDraft: domain.StaffProfileDraft{ID: uuid.New(), StaffID: staffID, Name: "Updated"},
		}
		h := newPortalHandler(&mockStaffAuthUsecase{}, mock, &mockStaffUsecase{})

		body := `{"name":"Updated","role":"Cast","bio":"bio","imageUrl":"","imageCropPosition":""}`
		req := httptest.NewRequest(http.MethodPut, "/portal/profile", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		testutil.SetStaffID(c, staffID)

		err := h.SaveMyProfileDraft(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("異常系: 不正なボディ → 400", func(t *testing.T) {
		h := newPortalHandler(&mockStaffAuthUsecase{}, &mockStaffPortalUsecase{}, &mockStaffUsecase{})

		req := httptest.NewRequest(http.MethodPut, "/portal/profile", strings.NewReader("bad"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		testutil.SetStaffID(c, staffID)

		_ = h.SaveMyProfileDraft(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestStaffPortalHandler_SubmitMyProfileDraft(t *testing.T) {
	staffID := uuid.New()
	draftID := uuid.New()

	t.Run("正常系: プロフィール提出 → 200", func(t *testing.T) {
		mock := &mockStaffPortalUsecase{
			profileDraft: domain.StaffProfileDraft{ID: draftID, Status: domain.DraftStatusPending},
		}
		h := newPortalHandler(&mockStaffAuthUsecase{}, mock, &mockStaffUsecase{})

		req := httptest.NewRequest(http.MethodPost, "/portal/profile/:id/submit", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(draftID.String())
		testutil.SetStaffID(c, staffID)

		err := h.SubmitMyProfileDraft(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("異常系: 不正なドラフトID → 400", func(t *testing.T) {
		h := newPortalHandler(&mockStaffAuthUsecase{}, &mockStaffPortalUsecase{}, &mockStaffUsecase{})

		req := httptest.NewRequest(http.MethodPost, "/portal/profile/:id/submit", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("bad")
		testutil.SetStaffID(c, staffID)

		_ = h.SubmitMyProfileDraft(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

// --- Schedule Draft Tests ---

func TestStaffPortalHandler_GetMyScheduleDraft(t *testing.T) {
	staffID := uuid.New()

	t.Run("正常系: スケジュール下書き取得 → 200", func(t *testing.T) {
		startTime, _ := time.Parse("15:04", "20:00")
		endTime, _ := time.Parse("15:04", "02:00")
		mock := &mockStaffPortalUsecase{
			scheduleDraft: domain.StaffScheduleDraft{
				ID: uuid.New(), StaffID: staffID, Status: domain.DraftStatusDraft,
				Items: []domain.ScheduleDraftItem{
					{ID: uuid.New(), DayOfWeek: 1, StartTime: startTime, EndTime: endTime},
				},
			},
		}
		h := newPortalHandler(&mockStaffAuthUsecase{}, mock, &mockStaffUsecase{})

		req := httptest.NewRequest(http.MethodGet, "/portal/schedule", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		testutil.SetStaffID(c, staffID)

		err := h.GetMyScheduleDraft(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp ScheduleDraftResponse
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp.Items, 1)
		assert.Equal(t, "20:00", resp.Items[0].StartTime)
	})
}

func TestStaffPortalHandler_SaveMyScheduleDraft(t *testing.T) {
	staffID := uuid.New()

	t.Run("正常系: スケジュール保存 → 200", func(t *testing.T) {
		mock := &mockStaffPortalUsecase{
			scheduleDraft: domain.StaffScheduleDraft{ID: uuid.New(), StaffID: staffID},
		}
		h := newPortalHandler(&mockStaffAuthUsecase{}, mock, &mockStaffUsecase{})

		body := `{"items":[{"dayOfWeek":1,"startTime":"20:00","endTime":"02:00"}]}`
		req := httptest.NewRequest(http.MethodPut, "/portal/schedule", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		testutil.SetStaffID(c, staffID)

		err := h.SaveMyScheduleDraft(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("異常系: 不正なボディ → 400", func(t *testing.T) {
		h := newPortalHandler(&mockStaffAuthUsecase{}, &mockStaffPortalUsecase{}, &mockStaffUsecase{})

		req := httptest.NewRequest(http.MethodPut, "/portal/schedule", strings.NewReader("bad"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		testutil.SetStaffID(c, staffID)

		_ = h.SaveMyScheduleDraft(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: 不正な時刻フォーマット → 400", func(t *testing.T) {
		h := newPortalHandler(&mockStaffAuthUsecase{}, &mockStaffPortalUsecase{}, &mockStaffUsecase{})

		body := `{"items":[{"dayOfWeek":1,"startTime":"invalid","endTime":"02:00"}]}`
		req := httptest.NewRequest(http.MethodPut, "/portal/schedule", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		testutil.SetStaffID(c, staffID)

		_ = h.SaveMyScheduleDraft(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestStaffPortalHandler_SubmitMyScheduleDraft(t *testing.T) {
	staffID := uuid.New()
	draftID := uuid.New()

	t.Run("正常系: スケジュール提出 → 200", func(t *testing.T) {
		mock := &mockStaffPortalUsecase{
			scheduleDraft: domain.StaffScheduleDraft{ID: draftID, Status: domain.DraftStatusPending},
		}
		h := newPortalHandler(&mockStaffAuthUsecase{}, mock, &mockStaffUsecase{})

		req := httptest.NewRequest(http.MethodPost, "/portal/schedule/:id/submit", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(draftID.String())
		testutil.SetStaffID(c, staffID)

		err := h.SubmitMyScheduleDraft(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("異常系: 不正なドラフトID → 400", func(t *testing.T) {
		h := newPortalHandler(&mockStaffAuthUsecase{}, &mockStaffPortalUsecase{}, &mockStaffUsecase{})

		req := httptest.NewRequest(http.MethodPost, "/portal/schedule/:id/submit", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("bad")
		testutil.SetStaffID(c, staffID)

		_ = h.SubmitMyScheduleDraft(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

// --- Image Management Tests ---

func TestStaffPortalHandler_ListMyImages(t *testing.T) {
	staffID := uuid.New()

	t.Run("正常系: 画像一覧を返す", func(t *testing.T) {
		staffMock := &mockStaffUsecase{
			staffDetail: domain.StaffWithSchedules{
				Staff:  domain.Staff{ID: staffID},
				Images: []domain.StaffImage{{ID: uuid.New(), StaffID: staffID, ImageURL: "/img.jpg"}},
			},
		}
		h := newPortalHandler(&mockStaffAuthUsecase{}, &mockStaffPortalUsecase{}, staffMock)

		req := httptest.NewRequest(http.MethodGet, "/portal/images", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		testutil.SetStaffID(c, staffID)

		err := h.ListMyImages(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("異常系: usecase エラー → 500", func(t *testing.T) {
		staffMock := &mockStaffUsecase{err: domain.ErrInternal}
		h := newPortalHandler(&mockStaffAuthUsecase{}, &mockStaffPortalUsecase{}, staffMock)

		req := httptest.NewRequest(http.MethodGet, "/portal/images", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		testutil.SetStaffID(c, staffID)

		_ = h.ListMyImages(c)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestStaffPortalHandler_DeleteMyImage(t *testing.T) {
	staffID := uuid.New()
	imageID := uuid.New()

	t.Run("正常系: 画像削除 → 204", func(t *testing.T) {
		staffMock := &mockStaffUsecase{
			staffDetail: domain.StaffWithSchedules{
				Staff:  domain.Staff{ID: staffID},
				Images: []domain.StaffImage{{ID: imageID, StaffID: staffID}},
			},
		}
		h := newPortalHandler(&mockStaffAuthUsecase{}, &mockStaffPortalUsecase{}, staffMock)

		req := httptest.NewRequest(http.MethodDelete, "/portal/images/:imageId", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("imageId")
		c.SetParamValues(imageID.String())
		testutil.SetStaffID(c, staffID)

		err := h.DeleteMyImage(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("異常系: imageId 空 → 400", func(t *testing.T) {
		h := newPortalHandler(&mockStaffAuthUsecase{}, &mockStaffPortalUsecase{}, &mockStaffUsecase{})

		req := httptest.NewRequest(http.MethodDelete, "/portal/images/:imageId", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("imageId")
		c.SetParamValues("")
		testutil.SetStaffID(c, staffID)

		_ = h.DeleteMyImage(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: 他人の画像 → 403", func(t *testing.T) {
		otherImageID := uuid.New()
		staffMock := &mockStaffUsecase{
			staffDetail: domain.StaffWithSchedules{
				Staff:  domain.Staff{ID: staffID},
				Images: []domain.StaffImage{{ID: imageID, StaffID: staffID}},
			},
		}
		h := newPortalHandler(&mockStaffAuthUsecase{}, &mockStaffPortalUsecase{}, staffMock)

		req := httptest.NewRequest(http.MethodDelete, "/portal/images/:imageId", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("imageId")
		c.SetParamValues(otherImageID.String())
		testutil.SetStaffID(c, staffID)

		_ = h.DeleteMyImage(c)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
}

func TestStaffPortalHandler_SetMyMainImage(t *testing.T) {
	staffID := uuid.New()
	imageID := uuid.New()

	t.Run("正常系: メイン画像設定 → 200", func(t *testing.T) {
		staffMock := &mockStaffUsecase{
			image: domain.StaffImage{ID: imageID, StaffID: staffID, IsMain: true},
		}
		h := newPortalHandler(&mockStaffAuthUsecase{}, &mockStaffPortalUsecase{}, staffMock)

		body := `{"imageId":"` + imageID.String() + `"}`
		req := httptest.NewRequest(http.MethodPut, "/portal/images/:imageId/main", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		testutil.SetStaffID(c, staffID)

		err := h.SetMyMainImage(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("異常系: 不正なボディ → 400", func(t *testing.T) {
		h := newPortalHandler(&mockStaffAuthUsecase{}, &mockStaffPortalUsecase{}, &mockStaffUsecase{})

		req := httptest.NewRequest(http.MethodPut, "/portal/images/:imageId/main", strings.NewReader("bad"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		testutil.SetStaffID(c, staffID)

		_ = h.SetMyMainImage(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestStaffPortalHandler_UpdateMyImageCropPosition(t *testing.T) {
	staffID := uuid.New()
	imageID := uuid.New()

	t.Run("正常系: クロップ位置更新 → 200", func(t *testing.T) {
		staffMock := &mockStaffUsecase{
			staffDetail: domain.StaffWithSchedules{
				Staff:  domain.Staff{ID: staffID},
				Images: []domain.StaffImage{{ID: imageID, StaffID: staffID}},
			},
			image: domain.StaffImage{ID: imageID, CropPosition: "30 70"},
		}
		h := newPortalHandler(&mockStaffAuthUsecase{}, &mockStaffPortalUsecase{}, staffMock)

		body := `{"cropPosition":"30 70"}`
		req := httptest.NewRequest(http.MethodPut, "/portal/images/:imageId/crop", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("imageId")
		c.SetParamValues(imageID.String())
		testutil.SetStaffID(c, staffID)

		err := h.UpdateMyImageCropPosition(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("異常系: imageId 空 → 400", func(t *testing.T) {
		h := newPortalHandler(&mockStaffAuthUsecase{}, &mockStaffPortalUsecase{}, &mockStaffUsecase{})

		body := `{"cropPosition":"50 50"}`
		req := httptest.NewRequest(http.MethodPut, "/portal/images/:imageId/crop", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("imageId")
		c.SetParamValues("")
		testutil.SetStaffID(c, staffID)

		_ = h.UpdateMyImageCropPosition(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: 他人の画像 → 403", func(t *testing.T) {
		otherImageID := uuid.New()
		staffMock := &mockStaffUsecase{
			staffDetail: domain.StaffWithSchedules{
				Staff:  domain.Staff{ID: staffID},
				Images: []domain.StaffImage{{ID: imageID, StaffID: staffID}},
			},
		}
		h := newPortalHandler(&mockStaffAuthUsecase{}, &mockStaffPortalUsecase{}, staffMock)

		body := `{"cropPosition":"50 50"}`
		req := httptest.NewRequest(http.MethodPut, "/portal/images/:imageId/crop", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("imageId")
		c.SetParamValues(otherImageID.String())
		testutil.SetStaffID(c, staffID)

		_ = h.UpdateMyImageCropPosition(c)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
}
