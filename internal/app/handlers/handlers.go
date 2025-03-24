package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/config"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/marshal"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/repository"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
)

type Handler struct {
	service  service.ServiceImpl
	config   config.Config
	logger   zap.Logger
	database *repository.SQLDatabase
}

func CreateGinHandler(service service.ServiceImpl, config config.Config, logger zap.Logger, database *repository.SQLDatabase) *Handler {
	return &Handler{
		service:  service,
		config:   config,
		logger:   logger,
		database: database,
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
				c.String(http.StatusConflict, handler.config.ShortURL+shortURL)
			} else {
				c.Writer.Header().Set("Content-Type", "text/plain; charset=UTF-8")
				c.String(http.StatusCreated, handler.config.ShortURL+shortURL)
			}
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
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			return
		}

		c.Header("Content-Type", "application/json")

		shortURL, err := handler.service.AddingURL(urlRequest.OriginalURL)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"result": handler.config.ShortURL + shortURL})
			return
		}
		c.Status(http.StatusCreated)

		urlResponse := marshal.URLResponse{ShortURL: handler.config.ShortURL + shortURL}

		if _, err := easyjson.MarshalToWriter(&urlResponse, c.Writer); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}
}

func (handler *Handler) HandlePingDB(database *repository.SQLDatabase) gin.HandlerFunc {
	return func(c *gin.Context) {

		defer database.Close()

		if database == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Server doesn't use database"})
			return
		}

		if err := database.Ping(); err != nil {
			http.Error(c.Writer, "database is not reachable", http.StatusInternalServerError)
			return
		}

		c.AbortWithStatus(http.StatusOK)
	}
}

type URLRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type URLResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func (handler *Handler) URLCreatorBatch(c *gin.Context) {
	defer c.Request.Body.Close()

	var requestURLs []URLRequest
	c.Header("Content-Type", "application/json")

	if err := json.NewDecoder(c.Request.Body).Decode(&requestURLs); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	responseURLs := make([]URLResponse, len(requestURLs))

	for i, requestURL := range requestURLs {
		shortURL, err := handler.service.AddingURL(requestURL.OriginalURL)
		if err != nil {
			c.String(http.StatusConflict, handler.config.ShortURL+shortURL)
			return
		}

		responseURLs[i] = URLResponse{
			CorrelationID: requestURL.CorrelationID,
			ShortURL:      handler.config.ShortURL + shortURL,
		}
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, responseURLs)
}
