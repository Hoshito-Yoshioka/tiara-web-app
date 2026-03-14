package usecase

import (
	"context"
	"tiara-web-app/backend/internal/domain"
)

// MenuRepository はメニューデータの永続化を抽象化するインターフェース。
type MenuRepository interface {
	ListMenuCategoriesWithItems(ctx context.Context) ([]domain.MenuCategoryWithItems, error)
	GetMenuCategoryByID(ctx context.Context, id string) (domain.MenuCategoryWithItems, error)
	CreateMenuCategory(ctx context.Context, input domain.CreateMenuCategoryInput) (domain.MenuCategory, error)
	UpdateMenuCategory(ctx context.Context, id string, input domain.UpdateMenuCategoryInput) (domain.MenuCategory, error)
	DeleteMenuCategory(ctx context.Context, id string) error
	CreateMenuItem(ctx context.Context, input domain.CreateMenuItemInput) (domain.MenuItem, error)
	UpdateMenuItem(ctx context.Context, id string, input domain.UpdateMenuItemInput) (domain.MenuItem, error)
	DeleteMenuItem(ctx context.Context, id string) error
}

// MenuUsecase はメニューに関するビジネスロジックを定義するインターフェース。
type MenuUsecase interface {
	ListMenuCategoriesWithItems(ctx context.Context) ([]domain.MenuCategoryWithItems, error)
	GetMenuCategoryByID(ctx context.Context, id string) (domain.MenuCategoryWithItems, error)
	CreateMenuCategory(ctx context.Context, input domain.CreateMenuCategoryInput) (domain.MenuCategory, error)
	UpdateMenuCategory(ctx context.Context, id string, input domain.UpdateMenuCategoryInput) (domain.MenuCategory, error)
	DeleteMenuCategory(ctx context.Context, id string) error
	CreateMenuItem(ctx context.Context, input domain.CreateMenuItemInput) (domain.MenuItem, error)
	UpdateMenuItem(ctx context.Context, id string, input domain.UpdateMenuItemInput) (domain.MenuItem, error)
	DeleteMenuItem(ctx context.Context, id string) error
}

type menuUsecase struct {
	menuRepo MenuRepository
}

func NewMenuUsecase(repo MenuRepository) MenuUsecase {
	return &menuUsecase{menuRepo: repo}
}

func (u *menuUsecase) ListMenuCategoriesWithItems(ctx context.Context) ([]domain.MenuCategoryWithItems, error) {
	return u.menuRepo.ListMenuCategoriesWithItems(ctx)
}

func (u *menuUsecase) GetMenuCategoryByID(ctx context.Context, id string) (domain.MenuCategoryWithItems, error) {
	return u.menuRepo.GetMenuCategoryByID(ctx, id)
}

func (u *menuUsecase) CreateMenuCategory(ctx context.Context, input domain.CreateMenuCategoryInput) (domain.MenuCategory, error) {
	return u.menuRepo.CreateMenuCategory(ctx, input)
}

func (u *menuUsecase) UpdateMenuCategory(ctx context.Context, id string, input domain.UpdateMenuCategoryInput) (domain.MenuCategory, error) {
	return u.menuRepo.UpdateMenuCategory(ctx, id, input)
}

func (u *menuUsecase) DeleteMenuCategory(ctx context.Context, id string) error {
	return u.menuRepo.DeleteMenuCategory(ctx, id)
}

func (u *menuUsecase) CreateMenuItem(ctx context.Context, input domain.CreateMenuItemInput) (domain.MenuItem, error) {
	return u.menuRepo.CreateMenuItem(ctx, input)
}

func (u *menuUsecase) UpdateMenuItem(ctx context.Context, id string, input domain.UpdateMenuItemInput) (domain.MenuItem, error) {
	return u.menuRepo.UpdateMenuItem(ctx, id, input)
}

func (u *menuUsecase) DeleteMenuItem(ctx context.Context, id string) error {
	return u.menuRepo.DeleteMenuItem(ctx, id)
}
