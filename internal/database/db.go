package database

import (
	"context"
	"errors"
	"sort"
	"sync"
	"urlshortener/internal/entities"
)

type InMemoryDatabase struct {
	shortUrlDB  map[string]entities.ShortURLDBData // Retrieval DB
	metricsDB   sync.Map                           // Retrieval DB
	longUrlDB   map[string]string                  // Duplicate Request DB
	repeatUrlDB map[string]bool                    // collision db
	topDomains  []entities.TopDomains
	mu          sync.RWMutex
}

func NewInMemoryDatabase() *InMemoryDatabase {
	shortUrlDB := make(map[string]entities.ShortURLDBData)
	repeatUrlDB := make(map[string]bool)
	longUrlDB := make(map[string]string)
	topDomains := make([]entities.TopDomains, 3)
	return &InMemoryDatabase{
		shortUrlDB:  shortUrlDB,
		repeatUrlDB: repeatUrlDB,
		longUrlDB:   longUrlDB,
		topDomains:  topDomains,
	}
}

type DB interface {
	AddData(ctx context.Context, key string, data entities.ShortURLDBData) error
	CheckDuplicateRequest(ctx context.Context, key string) error
	RetrieveData(ctx context.Context, data string) *entities.ShortURLDBData
	RetrieveDuplicateURL(ctx context.Context, data string) string
	RetrieveTop3Domain(ctx context.Context) []entities.TopDomains
	PopulateTop3Domain(ctx context.Context)
}

func (db *InMemoryDatabase) AddData(ctx context.Context, key string, data entities.ShortURLDBData) error {
	db.mu.Lock()
	if _, ok := db.shortUrlDB[key]; ok {
		return errors.New("URL is already shortened")
	}
	db.shortUrlDB[key] = data
	db.longUrlDB[data.LongURL] = key
	value, ok := db.metricsDB.Load(data.LongURLDomain)
	if !ok {
		db.metricsDB.Store(data.LongURLDomain, 1)
	} else {
		valueRead := value.(int) + 1
		db.metricsDB.Store(data.LongURLDomain, valueRead)
	}
	db.repeatUrlDB[data.ShortURl] = true
	db.mu.Unlock()
	return nil
}

func (db *InMemoryDatabase) CheckDuplicateRequest(ctx context.Context, key string) error {
	if _, ok := db.repeatUrlDB[key]; ok {
		errors.New("Duplicate Request")
	}
	return nil
}

func (db *InMemoryDatabase) RetrieveData(ctx context.Context, data string) *entities.ShortURLDBData {
	db.mu.RLock()
	defer db.mu.RUnlock()
	if result, ok := db.shortUrlDB[data]; ok {
		return &result
	}
	return nil
}

func (db *InMemoryDatabase) RetrieveDuplicateURL(ctx context.Context, data string) string {
	db.mu.RLock()
	defer db.mu.RUnlock()
	if result, ok := db.longUrlDB[data]; ok {
		return result
	}
	return ""
}

func (db *InMemoryDatabase) PopulateTop3Domain(ctx context.Context) {
	var domains []entities.TopDomains
	db.metricsDB.Range(func(key, value any) bool {
		if count, ok := value.(int); ok {
			domains = append(domains, entities.TopDomains{key.(string), count})
		}
		return true
	})
	sort.Slice(domains, func(i, j int) bool {
		return domains[i].Count > domains[j].Count
	})
	db.mu.Lock()
	defer db.mu.Unlock()
	db.topDomains = make([]entities.TopDomains, 0, 3)
	for i := 0; i < len(domains) && i < 3; i++ {
		db.topDomains = append(db.topDomains, entities.TopDomains{
			Domain: domains[i].Domain,
			Count:  domains[i].Count,
		})
	}
}

func (db *InMemoryDatabase) RetrieveTop3Domain(ctx context.Context) []entities.TopDomains {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return append([]entities.TopDomains{}, db.topDomains...)
}
