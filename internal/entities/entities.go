package entities

import "time"

type ShortenURLRequest struct {
	LongURL string `json:"longURL"`
	Domain  string `json:"domain"`
}

type ShortenURLResponse struct {
	ShortURl   string `json:"shortURL"`
	CreatedAt  time.Time
	ExpiryDate time.Time
}

type ShortURLDBData struct {
	LongURL       string
	Domain        string
	LongURLDomain string
	ShortURl      string
	CreatedAt     time.Time
	ExpiryDate    time.Time
}

type RedirectShortURLResponse struct {
	LongURl string
	Domain  string
}

type TopDomains struct {
	Domain string
	Count  int
}
