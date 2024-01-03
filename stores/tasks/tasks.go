package tasks

import (
	"context"

	"github.com/chihkaiyu/task-todo-api/models"
)

type Task interface {
	Create(ctx context.Context, name string) (*models.Task, error)
	Get(ctx context.Context, uuid string) (*models.Task, error)
	List(ctx context.Context) ([]*models.Task, error)
	Put(ctx context.Context, params *models.UpdateTaskParams) (*models.Task, error)
	Delete(ctx context.Context, uuid string) error
}
