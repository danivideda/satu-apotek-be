package service

import "github.com/danivideda/satu-apotek-be/internal/repository"

type OwnerService struct {
	ownersRepo repository.OwnersRepository
}

func NewOwnerService(ownersRepo repository.OwnersRepository) *OwnerService {
	return &OwnerService{
		ownersRepo: ownersRepo,
	}
}
