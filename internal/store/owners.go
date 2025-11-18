package store

import (
	"context"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
)

type OwnerStore struct {
	queries *dbsqlc.Queries
}

func (s *OwnerStore) Create(ctx context.Context, createOwnerParams dbsqlc.CreateOwnerParams) (*dbsqlc.CreateOwnerRow, error) {
	owner, err := s.queries.CreateOwner(ctx, createOwnerParams)
	if err != nil {
		return nil, err
	}
	return &owner, nil
}

func (s *OwnerStore) GetByID(ctx context.Context, id int) (*dbsqlc.Owner, error) {
	owner, err := s.queries.GetOwnerByID(ctx, int32(id))
	if err != nil {
		return nil, err
	}
	return &owner, nil
}

func (s *OwnerStore) GetByUsername(ctx context.Context, username string) (*dbsqlc.GetOwnerByUsernameRow, error) {
	owner, err := s.queries.GetOwnerByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return &owner, nil
}