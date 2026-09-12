package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
	"github.com/danivideda/satu-apotek-be/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/patrickmn/go-cache"
)

const (
	sessionIDA = "11111111-1111-1111-1111-111111111111"
	sessionIDB = "22222222-2222-2222-2222-222222222222"
)

type stubOwnerSessions struct {
	updateFn func(ctx context.Context, sessionID string, exp time.Time) (*dbsqlc.OwnerSession, error)
	createFn func(ctx context.Context, ownerID int64, exp time.Time) (*dbsqlc.OwnerSession, error)
}

func (s *stubOwnerSessions) Create(ctx context.Context, ownerID int64, exp time.Time) (*dbsqlc.OwnerSession, error) {
	if s.createFn != nil {
		return s.createFn(ctx, ownerID, exp)
	}
	return nil, errors.New("Create not implemented")
}

func (s *stubOwnerSessions) Update(ctx context.Context, sessionID string, exp time.Time) (*dbsqlc.OwnerSession, error) {
	if s.updateFn != nil {
		return s.updateFn(ctx, sessionID, exp)
	}
	return nil, errors.New("Update not implemented")
}

func (s *stubOwnerSessions) Get(context.Context, string) (*dbsqlc.OwnerSession, error) {
	return nil, errors.New("Get not implemented")
}

func (s *stubOwnerSessions) Delete(context.Context, string) (*dbsqlc.OwnerSession, error) {
	return nil, errors.New("Delete not implemented")
}

func (s *stubOwnerSessions) DeleteExpired(context.Context) (*[]dbsqlc.OwnerSession, error) {
	return nil, errors.New("DeleteExpired not implemented")
}

func mustUUID(s string) pgtype.UUID {
	var id pgtype.UUID
	if err := id.Scan(s); err != nil {
		panic(err)
	}
	return id
}

func newTestSession(ownerRepo repository.OwnerSessionsRepository, ownerCache *cache.Cache) *Session {
	return &Session{
		cacheStore:       &repository.CacheStore{OwnerSessions: ownerCache},
		ownerSessionRepo: ownerRepo,
		ownerSessionTTL:  24 * time.Hour,
	}
}

func TestResolveOwner_CacheHit(t *testing.T) {
	t.Parallel()

	ownerCache := cache.New(5*time.Minute, 10*time.Minute)
	ownerCache.SetDefault(sessionIDA, int64(42))

	var updateCalled bool
	svc := newTestSession(&stubOwnerSessions{
		updateFn: func(context.Context, string, time.Time) (*dbsqlc.OwnerSession, error) {
			updateCalled = true
			return nil, errors.New("Update should not be called on cache hit")
		},
	}, ownerCache)

	ownerID, resolvedID, expiresAt, err := svc.ResolveOwner(context.Background(), sessionIDA)
	if err != nil {
		t.Fatalf("ResolveOwner returned error: %v", err)
	}
	if ownerID != 42 {
		t.Fatalf("ownerID = %d, want 42", ownerID)
	}
	if resolvedID != sessionIDA {
		t.Fatalf("resolvedSessionID = %s, want %s", resolvedID, sessionIDA)
	}
	if !expiresAt.IsZero() {
		t.Fatalf("expiresAt = %v, want zero on cache hit", expiresAt)
	}
	if updateCalled {
		t.Fatal("Update was called on cache hit")
	}
}

func TestResolveOwner_CacheMissRenewsSession(t *testing.T) {
	t.Parallel()

	ownerCache := cache.New(5*time.Minute, 10*time.Minute)
	exp := time.Now().Add(24 * time.Hour)
	svc := newTestSession(&stubOwnerSessions{
		updateFn: func(_ context.Context, sessionID string, gotExp time.Time) (*dbsqlc.OwnerSession, error) {
			if sessionID != sessionIDA {
				t.Fatalf("Update sessionID = %s, want %s", sessionID, sessionIDA)
			}
			if gotExp.Before(time.Now()) {
				t.Fatal("Update expiry is in the past")
			}
			return &dbsqlc.OwnerSession{
				ID:        mustUUID(sessionIDB),
				OwnerID:   7,
				ExpiresAt: pgtype.Timestamptz{Time: exp, Valid: true},
			}, nil
		},
	}, ownerCache)

	ownerID, resolvedID, expiresAt, err := svc.ResolveOwner(context.Background(), sessionIDA)
	if err != nil {
		t.Fatalf("ResolveOwner returned error: %v", err)
	}
	if ownerID != 7 {
		t.Fatalf("ownerID = %d, want 7", ownerID)
	}
	if resolvedID != sessionIDB {
		t.Fatalf("resolvedSessionID = %s, want %s", resolvedID, sessionIDB)
	}
	if !expiresAt.Equal(exp) {
		t.Fatalf("expiresAt = %v, want %v", expiresAt, exp)
	}

	cached, found := ownerCache.Get(sessionIDB)
	if !found {
		t.Fatal("renewed session was not written to cache")
	}
	if cached.(int64) != 7 {
		t.Fatalf("cached ownerID = %v, want 7", cached)
	}
}

func TestResolveOwner_NotFound(t *testing.T) {
	t.Parallel()

	ownerCache := cache.New(5*time.Minute, 10*time.Minute)
	svc := newTestSession(&stubOwnerSessions{
		updateFn: func(context.Context, string, time.Time) (*dbsqlc.OwnerSession, error) {
			return nil, repository.ErrNotFound
		},
	}, ownerCache)

	_, _, _, err := svc.ResolveOwner(context.Background(), sessionIDA)
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("error = %v, want %v", err, repository.ErrNotFound)
	}
}

func TestResolveOwner_InvalidCacheType(t *testing.T) {
	t.Parallel()

	ownerCache := cache.New(5*time.Minute, 10*time.Minute)
	ownerCache.SetDefault(sessionIDA, "not-an-int64")
	svc := newTestSession(&stubOwnerSessions{}, ownerCache)

	_, _, _, err := svc.ResolveOwner(context.Background(), sessionIDA)
	if !errors.Is(err, ErrInvalidOwnerSessionCache) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidOwnerSessionCache)
	}
}
