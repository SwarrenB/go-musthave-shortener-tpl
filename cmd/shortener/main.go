package main

import (
	"log"

	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/config"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/logger"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/middleware"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/service"
	"github.com/gin-gonic/gin"
	compress "github.com/lf4096/gin-compress"
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
	router := gin.Default()
	handler := handlers.CreateGinHandler(service, *appConfig)
	router.Use(middleware.WithLogging)
	router.Use(compress.Compress())
	router.Use(middleware.Decompress())
	router.GET("/:id", handler.GinGetRequestHandler())
	router.POST("/api/shorten", handler.HandlePostJSON())
	router.POST("/", handler.GinPostRequestHandler())
	router.Run(appConfig.ServerAddress)
	return nil
}
