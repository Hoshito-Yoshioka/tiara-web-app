package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"tiara-web-app/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// mockStaffUsecase は StaffUsecase のテスト用モック。
type mockStaffUsecase struct {
	staffs       []domain.Staff
	staffDetail  domain.StaffWithSchedules
	staffDetails []domain.StaffWithSchedules
	image        domain.StaffImage
	err          error
}

func (m *mockStaffUsecase) ListStaffs(_ context.Context) ([]domain.Staff, error) {
	return m.staffs, m.err
}
func (m *mockStaffUsecase) GetStaffWithSchedules(_ context.Context, _ string) (domain.StaffWithSchedules, error) {
	return m.staffDetail, m.err
}
func (m *mockStaffUsecase) ListAllStaffsWithSchedules(_ context.Context) ([]domain.StaffWithSchedules, error) {
	return m.staffDetails, m.err
}
func (m *mockStaffUsecase) CreateStaff(_ context.Context, _ domain.CreateStaffInput) (domain.StaffWithSchedules, error) {
	return m.staffDetail, m.err
}
func (m *mockStaffUsecase) UpdateStaff(_ context.Context, _ string, _ domain.UpdateStaffInput) (domain.StaffWithSchedules, error) {
	return m.staffDetail, m.err
}
func (m *mockStaffUsecase) DeleteStaff(_ context.Context, _ string) error {
	return m.err
}
func (m *mockStaffUsecase) UploadStaffImage(_ context.Context, _ string, _ string, _ bool, _ int) (domain.StaffImage, error) {
	return m.image, m.err
}
func (m *mockStaffUsecase) DeleteStaffImage(_ context.Context, _ string) error {
	return m.err
}
func (m *mockStaffUsecase) SetMainImage(_ context.Context, _ string, _ string) (domain.StaffImage, error) {
	return m.image, m.err
}
func (m *mockStaffUsecase) UpdateImageCropPosition(_ context.Context, _ string, _ string) (domain.StaffImage, error) {
	return m.image, m.err
}

// ============================================================
// ListStaffs
// ============================================================

func TestStaffHandler_ListStaffs(t *testing.T) {
	t.Run("正常系: スタッフ一覧を返す", func(t *testing.T) {
		mock := &mockStaffUsecase{
			staffs: []domain.Staff{{ID: uuid.New(), Name: "Yuki"}},
		}
		h := NewStaffHandler(mock)

		req := httptest.NewRequest(http.MethodGet, "/staffs", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := h.ListStaffs(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("異常系: usecase エラー → 500", func(t *testing.T) {
		mock := &mockStaffUsecase{err: domain.ErrInternal}
		h := NewStaffHandler(mock)

		req := httptest.NewRequest(http.MethodGet, "/staffs", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		_ = h.ListStaffs(c)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

// ============================================================
// GetStaffWithSchedules
// ============================================================

func TestStaffHandler_GetStaffWithSchedules(t *testing.T) {
	staffID := uuid.New()

	t.Run("正常系: スタッフ詳細を返す", func(t *testing.T) {
		mock := &mockStaffUsecase{
			staffDetail: domain.StaffWithSchedules{
				Staff: domain.Staff{ID: staffID, Name: "Yuki"},
			},
		}
		h := NewStaffHandler(mock)

		req := httptest.NewRequest(http.MethodGet, "/staffs/:id", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(staffID.String())

		err := h.GetStaffWithSchedules(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("異常系: id パラメータが空 → 400", func(t *testing.T) {
		h := NewStaffHandler(&mockStaffUsecase{})

		req := httptest.NewRequest(http.MethodGet, "/staffs/:id", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("")

		_ = h.GetStaffWithSchedules(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: 存在しない → 404", func(t *testing.T) {
		mock := &mockStaffUsecase{err: domain.ErrNotFound}
		h := NewStaffHandler(mock)

		req := httptest.NewRequest(http.MethodGet, "/staffs/:id", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(uuid.New().String())

		_ = h.GetStaffWithSchedules(c)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

// ============================================================
// CreateStaff
// ============================================================

func TestStaffHandler_CreateStaff(t *testing.T) {
	t.Run("正常系: スタッフ作成 → 201", func(t *testing.T) {
		mock := &mockStaffUsecase{
			staffDetail: domain.StaffWithSchedules{
				Staff: domain.Staff{ID: uuid.New(), Name: "Yuki"},
			},
		}
		h := NewStaffHandler(mock)

		body := `{"shopId":"` + uuid.New().String() + `","name":"Yuki","role":"Cast","bio":"","imageUrl":"","imageCropPosition":"center","sortOrder":1,"schedules":[]}`
		req := httptest.NewRequest(http.MethodPost, "/admin/staffs", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := h.CreateStaff(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("異常系: 不正なボディ → 400", func(t *testing.T) {
		h := NewStaffHandler(&mockStaffUsecase{})

		req := httptest.NewRequest(http.MethodPost, "/admin/staffs", strings.NewReader("bad"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		_ = h.CreateStaff(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

// ============================================================
// UpdateStaff
// ============================================================

func TestStaffHandler_UpdateStaff(t *testing.T) {
	staffID := uuid.New()

	t.Run("正常系: スタッフ更新 → 200", func(t *testing.T) {
		mock := &mockStaffUsecase{
			staffDetail: domain.StaffWithSchedules{
				Staff: domain.Staff{ID: staffID, Name: "Updated"},
			},
		}
		h := NewStaffHandler(mock)

		body := `{"name":"Updated","role":"Cast","bio":"","imageUrl":"","imageCropPosition":"center","sortOrder":1,"schedules":[]}`
		req := httptest.NewRequest(http.MethodPut, "/admin/staffs/:id", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(staffID.String())

		err := h.UpdateStaff(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("異常系: id が空 → 400", func(t *testing.T) {
		h := NewStaffHandler(&mockStaffUsecase{})

		body := `{"name":"X","role":"","bio":"","imageUrl":"","imageCropPosition":"","sortOrder":0,"schedules":[]}`
		req := httptest.NewRequest(http.MethodPut, "/admin/staffs/:id", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("")

		_ = h.UpdateStaff(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

// ============================================================
// DeleteStaff
// ============================================================

func TestStaffHandler_DeleteStaff(t *testing.T) {
	t.Run("正常系: スタッフ削除 → 204", func(t *testing.T) {
		mock := &mockStaffUsecase{}
		h := NewStaffHandler(mock)

		req := httptest.NewRequest(http.MethodDelete, "/admin/staffs/:id", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(uuid.New().String())

		err := h.DeleteStaff(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("異常系: id が空 → 400", func(t *testing.T) {
		h := NewStaffHandler(&mockStaffUsecase{})

		req := httptest.NewRequest(http.MethodDelete, "/admin/staffs/:id", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("")

		_ = h.DeleteStaff(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: 存在しない → 404", func(t *testing.T) {
		mock := &mockStaffUsecase{err: domain.ErrNotFound}
		h := NewStaffHandler(mock)

		req := httptest.NewRequest(http.MethodDelete, "/admin/staffs/:id", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(uuid.New().String())

		_ = h.DeleteStaff(c)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

// ============================================================
// SetMainImage
// ============================================================

func TestStaffHandler_SetMainImage(t *testing.T) {
	staffID := uuid.New()
	imageID := uuid.New()

	t.Run("正常系: メイン画像設定 → 200", func(t *testing.T) {
		mock := &mockStaffUsecase{
			image: domain.StaffImage{ID: imageID, StaffID: staffID, IsMain: true},
		}
		h := NewStaffHandler(mock)

		body := `{"imageId":"` + imageID.String() + `"}`
		req := httptest.NewRequest(http.MethodPut, "/admin/staffs/:id/main-image", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(staffID.String())

		err := h.SetMainImage(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, true, resp["IsMain"])
	})

	t.Run("異常系: staffID が空 → 400", func(t *testing.T) {
		h := NewStaffHandler(&mockStaffUsecase{})

		body := `{"imageId":"` + imageID.String() + `"}`
		req := httptest.NewRequest(http.MethodPut, "/admin/staffs/:id/main-image", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("")

		_ = h.SetMainImage(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

// ============================================================
// UpdateImageCropPosition
// ============================================================

func TestStaffHandler_UpdateImageCropPosition(t *testing.T) {
	imageID := uuid.New()

	t.Run("正常系: クロップ位置更新 → 200", func(t *testing.T) {
		mock := &mockStaffUsecase{
			image: domain.StaffImage{ID: imageID, CropPosition: "top"},
		}
		h := NewStaffHandler(mock)

		body := `{"cropPosition":"top"}`
		req := httptest.NewRequest(http.MethodPut, "/admin/staffs/images/:imageId/crop", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("imageId")
		c.SetParamValues(imageID.String())

		err := h.UpdateImageCropPosition(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("異常系: imageId が空 → 400", func(t *testing.T) {
		h := NewStaffHandler(&mockStaffUsecase{})

		body := `{"cropPosition":"center"}`
		req := httptest.NewRequest(http.MethodPut, "/admin/staffs/images/:imageId/crop", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("imageId")
		c.SetParamValues("")

		_ = h.UpdateImageCropPosition(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
