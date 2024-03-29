//	@title			Task Todo API
//	@version		0.0.1
//	@description	Task Todo API server

//	@contact.name	Chih Kai Yu
//	@contact.email	kai.chihkaiyu@gmail.com

//	@host		localhost:8080
//	@BasePath	/api

//go:generate swag init -g ./cmd/api/main.go -o ./cmd/api/docs
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	bconfig "github.com/chihkaiyu/task-todo-api/base/config"
	"github.com/chihkaiyu/task-todo-api/base/server"
	"github.com/chihkaiyu/task-todo-api/cmd/api/api"
	"github.com/chihkaiyu/task-todo-api/cmd/api/config"
	"github.com/chihkaiyu/task-todo-api/middlewares"
	"github.com/chihkaiyu/task-todo-api/services/postgres"
	"github.com/chihkaiyu/task-todo-api/stores/tasks"

	_ "github.com/chihkaiyu/task-todo-api/cmd/api/docs"
)

func initLogger() zerolog.Logger {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.000000Z07:00"

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return logger
}

func main() {
	rootLogger := initLogger()
	rootCtx := rootLogger.WithContext(context.Background())

	cfg := config.Config{}
	if err := bconfig.Parse(&cfg); err != nil {
		rootLogger.Fatal().Err(err).Msg("bconfig.Parse failed")
	}

	dbPG, err := postgres.New(cfg.PostgresURI)
	if err != nil {
		rootLogger.Fatal().Msg("postgres.New failed")
	}

	// stores
	taskStore := tasks.New(dbPG)

	router := gin.New()
	router.Use(
		gin.CustomRecovery(middlewares.RecoveryHandle),
		middlewares.Cors(cfg.Env),
		requestid.New(),
		middlewares.Logger(rootCtx),
		middlewares.Stat(),
	)
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "healthy",
		})
	})
	if cfg.Debug {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	rg := router.Group("/")

	// routers
	api.NewTaskHandler(rg, taskStore)

	if err := server.Serve(fmt.Sprintf(":%s", cfg.Port), router); err != nil {
		rootLogger.Fatal().Err(err).Msg("server.Serve failed:")
	}
}
