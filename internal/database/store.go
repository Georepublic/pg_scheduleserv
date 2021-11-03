package database

import (
	"github.com/jackc/pgx/v4"
)

type Store struct {
	Querier
	db *pgx.Conn
}

func NewStore(db *pgx.Conn) *Store {
	return &Store{
		db:      db,
		Querier: New(db),
	}
}
