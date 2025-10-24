package store

import (
	"context"
	"log"

	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
)

func NewMockStorage() Storage {
	return Storage{}
}

type MockOwnerStore struct {
	mockData string
}

func (s *MockOwnerStore) GetOwnerByID(context.Context, int) (*dbsqlc.Owner, error) {
	log.Default().Print(s.mockData)
	return nil, nil
}

func (s *MockOwnerStore) CreateOwner(context.Context, dbsqlc.CreateOwnerParams) (*dbsqlc.CreateOwnerRow, error) {
	log.Default().Print(s.mockData)
	return nil, nil
}
