package db

import (
	"context"
	"tiara-web-app/backend/internal/domain"
	"tiara-web-app/backend/internal/usecase"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type menuRepository struct {
	q *Queries
}

func NewMenuRepository(q *Queries) usecase.MenuRepository {
	return &menuRepository{q: q}
}

// pgtype.UUID → uuid.UUID 変換ヘルパー
func toUUID(p pgtype.UUID) uuid.UUID {
	return uuid.UUID(p.Bytes)
}

// uuid.UUID → pgtype.UUID 変換ヘルパー
func toPgtypeUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

func toMenuCategory(r MenuCategory) domain.MenuCategory {
	return domain.MenuCategory{
		ID:          toUUID(r.ID),
		Name:        r.Name,
		Description: r.Description,
		SortOrder:   r.SortOrder,
		CreatedAt:   r.CreatedAt.Time,
		UpdatedAt:   r.UpdatedAt.Time,
	}
}

func toMenuItem(r MenuItem) domain.MenuItem {
	return domain.MenuItem{
		ID:          toUUID(r.ID),
		CategoryID:  toUUID(r.CategoryID),
		Name:        r.Name,
		Price:       r.Price,
		Description: r.Description,
		SortOrder:   r.SortOrder,
		CreatedAt:   r.CreatedAt.Time,
		UpdatedAt:   r.UpdatedAt.Time,
	}
}

// ListMenuCategoriesWithItems は全カテゴリとそれに属するアイテムを取得する。
func (r *menuRepository) ListMenuCategoriesWithItems(ctx context.Context) ([]domain.MenuCategoryWithItems, error) {
	cats, err := r.q.ListMenuCategories(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]domain.MenuCategoryWithItems, len(cats))
	for i, cat := range cats {
		items, err := r.q.ListMenuItemsByCategoryID(ctx, cat.ID)
		if err != nil {
			return nil, err
		}
		domainItems := make([]domain.MenuItem, len(items))
		for j, item := range items {
			domainItems[j] = toMenuItem(item)
		}
		result[i] = domain.MenuCategoryWithItems{
			Category: toMenuCategory(cat),
			Items:    domainItems,
		}
	}
	return result, nil
}

func (r *menuRepository) GetMenuCategoryByID(ctx context.Context, id string) (domain.MenuCategoryWithItems, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return domain.MenuCategoryWithItems{}, err
	}
	cat, err := r.q.GetMenuCategoryByID(ctx, toPgtypeUUID(uid))
	if err != nil {
		return domain.MenuCategoryWithItems{}, err
	}
	items, err := r.q.ListMenuItemsByCategoryID(ctx, cat.ID)
	if err != nil {
		return domain.MenuCategoryWithItems{}, err
	}
	domainItems := make([]domain.MenuItem, len(items))
	for i, item := range items {
		domainItems[i] = toMenuItem(item)
	}
	return domain.MenuCategoryWithItems{Category: toMenuCategory(cat), Items: domainItems}, nil
}

func (r *menuRepository) CreateMenuCategory(ctx context.Context, input domain.CreateMenuCategoryInput) (domain.MenuCategory, error) {
	cat, err := r.q.CreateMenuCategory(ctx, CreateMenuCategoryParams{
		Name:        input.Name,
		Description: input.Description,
		SortOrder:   input.SortOrder,
	})
	if err != nil {
		return domain.MenuCategory{}, err
	}
	return toMenuCategory(cat), nil
}

func (r *menuRepository) UpdateMenuCategory(ctx context.Context, id string, input domain.UpdateMenuCategoryInput) (domain.MenuCategory, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return domain.MenuCategory{}, err
	}
	cat, err := r.q.UpdateMenuCategory(ctx, UpdateMenuCategoryParams{
		ID:          toPgtypeUUID(uid),
		Name:        input.Name,
		Description: input.Description,
		SortOrder:   input.SortOrder,
	})
	if err != nil {
		return domain.MenuCategory{}, err
	}
	return toMenuCategory(cat), nil
}

func (r *menuRepository) DeleteMenuCategory(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.q.DeleteMenuCategory(ctx, toPgtypeUUID(uid))
}

func (r *menuRepository) CreateMenuItem(ctx context.Context, input domain.CreateMenuItemInput) (domain.MenuItem, error) {
	catUID, err := uuid.Parse(input.CategoryID)
	if err != nil {
		return domain.MenuItem{}, err
	}
	item, err := r.q.CreateMenuItem(ctx, CreateMenuItemParams{
		CategoryID:  toPgtypeUUID(catUID),
		Name:        input.Name,
		Price:       input.Price,
		Description: input.Description,
		SortOrder:   input.SortOrder,
	})
	if err != nil {
		return domain.MenuItem{}, err
	}
	return toMenuItem(item), nil
}

func (r *menuRepository) UpdateMenuItem(ctx context.Context, id string, input domain.UpdateMenuItemInput) (domain.MenuItem, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return domain.MenuItem{}, err
	}
	item, err := r.q.UpdateMenuItem(ctx, UpdateMenuItemParams{
		ID:          toPgtypeUUID(uid),
		Name:        input.Name,
		Price:       input.Price,
		Description: input.Description,
		SortOrder:   input.SortOrder,
	})
	if err != nil {
		return domain.MenuItem{}, err
	}
	return toMenuItem(item), nil
}

func (r *menuRepository) DeleteMenuItem(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.q.DeleteMenuItem(ctx, toPgtypeUUID(uid))
}
