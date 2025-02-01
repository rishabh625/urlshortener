package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/deatil/go-encoding/encoding"
	"math/rand"
	"net/url"
	"regexp"
	"time"
	"urlshortener/internal/database"
	"urlshortener/internal/entities"
)

type URLShortenService struct {
	db database.InMemoryDatabase
}

const MAXRANDOM = 10000000

var re = regexp.MustCompile("^https?:\\/\\/[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}(\\/.*)?$")
var FixDomain string

func NewURLShortenService(db database.InMemoryDatabase, domain string) *URLShortenService {
	FixDomain = domain
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
	if URL.Scheme != "http" && URL.Scheme != "https" {
		URL.Scheme = "https"
	}
	if URL.Host == "" {
		return nil, errors.New("Invalid URL")
	}
	if !re.MatchString(URL.String()) {
		return nil, errors.New("Invalid URL")
	}
	res := u.db.RetrieveDuplicateURL(ctx, URL.String())
	if res == "" {
		hash := u.GenerateHashOfURL(ctx, URL.String())
		var shortURL string
		if request.Domain == "" {
			shortURL = fmt.Sprintf("%s/%s", FixDomain, hash)
		} else {
			shortURL = fmt.Sprintf("%s/%s", request.Domain, hash)
		}
		response := entities.ShortenURLResponse{
			ShortURl:   shortURL,
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
		result.ShortURl = fmt.Sprintf("https://%s/%s", FixDomain, result.ShortURl)
	} else {
		result.ShortURl = fmt.Sprintf("https://%s/%s", result.Domain, result.ShortURl)
	}
	return &entities.ShortenURLResponse{
		ShortURl:   result.ShortURl,
		CreatedAt:  result.CreatedAt,
		ExpiryDate: result.ExpiryDate,
	}, nil
}

// GenerateHashOfURL : hashes long URL with https schema , for hashing it uses sha256, once
// sha256 sum is obtained first 6 character of sha256 is converted into hexa decimal
// approximately it returns 6*2 characters in hexa decimal format , this hexa value is base62 encoded.
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

func (u *URLShortenService) RetrieveTop3Domains(ctx context.Context) []entities.TopDomains {
	return u.db.RetrieveTop3Domain(ctx)
}
