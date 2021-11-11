/*GRP-GNU-AGPL******************************************************************

File: util.go

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
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/Georepublic/pg_scheduleserv/internal/api"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

func applyTestData(db_url string, filename string) {
	filepath := fmt.Sprintf("../scripts/%s", filename)
	cmd := exec.Command("psql", db_url, "-f", filepath)

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error executing query: %s", err)
	}
}

func setup(db_url string, filename string) (*api.Server, *pgx.Conn) {
	conn, err := pgx.Connect(context.Background(), db_url)
	if err != nil {
		logrus.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	server := api.NewServer(conn)
	logrus.Error(db_url)
	m, err := migrate.New("file://../migrations/", db_url)
	if err != nil {
		logrus.Printf("Unable to apply the migrations: %v\n", err)
		os.Exit(1)
	}
	m.Up()
	applyTestData(db_url, filename)
	return server, conn
}
