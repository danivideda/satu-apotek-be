package service

import "github.com/danivideda/satu-apotek-be/internal/repository"

type ownerService struct {
	ownersRepo repository.OwnersRepository
}

func NewOwnerService(ownersRepo repository.OwnersRepository) *ownerService {
	return &ownerService{
		ownersRepo: ownersRepo,
	}
}

func (s *ownerService) Create() error {
	return nil
}
