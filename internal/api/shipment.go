/*GRP-GNU-AGPL******************************************************************

File: shipment.go

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
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Georepublic/pg_scheduleserv/internal/database"
	"github.com/Georepublic/pg_scheduleserv/internal/util"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// CreateShipments godoc
// @Summary Create a new shipment
// @Description Create a new shipment with the input payload
// @Tags Shipment
// @Accept application/json
// @Produce application/json
// @Param project_id path int true "Project ID"
// @Param Shipment body database.CreateShipmentParams true "Create shipment"
// @Success 200 {object} database.Shipment
// @Failure 400 {object} util.MultiError
// @Router /projects/{project_id}/shipments [post]
func (server *Server) CreateShipment(w http.ResponseWriter, r *http.Request) {
	userInput := make(map[string]interface{})
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
			logrus.Error(err)
		}
	}

	// Add the project_id path variable
	vars := mux.Vars(r)
	userInput["project_id"] = vars["project_id"]

	// Validate the input type
	if err := util.ValidateInput(userInput, database.CreateShipmentParams{}); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	// Decode map[string]interface{} to struct
	userInputString, err := json.Marshal(userInput)
	if err != nil {
		logrus.Error(err)
	}
	shipment := database.CreateShipmentParams{}
	if err = json.Unmarshal(userInputString, &shipment); err != nil {
		logrus.Error(err)
	}

	logrus.Debugf("%v", userInput)
	logrus.Debugf("%+v", shipment)

	// Validate the struct
	if err := server.validate.Struct(shipment); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	created_shipment, err := server.DBCreateShipment(ctx, shipment)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusCreated, created_shipment)
}

// ListShipments godoc
// @Summary List shipments for a project
// @Description Get a list of shipments for a project with project_id
// @Tags Shipment
// @Accept application/json
// @Produce application/json
// @Param project_id path int true "Project ID"
// @Success 200 {object} []database.Shipment
// @Failure 400 {object} util.MultiError
// @Router /projects/{project_id}/shipments [get]
func (server *Server) ListShipments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	project_id, err := strconv.ParseInt(vars["project_id"], 10, 64)
	if err != nil {
		panic(err)
	}

	ctx := r.Context()
	created_shipment, err := server.DBListShipments(ctx, project_id)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusOK, created_shipment)
}

// GetShipment godoc
// @Summary Fetch a shipment
// @Description Fetch a shipment with its shipment_id
// @Tags Shipment
// @Accept application/json
// @Produce application/json
// @Param shipment_id path int true "Shipment ID"
// @Success 200 {object} database.Shipment
// @Failure 400 {object} util.MultiError
// @Failure 404 {object} util.NotFound
// @Router /shipments/{shipment_id} [get]
func (server *Server) GetShipment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shipment_id, err := strconv.ParseInt(vars["shipment_id"], 10, 64)
	if err != nil {
		panic(err)
	}

	ctx := r.Context()
	created_shipment, err := server.DBGetShipment(ctx, shipment_id)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusOK, created_shipment)
}

// UpdateShipment godoc
// @Summary Update a shipment
// @Description Update a shipment with its shipment_id
// @Tags Shipment
// @Accept application/json
// @Produce application/json
// @Param shipment_id path int true "Shipment ID"
// @Param Shipment body database.UpdateShipmentParams true "Update shipment"
// @Success 200 {object} database.Shipment
// @Failure 400 {object} util.MultiError
// @Failure 404 {object} util.NotFound
// @Router /shipments/{shipment_id} [patch]
func (server *Server) UpdateShipment(w http.ResponseWriter, r *http.Request) {
	userInput := make(map[string]interface{})
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
			logrus.Error(err)
		}
	}

	vars := mux.Vars(r)
	shipment_id, err := strconv.ParseInt(vars["shipment_id"], 10, 64)
	if err != nil {
		panic(err)
	}

	// Validate the input type
	if err := util.ValidateInput(userInput, database.UpdateShipmentParams{}); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	// Decode map[string]interface{} to struct
	userInputString, err := json.Marshal(userInput)
	if err != nil {
		logrus.Error(err)
	}
	shipment := database.UpdateShipmentParams{}
	if err = json.Unmarshal(userInputString, &shipment); err != nil {
		logrus.Error(err)
	}

	// Validate the struct
	if err := server.validate.Struct(shipment); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	created_shipment, err := server.DBUpdateShipment(ctx, shipment, shipment_id)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusOK, created_shipment)
}

// DeleteShipment godoc
// @Summary Delete a shipment
// @Description Delete a shipment with its shipment_id
// @Tags Shipment
// @Accept application/json
// @Produce application/json
// @Param shipment_id path int true "Shipment ID"
// @Success 200 {object} util.Success
// @Failure 400 {object} util.MultiError
// @Failure 404 {object} util.NotFound
// @Router /shipments/{shipment_id} [delete]
func (server *Server) DeleteShipment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shipment_id, err := strconv.ParseInt(vars["shipment_id"], 10, 64)
	if err != nil {
		panic(err)
	}

	ctx := r.Context()
	_, err = server.DBDeleteShipment(ctx, shipment_id)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusOK, nil)
}

// GetShipmentSchedule godoc
// @Summary Get the schedule for a shipment
// @Description Get the schedule for a shipment using shipment_id
// @Tags Shipment
// @Accept application/json
// @Produce text/calendar,application/json
// @Param shipment_id path int true "Shipment ID"
// @Success 200 {object} util.Schedule
// @Failure 400 {object} util.MultiError
// @Failure 404 {object} util.NotFound
// @Router /shipments/{shipment_id}/schedule [get]
func (server *Server) GetShipmentSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shipmentID, err := strconv.ParseInt(vars["shipment_id"], 10, 64)
	if err != nil {
		panic(err)
	}

	ctx := r.Context()
	schedule, err := server.DBGetScheduleShipment(ctx, shipmentID)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	switch r.Header.Get("Accept") {
	case "text/calendar":
		server.FormatICAL(w, http.StatusOK, schedule)
	case "application/json":
		server.FormatJSON(w, http.StatusOK, schedule)
	default:
		server.FormatICAL(w, http.StatusOK, schedule)
	}

}
