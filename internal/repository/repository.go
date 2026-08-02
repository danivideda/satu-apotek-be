package repository

import (
	"context"
	"time"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
	"github.com/danivideda/satu-apotek-be/internal/env"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ownerSessionTTL = env.GetString("OWNER_SESSION_TTL", "168h")

type Repository struct {
	Owners interface {
		GetByID(ctx context.Context, id int64) (*dbsqlc.Owner, error)
		Create(ctx context.Context, username, email, passwordHash string) (ownerID int64, ownerSessionID string, exp time.Time, err error)
		GetByUsername(ctx context.Context, username string) (*dbsqlc.GetOwnerByUsernameRow, error)
		GetByEmail(ctx context.Context, email string) (*dbsqlc.GetOwnerByEmailRow, error)
	}

	Users interface {
		Create(ctx context.Context, username, passwordHash string, pharmacyID int64) (*dbsqlc.User, error)
		GetByID(ctx context.Context, id int64) (*dbsqlc.User, error)
		GetByUsername(ctx context.Context, username string) (*dbsqlc.GetUserByUsernameRow, error)
		GetByPharmacyID(ctx context.Context, pharmacyID int64) (*[]dbsqlc.GetUserByPharmacyIDRow, error)
	}

	OwnerSessions interface {
		Create(ctx context.Context, ownerID int64, exp time.Time) (*dbsqlc.OwnerSession, error)
		Update(ctx context.Context, sessionID string, exp time.Time) (*dbsqlc.OwnerSession, error)
		Get(ctx context.Context, sessionID string) (*dbsqlc.OwnerSession, error)
		Delete(ctx context.Context, sessionID string) (*dbsqlc.OwnerSession, error)
		DeleteExpired(ctx context.Context) (*[]dbsqlc.OwnerSession, error)
	}

	UserSessions interface {
		Create(ctx context.Context, userID int64, exp time.Time) (*dbsqlc.UserSession, error)
		Update(ctx context.Context, sessionID string, exp time.Time) (*dbsqlc.UserSession, error)
		Get(ctx context.Context, sessionID string) (*dbsqlc.UserSession, error)
		Delete(ctx context.Context, sessionID string) (*dbsqlc.UserSession, error)
		DeleteExpired(ctx context.Context) (*[]dbsqlc.UserSession, error)
	}

	PharmacySessions interface {
		Create(ctx context.Context, pharmacyID int64, exp time.Time) (*dbsqlc.PharmacySession, error)
		Update(ctx context.Context, sessionID string, exp time.Time) (*dbsqlc.PharmacySession, error)
		Get(ctx context.Context, sessionID string) (*dbsqlc.PharmacySession, error)
		Delete(ctx context.Context, sessionID string) (*dbsqlc.PharmacySession, error)
		DeleteExpired(ctx context.Context) (*[]dbsqlc.PharmacySession, error)
	}

	Pharmacies interface {
		GetByIDForOwner(ctx context.Context, pharmacyID, ownerID int64) (*dbsqlc.Pharmacy, error)
		GetByID(ctx context.Context, pharmacyID int64) (*dbsqlc.Pharmacy, error)
		GetByAppIDForOwner(ctx context.Context, appID string, ownerID int64) (*dbsqlc.Pharmacy, error)
		GetCodeByCode(ctx context.Context, code string) (*dbsqlc.PharmacyCode, error)
		GetByOwnerID(ctx context.Context, ownerID int64) (*[]dbsqlc.Pharmacy, error)
		Create(ctx context.Context, ownerID int64, name string) (*dbsqlc.Pharmacy, error)
		UpsertCode(ctx context.Context, apotekID int64, code string) (*dbsqlc.PharmacyCode, error)
		GetCodeByID(ctx context.Context, apotekID int64) (*dbsqlc.PharmacyCode, error)
		DeleteExpiredCode(ctx context.Context) (*[]dbsqlc.PharmacyCode, error)
		DeleteCode(ctx context.Context, code string) (*dbsqlc.PharmacyCode, error)
	}

	CacheStore *CacheStore
}

func New(db *pgxpool.Pool, cs *CacheStore) Repository {
	q := dbsqlc.New(db)

	return Repository{
		Owners:           &ownersRepo{db: db, queries: q},
		Users:            &usersRepo{queries: q},
		OwnerSessions:    &ownerSessionsRepo{queries: q},
		UserSessions:     &userSessionsRepo{queries: q},
		PharmacySessions: &pharmacySessionsRepo{queries: q},
		Pharmacies:       &pharmaciesRepo{db: db, queries: q},
		CacheStore:       cs,
	}
}
