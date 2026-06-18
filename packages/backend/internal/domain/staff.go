package domain

import (
	"time"

	"github.com/google/uuid"
)

// Staff はスタッフのドメインモデル。
type Staff struct {
	ID                  uuid.UUID
	ShopID              uuid.UUID
	Name                string
	Role                string
	Bio                 string
	ImageURL            string
	ExternalScheduleURL string
	ImageCropPosition   string
	SortOrder           int
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// StaffSchedule はスタッフの出勤スケジュールのドメインモデル。
// DayOfWeek: 0=日, 1=月, 2=火, 3=水, 4=木, 5=金, 6=土
type StaffSchedule struct {
	ID        uuid.UUID
	StaffID   uuid.UUID
	DayOfWeek int
	StartTime time.Time
	EndTime   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// StaffWithSchedules はスタッフ情報に出勤スケジュールと画像を含めた集約モデル。
// StaffDetailView用に使用する。
type StaffWithSchedules struct {
	Staff     Staff
	Schedules []StaffSchedule
	Images    []StaffImage
}

// StaffImage はスタッフ画像のドメインモデル。
type StaffImage struct {
	ID           uuid.UUID
	StaffID      uuid.UUID
	ImageURL     string
	IsMain       bool
	SortOrder    int
	CropPosition string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// StaffRefreshToken はスタッフのリフレッシュトークンドメインモデル。
type StaffRefreshToken struct {
	ID        uuid.UUID
	StaffID   uuid.UUID
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

// CreateStaffInput はスタッフ作成時の入力構造体。
type CreateStaffInput struct {
	ShopID              string
	Name                string
	Role                string
	Bio                 string
	ImageURL            string
	ExternalScheduleURL string
	ImageCropPosition   string
	SortOrder           int
	Schedules           []ScheduleInput
}

// UpdateStaffInput はスタッフ更新時の入力構造体。
type UpdateStaffInput struct {
	Name                string
	Role                string
	Bio                 string
	ImageURL            string
	ExternalScheduleURL string
	ImageCropPosition   string
	SortOrder           int
	Schedules           []ScheduleInput
}

// ScheduleInput は出勤スケジュールの入力構造体。
type ScheduleInput struct {
	DayOfWeek int
	StartTime string
	EndTime   string
}

// Pagination はページネーションのメタ情報。
type Pagination struct {
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	TotalCount int `json:"totalCount"`
	TotalPages int `json:"totalPages"`
}

// PaginatedStaffs はページネーション付きスタッフ一覧。
type PaginatedStaffs struct {
	Data       []Staff    `json:"data"`
	Pagination Pagination `json:"pagination"`
}
