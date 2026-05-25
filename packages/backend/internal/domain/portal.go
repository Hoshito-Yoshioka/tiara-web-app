package domain

import (
	"time"

	"github.com/google/uuid"
)

// DraftStatus はドラフトの状態を表す値オブジェクト。
type DraftStatus string

// ドラフトステータス定数
const (
	DraftStatusDraft    DraftStatus = "draft"
	DraftStatusPending  DraftStatus = "pending"
	DraftStatusApproved DraftStatus = "approved"
	DraftStatusRejected DraftStatus = "rejected"
)

// IsValid は DraftStatus が有効な値かを検証する。
func (s DraftStatus) IsValid() bool {
	switch s {
	case DraftStatusDraft, DraftStatusPending, DraftStatusApproved, DraftStatusRejected:
		return true
	}
	return false
}

// IsEditable は下書きが編集・再申請可能な状態かを返す（draft または rejected）。
func (s DraftStatus) IsEditable() bool {
	return s == DraftStatusDraft || s == DraftStatusRejected
}

// String は DraftStatus の文字列表現を返す。
func (s DraftStatus) String() string {
	return string(s)
}

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
	Status            DraftStatus
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
	Status       DraftStatus
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
	Status       DraftStatus
	AdminComment string
}
