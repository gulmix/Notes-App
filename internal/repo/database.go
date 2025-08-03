package repo

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PGRepo struct {
	Mu   sync.Mutex
	Pool *pgxpool.Pool
}

func New(connString string) (*PGRepo, error) {
	db, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	return &PGRepo{Mu: sync.Mutex{}, Pool: db}, nil
}
