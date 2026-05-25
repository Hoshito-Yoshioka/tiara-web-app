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

// mockShopUsecase は ShopUsecase のテスト用モック。
type mockShopUsecase struct {
	shops []domain.Shop
	shop  domain.Shop
	err   error
}

func (m *mockShopUsecase) ListShops(_ context.Context) ([]domain.Shop, error) {
	return m.shops, m.err
}
func (m *mockShopUsecase) GetShopByID(_ context.Context, _ string) (domain.Shop, error) {
	return m.shop, m.err
}
func (m *mockShopUsecase) UpdateShop(_ context.Context, _ string, _ domain.UpdateShopInput) (domain.Shop, error) {
	return m.shop, m.err
}

// ============================================================
// ListShops
// ============================================================

func TestShopHandler_ListShops(t *testing.T) {
	t.Run("正常系: 店舗一覧を返す", func(t *testing.T) {
		mock := &mockShopUsecase{
			shops: []domain.Shop{{ID: uuid.New(), Name: "Tiara"}},
		}
		h := NewShopHandler(mock)

		req := httptest.NewRequest(http.MethodGet, "/shops", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := h.ListShops(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("異常系: usecase エラー → 500", func(t *testing.T) {
		mock := &mockShopUsecase{err: domain.ErrInternal}
		h := NewShopHandler(mock)

		req := httptest.NewRequest(http.MethodGet, "/shops", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		_ = h.ListShops(c)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

// ============================================================
// GetShopByID
// ============================================================

func TestShopHandler_GetShopByID(t *testing.T) {
	t.Run("正常系: 店舗詳細を返す", func(t *testing.T) {
		shopID := uuid.New()
		mock := &mockShopUsecase{
			shop: domain.Shop{ID: shopID, Name: "Tiara"},
		}
		h := NewShopHandler(mock)

		req := httptest.NewRequest(http.MethodGet, "/shops/:id", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(shopID.String())

		err := h.GetShopByID(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp ShopResponse
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Tiara", resp.Name)
	})

	t.Run("異常系: id パラメータが空 → 400", func(t *testing.T) {
		h := NewShopHandler(&mockShopUsecase{})

		req := httptest.NewRequest(http.MethodGet, "/shops/:id", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("")

		_ = h.GetShopByID(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: 存在しない店舗 → 404", func(t *testing.T) {
		mock := &mockShopUsecase{err: domain.ErrNotFound}
		h := NewShopHandler(mock)

		req := httptest.NewRequest(http.MethodGet, "/shops/:id", nil)
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(uuid.New().String())

		_ = h.GetShopByID(c)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

// ============================================================
// UpdateShop
// ============================================================

func TestShopHandler_UpdateShop(t *testing.T) {
	t.Run("正常系: 店舗更新 → 200", func(t *testing.T) {
		shopID := uuid.New()
		mock := &mockShopUsecase{
			shop: domain.Shop{ID: shopID, Name: "Updated Tiara"},
		}
		h := NewShopHandler(mock)

		body := `{"name":"Updated Tiara","address":"Tokyo","openingTime":"20:00","closingTime":"05:00"}`
		req := httptest.NewRequest(http.MethodPut, "/admin/shops/:id", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(shopID.String())

		err := h.UpdateShop(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("異常系: id パラメータが空 → 400", func(t *testing.T) {
		h := NewShopHandler(&mockShopUsecase{})

		body := `{"name":"X","address":"","openingTime":"","closingTime":""}`
		req := httptest.NewRequest(http.MethodPut, "/admin/shops/:id", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("")

		_ = h.UpdateShop(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: 不正なボディ → 400", func(t *testing.T) {
		h := NewShopHandler(&mockShopUsecase{})

		req := httptest.NewRequest(http.MethodPut, "/admin/shops/:id", strings.NewReader("bad"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(uuid.New().String())

		_ = h.UpdateShop(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: 存在しない店舗 → 404", func(t *testing.T) {
		mock := &mockShopUsecase{err: domain.ErrNotFound}
		h := NewShopHandler(mock)

		body := `{"name":"X","address":"","openingTime":"","closingTime":""}`
		req := httptest.NewRequest(http.MethodPut, "/admin/shops/:id", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(uuid.New().String())

		_ = h.UpdateShop(c)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
