package service

import (
	"context"
	"errors"
	"slices"

	"github.com/danivideda/satu-apotek-be/internal/repository"
)

func GetUsersFromPharmacyID(ctx context.Context, r repository.Repository, pharmacyID int64) (*[]repository.UserCache, error) {
	users, err := r.Users.GetByPharmacyID(ctx, pharmacyID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			users = nil
		} else {
			return nil, err
		}
	}
	var usersCache []repository.UserCache
	if users != nil {
		for _, user := range *users {
			userItem := repository.UserCache{
				ID:       user.ID,
				Username: user.Username,
			}
			usersCache = append(usersCache, userItem)
		}
	}
	return &usersCache, nil
}

func UserExistsInPharmacy(users []repository.UserCache, userID int64) bool {
	userExists := slices.ContainsFunc(users, func(u repository.UserCache) bool {
		if u.ID == userID {
			return true
		}
		return false
	})
	return userExists
}