package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_postRequestHandler(t *testing.T) {
	appConfig := config.CreateDefaultConfig()
	type args struct {
		code        int
		contentType string
	}
	tests := []struct {
		name       string
		body       string
		vocabulary map[string]string
		args       args
	}{
		{
			name:       "Normal case #1",
			body:       "https://yandex.practicum.ru",
			vocabulary: appConfig.Vocabulary,
			args: args{
				code:        http.StatusCreated,
				contentType: "text/plain; charset=UTF-8",
			},
		},
		{
			name:       "Error case #1",
			body:       "yandex.practicum",
			vocabulary: appConfig.Vocabulary,
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
			handler := ginPostRequestHandler(appConfig)
			handler(c)
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

var testUrls = []string{
	"http://practictum.yandex.ru",
	"http://vk.com",
	"https://yandex.ru",
	"https://github.com",
}

func fillVocabulary(vocabulary map[string]string) {
	for _, url := range testUrls {
		vocabulary[GenerateURL(len(url))] = url
	}
}

func Test_getRequestHandler(t *testing.T) {
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
	fillVocabulary(appConfig.Vocabulary)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for key, value := range appConfig.Vocabulary {
				request := httptest.NewRequest(http.MethodGet, key, nil)
				// создаём новый Recorder
				w := httptest.NewRecorder()
				// getRequestHandler(w, request)
				c, _ := gin.CreateTestContext(w)
				c.Request = request
				handler := ginGetRequestHandler(appConfig)
				handler(c)
				res := w.Result()
				// проверяем код ответа
				assert.Equal(t, test.args.code, res.StatusCode)

				defer res.Body.Close()

				assert.Equal(t, value, res.Header.Get("Location"))
				assert.Equal(t, strings.ToLower(test.args.contentType), strings.ToLower(res.Header.Get("Content-Type")))
			}
		})
	}
}
