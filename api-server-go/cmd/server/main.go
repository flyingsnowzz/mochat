package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mochat-api-server/internal/config"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/logger"
	"mochat-api-server/internal/pkg/storage"
	internalRedis "mochat-api-server/internal/redis"
	"mochat-api-server/internal/router"
	"mochat-api-server/internal/task"
)

var Version = "1.0.0"

func main() {
	fmt.Printf("MoChat API Server v%s starting...\n", Version)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := logger.Init(cfg.Log.Level, cfg.Log.Output); err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}

	if err := model.InitDB(cfg.DB); err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	defer model.CloseDB()

	if err := internalRedis.InitRedis(cfg.Redis); err != nil {
		log.Fatalf("Failed to connect redis: %v", err)
	}
	defer internalRedis.CloseRedis()

	storage.InitStorage(cfg.File, cfg.WeChat.APIBaseURL)

	task.InitScheduler()
	task.RegisterAllTasks(model.DB)
	task.StartScheduler()
	defer task.StopScheduler()

	r := router.SetupRouter(cfg, model.DB)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	go func() {
		fmt.Printf("Server listening on port %d\n", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Sugar.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Sugar.Errorf("Server forced to shutdown: %v", err)
	}

	logger.Sugar.Info("Server exited")
}
