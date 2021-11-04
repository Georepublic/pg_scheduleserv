/*GRP-GNU-AGPL******************************************************************

File: database_util.go

Copyright (C) 2021  Team Georepublic <info@georepublic.de>

Developer(s):
Copyright (C) 2021  Ashish Kumar <ashishkr23438@gmail.com>

-----

This file is part of pg_scheduleserv.

pg_scheduleserv is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

pg_scheduleserv is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with pg_scheduleserv.  If not, see <https://www.gnu.org/licenses/>.

******************************************************************GRP-GNU-AGPL*/

package e2etest

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
