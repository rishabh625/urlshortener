package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/deatil/go-encoding/encoding"
	"math/rand"
	"net/url"
	"time"
	"urlshortener/internal/database"
	"urlshortener/internal/entities"
)

type URLShortenService struct {
	db database.InMemoryDatabase
}

const MAXRANDOM = 10000000
const FixDomain = "shorty.tk"

func NewURLShortenService(db database.InMemoryDatabase) *URLShortenService {
	return &URLShortenService{db: db}
}

type URLShortener interface {
	ShortenURL(ctx context.Context, request entities.ShortenURLRequest) (*entities.ShortenURLResponse, error)
	RedirectURL(ctx context.Context, data string) *entities.RedirectShortURLResponse
	GenerateHashOfURL(ctx context.Context, URL string) string
}

func (u *URLShortenService) ShortenURL(ctx context.Context, request entities.ShortenURLRequest) (*entities.ShortenURLResponse, error) {
	// check if it's a valid url
	URL, err := url.ParseRequestURI(request.LongURL)
	if err != nil {
		return nil, err
	}
	res := u.db.RetrieveDuplicateURL(ctx, URL.String())
	if res == "" {
		hash := u.GenerateHashOfURL(ctx, URL.String())
		var shorturl string
		if request.Domain == "" {
			shorturl = fmt.Sprintf("https://%s/%s", FixDomain, hash)
		} else {
			shorturl = fmt.Sprintf("https://%s/%s", request.Domain, hash)
		}
		response := entities.ShortenURLResponse{
			ShortURl:   shorturl,
			CreatedAt:  time.Now(),
			ExpiryDate: time.Now().AddDate(0, 0, 7),
		}
		domain := FixDomain
		if request.Domain != "" {
			domain = request.Domain
		}
		dbData := entities.ShortURLDBData{
			LongURL:       request.LongURL,
			Domain:        domain,
			LongURLDomain: URL.Host,
			ShortURl:      hash,
			CreatedAt:     response.CreatedAt,
			ExpiryDate:    response.ExpiryDate,
		}
		err := u.db.AddData(ctx, hash, dbData)
		if err != nil {
			return nil, err
		}
		return &response, nil
	}
	result := u.db.RetrieveData(ctx, res)
	if result.Domain == "" {
		result.ShortURl = fmt.Sprintf("https://%s/%d", FixDomain, result.ShortURl)
	} else {
		result.ShortURl = fmt.Sprintf("https://%s/%d", result.Domain, result.ShortURl)
	}
	return &entities.ShortenURLResponse{
		ShortURl:   result.ShortURl,
		CreatedAt:  result.CreatedAt,
		ExpiryDate: result.ExpiryDate,
	}, nil
}

func (u *URLShortenService) GenerateHashOfURL(ctx context.Context, URL string) string {
	seed := ""
	for {
		hash := sha256.Sum256([]byte(URL + seed))
		shortHash := hash[:6]
		dst := make([]byte, hex.EncodedLen(len(shortHash)))
		hex.Encode(dst, shortHash)
		key := fmt.Sprintf("%x", dst)
		encodedData := encoding.FromString(key).Base62Encode()
		err := u.db.CheckDuplicateRequest(ctx, encodedData.String())
		if err != nil {
			seed = fmt.Sprintf("%d", time.Now().UnixNano()) + fmt.Sprintf("%d", rand.Int63n(MAXRANDOM))
		} else {
			return encodedData.String()
		}
	}
}

func (u *URLShortenService) RedirectURL(ctx context.Context, hash string) *entities.RedirectShortURLResponse {
	resp := u.db.RetrieveData(ctx, hash)
	if resp != nil && time.Now().Before(resp.ExpiryDate) {
		return &entities.RedirectShortURLResponse{
			LongURl: resp.LongURL,
			Domain:  resp.Domain,
		}
	}
	return nil
}
