package store

import (
	"context"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
)

type Storage struct {
	Owners interface {
		GetByID(ctx context.Context, id string) (*dbsqlc.Owner, error)
		Create(ctx context.Context, params dbsqlc.CreateOwnerParams) (*dbsqlc.CreateOwnerRow, error)
		GetByUsername(ctx context.Context, username string) (*dbsqlc.GetOwnerByUsernameRow, error)
	}

	Users interface {
		Create(ctx context.Context, username, passwordHash, pharmacyID string) (*dbsqlc.User, error)
	}

	RevokedSessions interface {
		Create(ctx context.Context, sessionID string) (*dbsqlc.RevokedSession, error)
		GetBySessionID(ctx context.Context, sessionID string) (*dbsqlc.RevokedSession, error)
	}

	Pharmacies interface {
		Create(ctx context.Context, ownerID string, name string) (*dbsqlc.Pharmacy, error)
	}
}

func NewStorage(db dbsqlc.DBTX) Storage {
	queries := dbsqlc.New(db)

	return Storage{
		Owners: &OwnerStore{queries: queries},
		Users: &UserStore{queries: queries},
		RevokedSessions: &RevokedSessionStore{queries: queries},
		Pharmacies: &PharmacyStore{queries: queries},
	}
}