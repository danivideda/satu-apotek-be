package store

import (
	"context"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type PharmacyStore struct {
	queries *dbsqlc.Queries
}

func (s *PharmacyStore) Create(ctx context.Context, ownerID string, name string) (*dbsqlc.Pharmacy, error) {

	var ownerUUID pgtype.UUID
	if err := ownerUUID.Scan(ownerID); err != nil {
		return nil, err
	}
	params := dbsqlc.CreatePharmacyParams{
		OwnerID: ownerUUID,
		Name:    pgtype.Text{String: name, Valid: true},
	}
	pharmacy, err := s.queries.CreatePharmacy(ctx, params)
	if err != nil {
		return nil, err
	}

	return &pharmacy, nil
}
