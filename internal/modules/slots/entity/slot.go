package entity

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID    uuid.UUID
	Name  string
	Color string
}

type Slot struct {
	ID          uuid.UUID
	Name        string
	Description string
	СategoryID  uuid.UUID
	UserID      uuid.UUID
	StartTime   time.Time
	EndTime     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
