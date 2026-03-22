package domain

import (
	"time"

	"github.com/google/uuid"
)

// ドラフトステータス定数
const (
	DraftStatusDraft    = "draft"
	DraftStatusPending  = "pending"
	DraftStatusApproved = "approved"
	DraftStatusRejected = "rejected"
)

// StaffAccount はスタッフポータル用のアカウント情報。
type StaffAccount struct {
	ID           uuid.UUID
	StaffID      uuid.UUID
	Username     string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// StaffProfileDraft はスタッフがポータルから提出するプロフィール変更の下書き。
type StaffProfileDraft struct {
	ID                uuid.UUID
	StaffID           uuid.UUID
	Name              string
	Role              string
	Bio               string
	ImageURL          string
	ImageCropPosition string
	Status            string
	AdminComment      string
	SubmittedAt       *time.Time
	ReviewedAt        *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// StaffScheduleDraft はスタッフがポータルから提出するスケジュール変更の下書き。
type StaffScheduleDraft struct {
	ID           uuid.UUID
	StaffID      uuid.UUID
	Status       string
	AdminComment string
	SubmittedAt  *time.Time
	ReviewedAt   *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Items        []ScheduleDraftItem
}

// ScheduleDraftItem はスケジュール下書きの個別の出勤枠。
// DayOfWeek: 0=日, 1=月, 2=火, 3=水, 4=木, 5=金, 6=土
type ScheduleDraftItem struct {
	ID        uuid.UUID
	DraftID   uuid.UUID
	DayOfWeek int
	StartTime time.Time
	EndTime   time.Time
}

// SaveProfileDraftInput はプロフィール下書き保存時の入力構造体。
type SaveProfileDraftInput struct {
	Name              string
	Role              string
	Bio               string
	ImageURL          string
	ImageCropPosition string
}

// ReviewDraftInput は管理者がドラフトをレビューする際の入力構造体。
type ReviewDraftInput struct {
	Status       string
	AdminComment string
}
