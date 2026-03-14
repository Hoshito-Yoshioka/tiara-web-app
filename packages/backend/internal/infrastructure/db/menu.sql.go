package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

// ============================================================
// MenuCategory queries
// ============================================================

const listMenuCategories = `
SELECT id, name, description, sort_order, created_at, updated_at
FROM menu_categories
ORDER BY sort_order ASC
`

func (q *Queries) ListMenuCategories(ctx context.Context) ([]MenuCategory, error) {
	rows, err := q.db.Query(ctx, listMenuCategories)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []MenuCategory
	for rows.Next() {
		var i MenuCategory
		if err := rows.Scan(
			&i.ID, &i.Name, &i.Description, &i.SortOrder,
			&i.CreatedAt, &i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

const getMenuCategoryByID = `
SELECT id, name, description, sort_order, created_at, updated_at
FROM menu_categories WHERE id = $1
`

func (q *Queries) GetMenuCategoryByID(ctx context.Context, id pgtype.UUID) (MenuCategory, error) {
	row := q.db.QueryRow(ctx, getMenuCategoryByID, id)
	var i MenuCategory
	err := row.Scan(&i.ID, &i.Name, &i.Description, &i.SortOrder, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

type CreateMenuCategoryParams struct {
	Name        string
	Description string
	SortOrder   int32
}

const createMenuCategory = `
INSERT INTO menu_categories (name, description, sort_order)
VALUES ($1, $2, $3)
RETURNING id, name, description, sort_order, created_at, updated_at
`

func (q *Queries) CreateMenuCategory(ctx context.Context, arg CreateMenuCategoryParams) (MenuCategory, error) {
	row := q.db.QueryRow(ctx, createMenuCategory, arg.Name, arg.Description, arg.SortOrder)
	var i MenuCategory
	err := row.Scan(&i.ID, &i.Name, &i.Description, &i.SortOrder, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

type UpdateMenuCategoryParams struct {
	ID          pgtype.UUID
	Name        string
	Description string
	SortOrder   int32
}

const updateMenuCategory = `
UPDATE menu_categories
SET name = $2, description = $3, sort_order = $4
WHERE id = $1
RETURNING id, name, description, sort_order, created_at, updated_at
`

func (q *Queries) UpdateMenuCategory(ctx context.Context, arg UpdateMenuCategoryParams) (MenuCategory, error) {
	row := q.db.QueryRow(ctx, updateMenuCategory, arg.ID, arg.Name, arg.Description, arg.SortOrder)
	var i MenuCategory
	err := row.Scan(&i.ID, &i.Name, &i.Description, &i.SortOrder, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const deleteMenuCategory = `DELETE FROM menu_categories WHERE id = $1`

func (q *Queries) DeleteMenuCategory(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteMenuCategory, id)
	return err
}

// ============================================================
// MenuItem queries
// ============================================================

const listMenuItemsByCategoryID = `
SELECT id, category_id, name, price, description, sort_order, created_at, updated_at
FROM menu_items WHERE category_id = $1
ORDER BY sort_order ASC
`

func (q *Queries) ListMenuItemsByCategoryID(ctx context.Context, categoryID pgtype.UUID) ([]MenuItem, error) {
	rows, err := q.db.Query(ctx, listMenuItemsByCategoryID, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []MenuItem
	for rows.Next() {
		var i MenuItem
		if err := rows.Scan(
			&i.ID, &i.CategoryID, &i.Name, &i.Price, &i.Description,
			&i.SortOrder, &i.CreatedAt, &i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

type CreateMenuItemParams struct {
	CategoryID  pgtype.UUID
	Name        string
	Price       string
	Description string
	SortOrder   int32
}

const createMenuItem = `
INSERT INTO menu_items (category_id, name, price, description, sort_order)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, category_id, name, price, description, sort_order, created_at, updated_at
`

func (q *Queries) CreateMenuItem(ctx context.Context, arg CreateMenuItemParams) (MenuItem, error) {
	row := q.db.QueryRow(ctx, createMenuItem, arg.CategoryID, arg.Name, arg.Price, arg.Description, arg.SortOrder)
	var i MenuItem
	err := row.Scan(&i.ID, &i.CategoryID, &i.Name, &i.Price, &i.Description, &i.SortOrder, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

type UpdateMenuItemParams struct {
	ID          pgtype.UUID
	Name        string
	Price       string
	Description string
	SortOrder   int32
}

const updateMenuItem = `
UPDATE menu_items
SET name = $2, price = $3, description = $4, sort_order = $5
WHERE id = $1
RETURNING id, category_id, name, price, description, sort_order, created_at, updated_at
`

func (q *Queries) UpdateMenuItem(ctx context.Context, arg UpdateMenuItemParams) (MenuItem, error) {
	row := q.db.QueryRow(ctx, updateMenuItem, arg.ID, arg.Name, arg.Price, arg.Description, arg.SortOrder)
	var i MenuItem
	err := row.Scan(&i.ID, &i.CategoryID, &i.Name, &i.Price, &i.Description, &i.SortOrder, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const deleteMenuItem = `DELETE FROM menu_items WHERE id = $1`

func (q *Queries) DeleteMenuItem(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteMenuItem, id)
	return err
}

const deleteMenuItemsByCategoryID = `DELETE FROM menu_items WHERE category_id = $1`

func (q *Queries) DeleteMenuItemsByCategoryID(ctx context.Context, categoryID pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteMenuItemsByCategoryID, categoryID)
	return err
}
