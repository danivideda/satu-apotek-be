package store

import "github.com/danivideda/satu-apotek-be/internal/dbsqlc"

type RevokedTokenStore struct {
	queries *dbsqlc.Queries
}

func (s *RevokedTokenStore) RevokeToken() {

}