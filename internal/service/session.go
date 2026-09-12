package service

import (
	"context"
	"time"

	"github.com/danivideda/satu-apotek-be/internal/config"
	"github.com/danivideda/satu-apotek-be/internal/repository"
)

type Session struct {
	cacheStore          *repository.CacheStore
	ownerSessionRepo    repository.OwnerSessionsRepository
	userSessionRepo     repository.UserSessionsRepository
	pharmacySessionRepo repository.PharmacySessionsRepository
	ownerSessionTTL     time.Duration
	userSessionTTL      time.Duration
	pharmacySessionTTL  time.Duration
}

func NewSession(repo repository.Repository, authConfig config.AuthConfig) *Session {
	return &Session{
		repo.CacheStore,
		repo.OwnerSessions,
		repo.UserSessions,
		repo.PharmacySessions,
		authConfig.OwnerSessionTTL,
		authConfig.UserSessionTTL,
		authConfig.PharmacySessionTTL,
	}
}

func (s *Session) NewOwner(ctx context.Context, ownerID int64) (sessionID string, expiresAt time.Time, err error) {
	exp := time.Now().Add(s.ownerSessionTTL)
	ownerSession, err := s.ownerSessionRepo.Create(ctx, ownerID, exp)
	if err != nil {
		return "", time.Time{}, err
	}
	ownerSessionID := ownerSession.ID.String()
	s.cacheStore.OwnerSessions.SetDefault(ownerSessionID, ownerID)
	return ownerSessionID, ownerSession.ExpiresAt.Time, nil
}

// ResolveOwner looks up an owner session for middleware.
// Cache hit returns the cached owner ID without rotating the session
// (ExpiresAt is zero because the cache only stores ownerID).
// Cache miss renews the session in the DB (new id + expiry), then caches it.
// Missing sessions return repository.ErrNotFound.
func (s *Session) ResolveOwner(ctx context.Context, sessionID string) (ownerID int64, resolvedSessionID string, expiresAt time.Time, err error) {
	if val, found := s.cacheStore.OwnerSessions.Get(sessionID); found {
		ownerID, ok := val.(int64)
		if !ok {
			return 0, "", time.Time{}, ErrInvalidOwnerSessionCache
		}
		return ownerID, sessionID, time.Time{}, nil
	}

	ownerSession, err := s.ownerSessionRepo.Update(ctx, sessionID, time.Now().Add(s.ownerSessionTTL))
	if err != nil {
		return 0, "", time.Time{}, err
	}

	resolvedSessionID = ownerSession.ID.String()
	s.cacheStore.OwnerSessions.SetDefault(resolvedSessionID, ownerSession.OwnerID)
	return ownerSession.OwnerID, resolvedSessionID, ownerSession.ExpiresAt.Time, nil
}
