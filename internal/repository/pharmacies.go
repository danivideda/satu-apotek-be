package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
	"github.com/danivideda/satu-apotek-be/internal/env"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sqids/sqids-go"
)

type pharmaciesRepo struct {
	db      *pgxpool.Pool
	queries *dbsqlc.Queries
}

func (r *pharmaciesRepo) GetByID(ctx context.Context, pharmacyID int64) (*dbsqlc.Pharmacy, error) {
	pharmacy, err := r.queries.GetPharmacyByID(ctx, pharmacyID)
	if err != nil {
		return nil, err
	}
	return &pharmacy, nil
}

func (r *pharmaciesRepo) GetByIDForOwner(ctx context.Context, pharmacyID, ownerID int64) (*dbsqlc.Pharmacy, error) {
	pharmacy, err := r.queries.GetPharmacyByIDForOwner(ctx, dbsqlc.GetPharmacyByIDForOwnerParams{
		ID:      pharmacyID,
		OwnerID: ownerID,
	})
	if err != nil {
		return nil, err
	}
	return &pharmacy, nil
}

func (r *pharmaciesRepo) GetByAppIDForOwner(ctx context.Context, appID string, ownerID int64) (*dbsqlc.Pharmacy, error) {
	pharmacy, err := r.queries.GetPharmacyByAppIDForOwner(ctx, dbsqlc.GetPharmacyByAppIDForOwnerParams{
		AppID:   appID,
		OwnerID: ownerID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &pharmacy, nil
}

func (r *pharmaciesRepo) GetCodeByCode(ctx context.Context, code string) (*dbsqlc.PharmacyCode, error) {
	pharmacyCode, err := r.queries.GetApotekCodeByCode(ctx, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &pharmacyCode, nil
}

func (r *pharmaciesRepo) GetByOwnerID(ctx context.Context, ownerID int64) (*[]dbsqlc.Pharmacy, error) {
	pharmacies, err := r.queries.GetPharmaciesByOwner(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	return &pharmacies, nil
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

	alphabets, err := qtx.GetAlphabets(ctx)
	if err != nil {
		err := fmt.Errorf("sqids config not set, %s", err)
		return nil, err
	}

	sq, err := sqids.New(sqids.Options{
		Alphabet:  alphabets,
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
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &pharmacyWithAppID, nil
}

func (r *pharmaciesRepo) GetCodeByID(ctx context.Context, pharmacyID int64) (*dbsqlc.PharmacyCode, error) {
	apotekCode, err := r.queries.GetApotekCode(ctx, pharmacyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &apotekCode, nil
}

func (r *pharmaciesRepo) UpsertCode(ctx context.Context, pharmacyID int64, code string) (*dbsqlc.PharmacyCode, error) {
	ttl, err := time.ParseDuration(env.GetString("CODE_TTL", "5m"))
	if err != nil {
		return nil, err
	}
	exp := pgtype.Timestamptz{
		Time:  time.Now().Add(ttl),
		Valid: true,
	}
	params := dbsqlc.UpsertApotekCodeParams{
		ApotekID:  pharmacyID,
		Code:      code,
		ExpiresAt: exp,
	}
	apotekCode, err := r.queries.UpsertApotekCode(ctx, params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, ErrDuplicateValue
			}
		}
		return nil, err
	}

	return &apotekCode, nil
}

func (r *pharmaciesRepo) DeleteCode(ctx context.Context, code string) (*dbsqlc.PharmacyCode, error) {
	pharmacyCode, err := r.queries.DeletePharmacyCode(ctx, code)
	if err != nil {
		return nil, err
	}

	return &pharmacyCode, nil
}

func (r *pharmaciesRepo) DeleteExpiredCode(ctx context.Context) (*[]dbsqlc.PharmacyCode, error) {
	apotekCode, err := r.queries.DeleteExpiredApotekCode(ctx)
	if err != nil {
		return nil, err
	}

	return &apotekCode, nil
}
