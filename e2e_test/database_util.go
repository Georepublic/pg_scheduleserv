package test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/Georepublic/pg_scheduleserv/internal/config"
	"github.com/jackc/pgx/v4"
)

func clone() (string, error) {
	name, err := randomDatabaseName()
	if err != nil {
		return "", fmt.Errorf("failed to generate random database name: %w", err)
	}

	ctx := context.Background()
	databaseName := "scheduler_test"
	q := fmt.Sprintf(`CREATE DATABASE "%s" WITH TEMPLATE "%s";`, name, databaseName)

	config, err := config.LoadConfig("../..")
	conn, err := pgx.Connect(context.Background(), config.DatabaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	if _, err := conn.Exec(ctx, q); err != nil {
		return "", fmt.Errorf("failed to clone template database: %w", err)
	}
	return name, nil
}

// NewDatabase creates a new database suitable for use in testing. It returns an
// established database connection and the configuration.
func NewDatabase(t *testing.T) string {
	newDatabaseName, err := clone()
	if err != nil {
		t.Fatal(err)
	}
	connectionURL := fmt.Sprintf("postgres://postgres:password@localhost:5432/%s", newDatabaseName)

	t.Cleanup(func() {
		ctx := context.Background()

		q := fmt.Sprintf(`DROP DATABASE IF EXISTS "%s" WITH (FORCE);`, newDatabaseName)

		config, err := config.LoadConfig("../..")
		conn, err := pgx.Connect(context.Background(), config.DatabaseURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			os.Exit(1)
		}

		if _, err := conn.Exec(ctx, q); err != nil {
			t.Errorf("failed to drop database %q: %s", newDatabaseName, err)
		}
	})

	return connectionURL
}

func randomDatabaseName() (string, error) {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
