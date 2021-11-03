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
// @Router /projects/{project_id}/shipments [post]
func (server *Server) CreateShipment(w http.ResponseWriter, r *http.Request) {
	userInput := make(map[string]interface{})
	if r.Body != nil {
		json.NewDecoder(r.Body).Decode(&userInput)
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
	json.Unmarshal(userInputString, &shipment)

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
// @Success 200 {object} database.Shipment
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

	server.FormatJSON(w, http.StatusCreated, created_shipment)
}

// GetShipment godoc
// @Summary Fetch a shipment
// @Description Fetch a shipment with its shipment_id
// @Tags Shipment
// @Accept application/json
// @Produce application/json
// @Param shipment_id path int true "Shipment ID"
// @Success 200 {object} database.Shipment
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

	server.FormatJSON(w, http.StatusCreated, created_shipment)
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
// @Router /shipments/{shipment_id} [patch]
func (server *Server) UpdateShipment(w http.ResponseWriter, r *http.Request) {
	userInput := make(map[string]interface{})
	if r.Body != nil {
		json.NewDecoder(r.Body).Decode(&userInput)
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
	json.Unmarshal(userInputString, &shipment)

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

	server.FormatJSON(w, http.StatusCreated, created_shipment)
}

// DeleteShipment godoc
// @Summary Delete a shipment
// @Description Delete a shipment with its shipment_id
// @Tags Shipment
// @Accept application/json
// @Produce application/json
// @Param shipment_id path int true "Shipment ID"
// @Success 200 {object} database.Shipment
// @Router /shipments/{shipment_id} [delete]
func (server *Server) DeleteShipment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shipment_id, err := strconv.ParseInt(vars["shipment_id"], 10, 64)
	if err != nil {
		panic(err)
	}

	ctx := r.Context()
	err = server.DBDeleteShipment(ctx, shipment_id)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusCreated, nil)
}
