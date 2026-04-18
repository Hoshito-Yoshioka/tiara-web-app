-- name: GetScheduleDraftByID :one
SELECT * FROM staff_schedule_drafts WHERE id = $1;

-- name: GetScheduleDraftByStaffID :one
SELECT * FROM staff_schedule_drafts WHERE staff_id = $1 AND status IN ('draft', 'pending', 'rejected') ORDER BY created_at DESC LIMIT 1;

-- name: ListPendingScheduleDrafts :many
SELECT * FROM staff_schedule_drafts WHERE status = 'pending' ORDER BY submitted_at ASC;

-- name: ListApprovedScheduleDrafts :many
SELECT * FROM staff_schedule_drafts WHERE status = 'approved' ORDER BY reviewed_at ASC;

-- name: CreateScheduleDraft :one
INSERT INTO staff_schedule_drafts (staff_id, status)
VALUES ($1, $2) RETURNING *;

-- name: UpdateScheduleDraftStatus :one
UPDATE staff_schedule_drafts SET status = $2 WHERE id = $1 RETURNING *;

-- name: SubmitScheduleDraft :one
UPDATE staff_schedule_drafts SET status = 'pending', submitted_at = NOW() WHERE id = $1 RETURNING *;

-- name: ReviewScheduleDraft :one
UPDATE staff_schedule_drafts SET status = $2, admin_comment = $3, reviewed_at = NOW() WHERE id = $1 RETURNING *;

-- name: DeleteScheduleDraft :exec
DELETE FROM staff_schedule_drafts WHERE id = $1;

-- name: ListScheduleDraftItems :many
SELECT * FROM staff_schedule_draft_items WHERE draft_id = $1 ORDER BY day_of_week ASC;

-- name: CreateScheduleDraftItem :one
INSERT INTO staff_schedule_draft_items (draft_id, day_of_week, start_time, end_time)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: DeleteScheduleDraftItems :exec
DELETE FROM staff_schedule_draft_items WHERE draft_id = $1;
