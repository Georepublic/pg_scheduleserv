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
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

type Server struct {
	conn     *pgx.Conn
	Router   *mux.Router
	validate *validator.Validate
	*database.Store
	*util.Formatter
}

func NewServer(conn *pgx.Conn) *Server {
	router := mux.NewRouter().StrictSlash(true)
	server := &Server{
		conn:      conn,
		Router:    router,
		validate:  util.NewValidator(),
		Store:     database.NewStore(conn),
		Formatter: util.NewFormatter(),
	}

	server.handleRoutes(router)
	serveSwagger(router)
	return server
}

func (server *Server) Start(port string) error {
	logrus.Info("Serving requests on port", port)
	return http.ListenAndServe(port, util.Logger(server.Router))
}

func (server *Server) handleRoutes(router *mux.Router) {
	// Use URLs without trailing slash

	// Projects endpoints
	router.HandleFunc("/projects", server.CreateProject).Methods("POST")
	router.HandleFunc("/projects", server.ListProjects).Methods("GET")
	router.HandleFunc("/projects/{project_id}", server.GetProject).Methods("GET")
	router.HandleFunc("/projects/{project_id}", server.UpdateProject).Methods("PATCH")
	router.HandleFunc("/projects/{project_id}", server.DeleteProject).Methods("DELETE")

	// Schedule related endpoints
	router.HandleFunc("/projects/{project_id}/schedule", server.GetSchedule).Methods("GET")
	router.HandleFunc("/projects/{project_id}/schedule", server.CreateSchedule).Methods("POST")
	router.HandleFunc("/projects/{project_id}/schedule", server.DeleteSchedule).Methods("DELETE")

	// Job endpoints
	router.HandleFunc("/projects/{project_id}/jobs", server.CreateJob).Methods("POST")
	router.HandleFunc("/projects/{project_id}/jobs", server.ListJobs).Methods("GET")
	router.HandleFunc("/jobs/{job_id}", server.GetJob).Methods("GET")
	router.HandleFunc("/jobs/{job_id}", server.UpdateJob).Methods("PATCH")
	router.HandleFunc("/jobs/{job_id}", server.DeleteJob).Methods("DELETE")
	router.HandleFunc("/jobs/{job_id}/schedule", server.GetJobSchedule).Methods("GET")

	// Job time windows endpoints
	router.HandleFunc("/jobs/{job_id}/time_windows", server.CreateJobTimeWindow).Methods("POST")
	router.HandleFunc("/jobs/{job_id}/time_windows", server.ListJobTimeWindows).Methods("GET")
	router.HandleFunc("/jobs/{job_id}/time_windows", server.DeleteJobTimeWindow).Methods("DELETE")

	// Shipment endpoints
	router.HandleFunc("/projects/{project_id}/shipments", server.CreateShipment).Methods("POST")
	router.HandleFunc("/projects/{project_id}/shipments", server.ListShipments).Methods("GET")
	router.HandleFunc("/shipments/{shipment_id}", server.GetShipment).Methods("GET")
	router.HandleFunc("/shipments/{shipment_id}", server.UpdateShipment).Methods("PATCH")
	router.HandleFunc("/shipments/{shipment_id}", server.DeleteShipment).Methods("DELETE")
	router.HandleFunc("/shipments/{shipment_id}/schedule", server.GetShipmentSchedule).Methods("GET")

	// Shipment time windows endpoints
	router.HandleFunc("/shipments/{shipment_id}/time_windows", server.CreateShipmentTimeWindow).Methods("POST")
	router.HandleFunc("/shipments/{shipment_id}/time_windows", server.ListShipmentTimeWindows).Methods("GET")
	router.HandleFunc("/shipments/{shipment_id}/time_windows", server.DeleteShipmentTimeWindow).Methods("DELETE")

	// Vehicle endpoints
	router.HandleFunc("/projects/{project_id}/vehicles", server.CreateVehicle).Methods("POST")
	router.HandleFunc("/projects/{project_id}/vehicles", server.ListVehicles).Methods("GET")
	router.HandleFunc("/vehicles/{vehicle_id}", server.GetVehicle).Methods("GET")
	router.HandleFunc("/vehicles/{vehicle_id}", server.UpdateVehicle).Methods("PATCH")
	router.HandleFunc("/vehicles/{vehicle_id}", server.DeleteVehicle).Methods("DELETE")
	router.HandleFunc("/vehicles/{vehicle_id}/schedule", server.GetVehicleSchedule).Methods("GET")

	// Vehicle breaks endpoints
	router.HandleFunc("/vehicles/{vehicle_id}/breaks", server.CreateBreak).Methods("POST")
	router.HandleFunc("/vehicles/{vehicle_id}/breaks", server.ListBreaks).Methods("GET")
	router.HandleFunc("/breaks/{break_id}", server.GetBreak).Methods("GET")
	router.HandleFunc("/breaks/{break_id}", server.UpdateBreak).Methods("PATCH")
	router.HandleFunc("/breaks/{break_id}", server.DeleteBreak).Methods("DELETE")

	// Break time windows endpoints
	router.HandleFunc("/breaks/{break_id}/time_windows", server.CreateBreakTimeWindow).Methods("POST")
	router.HandleFunc("/breaks/{break_id}/time_windows", server.ListBreakTimeWindows).Methods("GET")
	router.HandleFunc("/breaks/{break_id}/time_windows", server.DeleteBreakTimeWindow).Methods("DELETE")
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
