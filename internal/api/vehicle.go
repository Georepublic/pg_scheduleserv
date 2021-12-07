/*GRP-GNU-AGPL******************************************************************

File: vehicle.go

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

// CreateVehicles godoc
// @Summary Create a new vehicle
// @Description Create a new vehicle with the input payload
// @Tags Vehicle
// @Accept application/json
// @Produce application/json
// @Param project_id path int true "Project ID"
// @Param Vehicle body database.CreateVehicleParams true "Create vehicle"
// @Success 200 {object} database.Vehicle
// @Failure 400 {object} util.MultiError
// @Router /projects/{project_id}/vehicles [post]
func (server *Server) CreateVehicle(w http.ResponseWriter, r *http.Request) {
	userInput := make(map[string]interface{})
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
			logrus.Error(err)
		}
	}

	// Add the project_id path variable
	vars := mux.Vars(r)
	userInput["project_id"] = vars["project_id"]

	logrus.Debugf("%v", userInput)

	// Validate the input type
	if err := util.ValidateInput(userInput, database.CreateVehicleParams{}); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	// Decode map[string]interface{} to struct
	userInputString, err := json.Marshal(userInput)
	if err != nil {
		logrus.Error(err)
	}
	vehicle := database.CreateVehicleParams{}
	if err = json.Unmarshal(userInputString, &vehicle); err != nil {
		logrus.Error(err)
	}

	logrus.Debugf("%+v", vehicle)

	// Validate the struct
	if err := server.validate.Struct(vehicle); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	created_vehicle, err := server.DBCreateVehicle(ctx, vehicle)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusCreated, created_vehicle)
}

// ListVehicles godoc
// @Summary List vehicles for a project
// @Description Get a list of vehicles for a project with project_id
// @Tags Vehicle
// @Accept application/json
// @Produce application/json
// @Param project_id path int true "Project ID"
// @Success 200 {object} []database.Vehicle
// @Failure 400 {object} util.MultiError
// @Router /projects/{project_id}/vehicles [get]
func (server *Server) ListVehicles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	project_id, err := strconv.ParseInt(vars["project_id"], 10, 64)
	if err != nil {
		panic(err)
	}

	ctx := r.Context()
	created_vehicle, err := server.DBListVehicles(ctx, project_id)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusOK, created_vehicle)
}

// GetVehicle godoc
// @Summary Fetch a vehicle
// @Description Fetch a vehicle with its vehicle_id
// @Tags Vehicle
// @Accept application/json
// @Produce application/json
// @Param vehicle_id path int true "Vehicle ID"
// @Success 200 {object} database.Vehicle
// @Failure 400 {object} util.MultiError
// @Failure 404 {object} util.NotFound
// @Router /vehicles/{vehicle_id} [get]
func (server *Server) GetVehicle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vehicle_id, err := strconv.ParseInt(vars["vehicle_id"], 10, 64)
	if err != nil {
		panic(err)
	}

	ctx := r.Context()
	created_vehicle, err := server.DBGetVehicle(ctx, vehicle_id)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusOK, created_vehicle)
}

// UpdateVehicle godoc
// @Summary Update a vehicle
// @Description Update a vehicle with its vehicle_id
// @Tags Vehicle
// @Accept application/json
// @Produce application/json
// @Param vehicle_id path int true "Vehicle ID"
// @Param Vehicle body database.UpdateVehicleParams true "Update vehicle"
// @Success 200 {object} database.Vehicle
// @Failure 400 {object} util.MultiError
// @Failure 404 {object} util.NotFound
// @Router /vehicles/{vehicle_id} [patch]
func (server *Server) UpdateVehicle(w http.ResponseWriter, r *http.Request) {
	userInput := make(map[string]interface{})
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
			logrus.Error(err)
		}
	}

	vars := mux.Vars(r)
	vehicle_id, err := strconv.ParseInt(vars["vehicle_id"], 10, 64)
	if err != nil {
		panic(err)
	}

	// Validate the input type
	if err := util.ValidateInput(userInput, database.UpdateVehicleParams{}); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	// Decode map[string]interface{} to struct
	userInputString, err := json.Marshal(userInput)
	if err != nil {
		logrus.Error(err)
	}
	vehicle := database.UpdateVehicleParams{}
	if err = json.Unmarshal(userInputString, &vehicle); err != nil {
		logrus.Error(err)
	}

	// Validate the struct
	if err := server.validate.Struct(vehicle); err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()
	created_vehicle, err := server.DBUpdateVehicle(ctx, vehicle, vehicle_id)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusOK, created_vehicle)
}

// DeleteVehicle godoc
// @Summary Delete a vehicle
// @Description Delete a vehicle with its vehicle_id
// @Tags Vehicle
// @Accept application/json
// @Produce application/json
// @Param vehicle_id path int true "Vehicle ID"
// @Success 200 {object} util.Success
// @Failure 400 {object} util.MultiError
// @Failure 404 {object} util.NotFound
// @Router /vehicles/{vehicle_id} [delete]
func (server *Server) DeleteVehicle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vehicle_id, err := strconv.ParseInt(vars["vehicle_id"], 10, 64)
	if err != nil {
		panic(err)
	}

	ctx := r.Context()
	_, err = server.DBDeleteVehicle(ctx, vehicle_id)
	if err != nil {
		server.FormatJSON(w, http.StatusBadRequest, err)
		return
	}

	server.FormatJSON(w, http.StatusOK, nil)
}
