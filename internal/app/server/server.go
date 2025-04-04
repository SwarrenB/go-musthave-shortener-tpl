package server

import (
	"net/http"

	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/config"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/middleware"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/repository"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/service"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/urlgenerate"
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

func getRepository(config *config.Config, defaultRepo repository.URLRepository, database *repository.SQLDatabase, log zap.Logger) repository.URLRepository {
	if config.DatabaseDSN == "" {
		return defaultRepo
	}

	database.CreateTables(log)
	return database
}

func CreateServer(
	config *config.Config,
	repo repository.URLRepository,
	manager *repository.StateManager,
	log zap.Logger,
	database *repository.SQLDatabase,
) *Server {
	generator := urlgenerate.CreateURLGenerator()

	store := getRepository(config, repo, database, log)
	service := service.CreateShortenerService(store, generator, config)
	router := gin.Default()
	handler := handlers.CreateGinHandler(service, *config, log, database)
	router.Use(middleware.WithLogging(log))
	router.Use(compress.Compress())
	router.Use(middleware.Decompress())
	router.GET("/ping", handler.HandlePingDB(database))
	router.GET("/:id", handler.GinGetRequestHandler())
	router.POST("/api/shorten", handler.HandlePostJSON())
	router.POST("/", handler.GinPostRequestHandler())
	router.POST("/api/shorten/batch", handler.URLCreatorBatch)

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
