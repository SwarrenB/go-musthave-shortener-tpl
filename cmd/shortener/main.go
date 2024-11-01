package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/config"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/logger"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/middleware"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/repository"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/service"
	"github.com/gin-gonic/gin"
	compress "github.com/lf4096/gin-compress"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	logger.Initialize("Info")
	appConfig := config.CreateGeneralConfig()
	service := service.CreateShortenerService()
	repo := repository.CreateInMemoryURLRepository()
	stateManager := repository.CreateStateManager(appConfig, *logger.Log)
	router := gin.Default()
	handler := handlers.CreateGinHandler(service, *appConfig)
	router.Use(middleware.WithLogging)
	router.Use(compress.Compress())
	router.Use(middleware.Decompress())
	router.GET("/:id", handler.GinGetRequestHandler())
	router.POST("/api/shorten", handler.HandlePostJSON())
	router.POST("/", handler.GinPostRequestHandler())

	server := http.Server{
		Addr:    appConfig.ServerAddress,
		Handler: router.Handler(),
	}

	repoState, err := stateManager.LoadFromFile()
	if err != nil {
		logger.Log.Error("create repository state error")
	} else {
		err := repo.RestoreURLRepository(repoState)
		if err != nil {
			logger.Log.Error("restore repository state error")
		}
	}

	for k, v := range repoState.GetURLRepositoryState() {
		logger.Log.Info("state to load", zap.String("shortUrl", k), zap.String("origUrl", v))
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// router.Run(appConfig.ServerAddress)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	logger.Log.Info("Server started")
	logger.Log.Info("Repo after server started", zap.Any("repo", repo))

	<-stopChan
	logger.Log.Info("Shutdown signal received")

	repoState, err = repo.CreateURLRepository()
	if err != nil {
		logger.Log.Error("create repository state error")
	} else {
		if err := stateManager.SaveToFile(repoState); err != nil {
			logger.Log.Error("failed to store state to file", zap.String("file_storage_path", appConfig.FileStoragePath))
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Log.Error("Server forced to shut down")
	} else {
		logger.Log.Info("Server shut down gracefully")
	}

	return nil
}
