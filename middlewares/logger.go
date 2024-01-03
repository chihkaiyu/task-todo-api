package middlewares

import (
	"context"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func Logger(rootCtx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// NOTE: ignore healthy check log
		if path == "/" {
			c.Next()
			return
		}

		rid := requestid.Get(c)
		logger := zerolog.Ctx(rootCtx).With().
			Str("requestID", rid).
			Str("path", c.Request.URL.Path). // NOTE: don't use c.FullPath(), we need parameter in path
			Str("method", c.Request.Method).
			Logger()
		c.Request = c.Request.WithContext(logger.WithContext(rootCtx))

		c.Next()

		var l *zerolog.Event
		if c.Writer.Status() >= 300 {
			l = zerolog.Ctx(c.Request.Context()).Error()
		} else {
			l = zerolog.Ctx(c.Request.Context()).Info()
		}
		defer l.Send()

		l = l.Int("statusCode", c.Writer.Status()).
			Str("clientIP", c.ClientIP())
	}
}
