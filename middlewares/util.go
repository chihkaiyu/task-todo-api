package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/chihkaiyu/task-todo-api/models"
	"github.com/chihkaiyu/task-todo-api/services/metrics"
)

var met = metrics.New("api")

func JSON(c *gin.Context, code int, obj interface{}) {
	c.JSON(code, obj)
}

func Data(c *gin.Context, code int, contentType string, data []byte) {
	c.Data(code, contentType, data)
}

func Error(c *gin.Context, err interface{}) {
	var code int
	switch e := err.(type) {
	case models.BadRequestErr:
		code = http.StatusBadRequest
	case models.AuthorizationErr:
		code = http.StatusUnauthorized
	case models.ForbiddenErr:
		code = http.StatusForbidden
	case models.NotFoundErr:
		code = http.StatusNotFound
	case models.NotAllowedErr:
		code = http.StatusMethodNotAllowed
	case models.TooManyRequestErr:
		code = http.StatusTooManyRequests
	case models.ConflictErr:
		code = http.StatusConflict
	case error:
		code = http.StatusInternalServerError
		err = models.BaseError{Code: e.Error()}
	default:
		code = http.StatusInternalServerError
		err = models.InternalError
	}

	c.JSON(code, err)
}

func RecoveryHandle(c *gin.Context, err interface{}) {
	met.Counter("panic", 1, []metrics.Tag{})
}
