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

// mockShopRepository は ShopRepository のテスト用モック。
type mockShopRepository struct {
	shops []domain.Shop
	shop  domain.Shop
	err   error
}

func (m *mockShopRepository) ListShops(_ context.Context) ([]domain.Shop, error) {
	return m.shops, m.err
}

func (m *mockShopRepository) GetShopByID(_ context.Context, _ string) (domain.Shop, error) {
	return m.shop, m.err
}

func (m *mockShopRepository) UpdateShop(_ context.Context, _ string, _ domain.UpdateShopInput) (domain.Shop, error) {
	return m.shop, m.err
}

func TestShopUsecase_ListShops_Success(t *testing.T) {
	shops := []domain.Shop{
		{
			ID:   uuid.New(),
			Name: "Tiara",
		},
	}

	repo := &mockShopRepository{shops: shops}
	uc := NewShopUsecase(repo)

	result, err := uc.ListShops(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Tiara", result[0].Name)
}

func TestShopUsecase_ListShops_Error(t *testing.T) {
	repo := &mockShopRepository{err: errors.New("db error")}
	uc := NewShopUsecase(repo)

	result, err := uc.ListShops(context.Background())

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestShopUsecase_GetShopByID_Success(t *testing.T) {
	shop := domain.Shop{
		ID:      uuid.New(),
		Name:    "Tiara",
		Address: "Tokyo",
	}

	repo := &mockShopRepository{shop: shop}
	uc := NewShopUsecase(repo)

	result, err := uc.GetShopByID(context.Background(), shop.ID.String())

	assert.NoError(t, err)
	assert.Equal(t, shop.Name, result.Name)
	assert.Equal(t, shop.Address, result.Address)
}

func TestShopUsecase_GetShopByID_NotFound(t *testing.T) {
	repo := &mockShopRepository{err: errors.New("not found")}
	uc := NewShopUsecase(repo)

	_, err := uc.GetShopByID(context.Background(), uuid.New().String())

	assert.Error(t, err)
}

func TestShopUsecase_UpdateShop_Success(t *testing.T) {
	updated := domain.Shop{
		ID:          uuid.New(),
		Name:        "Tiara Updated",
		Address:     "Osaka",
		OpeningTime: time.Now(),
		ClosingTime: time.Now(),
	}

	repo := &mockShopRepository{shop: updated}
	uc := NewShopUsecase(repo)

	input := domain.UpdateShopInput{
		Name:    "Tiara Updated",
		Address: "Osaka",
	}

	result, err := uc.UpdateShop(context.Background(), updated.ID.String(), input)

	assert.NoError(t, err)
	assert.Equal(t, "Tiara Updated", result.Name)
}
