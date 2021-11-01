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
	"reflect"
	"strings"

	"github.com/Georepublic/pg_scheduleserv/internal/database"
	"github.com/Georepublic/pg_scheduleserv/internal/util"
	"github.com/go-openapi/runtime/middleware"
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

	validate := validator.New()
	// Get json tag name instead of actual struct field name
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	server := &Server{
		conn:      conn,
		router:    router,
		validate:  validate,
		Queries:   database.New(conn),
		Formatter: util.NewFormatter(),
	}

	server.handleRoutes(router)
	serveSwagger(router)
	return server
}

func (server *Server) Start(port string) error {
	logrus.Info("Serving requests on port", port)
	return http.ListenAndServe(port, util.Logger(server.router))
}

func (server *Server) handleRoutes(router *mux.Router) {
	// Use URLs without trailing slash

	// Projects endpoints
	router.HandleFunc("/projects", server.createProject).Methods("POST")
	router.HandleFunc("/projects", server.listProjects).Methods("GET")
	router.HandleFunc("/projects/{project_id}", server.getProject).Methods("GET")
	router.HandleFunc("/projects/{project_id}", server.updateProject).Methods("PATCH")
	router.HandleFunc("/projects/{project_id}", server.deleteProject).Methods("DELETE")

	// Job endpoints
	router.HandleFunc("/projects/{project_id}/jobs", server.createJob).Methods("POST")
	router.HandleFunc("/projects/{project_id}/jobs", server.listJobs).Methods("GET")
	router.HandleFunc("/jobs/{job_id}", server.getJob).Methods("GET")
	router.HandleFunc("/jobs/{job_id}", server.updateJob).Methods("PATCH")
	router.HandleFunc("/jobs/{job_id}", server.deleteJob).Methods("DELETE")

	// Job time windows endpoints
	router.HandleFunc("/jobs/{job_id}/time_windows", server.createJobTimeWindow).Methods("POST")
	router.HandleFunc("/jobs/{job_id}/time_windows", server.listJobTimeWindows).Methods("GET")
	router.HandleFunc("/jobs/{job_id}/time_windows/{tw_open}/{tw_close}", server.deleteJobTimeWindow).Methods("DELETE")

	// Shipment endpoints
	router.HandleFunc("/projects/{project_id}/shipments", server.createShipment).Methods("POST")
	router.HandleFunc("/projects/{project_id}/shipments", server.listShipments).Methods("GET")
	router.HandleFunc("/shipments/{shipment_id}", server.getShipment).Methods("GET")
	router.HandleFunc("/shipments/{shipment_id}", server.updateShipment).Methods("PATCH")
	router.HandleFunc("/shipments/{shipment_id}", server.deleteShipment).Methods("DELETE")

	// Shipment time windows endpoints
	router.HandleFunc("/shipments/{shipment_id}/time_windows", server.createShipmentTimeWindow).Methods("POST")
	router.HandleFunc("/shipments/{shipment_id}/time_windows", server.listShipmentTimeWindows).Methods("GET")
	router.HandleFunc("/shipments/{shipment_id}/time_windows/{kind}/{tw_open}/{tw_close}", server.deleteShipmentTimeWindow).Methods("DELETE")

	// Vehicle endpoints
	router.HandleFunc("/projects/{project_id}/vehicles", server.createVehicle).Methods("POST")
	router.HandleFunc("/projects/{project_id}/vehicles", server.listVehicles).Methods("GET")
	router.HandleFunc("/vehicles/{vehicle_id}", server.getVehicle).Methods("GET")
	router.HandleFunc("/vehicles/{vehicle_id}", server.updateVehicle).Methods("PATCH")
	router.HandleFunc("/vehicles/{vehicle_id}", server.deleteVehicle).Methods("DELETE")

	// Vehicle breaks endpoints
	router.HandleFunc("/vehicles/{vehicle_id}/breaks", server.createBreak).Methods("POST")
	router.HandleFunc("/vehicles/{vehicle_id}/breaks", server.listBreaks).Methods("GET")
	router.HandleFunc("/breaks/{vehicle_id}", server.getBreak).Methods("GET")
	router.HandleFunc("/breaks/{vehicle_id}", server.updateBreak).Methods("PATCH")
	router.HandleFunc("/breaks/{vehicle_id}", server.deleteBreak).Methods("DELETE")

	// Break time windows endpoints
	router.HandleFunc("/breaks/{break_id}/time_windows", server.createBreakTimeWindow).Methods("POST")
	router.HandleFunc("/breaks/{break_id}/time_windows", server.listBreakTimeWindows).Methods("GET")
	router.HandleFunc("/breaks/{break_id}/time_windows/{tw_open}/{tw_close}", server.deleteBreakTimeWindow).Methods("DELETE")
}

func serveSwagger(router *mux.Router) {
	redocOpts := middleware.RedocOpts{
		SpecURL: "./swagger.yaml",
		Path:    "/redoc",
	}
	redoc := middleware.Redoc(redocOpts, nil)
	swaggerOpts := middleware.SwaggerUIOpts{
		SpecURL: "./swagger.yaml",
		Path:    "/",
	}
	swagger := middleware.SwaggerUI(swaggerOpts, nil)

	router.Handle("/swagger.yaml", http.FileServer(http.Dir("./docs/")))
	router.Handle("/redoc", redoc)
	router.Handle("/", swagger)
}
