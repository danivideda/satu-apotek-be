package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres(addr string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(addr)
	if err != nil {
		return nil, err
	}

	dbpool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	if err := dbpool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return dbpool, nil
}