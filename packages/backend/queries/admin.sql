-- name: GetAdminByUsername :one
SELECT * FROM admins WHERE username = $1;
