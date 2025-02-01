package database

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"urlshortener/internal/entities"
)

func TestInMemoryDatabase_AddData(t *testing.T) {
	app := NewInMemoryDatabase()
	ctx := context.Background()
	data := entities.ShortURLDBData{
		ShortURl: "https://localhost:8080/test",
		LongURL:  "https://longURL.com/test",
	}
	app.AddData(ctx, "https://longURL.com/test", data)
	assert.Len(t, app.longUrlDB, 1)
	assert.Len(t, app.repeatUrlDB, 1)
	count := 0
	app.metricsDB.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	assert.Equal(t, count, 1)
	countOfArr := len(*app.topDomains)
	app.PopulateTop3Domain(ctx)
	assert.Equal(t, countOfArr, 3)
}

func TestInMemoryDatabase_CheckDuplicateRequest(t *testing.T) {
	app := NewInMemoryDatabase()
	ctx := context.Background()
	data := entities.ShortURLDBData{
		ShortURl: "https://localhost:8080/test",
		LongURL:  "https://longURL.com/test",
	}
	app.AddData(ctx, "https://longURL.com/test", data)
	assert.Len(t, app.longUrlDB, 1)

	data = entities.ShortURLDBData{
		ShortURl: "https://localhost:8080/test",
		LongURL:  "https://longURL.com/test",
	}
	app.AddData(ctx, "https://longURL.com/test", data)
	assert.Len(t, app.longUrlDB, 1)
}
