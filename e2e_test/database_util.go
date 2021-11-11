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
	"github.com/sirupsen/logrus"
)

func clone() (string, error) {
	name, err := randomDatabaseName()
	if err != nil {
		return "", fmt.Errorf("failed to generate random database name: %w", err)
	}

	ctx := context.Background()
	q := fmt.Sprintf(`CREATE DATABASE "%s"`, name)

	config, err := config.LoadConfig("..")
	if err != nil {
		logrus.Errorf("Cannot load config: %s", err)
	}
	connectionURL := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		config.DatabaseUser,
		config.DatabasePassword,
		config.DatabaseHost,
		config.DatabasePort,
		config.DatabaseName,
	)
	conn, err := pgx.Connect(context.Background(), connectionURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	if _, err := conn.Exec(ctx, q); err != nil {
		return "", fmt.Errorf("failed to clone template database: %w", err)
	}
	return name, nil
}

// NewTestDatabase creates a new database suitable for use in testing. It returns an
// established database connection and the configuration.
func NewTestDatabase(t *testing.T) string {
	newDatabaseName, err := clone()
	if err != nil {
		t.Fatal(err)
	}
	config, err := config.LoadConfig("..")
	if err != nil {
		logrus.Fatalf("Cannot load config: %s", err)
	}
	testConnectionURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.DatabaseUser,
		config.DatabasePassword,
		config.DatabaseHost,
		config.DatabasePort,
		newDatabaseName,
	)

	t.Cleanup(func() {
		ctx := context.Background()

		q := fmt.Sprintf(`DROP DATABASE IF EXISTS "%s" WITH (FORCE);`, newDatabaseName)

		connectionURL := fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			config.DatabaseUser,
			config.DatabasePassword,
			config.DatabaseHost,
			config.DatabasePort,
			config.DatabaseName,
		)
		conn, err := pgx.Connect(context.Background(), connectionURL)
		if err != nil {
			t.Errorf("Unable to connect to database: %v\n", err)
			os.Exit(1)
		}

		if _, err := conn.Exec(ctx, q); err != nil {
			t.Errorf("failed to drop database %q: %s", newDatabaseName, err)
		}
	})

	return testConnectionURL
}

func randomDatabaseName() (string, error) {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
