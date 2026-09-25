package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/danivideda/satu-apotek-be/internal/repository"
)

type Auth struct {
	Owner    OwnerAuth
	User     UserAuth
	Pharmacy PharmacyAuth
}

func NewAuth(repo repository.Repository, sessionSvc *Session) *Auth {
	return &Auth{
		Owner:    newOwnerAuth(repo.Owners, sessionSvc),
		User:     newUserAuth(repo.Users, sessionSvc),
		Pharmacy: &pharmacyAuth{},
	}
}

type OwnerAuth interface {
	Login(ctx context.Context, email, password string) (*OwnerSession, error)
	Logout(ctx context.Context, sessionID string) error
	Register(ctx context.Context, username, email, password string) error
}
type UserAuth interface {
	Login(ctx context.Context, userID int64, password string, users []repository.UserCacheValue) (*UserSession, error)
	Logout(sessionID string) error
}
type PharmacyAuth interface {
	Connect(code string) error
}

type ownerAuth struct {
	ownerRepo  repository.OwnersRepository
	sessionSvc *Session
}

func newOwnerAuth(
	ownerRepo repository.OwnersRepository,
	session *Session,
) *ownerAuth {
	return &ownerAuth{ownerRepo, session}
}

func (a *ownerAuth) Login(ctx context.Context, email, password string) (*OwnerSession, error) {
	owner, err := a.ownerRepo.GetByEmail(ctx, email)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			return nil, fmt.Errorf("%w: owner %w", err, ErrInvalidCredentials)
		default:
			return nil, err
		}
	}

	match, err := argon2id.ComparePasswordAndHash(password, owner.PasswordHash)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, fmt.Errorf("owner %w", ErrInvalidCredentials)
	}

	ownerSession, err := a.sessionSvc.NewOwner(ctx, owner.ID)
	if err != nil {
		return nil, err
	}

	return ownerSession, nil
}
func (a *ownerAuth) Logout(ctx context.Context, sessionID string) error {
	return a.sessionSvc.DeleteOwner(ctx, sessionID)
}
func (a *ownerAuth) Register(ctx context.Context, username, email, password string) error {
	passwordHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return err
	}
	_, err = a.ownerRepo.Create(ctx, username, email, passwordHash)
	if err != nil {
		return err
	}

	return nil
}

type userAuth struct {
	userRepo   repository.UsersRepository
	sessionSvc *Session
}

func newUserAuth(
	userRepo repository.UsersRepository,
	session *Session,
) *userAuth {
	return &userAuth{userRepo, session}
}

func (a *userAuth) Login(ctx context.Context, userID int64, password string, users []repository.UserCacheValue) (*UserSession, error) {
	// Check UserID exist in Pharmacy
	if !UserExistsInPharmacy(users, userID) {
		return nil, ErrUserForbidden
	}

	// Get User
	user, err := a.userRepo.GetByID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			return nil, fmt.Errorf("%w: user doesn't exist", err)
		default:
			return nil, err
		}
	}

	// check password hash
	match, err := argon2id.ComparePasswordAndHash(password, user.PasswordHash)
	if err != nil {
		return nil, err
	}

	if !match {
		return nil, ErrInvalidCredentials
	}

	newExp := time.Now().Add(a.sessionSvc.ownerSessionTTL)
	newSession, err := a.sessionSvc.userSessionRepo.Create(ctx, userID, newExp)
	if err != nil {
		return nil, err
	}
	userCacheValue := repository.UserCacheValue{
		ID:       newSession.UserID,
		Username: user.Username,
	}
	a.sessionSvc.cacheStore.UserSessions.SetDefault(newSession.ID.String(), userCacheValue)

	userSession := &UserSession{
		ID:     newSession.ID.String(),
		Exp:    newSession.ExpiresAt.Time,
		UserID: newSession.UserID,
	}
	return userSession, nil
}
func (a *userAuth) Logout(sessionID string) error { return nil }

type pharmacyAuth struct{}

func (a *pharmacyAuth) Connect(code string) error { return nil }
