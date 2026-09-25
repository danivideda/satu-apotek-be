package service

import (
	"context"
	"errors"
	"slices"

	"github.com/danivideda/satu-apotek-be/internal/repository"
)

func GetUsersFromPharmacyID(ctx context.Context, users repository.UsersRepository, pharmacyID int64) (*[]repository.UserCacheValue, error) {
	rows, err := users.GetByPharmacyID(ctx, pharmacyID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			rows = nil
		} else {
			return nil, err
		}
	}
	var usersCache []repository.UserCacheValue
	if rows != nil {
		for _, user := range *rows {
			userItem := repository.UserCacheValue{
				ID:       user.ID,
				Username: user.Username,
			}
			usersCache = append(usersCache, userItem)
		}
	}
	return &usersCache, nil
}

func UserExistsInPharmacy(users []repository.UserCacheValue, userID int64) bool {
	userExists := slices.ContainsFunc(users, func(u repository.UserCacheValue) bool {
		if u.ID == userID {
			return true
		}
		return false
	})
	return userExists
}
