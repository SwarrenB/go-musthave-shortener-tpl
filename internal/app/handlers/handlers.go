package handlers

import (
	"io"
	"net/http"

	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/config"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/marshal"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/service"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/utils"
	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
)

type Handler struct {
	service service.ServiceImpl
	config  config.Config
	logger  zap.Logger
}

func CreateGinHandler(service service.ServiceImpl, config config.Config, logger zap.Logger) *Handler {
	return &Handler{
		service: service,
		config:  config,
		logger:  logger,
	}
}

func (handler *Handler) GinPostRequestHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		isURL := string(body[0:4]) == "http"
		if err != nil || !isURL {
			c.String(http.StatusBadRequest, "URL is invalid.")
			return
		} else {
			shortURL, err := handler.service.AddingURL(string(body))
			if err != nil {
				c.String(http.StatusBadRequest, err.Error())
				return
			}
			c.Writer.Header().Set("Content-Type", "text/plain; charset=UTF-8")
			c.String(http.StatusCreated, handler.config.ShortURL+shortURL)
			return
		}
	}
}

func (handler *Handler) GinGetRequestHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Request.URL.Path
		val, ok := handler.service.GetOriginalURL(id)
		if ok == nil {
			c.Writer.Header().Set("Content-Type", "text/plain; charset=UTF-8")
			c.Writer.Header().Set("Location", val)
			c.AbortWithStatus(http.StatusTemporaryRedirect)
			return
		} else {
			c.String(http.StatusBadRequest, "This URL does not exist in vocabulary.")
			return
		}
	}
}

func (handler *Handler) HandlePostJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		urlRequest := marshal.URLRequest{OriginalURL: ""}

		if err := easyjson.UnmarshalFromReader(c.Request.Body, &urlRequest); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		shortURL, err := handler.service.AddingURL(urlRequest.OriginalURL)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Status(http.StatusCreated)

		urlResponse := marshal.URLResponse{ShortURL: handler.config.ShortURL + shortURL}

		if _, err = easyjson.MarshalToWriter(&urlResponse, c.Writer); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}
}

func (handler *Handler) HandlePingDB() gin.HandlerFunc {
	return func(c *gin.Context) {
		database, err := utils.NewPG(handler.config.DatabaseDSN, handler.logger)
		if err != nil {
			http.Error(c.Writer, "cannot open database", http.StatusInternalServerError)

			return
		}
		defer database.Close()

		if err := database.Ping(); err != nil {
			http.Error(c.Writer, "database is not reachable", http.StatusInternalServerError)

			return
		}

		c.AbortWithStatus(http.StatusOK)
	}
}
