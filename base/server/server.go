package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/chihkaiyu/task-todo-api/base/goroutine"
)

const (
	shutdownTimeout = 10 * time.Second
)

func Serve(addr string, router *gin.Engine) error {
	srv := http.Server{
		Addr:    addr,
		Handler: router,
	}

	router.GET("/metrics", prometheusHandler())

	srvCh := make(chan error, 1)
	goroutine.Go(func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			srvCh <- err
		}
		srvCh <- nil
	})

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-srvCh:
		return err
	case <-shutdownCh:
		timeoutCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := srv.Shutdown(timeoutCtx); err != nil {
			return err
		}
		return nil
	}
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
