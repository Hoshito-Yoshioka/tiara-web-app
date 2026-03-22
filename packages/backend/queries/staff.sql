-- name: ListStaffs :many
SELECT * FROM staffs ORDER BY sort_order ASC;

-- name: GetStaffByID :one
SELECT * FROM staffs WHERE id = $1;

-- name: ListStaffsByShopID :many
SELECT * FROM staffs WHERE shop_id = $1 ORDER BY sort_order ASC;

-- name: ListSchedulesByStaffID :many
SELECT * FROM staff_schedules WHERE staff_id = $1 ORDER BY day_of_week ASC;

-- name: ListAllSchedules :many
SELECT * FROM staff_schedules ORDER BY staff_id, day_of_week ASC;

-- name: CreateStaff :one
INSERT INTO staffs (shop_id, name, role, bio, image_url, image_crop_position, sort_order) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: UpdateStaff :one
UPDATE staffs SET name = $2, role = $3, bio = $4, image_url = $5, image_crop_position = $6, sort_order = $7 WHERE id = $1 RETURNING *;

-- name: DeleteStaff :exec
DELETE FROM staffs WHERE id = $1;

-- name: DeleteSchedulesByStaffID :exec
DELETE FROM staff_schedules WHERE staff_id = $1;

-- name: CreateSchedule :one
INSERT INTO staff_schedules (staff_id, day_of_week, start_time, end_time) VALUES ($1, $2, $3, $4) RETURNING *;
