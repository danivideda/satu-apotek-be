package repository

import (
	"time"

	"github.com/danivideda/satu-apotek-be/internal/env"
	"github.com/patrickmn/go-cache"
)

type CacheStore struct {
	OwnerSessions    *cache.Cache
	UserSessions     *cache.Cache
	PharmacySessions *cache.Cache
}

type PharmacySessionCacheValue struct {
	PharmacyID int64
	Name       string
	Users      []UserCache
}

type UserCache struct {
	ID       int64
	Username string
}

func NewCacheStore() (*CacheStore, error) {

	cacheTTL, err := time.ParseDuration(env.GetString("CACHE_SESSION_TTL", "5m"))
	if err != nil {
		return nil, err
	}

	cs := &CacheStore{
		OwnerSessions:    cache.New(cacheTTL, 2*cacheTTL),
		UserSessions:     cache.New(cacheTTL, 2*cacheTTL),
		PharmacySessions: cache.New(cacheTTL, 2*cacheTTL),
	}

	return cs, nil
}
