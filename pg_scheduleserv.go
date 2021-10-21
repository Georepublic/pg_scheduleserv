/*GRP-GNU-AGPL******************************************************************

File: pg_scheduleserv.go

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

package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/Georepublic/pg_scheduleserv/internal/api"
	"github.com/Georepublic/pg_scheduleserv/internal/config"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05 -0700",
		FullTimestamp:   true,
		CallerPrettyfier: func(frame *runtime.Frame) (string, string) {
			s := strings.Split(frame.Function, ".")
			funcname := s[len(s)-1] + "()"
			s = strings.Split(frame.File, "pg_scheduleserv/")
			filename := s[len(s)-1] + ":" + strconv.Itoa(frame.Line)
			return funcname, filename
		},
	})
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
	config, err := config.LoadConfig(".")
	if err != nil {
		logrus.Error("Cannot load config:", err)
	}
	conn, err := pgx.Connect(context.Background(), config.DatabaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	server := api.NewServer(conn)
	logrus.Error(server.Start(config.ServerBindAddress))
}
