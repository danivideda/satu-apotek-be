package store

import (
	"context"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
)

type OwnerStore struct {
	queries *dbsqlc.Queries
}

func (s *OwnerStore) CreateOwner(ctx context.Context, createOwnerParams dbsqlc.CreateOwnerParams) (*dbsqlc.CreateOwnerRow, error) {
	return nil, nil
}

func (s *OwnerStore) GetOwnerByID(ctx context.Context, id int) (*dbsqlc.Owner, error) {
	owner, err := s.queries.GetOwnerByID(ctx, int32(id))

	return &owner, err
}
