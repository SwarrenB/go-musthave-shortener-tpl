package main

import (
	"log"
	"os"

	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/config"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/service"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func run() error {
	appConfig := config.CreateGeneralConfig()
	service := service.CreateShortenerService()
	router := gin.Default()
	handler := handlers.CreateGinHandler(service, *appConfig)
	router.GET("/:id", handler.GinGetRequestHandler())
	router.POST("/", handler.GinPostRequestHandler())
	router.Run(appConfig.ServerAddress)
	return nil
}
