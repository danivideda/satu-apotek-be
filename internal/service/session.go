package service

import (
	"context"
	"errors"
	"fmt"
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

type OwnerSession struct {
	ID      string
	Exp     time.Time
	OwnerID int64
}

func (s *Session) NewOwner(ctx context.Context, ownerID int64) (*OwnerSession, error) {
	exp := time.Now().Add(s.ownerSessionTTL)
	session, err := s.ownerSessionRepo.Create(ctx, ownerID, exp)
	if err != nil {
		return nil, err
	}
	s.cacheStore.OwnerSessions.SetDefault(session.ID.String(), ownerID)

	newOwnerSession := &OwnerSession{
		OwnerID: ownerID,
		ID:      session.ID.String(),
		Exp:     exp,
	}

	return newOwnerSession, nil
}

// Return the same Owner SessionID if still exist in cache, otherwise checks the DB and rotate sessionID if still valid / not expired
func (s *Session) ResolveOwner(ctx context.Context, sessionID string) (ownerSession *OwnerSession, rotated bool, err error) {
	newExp := time.Now().Add(s.ownerSessionTTL)

	// 1. Check session if still in cache
	if val, found := s.cacheStore.OwnerSessions.Get(sessionID); found {
		ownerID, ok := val.(int64)
		if !ok {
			s.cacheStore.OwnerSessions.Delete(sessionID)
			return nil, false, fmt.Errorf("%w: int64 expected, got %T", ErrTypeAssertionFailed, ownerID)
		}
		ownerSession := &OwnerSession{
			OwnerID: ownerID,
			ID:      sessionID,
			Exp:     time.Time{},
		}
		return ownerSession, false, nil
	}

	// 2. check session in DB
	ownerSessionItem, err := s.ownerSessionRepo.Update(ctx, sessionID, newExp)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, false, fmt.Errorf("%w: owner %w", err, ErrInvalidSession)
		}
		return nil, false, err
	}

	ownerSession = &OwnerSession{
		OwnerID: ownerSessionItem.OwnerID,
		ID:      ownerSessionItem.ID.String(),
		Exp:     newExp,
	}
	s.cacheStore.OwnerSessions.SetDefault(ownerSession.ID, ownerSession.OwnerID)

	return ownerSession, true, nil
}

func (s *Session) DeleteOwner(ctx context.Context, sessionID string) error {
	_, err := s.ownerSessionRepo.Delete(ctx, sessionID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			s.cacheStore.OwnerSessions.Delete(sessionID)
			return fmt.Errorf("%w: owner %w", err, ErrInvalidSession)
		default:
			return err
		}
	}

	s.cacheStore.OwnerSessions.Delete(sessionID)
	return nil
}

type UserSession struct {
	ID     string
	Exp    time.Time
	UserID int64
}

