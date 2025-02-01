package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"urlshortener/internal/database"
	"urlshortener/internal/entities"
	"urlshortener/internal/service"
)

type App struct {
	service *service.URLShortenService
}

func NewApp() *App {
	db := database.NewInMemoryDatabase()
	s := service.NewURLShortenService(*db)
	app := App{
		service: s,
	}
	return &app
}

func (a *App) RedirectHandler() http.HandlerFunc {
	f := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			id := request.URL.Path[1:]
			if id == "" {
				writer.WriteHeader(http.StatusBadRequest)
			}
			ctx := context.Background()
			resp := a.service.RedirectURL(ctx, id)
			if resp != nil {
				if resp.Domain == "" || request.URL.Host == resp.Domain {
					http.RedirectHandler(resp.LongURl, http.StatusPermanentRedirect)
					return
				}
				writer.WriteHeader(http.StatusBadRequest)
			}
		default:
			writer.WriteHeader(http.StatusBadRequest)
		}
	})
	return f
}

func (a *App) GenerateShortURL() http.HandlerFunc {
	f := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodPost:
			req := entities.ShortenURLRequest{}
			data, err := io.ReadAll(request.Body)
			if err != nil {
				writer.WriteHeader(http.StatusBadRequest)
				writer.Write([]byte(err.Error()))
				return
			}
			err = json.Unmarshal(data, &req)
			if err != nil {
				writer.WriteHeader(http.StatusBadRequest)
				writer.Write([]byte(err.Error()))
				return
			}
			ctx := context.Background()
			resp, err := a.service.ShortenURL(ctx, req)
			repData, err := json.Marshal(resp)
			_, err = writer.Write(repData)
			//writer.WriteHeader(http.StatusCreated)
			if err != nil {
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
		default:
			writer.WriteHeader(http.StatusBadRequest)
		}
	})
	return f
}
