package middlewares

import (
	"github.com/gin-gonic/gin"

	"github.com/chihkaiyu/task-todo-api/services/metrics"
)

func Stat() gin.HandlerFunc {
	return func(c *gin.Context) {
		ender := met.Time("response_time_seconds", []metrics.Tag{
			{
				Name:  "method",
				Value: c.Request.Method,
			},
			{
				Name:  "path",
				Value: c.FullPath(),
			},
		})
		c.Next()
		ender.End()
	}
}
