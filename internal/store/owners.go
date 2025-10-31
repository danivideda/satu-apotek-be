package store

import (
	"context"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
)

type OwnerStore struct {
	queries *dbsqlc.Queries
}

func (s *OwnerStore) CreateOwner(ctx context.Context, createOwnerParams dbsqlc.CreateOwnerParams) (*dbsqlc.CreateOwnerRow, error) {
	owner, err := s.queries.CreateOwner(ctx, createOwnerParams)
	return &owner, err
}

func (s *OwnerStore) GetOwnerByID(ctx context.Context, id int) (*dbsqlc.Owner, error) {
	owner, err := s.queries.GetOwnerByID(ctx, int32(id))
	return &owner, err
}

func (s *OwnerStore) GetOwnerByUsername(ctx context.Context, username string) (*dbsqlc.GetOwnerByUsernameRow, error) {
	owner, err := s.queries.GetOwnerByUsername(ctx, username)
	return &owner, err
}