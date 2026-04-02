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
	}

	Users interface {
		Create(ctx context.Context, username, passwordHash string, pharmacyID int64) (*dbsqlc.User, error)
	}

	OwnerSessions interface {
		Create(ctx context.Context, ownerID int64, exp time.Time) (*dbsqlc.OwnerSession, error)
		Update(ctx context.Context, sessionID string, exp time.Time) (*dbsqlc.OwnerSession, error)
		Get(ctx context.Context, sessionID string) (*dbsqlc.OwnerSession, error)
		Delete(ctx context.Context, sessionID string) (*dbsqlc.OwnerSession, error)
		DeleteExpired(ctx context.Context) (*[]dbsqlc.OwnerSession, error)
	}

	Pharmacies interface {
		GetByID(ctx context.Context, pharmacyID int64) (*dbsqlc.Pharmacy, error)
		GetByAppID(ctx context.Context, appID string) (*dbsqlc.Pharmacy, error)
		GetByCode(ctx context.Context, code string) (*dbsqlc.PharmacyCode, error)
		GetByOwnerID(ctx context.Context, ownerID int64) (*[]dbsqlc.Pharmacy, error)
		Create(ctx context.Context, ownerID int64, name string) (*dbsqlc.Pharmacy, error)
		UpsertCode(ctx context.Context, apotekID int64, code string) (*dbsqlc.PharmacyCode, error)
		GetCodeByID(ctx context.Context, apotekID int64) (*dbsqlc.PharmacyCode, error)
		DeleteExpired(ctx context.Context) (*[]dbsqlc.PharmacyCode, error)
	}

	CacheStore *CacheStore
}

func New(db *pgxpool.Pool, cs *CacheStore) Repository {
	q := dbsqlc.New(db)

	return Repository{
		Owners:        &ownersRepo{db: db, queries: q},
		Users:         &usersRepo{queries: q},
		OwnerSessions: &ownerSessionsRepo{queries: q},
		Pharmacies:    &pharmaciesRepo{db: db, queries: q},
		CacheStore:    cs,
	}
}
