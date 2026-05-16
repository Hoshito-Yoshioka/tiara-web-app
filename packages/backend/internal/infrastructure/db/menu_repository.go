package db

import (
	"context"
	"tiara-web-app/backend/internal/domain"
	"tiara-web-app/backend/internal/usecase"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type menuRepository struct {
	q    *Queries
	pool *pgxpool.Pool
}

func NewMenuRepository(q *Queries, pool *pgxpool.Pool) usecase.MenuRepository {
	return &menuRepository{q: q, pool: pool}
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
	pgtypeUID := toPgtypeUUID(uid)

	// トランザクション開始（sort_order スワップ対応）
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.MenuCategory{}, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck // rollback after commit is a no-op

	qtx := r.q.WithTx(tx)

	// 現在のカテゴリを取得
	current, err := qtx.GetMenuCategoryByID(ctx, pgtypeUID)
	if err != nil {
		return domain.MenuCategory{}, err
	}

	// sort_order が変更された場合、他のカテゴリと入れ替え
	if current.SortOrder != input.SortOrder {
		_ = qtx.SwapMenuCategorySortOrder(ctx, SwapMenuCategorySortOrderParams{
			SortOrder:   input.SortOrder,
			SortOrder_2: current.SortOrder,
			ID:          pgtypeUID,
		})
	}

	cat, err := qtx.UpdateMenuCategory(ctx, UpdateMenuCategoryParams{
		ID:          pgtypeUID,
		Name:        input.Name,
		Description: input.Description,
		SortOrder:   input.SortOrder,
	})
	if err != nil {
		return domain.MenuCategory{}, err
	}

	err = tx.Commit(ctx)
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
	pgtypeUID := toPgtypeUUID(uid)

	// トランザクション開始（sort_order スワップ対応）
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.MenuItem{}, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck // rollback after commit is a no-op

	qtx := r.q.WithTx(tx)

	// 現在のアイテムの sort_order と category_id を取得するための SQL（sqlc クエリではなく直接取得）
	var currentSortOrder int32
	var currentCategoryID pgtype.UUID
	err = tx.QueryRow(ctx, "SELECT sort_order, category_id FROM menu_items WHERE id = $1", pgtypeUID).Scan(&currentSortOrder, &currentCategoryID)
	if err != nil {
		return domain.MenuItem{}, err
	}

	// sort_order が変更された場合、同一カテゴリ内の他のアイテムと入れ替え
	if currentSortOrder != input.SortOrder {
		_ = qtx.SwapMenuItemSortOrder(ctx, SwapMenuItemSortOrderParams{
			CategoryID:  currentCategoryID,
			SortOrder:   input.SortOrder,
			SortOrder_2: currentSortOrder,
			ID:          pgtypeUID,
		})
	}

	item, err := qtx.UpdateMenuItem(ctx, UpdateMenuItemParams{
		ID:          pgtypeUID,
		Name:        input.Name,
		Price:       input.Price,
		Description: input.Description,
		SortOrder:   input.SortOrder,
	})
	if err != nil {
		return domain.MenuItem{}, err
	}

	err = tx.Commit(ctx)
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
