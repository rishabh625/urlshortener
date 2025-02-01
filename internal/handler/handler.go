package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"urlshortener/internal/database"
	"urlshortener/internal/entities"
	"urlshortener/internal/service"
)

type App struct {
	service *service.URLShortenService
}

func NewApp(domain string) *App {
	db := database.NewInMemoryDatabase()
	s := service.NewURLShortenService(*db, domain)
	app := App{
		service: s,
	}
	go s.PopulateTopDomains()
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
			if resp == nil {
				writer.WriteHeader(http.StatusServiceUnavailable)
			}
			if resp != nil {
				fmt.Println("Redirecting")
				http.Redirect(writer, request, resp.LongURl, http.StatusPermanentRedirect)
				return
			}
			writer.WriteHeader(http.StatusBadRequest)
			return
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
				_, _ = writer.Write([]byte(err.Error()))
				return
			}
			err = json.Unmarshal(data, &req)
			if err != nil {
				writer.WriteHeader(http.StatusBadRequest)
				_, _ = writer.Write([]byte(err.Error()))
				return
			}
			ctx := context.Background()
			resp, err := a.service.ShortenURL(ctx, req)
			repData, err := json.Marshal(resp)
			_, err = writer.Write(repData)
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

func (a *App) Top3Domains() http.HandlerFunc {
	f := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			ctx := context.Background()
			res := a.service.RetrieveTop3Domains(ctx)
			data, err := json.Marshal(res)
			if err != nil {
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
			_, err = writer.Write(data)
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
