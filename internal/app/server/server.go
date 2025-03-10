package server

import (
	"net/http"

	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/config"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/middleware"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/repository"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/service"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/urlgenerate"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/utils"
	"github.com/gin-gonic/gin"
	compress "github.com/lf4096/gin-compress"
	"go.uber.org/zap"
)

type Server struct {
	*http.Server
	config  *config.Config
	repo    repository.URLRepository
	manager *repository.StateManager
	log     zap.Logger
}

func CreateServer(
	config *config.Config,
	repo repository.URLRepository,
	manager *repository.StateManager,
	log zap.Logger,
	database *utils.SQLDatabase,
) *Server {
	generator := urlgenerate.CreateURLGenerator()
	service := service.CreateShortenerService(repo, generator, config)

	router := gin.Default()
	handler := handlers.CreateGinHandler(service, *config, log, database)
	router.Use(middleware.WithLogging(log))
	router.Use(compress.Compress())
	router.Use(middleware.Decompress())
	router.GET("/ping", handler.HandlePingDB(database))
	router.GET("/:id", handler.GinGetRequestHandler())
	router.POST("/api/shorten", handler.HandlePostJSON())
	router.POST("/", handler.GinPostRequestHandler())

	server := http.Server{
		Addr:    config.ServerAddress,
		Handler: router,
	}
	return &Server{
		Server:  &server,
		config:  config,
		repo:    repo,
		manager: manager,
		log:     log,
	}
}
