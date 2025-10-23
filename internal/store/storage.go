package store

import (
	"context"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
)

type Storage struct {
	Owners interface {
		GetOwnerByID(context.Context, int) (*dbsqlc.Owner, error)
		CreateOwner(context.Context, dbsqlc.CreateOwnerParams) (*dbsqlc.CreateOwnerRow, error)
	}
}

func NewStorage(db dbsqlc.DBTX) Storage {
	return Storage{
		Owners: &OwnerStore{db: db},
	}
}