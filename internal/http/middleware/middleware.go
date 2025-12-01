package middleware

import "github.com/patrickmn/go-cache"

type AppMiddleware struct {
	cache *cache.Cache
}

func New(c *cache.Cache) AppMiddleware {
	return AppMiddleware{
		cache: c,
	}
}