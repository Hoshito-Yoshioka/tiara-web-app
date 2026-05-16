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

// mockMenuUsecase は MenuUsecase のテスト用モック。
type mockMenuUsecase struct {
	categories     []domain.MenuCategoryWithItems
	category       domain.MenuCategory
	categoryDetail domain.MenuCategoryWithItems
	item           domain.MenuItem
	err            error
}

func (m *mockMenuUsecase) ListMenuCategoriesWithItems(_ context.Context) ([]domain.MenuCategoryWithItems, error) {
	return m.categories, m.err
}
func (m *mockMenuUsecase) GetMenuCategoryByID(_ context.Context, _ string) (domain.MenuCategoryWithItems, error) {
	return m.categoryDetail, m.err
}
func (m *mockMenuUsecase) CreateMenuCategory(_ context.Context, _ domain.CreateMenuCategoryInput) (domain.MenuCategory, error) {
	return m.category, m.err
}
func (m *mockMenuUsecase) UpdateMenuCategory(_ context.Context, _ string, _ domain.UpdateMenuCategoryInput) (domain.MenuCategory, error) {
	return m.category, m.err
}
func (m *mockMenuUsecase) DeleteMenuCategory(_ context.Context, _ string) error {
	return m.err
}
func (m *mockMenuUsecase) CreateMenuItem(_ context.Context, _ domain.CreateMenuItemInput) (domain.MenuItem, error) {
	return m.item, m.err
}
func (m *mockMenuUsecase) UpdateMenuItem(_ context.Context, _ string, _ domain.UpdateMenuItemInput) (domain.MenuItem, error) {
	return m.item, m.err
}
func (m *mockMenuUsecase) DeleteMenuItem(_ context.Context, _ string) error {
	return m.err
}

// ============================================================
// ListMenuCategoriesWithItems
// ============================================================

func TestMenuHandler_ListMenuCategoriesWithItems(t *testing.T) {
	catID := uuid.New()

	t.Run("正常系: カテゴリ＋アイテム一覧を返す", func(t *testing.T) {
		mock := &mockMenuUsecase{
			categories: []domain.MenuCategoryWithItems{
				{
					Category: domain.MenuCategory{ID: catID, Name: "Cocktails"},
					Items:    []domain.MenuItem{{ID: uuid.New(), Name: "Mojito", Price: "800"}},
				},
			},
		}
		h := NewMenuHandler(mock)

		req := httptest.NewRequest(http.MethodGet, "/menus", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := h.ListMenuCategoriesWithItems(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("異常系: usecase エラー → 500", func(t *testing.T) {
		mock := &mockMenuUsecase{err: domain.ErrInternal}
		h := NewMenuHandler(mock)

		req := httptest.NewRequest(http.MethodGet, "/menus", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		_ = h.ListMenuCategoriesWithItems(c)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

// ============================================================
// CreateMenuCategory
// ============================================================

func TestMenuHandler_CreateMenuCategory(t *testing.T) {
	t.Run("正常系: カテゴリ作成 → 201", func(t *testing.T) {
		mock := &mockMenuUsecase{
			category: domain.MenuCategory{ID: uuid.New(), Name: "Whisky"},
		}
		h := NewMenuHandler(mock)

		body := `{"name":"Whisky","description":"Single malt","sortOrder":1}`
		req := httptest.NewRequest(http.MethodPost, "/admin/menu/categories", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := h.CreateMenuCategory(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp domain.MenuCategory
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Whisky", resp.Name)
	})

	t.Run("異常系: 不正なボディ → 400", func(t *testing.T) {
		h := NewMenuHandler(&mockMenuUsecase{})

		req := httptest.NewRequest(http.MethodPost, "/admin/menu/categories", strings.NewReader("invalid"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		_ = h.CreateMenuCategory(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: usecase エラー → handleError", func(t *testing.T) {
		mock := &mockMenuUsecase{err: domain.ErrConflict}
		h := NewMenuHandler(mock)

		body := `{"name":"Whisky","description":"","sortOrder":1}`
		req := httptest.NewRequest(http.MethodPost, "/admin/menu/categories", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		_ = h.CreateMenuCategory(c)
		assert.Equal(t, http.StatusConflict, rec.Code)
	})
}

// ============================================================
// UpdateMenuCategory
// ============================================================

func TestMenuHandler_UpdateMenuCategory(t *testing.T) {
	t.Run("正常系: カテゴリ更新 → 200", func(t *testing.T) {
		mock := &mockMenuUsecase{
			category: domain.MenuCategory{ID: uuid.New(), Name: "Updated"},
		}
		h := NewMenuHandler(mock)

		body := `{"name":"Updated","description":"desc","sortOrder":2}`
		req := httptest.NewRequest(http.MethodPut, "/admin/menu/categories/:id", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(uuid.New().String())

		err := h.UpdateMenuCategory(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("異常系: 存在しない ID → 404", func(t *testing.T) {
		mock := &mockMenuUsecase{err: domain.ErrNotFound}
		h := NewMenuHandler(mock)

		body := `{"name":"X","description":"","sortOrder":1}`
		req := httptest.NewRequest(http.MethodPut, "/admin/menu/categories/:id", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(uuid.New().String())

		_ = h.UpdateMenuCategory(c)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

// ============================================================
// DeleteMenuCategory
// ============================================================

func TestMenuHandler_DeleteMenuCategory(t *testing.T) {
	t.Run("正常系: カテゴリ削除 → 204", func(t *testing.T) {
		mock := &mockMenuUsecase{}
		h := NewMenuHandler(mock)

		req := httptest.NewRequest(http.MethodDelete, "/admin/menu/categories/:id", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(uuid.New().String())

		err := h.DeleteMenuCategory(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("異常系: 存在しない ID → 404", func(t *testing.T) {
		mock := &mockMenuUsecase{err: domain.ErrNotFound}
		h := NewMenuHandler(mock)

		req := httptest.NewRequest(http.MethodDelete, "/admin/menu/categories/:id", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(uuid.New().String())

		_ = h.DeleteMenuCategory(c)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

// ============================================================
// CreateMenuItem
// ============================================================

func TestMenuHandler_CreateMenuItem(t *testing.T) {
	t.Run("正常系: アイテム作成 → 201", func(t *testing.T) {
		mock := &mockMenuUsecase{
			item: domain.MenuItem{ID: uuid.New(), Name: "Mojito", Price: "800"},
		}
		h := NewMenuHandler(mock)

		body := `{"categoryId":"` + uuid.New().String() + `","name":"Mojito","price":"800","description":"","sortOrder":1}`
		req := httptest.NewRequest(http.MethodPost, "/admin/menu/items", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := h.CreateMenuItem(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("異常系: 不正なボディ → 400", func(t *testing.T) {
		h := NewMenuHandler(&mockMenuUsecase{})

		req := httptest.NewRequest(http.MethodPost, "/admin/menu/items", strings.NewReader("bad"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		_ = h.CreateMenuItem(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

// ============================================================
// UpdateMenuItem
// ============================================================

func TestMenuHandler_UpdateMenuItem(t *testing.T) {
	t.Run("正常系: アイテム更新 → 200", func(t *testing.T) {
		mock := &mockMenuUsecase{
			item: domain.MenuItem{ID: uuid.New(), Name: "Updated"},
		}
		h := NewMenuHandler(mock)

		body := `{"name":"Updated","price":"900","description":"desc","sortOrder":1}`
		req := httptest.NewRequest(http.MethodPut, "/admin/menu/items/:id", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(uuid.New().String())

		err := h.UpdateMenuItem(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

// ============================================================
// DeleteMenuItem
// ============================================================

func TestMenuHandler_DeleteMenuItem(t *testing.T) {
	t.Run("正常系: アイテム削除 → 204", func(t *testing.T) {
		mock := &mockMenuUsecase{}
		h := NewMenuHandler(mock)

		req := httptest.NewRequest(http.MethodDelete, "/admin/menu/items/:id", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(uuid.New().String())

		err := h.DeleteMenuItem(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})
}
