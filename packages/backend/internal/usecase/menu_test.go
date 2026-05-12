package usecase

import (
	"context"
	"errors"
	"testing"

	"tiara-web-app/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// mockMenuRepository は MenuRepository のテスト用モック。
type mockMenuRepository struct {
	categories   []domain.MenuCategoryWithItems
	category     domain.MenuCategoryWithItems
	menuCategory domain.MenuCategory
	menuItem     domain.MenuItem
	err          error
}

func (m *mockMenuRepository) ListMenuCategoriesWithItems(_ context.Context) ([]domain.MenuCategoryWithItems, error) {
	return m.categories, m.err
}

func (m *mockMenuRepository) GetMenuCategoryByID(_ context.Context, _ string) (domain.MenuCategoryWithItems, error) {
	return m.category, m.err
}

func (m *mockMenuRepository) CreateMenuCategory(_ context.Context, _ domain.CreateMenuCategoryInput) (domain.MenuCategory, error) {
	return m.menuCategory, m.err
}

func (m *mockMenuRepository) UpdateMenuCategory(_ context.Context, _ string, _ domain.UpdateMenuCategoryInput) (domain.MenuCategory, error) {
	return m.menuCategory, m.err
}

func (m *mockMenuRepository) DeleteMenuCategory(_ context.Context, _ string) error {
	return m.err
}

func (m *mockMenuRepository) CreateMenuItem(_ context.Context, _ domain.CreateMenuItemInput) (domain.MenuItem, error) {
	return m.menuItem, m.err
}

func (m *mockMenuRepository) UpdateMenuItem(_ context.Context, _ string, _ domain.UpdateMenuItemInput) (domain.MenuItem, error) {
	return m.menuItem, m.err
}

func (m *mockMenuRepository) DeleteMenuItem(_ context.Context, _ string) error {
	return m.err
}

func TestMenuUsecase_ListMenuCategoriesWithItems_Success(t *testing.T) {
	categories := []domain.MenuCategoryWithItems{
		{
			Category: domain.MenuCategory{
				ID:   uuid.New(),
				Name: "Cocktails",
			},
			Items: []domain.MenuItem{
				{ID: uuid.New(), Name: "Mojito", Price: "¥1,200"},
			},
		},
	}

	repo := &mockMenuRepository{categories: categories}
	uc := NewMenuUsecase(repo)

	result, err := uc.ListMenuCategoriesWithItems(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Cocktails", result[0].Category.Name)
	assert.Len(t, result[0].Items, 1)
	assert.Equal(t, "Mojito", result[0].Items[0].Name)
}

func TestMenuUsecase_ListMenuCategoriesWithItems_Error(t *testing.T) {
	repo := &mockMenuRepository{err: errors.New("db error")}
	uc := NewMenuUsecase(repo)

	result, err := uc.ListMenuCategoriesWithItems(context.Background())

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestMenuUsecase_CreateMenuCategory_Success(t *testing.T) {
	cat := domain.MenuCategory{
		ID:   uuid.New(),
		Name: "Whisky",
	}

	repo := &mockMenuRepository{menuCategory: cat}
	uc := NewMenuUsecase(repo)

	result, err := uc.CreateMenuCategory(context.Background(), domain.CreateMenuCategoryInput{
		Name:      "Whisky",
		SortOrder: 1,
	})

	assert.NoError(t, err)
	assert.Equal(t, "Whisky", result.Name)
}

func TestMenuUsecase_DeleteMenuCategory_Success(t *testing.T) {
	repo := &mockMenuRepository{}
	uc := NewMenuUsecase(repo)

	err := uc.DeleteMenuCategory(context.Background(), uuid.New().String())

	assert.NoError(t, err)
}

func TestMenuUsecase_CreateMenuItem_Success(t *testing.T) {
	item := domain.MenuItem{
		ID:    uuid.New(),
		Name:  "Old Fashioned",
		Price: "¥1,500",
	}

	repo := &mockMenuRepository{menuItem: item}
	uc := NewMenuUsecase(repo)

	result, err := uc.CreateMenuItem(context.Background(), domain.CreateMenuItemInput{
		CategoryID: uuid.New().String(),
		Name:       "Old Fashioned",
		Price:      "¥1,500",
	})

	assert.NoError(t, err)
	assert.Equal(t, "Old Fashioned", result.Name)
}
