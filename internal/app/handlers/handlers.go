package handlers

import (
	"io"
	"net/http"

	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/config"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/marshal"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"
)

type GinHandler struct {
	service service.ServiceImpl
	config  config.Config
}

func CreateGinHandler(service service.ServiceImpl, config config.Config) *GinHandler {
	return &GinHandler{
		service: service,
		config:  config,
	}
}

func (handler *GinHandler) GinPostRequestHandler() gin.HandlerFunc {
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

func (handler *GinHandler) GinGetRequestHandler() gin.HandlerFunc {
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

func (handler *GinHandler) HandlePostJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		urlRequest := marshal.URLRequest{}

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
