package entity

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID  `db:"id"`
	Name        string     `db:"name"`
	Description string     `db:"description"`
	IsCompleted bool       `db:"is_completed"`
	DueDate     *time.Time `db:"due_date"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}
