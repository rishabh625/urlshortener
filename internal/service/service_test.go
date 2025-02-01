package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
	"time"
	"urlshortener/internal/database"
	"urlshortener/internal/entities"
)

// Test Shortening of URL for valid url
func TestURLShortenService_ShortenURL_VALIDURL(t *testing.T) {
	db := database.NewInMemoryDatabase()
	app := NewURLShortenService(*db, "http://localhost:8080")
	ctx := context.Background()
	req := entities.ShortenURLRequest{
		LongURL: "https://www.reddit.com/r/Fedora/",
	}
	resp, err := app.ShortenURL(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

// Test Shortening of URL for in valid url
func TestURLShortenService_ShortenURL_INVALIDURL(t *testing.T) {
	db := database.NewInMemoryDatabase()
	app := NewURLShortenService(*db, "http://localhost:8080")
	ctx := context.Background()
	req := entities.ShortenURLRequest{
		LongURL: "https://invalid-url.",
	}
	_, err := app.ShortenURL(ctx, req)
	assert.Error(t, err)
}

// Test Shortening of URL for empty url
func TestURLShortenService_ShortenURL_BadRequest(t *testing.T) {
	db := database.NewInMemoryDatabase()
	app := NewURLShortenService(*db, "http://localhost:8080")
	ctx := context.Background()
	req := entities.ShortenURLRequest{
		LongURL: "",
	}
	_, err := app.ShortenURL(ctx, req)
	assert.Error(t, err)
}

// Test Shortening of URL if same url shortened twice
func TestURLShortenService_ShortenURL_IFSAMELongURLPassed(t *testing.T) {
	db := database.NewInMemoryDatabase()
	app := NewURLShortenService(*db, "http://localhost:8080")
	ctx := context.Background()
	req := entities.ShortenURLRequest{
		LongURL: "https://www.reddit.com/r/Fedora/",
	}
	resp, err := app.ShortenURL(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	resp2, err := app.ShortenURL(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp2)
	assert.Equal(t, resp.ShortURl, resp2.ShortURl)
	assert.Equal(t, resp.ExpiryDate, resp2.ExpiryDate)
}

// Test Shortening of URL for very long url
func TestURLShortenService_ShortenURL_LongURL(t *testing.T) {
	db := database.NewInMemoryDatabase()
	app := NewURLShortenService(*db, "http://localhost:8080")
	ctx := context.Background()
	req := entities.ShortenURLRequest{
		LongURL: "https://accounts.google.com/signin/oauth/consent?authuser=1&part=AJi8hAO8Y7aSI6oqSN5ZFgylTOD4-8IHTf--L_vfW1OSSofRtawOX1M9kDBKtErYlCQmc21gn5uS8Zn1Sxgv9-KAhHv5EKWSgsqjN094rNbw-W0JEs1Fa9k36h75095hQ2ApvCv1EIOxioCxU2VkEa_OxDjGTHHLoBtyO_ZoQiTejVkVXFvAfhq0qlwLC7LODsFpGdKulf7y5I-F-ou3bPh7cZWy0yxGVWGJsTsqBXTeoZDWTRViNhjFUV42ayZKT0l1Fh_-MK96Keig_CuJEhitIrGubdXNwofnwFTf7QH2yo7s4xyFy1lDwJzLdDLxXzOdvhNWCNZfS80aVg2pqcVp4ftmUm5bLFN4U8dB-2GNjsUx9SFFiU8PRPCcmvqT4Jq68HojYPsPIG-SW6a_-lPfXAsJBzLc39ammfnREiLrLm6aFio5zA8qfTuRlMjo5l_SSyo1D33AxPbnoFzi81ks-6hokEErhKCIuQUiSGmKCIk71TkN5aA&flowName=GeneralOAuthFlow&as=S-1540822420%3A1738341489393397&client_id=993576537952-o63tbj4issluoheejqdfan468foht25p.apps.googleusercontent.com&pli=1&rapt=AEjHL4PQqC0VtqcT3S19Kv8Ox5Y7I8_pMD_qdtK4S8BK0OcXG9K7GCzD1rircLh9JedQYt79xrSpWVLXsiLKp-eXFjr3UkNkdA#",
	}
	resp, err := app.ShortenURL(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

// Test Shortening of URL for URL containing chinese character
func TestURLShortenService_ShortenURL_NonEnglishCharacter(t *testing.T) {
	db := database.NewInMemoryDatabase()
	app := NewURLShortenService(*db, "http://localhost:8080")
	ctx := context.Background()
	req := entities.ShortenURLRequest{
		LongURL: "https://zh.wikipedia.org/wiki/%E7%99%BE%E5%BA%A6",
	}
	resp, err := app.ShortenURL(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

// Test Redirection of URL once shortened
func TestURLShortenService_RedirectURL(t *testing.T) {
	// Shorten a URL
	db := database.NewInMemoryDatabase()
	app := NewURLShortenService(*db, "http://localhost:8080")
	ctx := context.Background()
	req := entities.ShortenURLRequest{
		LongURL: "https://www.reddit.com/r/Fedora/",
	}
	resp, err := app.ShortenURL(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// Redirect the URL
	shortURL := strings.Split(resp.ShortURl, "/")
	response := app.RedirectURL(ctx, shortURL[len(shortURL)-1])
	assert.Equal(t, req.LongURL, response.LongURl)
}

// Test Redirection of URL once shortened
func TestURLShortenService_RedirectURLNonExistent(t *testing.T) {
	db := database.NewInMemoryDatabase()
	app := NewURLShortenService(*db, "http://localhost:8080")
	ctx := context.Background()
	// Redirect the URL
	shortURL := strings.Split("https://www.reddit.com/hdjdjknkdnkj", "/")
	response := app.RedirectURL(ctx, shortURL[len(shortURL)-1])
	assert.Nil(t, response)
}

// Test Redirection of URL once if it is not shortened by this application
func TestURLShortenService_RetrieveTop3Domain(t *testing.T) {
	// Redirect the URL
	// Redit - 2
	db := database.NewInMemoryDatabase()
	app := NewURLShortenService(*db, "http://localhost:8080")
	ctx := context.Background()
	req := entities.ShortenURLRequest{
		LongURL: "https://www.reddit.com/r/Fedora/",
	}
	resp, err := app.ShortenURL(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	req = entities.ShortenURLRequest{
		LongURL: "https://www.reddit.com/r/Fedora/link2",
	}
	resp, err = app.ShortenURL(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	// wikepedia 2
	resp, err = app.ShortenURL(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	req = entities.ShortenURLRequest{
		LongURL: "https://www.wikepedia.com/r/Fedora/link2",
	}
	resp, err = app.ShortenURL(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	resp, err = app.ShortenURL(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	req = entities.ShortenURLRequest{
		LongURL: "https://www.wikepedia.com/r/Fedora/link3",
	}
	resp, err = app.ShortenURL(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	resp, err = app.ShortenURL(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// amazon - 2
	req = entities.ShortenURLRequest{
		LongURL: "https://www.amazon.com/r/Fedora/",
	}
	resp, err = app.ShortenURL(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	resp, err = app.ShortenURL(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	req = entities.ShortenURLRequest{
		LongURL: "https://www.amazon.com/r/Fedora/link2",
	}
	resp, err = app.ShortenURL(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// baeldung 1
	req = entities.ShortenURLRequest{
		LongURL: "https://www.baeldung.com/r/Fedora/",
	}
	resp, err = app.ShortenURL(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// Run Cron to populate db
	go app.PopulateTopDomains()
	time.Sleep(2 * time.Second)
	dom := app.RetrieveTop3Domains(ctx)
	assert.Len(t, dom, 3)
	domains := []entities.TopDomains{
		{"www.amazon.com", 2},
		{"www.wikepedia.com", 2},
		{"www.reddit.com", 2},
	}

	// Loop through slice and check if domain matches any in domainsToCheck
	for _, domain := range domains {
		flag := false
		for _, check := range dom {
			if reflect.DeepEqual(domain, check) {
				flag = true
			}
		}
		assert.True(t, flag, "for "+domain.Domain+" having count :", domain.Count)
	}
}

func TestTop3Domain(t *testing.T) {
	// Redirect the URL
	db := database.NewInMemoryDatabase()
	app := NewURLShortenService(*db, "http://localhost:8080")
	ctx := context.Background()
	shortURL := strings.Split("https://www.reddit.com/r", "/")
	response := app.RedirectURL(ctx, shortURL[len(shortURL)-1])
	assert.Nil(t, response)
}
