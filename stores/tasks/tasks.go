package tasks

import (
	"context"

	"github.com/chihkaiyu/task-todo-api/models"
)

var (
	ErrTaskNotFound = models.NotFoundErr{Code: "TASK_NOT_FOUND"}
	ErrInvalidID    = models.BadRequestErr{Code: "INVALID_ID"}
)

type ListTaskOption struct {
	WithDeleted bool
}

type ListTaskOptionFunc func(*ListTaskOption)

func WithDeleted() ListTaskOptionFunc {
	return func(to *ListTaskOption) {
		to.WithDeleted = true
	}
}

type Task interface {
	Create(ctx context.Context, name string) (*models.Task, error)
	Get(ctx context.Context, id string) (*models.Task, error)
	List(ctx context.Context, opts ...ListTaskOptionFunc) ([]*models.Task, error)
	Put(ctx context.Context, id string, params *models.PutTaskParams) (*models.Task, error)
	Delete(ctx context.Context, id string) error
}
