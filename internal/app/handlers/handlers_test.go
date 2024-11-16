package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/config"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/logger"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/repository"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/service"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/urlgenerate"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ginPostRequestHandler(t *testing.T) {
	appConfig := config.CreateDefaultConfig()
	repo := repository.CreateInMemoryURLRepository()
	generator := urlgenerate.CreateURLGenerator()
	service := service.CreateShortenerService(repo, generator, appConfig)
	log := logger.CreateLogger("Info").GetLogger()
	handler := CreateGinHandler(service, *appConfig, log)
	type args struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		body string
		args args
	}{
		{
			name: "Normal case #1",
			body: "https://yandex.practicum.ru",
			args: args{
				code:        http.StatusCreated,
				contentType: "text/plain; charset=UTF-8",
			},
		},
		{
			name: "Error case #1",
			body: "yandex.practicum",
			args: args{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=UTF-8",
			},
		},
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.body))
			// создаём новый Recorder
			w := httptest.NewRecorder()
			// postRequestHandler(w, request)
			c, _ := gin.CreateTestContext(w)
			c.Request = request
			h := handler.GinPostRequestHandler()
			h(c)
			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.args.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, strings.ToLower(test.args.contentType), strings.ToLower(res.Header.Get("Content-Type")))
		})
	}
}

func Test_ginGetRequestHandler(t *testing.T) {
	type args struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Normal case #1",
			args: args{
				code:        http.StatusTemporaryRedirect,
				contentType: "text/plain; charset=UTF-8",
			},
		},
		// TODO: Add test cases.
	}
	appConfig := config.CreateDefaultConfig()
	repo := repository.CreateInMemoryURLRepository()
	generator := urlgenerate.CreateURLGenerator()
	log := logger.CreateLogger("Info").GetLogger()
	service := service.CreateShortenerService(repo, generator, appConfig)
	handler := CreateGinHandler(service, *appConfig, log)
	originalURL := "http://practictum.yandex.ru"
	shortURL, _ := handler.service.AddingURL(originalURL)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, shortURL, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			// getRequestHandler(w, request)
			c, _ := gin.CreateTestContext(w)
			c.Request = request
			h := handler.GinGetRequestHandler()
			h(c)
			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.args.code, res.StatusCode)

			defer res.Body.Close()

			assert.Equal(t, originalURL, res.Header.Get("Location"))
			assert.Equal(t, strings.ToLower(test.args.contentType), strings.ToLower(res.Header.Get("Content-Type")))
		})
	}
}

func TestGinHandler_HandlePostJSON(t *testing.T) {
	tests := []struct {
		name    string
		handler *Handler
		want    gin.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.handler.HandlePostJSON(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GinHandler.HandlePostJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
