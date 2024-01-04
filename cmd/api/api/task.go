package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	mw "github.com/chihkaiyu/task-todo-api/middlewares"
	"github.com/chihkaiyu/task-todo-api/models"
	"github.com/chihkaiyu/task-todo-api/stores/tasks"
)

type taskHandler struct {
	taskStore tasks.Task
}

func NewTaskHandler(taskRG *gin.RouterGroup, taskStore tasks.Task) {
	th := taskHandler{
		taskStore: taskStore,
	}

	taskRG.GET("/tasks", th.listTask)
	taskRG.POST("/task", th.createTask)
	taskRG.PUT("task/:id", th.putTask)
	taskRG.DELETE("/task/:id", th.deleteTask)
}

// @Summary List tasks
// @Tags task
// @Accept json
// @Produce json
// @Success 200 {object} models.ListTaskResp
// @Failure 400 {object} models.BaseError
// @Failure 500 {object} models.BaseError
// @Router /tasks [get]
func (th *taskHandler) listTask(c *gin.Context) {
	ctx := c.Request.Context()

	tasks, err := th.taskStore.List(ctx)
	if err != nil {
		mw.Error(c, err)
		return
	}

	dt := make([]*models.DisplayTask, len(tasks))
	for i, t := range tasks {
		dt[i] = t.Parse()
	}

	mw.JSON(c, http.StatusOK, models.ListTaskResp{
		Result: dt,
	})
}

// @Summary Create task
// @Tags task
// @Accept json
// @Produce json
// @Param CreateTaskParams body models.CreateTaskParams true "parameters for creating task"
// @Success 200 {object} models.CreateTaskResp
// @Failure 400 {object} models.BaseError
// @Failure 500 {object} models.BaseError
// @Router /task [post]
func (th *taskHandler) createTask(c *gin.Context) {
	ctx := c.Request.Context()

	params := models.CreateTaskParams{}
	if err := c.ShouldBindJSON(&params); err != nil {
		mw.Error(c, err)
		return
	}

	task, err := th.taskStore.Create(ctx, params.Name)
	if err != nil {
		mw.Error(c, err)
		return
	}

	mw.JSON(c, http.StatusCreated, models.CreateTaskResp{
		Result: task.Parse(),
	})
}

// @Summary Put task
// @Tags task
// @Accept json
// @Produce json
// @Param id path string true "task's ID"
// @Param PutTaskParams body models.PutTaskParams true "parameters for updating task"
// @Success 200 {object} models.PutTaskResp
// @Failure 400 {object} models.BaseError
// @Failure 500 {object} models.BaseError
// @Router /task/{id} [put]
func (th *taskHandler) putTask(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	params := models.PutTaskParams{}
	if err := c.ShouldBindJSON(&params); err != nil {
		mw.Error(c, err)
		return
	}

	task, err := th.taskStore.Put(ctx, id, &params)
	if err != nil {
		mw.Error(c, err)
		return
	}

	mw.JSON(c, http.StatusOK, models.PutTaskResp{
		Result: task.Parse(),
	})
}

// @Summary Delete task
// @Tags task
// @Accept json
// @Produce json
// @Param id path string true "task's ID"
// @Success 200 {object} string
// @Failure 400 {object} models.BaseError
// @Failure 500 {object} models.BaseError
// @Router /task/{id} [delete]
func (th *taskHandler) deleteTask(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	if err := th.taskStore.Delete(ctx, id); err != nil {
		mw.Error(c, err)
		return
	}

	mw.JSON(c, http.StatusOK, gin.H{})
}
