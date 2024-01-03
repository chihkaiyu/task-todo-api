package models

import "github.com/google/uuid"

type Task struct {
	ID     int       `db:"id"`
	UUID   uuid.UUID `db:"uuid"`
	Name   string    `db:"name"`
	Status int       `db:"status"`
}

type DisplayTask struct {
	UUID   uuid.UUID `json:"uuid"`
	Name   string    `json:"name"`
	Status int       `json:"status"`
}

func (t *Task) Parse() *DisplayTask {
	return &DisplayTask{
		UUID:   t.UUID,
		Name:   t.Name,
		Status: t.Status,
	}
}

type UpdateTaskParams struct {
	Name   string `json:"name"`
	Status int    `json:"status"`
}
