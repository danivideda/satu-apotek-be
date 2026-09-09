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
