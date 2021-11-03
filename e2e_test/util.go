package test

import (
	"context"
	"fmt"
	"os"

	"github.com/Georepublic/pg_scheduleserv/internal/api"
	"github.com/jackc/pgx/v4"
)

func setup(db_url string) (*api.Server, *pgx.Conn) {
	conn, err := pgx.Connect(context.Background(), db_url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	server := api.NewServer(conn)
	return server, conn
}
