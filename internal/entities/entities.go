package entities

import "time"

type ShortenURLRequest struct {
	LongURL string
	Domain  string
}

type ShortenURLResponse struct {
	ShortURl   string
	CreatedAt  time.Time
	ExpiryDate time.Time
	Error      string
}

type ShortURLDBData struct {
	LongURL    string
	Domain     string
	ShortURl   string
	CreatedAt  time.Time
	ExpiryDate time.Time
}

type RedirectShortURLRequest struct {
	ShortURl string
}

type RedirectShortURLResponse struct {
	LongURl string
}
