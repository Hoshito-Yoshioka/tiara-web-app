-- name: ListMenuCategories :many
SELECT * FROM menu_categories ORDER BY sort_order ASC;

-- name: GetMenuCategoryByID :one
SELECT * FROM menu_categories WHERE id = $1;

-- name: CreateMenuCategory :one
INSERT INTO menu_categories (name, description, sort_order)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateMenuCategory :one
UPDATE menu_categories
SET name = $2, description = $3, sort_order = $4
WHERE id = $1
RETURNING *;

-- name: DeleteMenuCategory :exec
DELETE FROM menu_categories WHERE id = $1;

-- name: ListMenuItemsByCategoryID :many
SELECT * FROM menu_items WHERE category_id = $1 ORDER BY sort_order ASC;

-- name: ListAllMenuItems :many
SELECT * FROM menu_items ORDER BY sort_order ASC;

-- name: CreateMenuItem :one
INSERT INTO menu_items (category_id, name, price, description, sort_order)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateMenuItem :one
UPDATE menu_items
SET name = $2, price = $3, description = $4, sort_order = $5
WHERE id = $1
RETURNING *;

-- name: DeleteMenuItem :exec
DELETE FROM menu_items WHERE id = $1;

-- name: DeleteMenuItemsByCategoryID :exec
DELETE FROM menu_items WHERE category_id = $1;

-- name: SwapMenuCategorySortOrder :exec
UPDATE menu_categories SET sort_order = $2 WHERE sort_order = $1 AND id != $3;

-- name: SwapMenuItemSortOrder :exec
UPDATE menu_items SET sort_order = $3 WHERE category_id = $1 AND sort_order = $2 AND id != $4;
