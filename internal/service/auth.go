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

func NewAuth(repo repository.Repository, session *Session) *Auth {
	return &Auth{
		Owner:    newOwnerAuth(repo.Owners, session),
		User:     newUserAuth(repo.Users, session),
		Pharmacy: newPharmacyAuth(repo.Pharmacies, repo.Users, session),
	}
}

type OwnerAuth interface {
	Login(ctx context.Context, email, password string) (*OwnerSession, error)
	Logout(ctx context.Context, sessionID string) error
	Register(ctx context.Context, username, email, password string) error
}
type UserAuth interface {
	Login(ctx context.Context, userID int64, password string, users []repository.UserCacheValue) (*UserSession, error)
	Logout(ctx context.Context, sessionID string) error
}
type PharmacyAuth interface {
	Connect(ctx context.Context, code string) (*PharmacySession, error)
}

type ownerAuth struct {
	ownerRepo repository.OwnersRepository
	session   *Session
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

	return a.session.NewOwner(ctx, owner.ID)
}
func (a *ownerAuth) Logout(ctx context.Context, sessionID string) error {
	return a.session.DeleteOwner(ctx, sessionID)
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
	userRepo repository.UsersRepository
	session  *Session
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

	return a.session.NewUser(ctx, user.ID, user.Username)
}
func (a *userAuth) Logout(ctx context.Context, sessionID string) error {
	return a.session.DeleteUser(ctx, sessionID)
}

type pharmacyAuth struct {
	pharmaciesRepo repository.PharmaciesRepository
	usersRepo      repository.UsersRepository
	session        *Session
}

func newPharmacyAuth(
	pharmaciesRepo repository.PharmaciesRepository,
	usersRepo repository.UsersRepository,
	session *Session,
) *pharmacyAuth {
	return &pharmacyAuth{pharmaciesRepo, usersRepo, session}
}

func (a *pharmacyAuth) Connect(ctx context.Context, code string) (*PharmacySession, error) {
	// Check the code
	pharmacyCode, err := a.pharmaciesRepo.GetCodeByCode(ctx, code)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			return nil, ErrPharmacyCodeExpiredOrNotFound
		default:
			return nil, err
		}
	}

	//--check if code expired
	if time.Now().After(pharmacyCode.ExpiresAt.Time) {
		return nil, ErrPharmacyCodeExpiredOrNotFound
	}

	// Delete immediately after being used
	// doesn't matter if the session creation later have some error
	if _, err := a.pharmaciesRepo.DeleteCode(ctx, pharmacyCode.Code); err != nil {
		return nil, err
	}

	newPharmacySession, err := a.session.NewPharmacy(ctx, pharmacyCode.ApotekID)
	if err != nil {
		return nil, err
	}

	return newPharmacySession, nil
}
