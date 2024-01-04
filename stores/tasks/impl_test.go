package tasks

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	bdocker "github.com/chihkaiyu/task-todo-api/base/docker"
	"github.com/chihkaiyu/task-todo-api/models"
	"github.com/chihkaiyu/task-todo-api/services/postgres"
)

var (
	mockCTX   = context.Background()
	mockNow   = time.Now().UTC()
	mockUUID  = uuid.New()
	mockUUID2 = uuid.New()
)

type mockFuncs struct {
	mock.Mock
}

func (m *mockFuncs) timeNow() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

type taskSuite struct {
	suite.Suite
	taskStore    *impl
	db           *sqlx.DB
	postgresPort string

	mockFuncs *mockFuncs
}

func TestTaskSuite(t *testing.T) {
	suite.Run(t, new(taskSuite))
}

func (s *taskSuite) SetupSuite() {
	ports, err := bdocker.RunExternal([]string{"postgres"})
	s.Require().NoError(err)
	s.postgresPort = ports[0]
}

func (s *taskSuite) TearDownSuite() {
	s.NoError(bdocker.RemoveExternal())
}

func (s *taskSuite) SetupTest() {
	createDB("redreamer", s.postgresPort)
	create("redreamer", s.postgresPort)

	db, err := postgres.New(fmt.Sprintf("postgres://postgres@localhost:%s/redreamer?sslmode=disable", s.postgresPort))
	s.Require().NoError(err)
	s.db = db
	s.mockFuncs = new(mockFuncs)
	s.taskStore = New(s.db).(*impl)

	// mock functions
	timeNow = s.mockFuncs.timeNow
}

func (s *taskSuite) TearDownTest() {
	s.mockFuncs.AssertExpectations(s.T())

	s.db.Close()
	s.Require().NoError(bdocker.ClearPostgres(s.postgresPort))
}

func createDB(name, port string) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://postgres@localhost:%s/?sslmode=disable", port))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE " + name)
	if err != nil {
		panic(err)
	}
}

func create(name, port string) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://postgres@localhost:%s/%s?sslmode=disable", port, name))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	migrations := &migrate.FileMigrationSource{
		Dir: "../../infra/databases/api/migrations",
	}

	_, err = migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

type createTaskOption struct {
	id uuid.UUID
}

type createTaskOptionFunc func(*createTaskOption)

func createWithID(id uuid.UUID) createTaskOptionFunc {
	return func(cto *createTaskOption) {
		cto.id = id
	}
}

func (s *taskSuite) createTask(opts ...createTaskOptionFunc) *models.Task {
	task := &models.Task{
		ID:        mockUUID,
		Name:      "mock-task-name",
		Status:    0,
		CreatedAt: mockNow,
		UpdatedAt: mockNow,
		DeletedAt: pq.NullTime{},
	}

	opt := createTaskOption{}
	for _, f := range opts {
		f(&opt)
	}
	if opt.id != uuid.Nil {
		task.ID = opt.id
	}

	insertSQL := "INSERT INTO tasks (id, name, status) VALUES (:id, :name, :status)"
	_, err := s.db.NamedExec(insertSQL, task)
	s.Require().NoError(err)
	return task
}

func (s *taskSuite) deleteTask(id uuid.UUID) {
	deleteSQL := "UPDATE tasks SET deleted_at=$1 WHERE id=$2"
	_, err := s.db.Exec(deleteSQL, mockNow, id)
	s.Require().NoError(err)
}

func (s *taskSuite) TestGet() {
	tests := []struct {
		desc   string
		id     string
		expErr error
	}{
		{
			desc:   "get normally",
			id:     mockUUID.String(),
			expErr: nil,
		},
		{
			desc:   "not found",
			id:     mockUUID2.String(),
			expErr: sql.ErrNoRows,
		},
		{
			desc:   "invalid id",
			id:     "invalid-uuid",
			expErr: ErrInvalidID,
		},
	}

	s.TearDownTest()
	for _, test := range tests {
		s.SetupTest()

		expected := s.createTask()

		act, err := s.taskStore.Get(mockCTX, test.id)
		if test.expErr != nil {
			s.Require().EqualError(err, test.expErr.Error(), test.desc)
		} else {
			s.Require().NoError(err, test.desc)

			s.Require().Equal(expected.ID, act.ID, test.desc)
			s.Require().Equal(expected.Name, act.Name, test.desc)
			s.Require().Equal(expected.Status, act.Status, test.desc)
		}

		s.TearDownTest()
	}
}

func (s *taskSuite) TestCreate() {
	tests := []struct {
		desc     string
		mockFunc func()
		name     string
	}{
		{
			desc: "create normally",
			mockFunc: func() {
				s.mockFuncs.On("timeNow").Return(mockNow).Once()
			},
			name: "mock-task-name",
		},
	}

	s.TearDownTest()
	for _, test := range tests {
		s.SetupTest()

		test.mockFunc()

		expected, err := s.taskStore.Create(mockCTX, test.name)
		s.Require().NoError(err, test.desc)

		act, err := s.taskStore.Get(mockCTX, expected.ID.String())
		s.Require().NoError(err, test.desc)

		s.Require().Equal(expected.ID, act.ID, test.desc)
		s.Require().Equal(expected.Name, act.Name, test.desc)
		s.Require().Equal(expected.Status, act.Status, test.desc)

		s.TearDownTest()
	}
}

func (s *taskSuite) TestList() {
	tests := []struct {
		desc     string
		mockFunc func()
		opts     []ListTaskOptionFunc
		expNum   int
	}{
		{
			desc: "list normally",
			mockFunc: func() {
				for i := 0; i < 3; i++ {
					s.createTask(createWithID(uuid.New()))
				}
			},
			opts:   []ListTaskOptionFunc{},
			expNum: 3,
		},
		{
			desc: "list all without deleted",
			mockFunc: func() {
				for i := 0; i < 3; i++ {
					s.createTask(createWithID(uuid.New()))
				}
				for i := 0; i < 2; i++ {
					t := s.createTask(createWithID(uuid.New()))
					s.deleteTask(t.ID)
				}
			},
			opts:   []ListTaskOptionFunc{},
			expNum: 3,
		},
		{
			desc: "list all with deleted",
			mockFunc: func() {
				for i := 0; i < 3; i++ {
					s.createTask(createWithID(uuid.New()))
				}
				for i := 0; i < 2; i++ {
					t := s.createTask(createWithID(uuid.New()))
					s.deleteTask(t.ID)
				}
			},
			opts:   []ListTaskOptionFunc{WithDeleted()},
			expNum: 5,
		},
	}

	s.TearDownTest()
	for _, test := range tests {
		s.SetupTest()

		test.mockFunc()

		tasks, err := s.taskStore.List(mockCTX, test.opts...)
		s.Require().NoError(err, test.desc)
		s.Require().Len(tasks, test.expNum, test.desc)

		s.TearDownTest()
	}
}

func (s *taskSuite) TestPut() {
	tests := []struct {
		desc     string
		mockFunc func()
		id       string
		params   *models.PutTaskParams
		expTask  *models.Task
		expErr   error
	}{
		{
			desc: "put normally",
			mockFunc: func() {
				s.createTask()
				s.mockFuncs.On("timeNow").Return(mockNow.Add(7 * time.Minute)).Once()
			},
			id: mockUUID.String(),
			params: &models.PutTaskParams{
				Name:   "updated-task-name",
				Status: 1,
			},
			expTask: &models.Task{
				ID:        mockUUID,
				Name:      "updated-task-name",
				Status:    1,
				UpdatedAt: mockNow.Add(7 * time.Minute),
			},
			expErr: nil,
		},
		{
			desc:     "invalid uuid",
			mockFunc: func() {},
			id:       "mock-invalid-id",
			params: &models.PutTaskParams{
				Name:   "updated-task-name",
				Status: 1,
			},
			expTask: nil,
			expErr:  ErrInvalidID,
		},
		{
			desc: "put non-exist task",
			mockFunc: func() {
				s.createTask()
				s.mockFuncs.On("timeNow").Return(mockNow.Add(7 * time.Minute)).Once()
			},
			id: mockUUID2.String(),
			params: &models.PutTaskParams{
				Name:   "updated-task-name",
				Status: 1,
			},
			expTask: nil,
			expErr:  ErrTaskNotFound,
		},
	}

	s.TearDownTest()
	for _, test := range tests {
		s.SetupTest()

		test.mockFunc()
		updated, err := s.taskStore.Put(mockCTX, test.id, test.params)
		if test.expErr != nil {
			s.Require().EqualError(err, test.expErr.Error(), test.desc)
		} else {
			s.Require().NoError(err, test.desc)

			s.Require().Equal(test.expTask.ID, updated.ID, test.desc)
			s.Require().Equal(test.expTask.Name, updated.Name, test.desc)
			s.Require().Equal(test.expTask.Status, updated.Status, test.desc)
			s.Require().Equal(test.expTask.UpdatedAt, updated.UpdatedAt, test.desc)
		}

		s.TearDownTest()
	}
}

func (s *taskSuite) TestDelete() {
	tests := []struct {
		desc     string
		mockFunc func()
		id       string
		deleted  bool
		expErr   error
	}{
		{
			desc: "delete normally",
			mockFunc: func() {
				s.createTask()
				s.mockFuncs.On("timeNow").Return(mockNow.Add(7 * time.Minute)).Once()
			},
			id:      mockUUID.String(),
			deleted: true,
			expErr:  nil,
		},
		{
			desc: "invalid id",
			mockFunc: func() {
			},
			id:      "mock-invalid-id",
			deleted: false,
			expErr:  ErrInvalidID,
		},
		{
			desc: "delete non-exist task",
			mockFunc: func() {
				s.createTask()
				s.mockFuncs.On("timeNow").Return(mockNow.Add(7 * time.Minute)).Once()
			},
			id:      mockUUID2.String(),
			deleted: false,
			expErr:  nil,
		},
	}

	s.TearDownTest()
	for _, test := range tests {
		s.SetupTest()

		test.mockFunc()
		err := s.taskStore.Delete(mockCTX, test.id)
		if test.expErr != nil {
			s.Require().EqualError(err, test.expErr.Error(), test.desc)
		} else {
			s.Require().NoError(err, test.desc)

			if test.deleted {
				deleted, err := s.taskStore.Get(mockCTX, test.id)
				s.Require().NoError(err, test.desc)
				s.Require().NotNil(deleted.DeletedAt, test.desc)
			}
		}

		s.TearDownTest()
	}
}
