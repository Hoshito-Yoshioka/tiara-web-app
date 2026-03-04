package domain

import (
	"time"

	"github.com/google/uuid"
)

// Admin は管理者のドメインモデル。
type Admin struct {
	ID           uuid.UUID
	Username     string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
