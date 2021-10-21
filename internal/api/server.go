/*GRP-GNU-AGPL******************************************************************

File: server.go

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

package api

import (
	"net/http"

	"github.com/Georepublic/pg_scheduleserv/internal/database"
	"github.com/Georepublic/pg_scheduleserv/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

type Server struct {
	conn     *pgx.Conn
	router   *mux.Router
	validate *validator.Validate
	*database.Queries
	*util.Formatter
}

func NewServer(conn *pgx.Conn) *Server {
	router := mux.NewRouter().StrictSlash(true)

	server := &Server{
		conn:      conn,
		router:    router,
		validate:  validator.New(),
		Queries:   database.New(conn),
		Formatter: util.NewFormatter(),
	}

	router.HandleFunc("/projects/{project_id}/jobs", server.createJob).Methods("POST")
	router.HandleFunc("/projects/", server.createProject).Methods("POST")

	return server
}

func (server *Server) Start(port string) error {
	logrus.Info("Serving requests on port", port)
	return http.ListenAndServe(port, util.Logger(server.router))
}
