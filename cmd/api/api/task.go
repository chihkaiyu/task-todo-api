package api

import (
	"github.com/gin-gonic/gin"

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
}

// @Summary Create task
// @Tags task
// @Accept json
// @Produce json
// @Success 200 {object} models.CreateTaskResp
// @Failure 400 {object} models.BaseError
// @Failure 500 {object} models.BaseError
// @Router /task [post]
func (th *taskHandler) createTask(c *gin.Context) {
}

// @Summary Put task
// @Tags task
// @Accept json
// @Produce json
// @Success 200 {object} models.PutTaskResp
// @Failure 400 {object} models.BaseError
// @Failure 500 {object} models.BaseError
// @Router /task/{id} [put]
func (th *taskHandler) putTask(c *gin.Context) {
}

// @Summary Delete task
// @Tags task
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Failure 400 {object} models.BaseError
// @Failure 500 {object} models.BaseError
// @Router /task/{id} [delete]
func (th *taskHandler) deleteTask(c *gin.Context) {
}
