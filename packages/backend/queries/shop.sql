-- name: ListShops :many
SELECT * FROM shops;

-- name: GetShopByID :one
SELECT * FROM shops WHERE id = $1;
