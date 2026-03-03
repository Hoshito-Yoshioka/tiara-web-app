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
