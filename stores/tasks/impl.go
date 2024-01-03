package tasks

import (
	"context"
	"time"

	"github.com/chihkaiyu/task-todo-api/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rs/zerolog"
)

var timeNow = time.Now

type impl struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Task {
	return &impl{
		db: db,
	}
}

func (im *impl) Create(ctx context.Context, name string) (*models.Task, error) {
	s := "INSERT INTO tasks (id, name, status, created_at, updated_at)\n" +
		"VALUES (:id, :name, :status, :created_at, :updated_at)"
	now := timeNow().UTC()
	task := &models.Task{
		ID:        uuid.New(),
		Name:      name,
		Status:    0,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: pq.NullTime{},
	}
	_, err := im.db.NamedExec(s, task)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("im.db.NamedExec failed")
		return nil, err
	}

	return task, nil
}

func (im *impl) Get(ctx context.Context, uuidStr string) (*models.Task, error) {
	uid, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, err
	}

	s := "SELECT id, name, status FROM tasks WHERE id=$1"
	task := &models.Task{}
	if err := im.db.Get(task, s, uid); err != nil {
		return nil, err
	}

	return task, nil
}

func (im *impl) List(ctx context.Context, opts ...ListTaskOptionFunc) ([]*models.Task, error) {
	opt := ListTaskOption{}
	for _, f := range opts {
		f(&opt)
	}
	s := "SELECT id, name, status FROM tasks\n"
	if !opt.WithDeleted {
		s += "WHERE deleted_at IS NULL"
	}
	tasks := []*models.Task{}
	if err := im.db.Select(&tasks, s); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (im *impl) Put(ctx context.Context, uuidStr string, params *models.UpdateTaskParams) (*models.Task, error) {
	uid, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, err
	}

	s := "UPDATE tasks SET name=:name, status=:status, updated_at=:updated_at WHERE id=:id RETURNING name, status"
	now := timeNow().UTC()
	task := &models.Task{
		ID:        uid,
		Name:      params.Name,
		Status:    params.Status,
		UpdatedAt: now,
	}
	row, err := im.db.NamedQuery(s, task)
	if err != nil {
		return nil, err
	}

	updated := &models.Task{}
	for row.Next() {
		if err := row.StructScan(updated); err != nil {
			return nil, err
		}
	}
	if updated.ID == uuid.Nil {
		return nil, ErrTaskNotFound
	}

	return updated, nil
}

func (im *impl) Delete(ctx context.Context, uuidStr string) error {
	uid, err := uuid.Parse(uuidStr)
	if err != nil {
		return err
	}

	now := timeNow().UTC()
	s := "UPDATE tasks SET deleted_at=$1 WHERE id=$2"
	_, err = im.db.Exec(s, now, uid)
	if err != nil {
		return err
	}
	return nil
}
