package repository

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type CacheStore struct {
	OwnerSessions  *cache.Cache
	UserSessions   *cache.Cache
	ApotekSessions *cache.Cache
}

func NewCacheStore() (*CacheStore, error) {

	cs := &CacheStore{
		OwnerSessions:  cache.New(5*time.Minute, 10*time.Minute),
		UserSessions:   cache.New(5*time.Minute, 10*time.Minute),
		ApotekSessions: cache.New(5*time.Minute, 10*time.Minute),
	}

	return cs, nil
}