package store

import (
	"context"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
	"github.com/jackc/pgx/v5/pgtype"
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

func (s *OwnerStore) GetByID(ctx context.Context, id string) (*dbsqlc.Owner, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan("hello"); err != nil {
		return nil, err
	}
	owner, err := s.queries.GetOwnerByID(ctx, uuid)
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