package store

import (
	"context"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
)

type OwnerStore struct {
	db dbsqlc.DBTX
}

func (s *OwnerStore) CreateOwner(ctx context.Context, createOwnerParams dbsqlc.CreateOwnerParams) (*dbsqlc.Owner, error) {
	return nil, nil
}

func (s *OwnerStore) GetOwnerByID(ctx context.Context, id int) (*dbsqlc.Owner, error) {
	return nil, nil
}
