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

	ownerSessionTTL    time.Duration
	userSessionTTL     time.Duration
	pharmacySessionTTL time.Duration

	usersRepo      repository.UsersRepository
	pharmaciesRepo repository.PharmaciesRepository
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
		repo.Users,
		repo.Pharmacies,
	}
}

type OwnerSession struct {
	ID      string
	Exp     time.Time
	OwnerID int64
}

func (s *Session) NewOwner(ctx context.Context, ownerID int64) (*OwnerSession, error) {
	newExp := time.Now().Add(s.ownerSessionTTL)
	newSession, err := s.ownerSessionRepo.Create(ctx, ownerID, newExp)
	if err != nil {
		return nil, err
	}
	s.cacheStore.OwnerSessions.SetDefault(newSession.ID.String(), ownerID)

	newOwnerSession := &OwnerSession{
		OwnerID: ownerID,
		ID:      newSession.ID.String(),
		Exp:     newExp,
	}

	return newOwnerSession, nil
}

// Return the same Owner SessionID if still exist in cache, otherwise checks the DB and rotate sessionID if still valid / not expired
func (s *Session) ResolveOwner(ctx context.Context, sessionID string) (newOwnerSession *OwnerSession, rotated bool, err error) {
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

	newOwnerSession = &OwnerSession{
		OwnerID: ownerSessionItem.OwnerID,
		ID:      ownerSessionItem.ID.String(),
		Exp:     newExp,
	}
	s.cacheStore.OwnerSessions.SetDefault(newOwnerSession.ID, newOwnerSession.OwnerID)

	return newOwnerSession, true, nil
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
	ID       string
	Exp      time.Time
	UserID   int64
	Username string
}

func (s *Session) NewUser(ctx context.Context, userID int64, username string) (*UserSession, error) {
	newExp := time.Now().Add(s.userSessionTTL)
	newSession, err := s.userSessionRepo.Create(ctx, userID, newExp)
	if err != nil {
		return nil, err
	}
	userCacheValue := repository.UserCacheValue{
		ID:       newSession.UserID,
		Username: username,
	}
	s.cacheStore.UserSessions.SetDefault(newSession.ID.String(), userCacheValue)

	newUserSession := &UserSession{
		ID:       newSession.ID.String(),
		Exp:      newSession.ExpiresAt.Time,
		UserID:   newSession.UserID,
		Username: username,
	}
	return newUserSession, nil
}

// Return the same User SessionID if still exist in cache, otherwise checks the DB and rotate sessionID if still valid / not expired
func (s *Session) ResolveUser(ctx context.Context, sessionID string) (newUserSession *UserSession, rotated bool, err error) {
	newExp := time.Now().Add(s.userSessionTTL)

	// 1. Check session if still in cache
	if val, found := s.cacheStore.UserSessions.Get(sessionID); found {
		userCache, ok := val.(repository.UserCacheValue)
		if !ok {
			s.cacheStore.UserSessions.Delete(sessionID)
			return nil, false, fmt.Errorf("%w: UserCacheValue expected, got %T", ErrTypeAssertionFailed, userCache)
		}
		userSession := &UserSession{
			ID:       sessionID,
			Exp:      time.Time{},
			UserID:   userCache.ID,
			Username: userCache.Username,
		}
		return userSession, false, nil
	}

	// 2. check session in DB
	userSessionItem, err := s.userSessionRepo.Update(ctx, sessionID, newExp)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, false, fmt.Errorf("%w: user %w", err, ErrInvalidSession)
		}
		return nil, false, err
	}

	user, err := s.usersRepo.GetByID(ctx, userSessionItem.UserID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			return nil, false, fmt.Errorf("%w: session exist but user not found", err)
		default:
			return nil, false, err
		}
	}
	userCache := repository.UserCacheValue{
		ID:       userSessionItem.UserID,
		Username: user.Username,
	}
	newUserSession = &UserSession{
		UserID: userSessionItem.UserID,
		ID:     userSessionItem.ID.String(),
		Exp:    newExp,
	}
	s.cacheStore.UserSessions.SetDefault(newUserSession.ID, userCache)

	return newUserSession, true, nil
}

func (s *Session) DeleteUser(ctx context.Context, sessionID string) error {
	_, err := s.userSessionRepo.Delete(ctx, sessionID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			s.cacheStore.UserSessions.Delete(sessionID)
			return fmt.Errorf("%w: user %w", err, ErrInvalidSession)
		default:
			return err
		}
	}

	s.cacheStore.UserSessions.Delete(sessionID)
	return nil
}

type PharmacySession struct {
	ID  string
	Exp time.Time
	repository.PharmacyCacheValue
}

func (s *Session) NewPharmacy(ctx context.Context, pharmacyID int64) (*PharmacySession, error) {
	// create pharmacy session
	exp := time.Now().Add(s.pharmacySessionTTL)
	pharmacySessionItem, err := s.pharmacySessionRepo.Create(ctx, pharmacyID, exp)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			return nil, fmt.Errorf("%w: pharmacyID not found when creating session", err)
		default:
			return nil, err
		}
	}

	// Gather information for session creation
	// embed the Users belong in a Pharmacy session (this will be used in login-page)
	users, err := s.usersRepo.GetByPharmacyID(ctx, pharmacyID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			users = nil
		default:
			return nil, err
		}
	}

	var usersCache []repository.UserCacheValue
	if users != nil {
		for _, user := range *users {
			userItem := repository.UserCacheValue{
				ID:       user.ID,
				Username: user.Username,
			}
			usersCache = append(usersCache, userItem)
		}
	}

	pharmacy, err := s.pharmaciesRepo.GetByID(ctx, pharmacyID)
	if err != nil {
		return nil, err
	}

	s.cacheStore.PharmacySessions.SetDefault(pharmacySessionItem.ID.String(), repository.PharmacyCacheValue{
		PharmacyID: pharmacyID,
		Name:       pharmacy.Name,
		Users:      usersCache,
	})

	newPharmacySession := &PharmacySession{
		ID:  pharmacySessionItem.ID.String(),
		Exp: exp,
	}

	return newPharmacySession, nil
}

// Return the same Pharmacy SessionID if still exist in cache, otherwise checks the DB and rotate sessionID if still valid / not expired
func (s *Session) ResolvePharmacy(ctx context.Context, sessionID string) (newPharmacySession *PharmacySession, rotated bool, err error) {
	newExp := time.Now().Add(s.pharmacySessionTTL)

	// 1. Check session if still in cache
	if val, found := s.cacheStore.PharmacySessions.Get(sessionID); found {
		pharmacyCache, ok := val.(repository.PharmacyCacheValue)
		if !ok {
			s.cacheStore.PharmacySessions.Delete(sessionID)
			return nil, false, fmt.Errorf("%w: PharmacyCacheValue expected, got %T", ErrTypeAssertionFailed, pharmacyCache)
		}
		pharmacySession := &PharmacySession{
			ID:  sessionID,
			Exp: time.Time{},
			PharmacyCacheValue: pharmacyCache,
		}
		return pharmacySession, false, nil
	}

	// 2. check session in DB
	pharmacySessionItem, err := s.pharmacySessionRepo.Update(ctx, sessionID, newExp)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, false, fmt.Errorf("%w: pharmachy %w", err, ErrInvalidSession)
		}
		return nil, false, err
	}

	// Get pharmacy name
	pharmacyItem, err := s.pharmaciesRepo.GetByID(ctx, pharmacySessionItem.PharmacyID)
	if err != nil {
		return nil, false, err
	}
	// Get users
	users, err := s.usersRepo.GetByPharmacyID(ctx, pharmacySessionItem.PharmacyID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			users = nil
		} else {
			return nil, false, err
		}
	}

	var usersCache []repository.UserCacheValue
	if users != nil {
		for _, user := range *users {
			userItem := repository.UserCacheValue{
				ID:       user.ID,
				Username: user.Username,
			}

			usersCache = append(usersCache, userItem)
		}
	}
	pharmacyCache := repository.PharmacyCacheValue{
		PharmacyID: pharmacyItem.ID,
		Name:       pharmacyItem.Name,
		Users:      usersCache,
	}

	newPharmacySession = &PharmacySession{
		ID:  pharmacySessionItem.ID.String(),
		Exp: newExp,
		PharmacyCacheValue: pharmacyCache,
	}

	s.cacheStore.PharmacySessions.SetDefault(newPharmacySession.ID, pharmacyCache)

	return newPharmacySession, true, nil
}
