-- name: ListShops :many
SELECT * FROM shops;

-- name: GetShopByID :one
SELECT * FROM shops WHERE id = $1;

-- name: CreateShop :one
INSERT INTO shops (name, address, opening_time, closing_time) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateShop :one
UPDATE shops SET name = $2, address = $3, opening_time = $4, closing_time = $5 WHERE id = $1 RETURNING *;
