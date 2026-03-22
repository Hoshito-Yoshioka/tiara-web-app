-- name: ListImagesByStaffID :many
SELECT * FROM staff_images WHERE staff_id = $1 ORDER BY is_main DESC, sort_order ASC;

-- name: ListAllStaffImages :many
SELECT * FROM staff_images ORDER BY staff_id, is_main DESC, sort_order ASC;

-- name: CreateStaffImage :one
INSERT INTO staff_images (staff_id, image_url, is_main, sort_order) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateStaffImage :one
UPDATE staff_images SET image_url = $2, is_main = $3, sort_order = $4 WHERE id = $1 RETURNING *;

-- name: DeleteStaffImage :exec
DELETE FROM staff_images WHERE id = $1;

-- name: DeleteImagesByStaffID :exec
DELETE FROM staff_images WHERE staff_id = $1;

-- name: ClearMainFlagByStaffID :exec
UPDATE staff_images SET is_main = false WHERE staff_id = $1;

-- name: SetMainImage :one
UPDATE staff_images SET is_main = true WHERE id = $1 RETURNING *;
