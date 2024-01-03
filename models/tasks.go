package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Task struct {
	// TODO:
	PK        int         `db:"pk"`
	ID        uuid.UUID   `db:"id"`
	Name      string      `db:"name"`
	Status    int         `db:"status"`
	CreatedAt time.Time   `db:"created_at"`
	UpdatedAt time.Time   `db:"updated_at"`
	DeletedAt pq.NullTime `db:"deleted_at"`
}

type DisplayTask struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Status int       `json:"status"`
}

func (t *Task) Parse() *DisplayTask {
	return &DisplayTask{
		ID:     t.ID,
		Name:   t.Name,
		Status: t.Status,
	}
}

type UpdateTaskParams struct {
	Name   string `json:"name"`
	Status int    `json:"status"`
}