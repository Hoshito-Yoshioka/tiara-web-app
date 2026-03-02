package domain

import (
	"time"

	"github.com/google/uuid"
)

type Shop struct {
	ID         uuid.UUID
	Name       string
	Address    string
	OpeningTime time.Time
	ClosingTime time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
