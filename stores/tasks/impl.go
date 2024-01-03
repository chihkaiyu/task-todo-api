package tasks

import (
	"context"
	"errors"

	"github.com/chihkaiyu/task-todo-api/models"
)

type impl struct{}

func New() Task {
	return &impl{}
}

func (im *impl) Create(ctx context.Context, name string) (*models.Task, error) {
	return nil, errors.New("not implemented")
}

func (im *impl) Get(ctx context.Context, uuid string) (*models.Task, error) {
	return nil, errors.New("not implemented")
}

func (im *impl) List(ctx context.Context) ([]*models.Task, error) {
	return nil, errors.New("not implemented")
}

func (im *impl) Put(ctx context.Context, params *models.UpdateTaskParams) (*models.Task, error) {
	return nil, errors.New("not implemented")
}

func (im *impl) Delete(ctx context.Context, uuid string) error {
	return errors.New("not implemented")
}
