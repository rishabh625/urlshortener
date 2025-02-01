package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"urlshortener/internal/handler"
	"urlshortener/internal/middleware"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <fixDomain>")
	}
	fixDomain := "http://localhost:8080"
	if len(os.Args) > 1 {
		fixDomain = os.Args[1]
	}
	log.Printf("Server to be started at 8080")
	app := handler.NewApp(fixDomain)
	mux := http.NewServeMux()
	mux.HandleFunc("/shortURL", app.GenerateShortURL())
	mux.HandleFunc("/{id}", app.RedirectHandler())
	mux.HandleFunc("/metrics", app.Top3Domains())
	srv := http.Server{
		Addr:    ":8080",
		Handler: middleware.LatencyMiddleware(mux),
	}
	log.Printf("Server is serving on 8080" +
		"/shortURL - shorten the URL" +
		"/{id} - for redirection" +
		"/metrics -  for top3 Domains")

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %s", err)
	}
}
