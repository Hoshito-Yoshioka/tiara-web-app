-- name: GetProfileDraftByID :one
SELECT * FROM staff_profile_drafts WHERE id = $1;

-- name: GetProfileDraftByStaffID :one
SELECT * FROM staff_profile_drafts WHERE staff_id = $1 AND status IN ('draft', 'pending', 'rejected') ORDER BY created_at DESC LIMIT 1;

-- name: ListPendingProfileDrafts :many
SELECT * FROM staff_profile_drafts WHERE status = 'pending' ORDER BY submitted_at ASC;

-- name: CreateProfileDraft :one
INSERT INTO staff_profile_drafts (staff_id, name, role, bio, image_url, external_schedule_url, image_crop_position, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: UpdateProfileDraft :one
UPDATE staff_profile_drafts
SET name = $2, role = $3, bio = $4, image_url = $5, external_schedule_url = $6, image_crop_position = $7, status = $8
WHERE id = $1 AND updated_at = $9 RETURNING *;

-- name: SubmitProfileDraft :one
UPDATE staff_profile_drafts SET status = 'pending', submitted_at = NOW() WHERE id = $1 AND updated_at = $2 RETURNING *;

-- name: ReviewProfileDraft :one
UPDATE staff_profile_drafts SET status = $2, admin_comment = $3, reviewed_at = NOW() WHERE id = $1 AND updated_at = $4 RETURNING *;

-- name: UpdateProfileDraftContent :one
UPDATE staff_profile_drafts
SET name = $2, role = $3, bio = $4, image_url = $5, external_schedule_url = $6, image_crop_position = $7
WHERE id = $1 AND updated_at = $8 RETURNING *;

-- name: DeleteProfileDraft :exec
DELETE FROM staff_profile_drafts WHERE id = $1;
