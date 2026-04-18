-- name: GetStaffAccountByUsername :one
SELECT * FROM staff_accounts WHERE username = $1;

-- name: GetStaffAccountByStaffID :one
SELECT * FROM staff_accounts WHERE staff_id = $1;

-- name: CreateStaffAccount :one
INSERT INTO staff_accounts (staff_id, username, password_hash)
VALUES ($1, $2, $3) RETURNING *;

-- name: ListStaffAccounts :many
SELECT * FROM staff_accounts ORDER BY created_at ASC;

-- name: UpdateStaffAccount :one
UPDATE staff_accounts SET username = $2, password_hash = $3 WHERE id = $1 RETURNING *;

-- name: DeleteStaffAccount :exec
DELETE FROM staff_accounts WHERE id = $1;
