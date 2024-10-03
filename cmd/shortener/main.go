package main

import (
	"io"
	"net/http"

	"math/rand/v2"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var vocabulary = make(map[string]string)

func GenerateURL(n int) []byte {
	b := "http://localhost:8080/"
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
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, postRequestHandler)
	mux.HandleFunc(`/{id}`, getRequestHandler)
	return http.ListenAndServe(`:8080`, mux)
}

func checkCorrectRequest(w http.ResponseWriter, r *http.Request) {
	correctRequest := r.Method != http.MethodPost && r.Method != http.MethodGet && r.Header.Get("Content-Type") != "text/plain"
	if correctRequest {
		http.Error(w, "This request is not allowed.", http.StatusBadRequest)
	}
}

func postRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		checkCorrectRequest(w, r)
		body, _ := io.ReadAll(r.Body)
		isURL := string(body[0:4])
		if isURL != "http" {
			http.Error(w, "URL is invalid.", http.StatusBadRequest)
			return
		} else {
			result := GenerateURL(rand.IntN(int(len(body))))
			vocabulary[string(`/{result}`)] = string(body)
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(http.StatusCreated)
			w.Write(result)
			return
		}
	}
}

func getRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		checkCorrectRequest(w, r)
		id := r.URL.Path
		key, ok := vocabulary[id]
		if ok {
			w.Header().Add("Content-Type", "text/plain")
			w.Header().Add("Location", key)
			w.WriteHeader(307)
		} else {
			http.Error(w, "This request is not allowed.", http.StatusBadRequest)
		}
		return
	}
}
