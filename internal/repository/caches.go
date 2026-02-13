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
	
	ownerSessionTTL, err := time.ParseDuration(env.GetString("OWNER_SESSION_TTL", ""))
	if err != nil {
		return nil, err
	}

	cs := &CacheStore{
		OwnerSessions:  cache.New(ownerSessionTTL, 10*time.Minute),
		UserSessions:   cache.New(5*time.Minute, 10*time.Minute),
		ApotekSessions: cache.New(5*time.Minute, 10*time.Minute),
	}

	return cs, nil
}