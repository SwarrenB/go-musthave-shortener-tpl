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
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/logger"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/repository"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/server"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	logger := logger.CreateLogger("Info").GetLogger()
	appConfig := config.CreateGeneralConfig()
	repo := repository.CreateInMemoryURLRepository()
	stateManager := repository.CreateStateManager(appConfig, logger)
	server := server.CreateServer(appConfig, repo, stateManager, logger)

	repoState, err := stateManager.LoadFromFile()
	if err != nil {
		logger.Error("create repository state error")
	} else {
		err := repo.RestoreURLRepository(repoState)
		if err != nil {
			logger.Error("restore repository state error")
		}
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// router.Run(appConfig.ServerAddress)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	logger.Info("Server started")

	<-stopChan
	logger.Info("Shutdown signal received")

	repoState, err = repo.CreateURLRepository()
	if err != nil {
		logger.Error("create repository state error")
	} else {
		if err := stateManager.SaveToFile(repoState); err != nil {
			logger.Error("failed to store state to file", zap.String("file_storage_path", appConfig.FileStoragePath))
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shut down")
	} else {
		logger.Info("Server shut down gracefully")
	}

	return nil
}
