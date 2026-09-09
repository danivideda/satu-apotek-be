package service

import (
	"context"
	"errors"
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
		User:     &userAuth{},
		Pharmacy: &pharmacyAuth{},
	}
}

type OwnerAuth interface {
	Login(ctx context.Context, email, password string) (ownerSessionID string, expiresAt time.Time, err error)
	Logout(ctx context.Context, sessionID string) error
	Register(ctx context.Context, username, email, password string) error
}
type UserAuth interface {
	Login(username, password string) (sessionID string, err error)
	Logout(sessionID string) error
}
type PharmacyAuth interface {
	Connect(code string) error
}

type ownerAuth struct {
	ownerRepo        repository.OwnersRepository
	sessionSvc       *Session
}

func newOwnerAuth(
	ownerRepo repository.OwnersRepository,
	session *Session,
) *ownerAuth {
	return &ownerAuth{ownerRepo, session}
}

func (a *ownerAuth) Login(ctx context.Context, email, password string) (ownerSessionID string, expiresAt time.Time, err error) {
	owner, err := a.ownerRepo.GetByEmail(ctx, email)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			return "", time.Time{}, ErrInvalidCredentials
		default:
			return "", time.Time{}, err
		}
	}

	match, err := argon2id.ComparePasswordAndHash(password, owner.PasswordHash)
	if err != nil {
		return "", time.Time{}, err
	}
	if !match {
		return "", time.Time{}, ErrInvalidCredentials
	}

	ownerSessionID, expiresAt, err = a.sessionSvc.NewOwner(ctx, owner.ID)
	if err != nil {
		return "", time.Time{}, err
	}

	return ownerSessionID, expiresAt, nil
}
func (a *ownerAuth) Logout(ctx context.Context, sessionID string) error {

	return nil
}
func (a *ownerAuth) Register(ctx context.Context, username, email, password string) error {
	return nil
}

type userAuth struct{}

func (a *userAuth) Login(username, password string) (string, error) { return "", nil }
func (a *userAuth) Logout(session_id string) error                  { return nil }

type pharmacyAuth struct{}

func (a *pharmacyAuth) Connect(code string) error { return nil }
