package main

import (
	"io"
	"math/rand/v2"
	"net/http"
	"strings"

	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/config"
	"github.com/gin-gonic/gin"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var vocabulary = make(map[string]string)

func GenerateURL(n int, defaultURL string) []byte {
	b := defaultURL
	for i := 0; i < n; i++ {
		b += string(letters[rand.IntN(len(letters))])
	}
	return []byte(b)
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	appConfig := config.CreateConfigWithFlags()
	router := gin.Default()
	router.GET("/{id}", ginGetRequestHandler(appConfig))
	router.POST("/", ginPostRequestHandler(appConfig))
	// mux := http.NewServeMux()
	// mux.HandleFunc(`/`, postRequestHandler)
	// mux.HandleFunc(`/:id`, getRequestHandler)
	// return http.ListenAndServe(`:8080`, mux)
	router.Run(":8080")
	return nil
}

// func checkCorrectRequest(w http.ResponseWriter, r *http.Request) {
// 	correctRequest := r.Method != http.MethodPost && r.Method != http.MethodGet
// 	if correctRequest {
// 		http.Error(w, "This request is not allowed.", http.StatusBadRequest)
// 	}
// }

func ginPostRequestHandler(appConfig *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		isURL := string(body[0:4]) == "http"
		if err != nil || !isURL {
			c.String(http.StatusBadRequest, "URL is invalid.")
			return
		} else {
			result := GenerateURL(rand.IntN(int(len(body))), appConfig.GetConfigURL())
			vocabulary[string(result)] = string(body)
			c.Writer.Header().Set("Content-Type", "text/plain; charset=UTF-8")
			c.String(http.StatusCreated, string(result))
			return
		}
	}
}

// func postRequestHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodPost {
// 		checkCorrectRequest(w, r)
// 		body, err := io.ReadAll(r.Body)
// 		isURL := string(body[0:4])
// 		if err != nil || isURL != "http" {
// 			http.Error(w, "URL is invalid.", http.StatusBadRequest)
// 			return
// 		} else {
// 			result := GenerateURL(rand.IntN(int(len(body))))
// 			vocabulary[string(result)] = string(body)
// 			w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
// 			w.WriteHeader(http.StatusCreated)
// 			w.Write(result)
// 			return
// 		}
// 	}
// }

func ginGetRequestHandler(appConfig *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Request.URL.Path
		val, ok := vocabulary[appConfig.GetConfigURL()+strings.TrimPrefix(id, "/")]
		if ok {
			c.Writer.Header().Set("Content-Type", "text/plain; charset=UTF-8")
			c.Writer.Header().Set("Location", val)
			c.AbortWithStatus(http.StatusTemporaryRedirect)
		} else {
			c.String(http.StatusBadRequest, "This URL does not exist in vocabulary.")
		}
	}
}

// func getRequestHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodGet {
// 		checkCorrectRequest(w, r)
// 		id := "http://localhost:8080" + r.URL.Path
// 		val, ok := vocabulary[id]
// 		if ok {
// 			w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
// 			w.Header().Set("Location", val)
// 			w.WriteHeader(307)
// 		} else {
// 			http.Error(w, "This URL does not exist in vocabulary.", http.StatusBadRequest)
// 		}
// 		return
// 	}
// }
