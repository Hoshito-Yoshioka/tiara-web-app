package domain

import (
	"time"

	"github.com/google/uuid"
)

// Staff はスタッフのドメインモデル。
type Staff struct {
	ID        uuid.UUID
	ShopID    uuid.UUID
	Name      string
	Role      string
	Bio       string
	ImageURL  string
	SortOrder int
	CreatedAt time.Time
	UpdatedAt time.Time
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

// StaffWithSchedules はスタッフ情報に出勤スケジュールを含めた集約モデル。
// StaffDetailView用に使用する。
type StaffWithSchedules struct {
	Staff     Staff
	Schedules []StaffSchedule
}
