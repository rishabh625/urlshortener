package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"urlshortener/internal/handler"
	"urlshortener/internal/middleware"
)

func main() {
	fmt.Println("Server to be started at 8080")
	app := handler.NewApp()
	mux := http.NewServeMux()
	mux.HandleFunc("/shortURL", app.GenerateShortURL())
	mux.HandleFunc("/{id}", app.RedirectHandler())
	srv := http.Server{
		Addr:    ":8080",
		Handler: middleware.LatencyMiddleware(mux),
	}
	log.Printf("Server")
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %s", err)
	}
}
