package main

import (
	"net/http"
	"rate-limiter/fixedWindow"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", fixedWindow.Middleware(http.HandlerFunc(handler)))
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
		return
	}
}
