package domain

import (
	"time"

	"github.com/google/uuid"
)

type Shop struct {
	ID          uuid.UUID
	Name        string
	Address     string
	OpeningTime time.Time
	ClosingTime time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// UpdateShopInput は店舗更新時の入力構造体。
type UpdateShopInput struct {
	Name        string
	Address     string
	OpeningTime string
	ClosingTime string
}
