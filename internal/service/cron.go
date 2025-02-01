package service

import (
	"context"
	"time"
)

func (u *URLShortenService) PopulateTopDomains() {
	for {
		ctx := context.Background()
		u.db.PopulateTop3Domain(ctx)
		time.Sleep(2 * time.Second)
	}
}
