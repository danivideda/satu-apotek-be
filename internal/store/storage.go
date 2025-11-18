package store

import (
	"context"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
)

type Storage struct {
	Owners interface {
		GetByID(context.Context, int) (*dbsqlc.Owner, error)
		Create(context.Context, dbsqlc.CreateOwnerParams) (*dbsqlc.CreateOwnerRow, error)
		GetByUsername(context.Context, string) (*dbsqlc.GetOwnerByUsernameRow, error)
	}

	RevokedSessions interface {
		Create(context.Context, string) (*dbsqlc.RevokedSession, error)
		GetBySessionID(context.Context, string) (*dbsqlc.RevokedSession, error)
	}
}

func NewStorage(db dbsqlc.DBTX) Storage {
	queries := dbsqlc.New(db)

	return Storage{
		Owners: &OwnerStore{queries: queries},
		RevokedSessions: &RevokedSessionStore{queries: queries},
	}
}