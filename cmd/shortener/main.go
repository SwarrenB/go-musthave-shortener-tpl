package main

import (
	"io"
	"net/http"
	"strings"

	"golang.org/x/exp/rand"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var vocabulary = make(map[string]string)

func GenerateURL(n int) []byte {
	b := "http://localhost:8080/"
	for i := 0; i < n; i++ {
		b += string(letters[rand.Intn(len(letters))])
	}
	return []byte(b)
}

// функция main вызывается автоматически при запуске приложения
func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

// функция run будет полезна при инициализации зависимостей сервера перед запуском
func run() error {
	return http.ListenAndServe(`:8080`, middleware(http.HandlerFunc(webhook)))
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		next.ServeHTTP(w, r)
	})
}

// функция webhook — обработчик HTTP-запроса
func webhook(w http.ResponseWriter, r *http.Request) {
	correctRequest := r.Method != http.MethodPost && r.Method != http.MethodGet && r.Header.Get("Content-Type") != "text/plain"
	if correctRequest {
		http.Error(w, "This request is not allowed.", http.StatusBadRequest)
		return
	} else if r.Method == http.MethodGet {
		id := r.Host[strings.LastIndex(r.Host, "/"):]
		key, val := vocabulary[id]
		if val {
			w.Header().Add("Location", key)
			w.WriteHeader(307)
		} else {
			http.Error(w, "This request is not allowed.", http.StatusBadRequest)
		}
		return
	} else {
		body, _ := io.ReadAll(r.Body)
		isURL := string(body[0:4])
		if isURL != "http" {
			http.Error(w, "URL is invalid.", http.StatusBadRequest)
			return
		} else {
			result := GenerateURL(5)
			vocabulary[string(body)] = string(result)
			w.Write(result)
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(http.StatusCreated)
		}
	}
}
