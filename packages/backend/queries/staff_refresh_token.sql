-- name: CreateRefreshToken :one
INSERT INTO staff_refresh_tokens (staff_id, token, expires_at)
VALUES ($1, $2, $3) RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM staff_refresh_tokens WHERE token = $1;

-- name: DeleteRefreshToken :exec
DELETE FROM staff_refresh_tokens WHERE token = $1;

-- name: DeleteRefreshTokensByStaffID :exec
DELETE FROM staff_refresh_tokens WHERE staff_id = $1;

-- name: DeleteExpiredRefreshTokens :exec
DELETE FROM staff_refresh_tokens WHERE expires_at < NOW();
