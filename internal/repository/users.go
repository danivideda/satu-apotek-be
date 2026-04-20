package repository

import (
	"context"
	"fmt"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
)

type usersRepo struct {
	queries *dbsqlc.Queries
}

func (r *usersRepo) Create(ctx context.Context, username, passwordHash string, pharmacyID int64) (*dbsqlc.User, error) {
	params := dbsqlc.CreateUserParams{
		Username:     username,
		PasswordHash: passwordHash,
		PharmacyID:   pharmacyID,
	}
	user, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		if isDuplicateError(err) {
			return nil, ErrDuplicateValue
		}
		return nil, err
	}

	return &user, nil
}

func (r *usersRepo) GetByID(ctx context.Context, id int64) (*dbsqlc.User, error) {
	user, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if isNotFoundError(err) {
			return nil, fmt.Errorf("%w: user doesn't exist", ErrNotFound)
		}
		return nil, err
	}
	return &user, nil
}

func (r *usersRepo) GetByUsername(ctx context.Context, username string) (*dbsqlc.GetUserByUsernameRow, error) {
	user, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if isNotFoundError(err) {
			return nil, fmt.Errorf("%w: user doesn't exist", ErrNotFound)
		}
		return nil, err
	}
	return &user, nil
}

func (r *usersRepo) GetByPharmacyID(ctx context.Context, pharmacyID int64) (*[]dbsqlc.GetUserByPharmacyIDRow, error) {
	users, err := r.queries.GetUserByPharmacyID(ctx, pharmacyID)
	if err != nil {
		if isNotFoundError(err) {
			return nil, fmt.Errorf("%w: user doesn't exist", ErrNotFound)
		}
		return nil, err
	}
	return &users, nil
}
