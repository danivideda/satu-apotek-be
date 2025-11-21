package repository

import (
	"context"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
)

type Repository struct {
	Owners interface {
		GetByID(ctx context.Context, id string) (*dbsqlc.Owner, error)
		Create(ctx context.Context, username, email, passwordHash string) (*dbsqlc.CreateOwnerRow, error)
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

func NewRepository(db dbsqlc.DBTX) Repository {
	queries := dbsqlc.New(db)

	return Repository{
		Owners: &OwnersRepo{queries: queries},
		Users: &UsersRepo{queries: queries},
		RevokedSessions: &RevokedSessionsRepo{queries: queries},
		Pharmacies: &PharmaciesRepo{queries: queries},
	}
}