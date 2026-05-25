package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"tiara-web-app/backend/internal/domain"
	"tiara-web-app/backend/internal/testutil"

	"github.com/google/uuid"
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
		shop := testutil.NewShop()
		mock := &mockShopUsecase{
			shops: []domain.Shop{shop},
		}
		h := NewShopHandler(mock)

		c, rec := testutil.NewEchoContext(http.MethodGet, "/shops")

		err := h.ListShops(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("異常系: usecase エラー → 500", func(t *testing.T) {
		mock := &mockShopUsecase{err: domain.ErrInternal}
		h := NewShopHandler(mock)

		c, rec := testutil.NewEchoContext(http.MethodGet, "/shops")

		_ = h.ListShops(c)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

// ============================================================
// GetShopByID
// ============================================================

func TestShopHandler_GetShopByID(t *testing.T) {
	t.Run("正常系: 店舗詳細を返す", func(t *testing.T) {
		shop := testutil.NewShop()
		mock := &mockShopUsecase{shop: shop}
		h := NewShopHandler(mock)

		c, rec := testutil.NewEchoContext(http.MethodGet, "/shops/:id")
		testutil.SetPathParams(c, "id", shop.ID.String())

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

		c, rec := testutil.NewEchoContext(http.MethodGet, "/shops/:id")
		testutil.SetPathParams(c, "id", "")

		_ = h.GetShopByID(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: 存在しない店舗 → 404", func(t *testing.T) {
		mock := &mockShopUsecase{err: domain.ErrNotFound}
		h := NewShopHandler(mock)

		c, rec := testutil.NewEchoContext(http.MethodGet, "/shops/:id")
		testutil.SetPathParams(c, "id", uuid.New().String())

		_ = h.GetShopByID(c)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

// ============================================================
// UpdateShop
// ============================================================

func TestShopHandler_UpdateShop(t *testing.T) {
	t.Run("正常系: 店舗更新 → 200", func(t *testing.T) {
		shop := testutil.NewShop()
		shop.Name = "Updated Tiara"
		mock := &mockShopUsecase{shop: shop}
		h := NewShopHandler(mock)

		body := `{"name":"Updated Tiara","address":"Tokyo","openingTime":"20:00","closingTime":"05:00"}`
		c, rec := testutil.NewEchoContext(http.MethodPut, "/admin/shops/:id", body)
		testutil.SetPathParams(c, "id", shop.ID.String())

		err := h.UpdateShop(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("異常系: id パラメータが空 → 400", func(t *testing.T) {
		h := NewShopHandler(&mockShopUsecase{})

		body := `{"name":"X","address":"","openingTime":"","closingTime":""}`
		c, rec := testutil.NewEchoContext(http.MethodPut, "/admin/shops/:id", body)
		testutil.SetPathParams(c, "id", "")

		_ = h.UpdateShop(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: 不正なボディ → 400", func(t *testing.T) {
		h := NewShopHandler(&mockShopUsecase{})

		c, rec := testutil.NewEchoContext(http.MethodPut, "/admin/shops/:id", "bad")
		testutil.SetPathParams(c, "id", uuid.New().String())

		_ = h.UpdateShop(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: 存在しない店舗 → 404", func(t *testing.T) {
		mock := &mockShopUsecase{err: domain.ErrNotFound}
		h := NewShopHandler(mock)

		body := `{"name":"X","address":"","openingTime":"","closingTime":""}`
		c, rec := testutil.NewEchoContext(http.MethodPut, "/admin/shops/:id", body)
		testutil.SetPathParams(c, "id", uuid.New().String())

		_ = h.UpdateShop(c)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
