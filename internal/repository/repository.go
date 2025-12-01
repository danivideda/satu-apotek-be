package repository

import (
	"context"
	"time"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
)

type Repository struct {
	Owners interface {
		GetByID(ctx context.Context, id int64) (*dbsqlc.Owner, error)
		Create(ctx context.Context, username, email, passwordHash string) (*dbsqlc.CreateOwnerRow, error)
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
	}

	Pharmacies interface {
		Create(ctx context.Context, ownerID int64, name string) (*dbsqlc.Pharmacy, error)
	}

	ApotekCode interface {
		Get(ctx context.Context, apotekID int64) (*dbsqlc.ApotekCode, error)
		GetByCode(ctx context.Context, code string) (*dbsqlc.ApotekCode, error)
		Upsert(ctx context.Context, apotekID int64, code string) (*dbsqlc.ApotekCode, error)
		DeleteExpired(ctx context.Context) (*[]dbsqlc.ApotekCode, error)
	}
}

func New(db dbsqlc.DBTX) Repository {
	queries := dbsqlc.New(db)

	return Repository{
		Owners:        &ownersRepo{queries: queries},
		Users:         &usersRepo{queries: queries},
		OwnerSessions: &ownerSessionsRepo{queries: queries},
		Pharmacies:    &pharmaciesRepo{queries: queries},
		ApotekCode:    &apotekCodesRepo{queries: queries},
	}
}
