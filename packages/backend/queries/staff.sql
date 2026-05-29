-- name: ListStaffs :many
SELECT * FROM staffs ORDER BY sort_order ASC;

-- name: ListStaffsPaginated :many
SELECT * FROM staffs ORDER BY sort_order ASC LIMIT $1 OFFSET $2;

-- name: CountStaffs :one
SELECT count(*) FROM staffs;

-- name: GetStaffByID :one
SELECT * FROM staffs WHERE id = $1;

-- name: ListStaffsByShopID :many
SELECT * FROM staffs WHERE shop_id = $1 ORDER BY sort_order ASC;

-- name: ListSchedulesByStaffID :many
SELECT * FROM staff_schedules WHERE staff_id = $1 ORDER BY day_of_week ASC;

-- name: ListAllSchedules :many
SELECT * FROM staff_schedules ORDER BY staff_id, day_of_week ASC;

-- name: CreateStaff :one
INSERT INTO staffs (shop_id, name, role, bio, image_url, external_schedule_url, image_crop_position, sort_order) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: UpdateStaff :one
UPDATE staffs SET name = $2, role = $3, bio = $4, image_url = $5, external_schedule_url = $6, image_crop_position = $7, sort_order = $8 WHERE id = $1 RETURNING *;

-- name: DeleteStaff :exec
DELETE FROM staffs WHERE id = $1;

-- name: DeleteSchedulesByStaffID :exec
DELETE FROM staff_schedules WHERE staff_id = $1;

-- name: CreateSchedule :one
INSERT INTO staff_schedules (staff_id, day_of_week, start_time, end_time) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: SwapStaffSortOrder :exec
UPDATE staffs SET sort_order = $3 WHERE shop_id = $1 AND sort_order = $2 AND id != $4;
