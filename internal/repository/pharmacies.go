package repository

import (
	"context"
	"fmt"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
	"github.com/danivideda/satu-apotek-be/internal/env"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sqids/sqids-go"
)

type pharmaciesRepo struct {
	db      *pgxpool.Pool
	queries *dbsqlc.Queries
}

func (r *pharmaciesRepo) Create(ctx context.Context, ownerID int64, name string) (*dbsqlc.Pharmacy, error) {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	qtx := r.queries.WithTx(tx)

	params := dbsqlc.CreatePharmacyParams{OwnerID: ownerID, Name: name}
	pharmacy, err := qtx.CreatePharmacy(ctx, params)
	if err != nil {
		return nil, err
	}

	sq, err := sqids.New(sqids.Options{
		Alphabet:  env.GetString("SQIDS_LIST", "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"),
		MinLength: 5,
	})
	if err != nil {
		return nil, err
	}

	appID, err := sq.Encode([]uint64{uint64(pharmacy.ID)})
	if err != nil {
		return nil, err
	}

	params2 := dbsqlc.InsertAppIDParams{
		AppID: appID,
		ID:    pharmacy.ID,
	}
	pharmacyWithAppID, err := qtx.InsertAppID(ctx, params2)
	if err != nil {
		fmt.Println("ERR: no 2")
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &pharmacyWithAppID, nil
}

func (r *pharmaciesRepo) GetByOwner(ctx context.Context, ownerID int64) (*[]dbsqlc.Pharmacy, error) {
	pharmacies, err := r.queries.GetPharmaciesByOwner(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	return &pharmacies, nil
}
