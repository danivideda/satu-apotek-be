package repository

import (
	"time"

	"github.com/danivideda/satu-apotek-be/internal/env"
	"github.com/patrickmn/go-cache"
)

type CacheStore struct {
	OwnerSessions  *cache.Cache
	UserSessions   *cache.Cache
	ApotekSessions *cache.Cache
}

func NewCacheStore() (*CacheStore, error) {

	cacheTTL, err := time.ParseDuration(env.GetString("CACHE_SESSION_TTL", "5m"))
	if err != nil {
		return nil, err
	}

	cs := &CacheStore{
		OwnerSessions:  cache.New(cacheTTL, 2*cacheTTL),
		UserSessions:   cache.New(cacheTTL, 2*cacheTTL),
		ApotekSessions: cache.New(cacheTTL, 2*cacheTTL),
	}

	return cs, nil
}
