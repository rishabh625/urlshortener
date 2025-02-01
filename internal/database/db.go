package database

import (
	"context"
	"errors"
	"strings"
	"urlshortener/internal/entities"
)

type InMemoryDatabase struct {
	shortUrlDB  map[string]entities.ShortURLDBData
	repeatUrlDB map[string]bool
}

func NewInMemoryDatabase() *InMemoryDatabase {
	shortUrlDB := make(map[string]entities.ShortURLDBData)
	repeatUrlDB := make(map[string]bool)
	return &InMemoryDatabase{
		shortUrlDB:  shortUrlDB,
		repeatUrlDB: repeatUrlDB,
	}
}

type DB interface {
	AddData(ctx context.Context, key string, data entities.ShortURLDBData) error
	CheckDuplicateRequest(ctx context.Context, key string) error
	RetrieveData(ctx context.Context, data string) *entities.ShortURLDBData
}

func (db *InMemoryDatabase) AddData(ctx context.Context, key string, data entities.ShortURLDBData) error {
	if _, ok := db.shortUrlDB[key]; ok {
		return errors.New("URL is already shortened")
	}
	db.shortUrlDB[key] = data
	db.repeatUrlDB[strings.Trim(data.ShortURl, "https://")] = true
	return nil
}

func (db *InMemoryDatabase) CheckDuplicateRequest(ctx context.Context, key string) error {
	if _, ok := db.repeatUrlDB[key]; ok {
		errors.New("Duplicate Request")
	}
	return nil
}

func (db *InMemoryDatabase) RetrieveData(ctx context.Context, data string) *entities.ShortURLDBData {
	if result, ok := db.shortUrlDB[data]; ok {
		return &result
	}
	return nil
}
