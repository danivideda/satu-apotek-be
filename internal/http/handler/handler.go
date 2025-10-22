package handler

import "github.com/danivideda/satu-apotek-be/internal/dbsqlc"

type Handler struct {
	store dbsqlc.Querier
}

func New(db dbsqlc.DBTX) Handler {
	return Handler{
		store: dbsqlc.New(db),
	}
}
